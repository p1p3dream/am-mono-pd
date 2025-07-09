package worker

import (
	"path"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/zeebo/xxh3"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type EnsureDataFileObjectInput struct {
	DataFileType entities.DataFileType
	DirectoryId  *uuid.UUID
	FileSize     int64
	Meta         map[string]any
	ParentFileId *uuid.UUID
	Path         string
	Priorities   []int32
	WorkerId     *uuid.UUID
}

type EnsureDataFileObjectOutput struct {
	Entity *entities.DataFileObject
}

// EnsureDataFileObject checks if a DataFileObject exists in the database
// and returns it, or creates a new one otherwise.
func (dom *domain) EnsureDataFileObject(r *arc.Request, in *EnsureDataFileObjectInput) (*EnsureDataFileObjectOutput, error) {
	// Strip the leading slash from the path to ensure compatibility
	// with previous versions when matching the hash on the database.
	hash := val.ByteArray16ToSlice(
		xxh3.HashString128(strings.TrimPrefix(in.Path, "/")).Bytes(),
	)

	selectObjectOut, err := dom.SelectDataFileObject(r, &SelectDataFileObjectInput{
		DirectoryId:  in.DirectoryId,
		ParentFileId: in.ParentFileId,
		FileType:     in.DataFileType,
		Hash:         hash,
		Meta:         in.Meta,
	})
	if err != nil {
		return nil, errors.Forward(err, "1e7bc1e7-deaa-4bcf-a733-515aa0522308")
	}

	out := &EnsureDataFileObjectOutput{}

	if selectObjectOut.Entity != nil {
		out.Entity = selectObjectOut.Entity
		return out, nil
	}

	id, err := val.NewUUID7()
	if err != nil {
		return nil, errors.Forward(err, "8e7f54c7-1245-4049-94d6-5af3b4f7706d")
	}

	now := time.Now()

	insertObjectOut, err := dom.InsertDataFileObject(r, &InsertDataFileObjectInput{
		Entity: &entities.DataFileObject{
			Id:           id,
			CreatedAt:    now,
			UpdatedAt:    now,
			Meta:         in.Meta,
			DirectoryId:  in.DirectoryId,
			ParentFileId: in.ParentFileId,
			FileType:     in.DataFileType,
			Hash:         hash,
			Status:       entities.DataFileObjectStatusToDo,
			FileDir:      path.Dir(in.Path),
			FileName:     path.Base(in.Path),
			FileSize:     &in.FileSize,
			Priorities:   in.Priorities,
			WorkerId:     in.WorkerId,
		},
		ReturnEntity: true,
	})
	if err != nil {
		return nil, errors.Forward(err, "15d7f140-7f3d-46b9-ae9f-17d14cd3d095")
	}

	out.Entity = insertObjectOut.Entity

	return out, nil
}

type InsertDataFileObjectInput struct {
	Entity       *entities.DataFileObject
	ReturnEntity bool
}

type InsertDataFileObjectOutput struct {
	Entity *entities.DataFileObject
}

func (dom *domain) InsertDataFileObject(r *arc.Request, in *InsertDataFileObjectInput) (*InsertDataFileObjectOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "9f01b149-c492-4344-8dec-a9766d6b6de8",
			Code:   errors.Code_INTERNAL,
			Detail: "Missing input.",
		}
	}

	obj := in.Entity

	if obj == nil {
		return nil, &errors.Object{
			Id:     "4cfe9efd-080d-49eb-be48-96162c8164ae",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing DataFileObject.",
		}
	}

	if len(obj.Hash) == 0 {
		return nil, &errors.Object{
			Id:     "6f8acc2c-567b-4b44-86bd-20ba59ff2849",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Hash.",
		}
	}

	if obj.FileType == 0 {
		return nil, &errors.Object{
			Id:     "a6b72469-9df4-4127-b176-68ee86da731e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing FileType.",
		}
	}

	if obj.Status == 0 {
		return nil, &errors.Object{
			Id:     "62208df3-5381-4843-94ca-22919beafc53",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing FileType.",
		}
	}

	obj.FileDir = strings.TrimSpace(obj.FileDir)

	if obj.FileDir == "" {
		return nil, &errors.Object{
			Id:     "f642f63d-a526-417a-9ece-140c86ab50c9",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing FileType.",
		}
	}

	obj.FileName = strings.TrimSpace(obj.FileName)

	if obj.FileName == "" {
		return nil, &errors.Object{
			Id:     "db74b226-98da-4082-837b-d6efcf57bf40",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing FileType.",
		}
	}

	if obj.Id == uuid.Nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "daafd239-a58f-4909-8025-7a106660c69b")
		}
		obj.Id = id
	}

	if obj.CreatedAt.IsZero() {
		obj.CreatedAt = time.Now()
	}

	if obj.UpdatedAt.IsZero() {
		obj.UpdatedAt = time.Now()
	}

	insertRecordOut, err := dom.repository.InsertDataFileObjectRecord(r, &InsertDataFileObjectRecordInput{
		Record:            obj,
		DoNotReturnRecord: !in.ReturnEntity,
	})
	if err != nil {
		return nil, errors.Forward(err, "8183d778-c9a8-413a-b123-276efb5fa007")
	}

	out := &InsertDataFileObjectOutput{
		Entity: insertRecordOut.Record,
	}

	return out, nil
}

type InsertDataFileObjectRecordInput struct {
	Record *entities.DataFileObject

	DoNotReturnRecord bool
}

type InsertDataFileObjectRecordOutput struct {
	Record *entities.DataFileObject
}

func (rep *repository) InsertDataFileObjectRecord(r *arc.Request, in *InsertDataFileObjectRecordInput) (*InsertDataFileObjectRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("data_file_objects").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"directory_id",
			"parent_file_id",
			"file_type",
			"hash",
			"status",
			"record_count",
			"file_dir",
			"file_name",
			"file_size",
			"priorities",
			"worker_id",
		).
		Values(
			in.Record.Id,
			in.Record.CreatedAt,
			in.Record.UpdatedAt,
			in.Record.Meta,
			in.Record.DirectoryId,
			in.Record.ParentFileId,
			in.Record.FileType,
			in.Record.Hash,
			in.Record.Status,
			in.Record.RecordCount,
			in.Record.FileDir,
			in.Record.FileName,
			in.Record.FileSize,
			in.Record.Priorities,
			in.Record.WorkerId,
		)

	if !in.DoNotReturnRecord {
		builder = builder.Suffix(`
			RETURNING
				id,
				created_at,
				updated_at,
				meta,
				directory_id,
				parent_file_id,
				file_type,
				hash,
				status,
				record_count,
				file_dir,
				file_name,
				file_size,
				priorities,
				worker_id
		`)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "517e8fdb-fe74-4fec-add8-ab07a79fc005",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "a6a1f1b7-63f9-4d32-b3c5-1aa1af67e2c3")
	}

	out := new(InsertDataFileObjectRecordOutput)

	if !in.DoNotReturnRecord {
		record := new(entities.DataFileObject)

		if err := row.Scan(
			&record.Id,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.Meta,
			&record.DirectoryId,
			&record.ParentFileId,
			&record.FileType,
			&record.Hash,
			&record.Status,
			&record.RecordCount,
			&record.FileDir,
			&record.FileName,
			&record.FileSize,
			&record.Priorities,
			&record.WorkerId,
		); err != nil {
			return nil, &errors.Object{
				Id:     "5e5723db-cb4a-4d25-afde-7917a3a505e6",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to insert row.",
				Cause:  err.Error(),
			}
		}

		out.Record = record
	}

	return out, nil
}

type SelectDataFileObjectInput struct {
	DirectoryId  *uuid.UUID
	FileType     entities.DataFileType
	Hash         []byte
	Meta         map[string]any
	ParentFileId *uuid.UUID
}

type SelectDataFileObjectOutput struct {
	Entity *entities.DataFileObject
}

func (dom *domain) SelectDataFileObject(r *arc.Request, in *SelectDataFileObjectInput) (*SelectDataFileObjectOutput, error) {
	selectDataFileObjectOut, err := dom.repository.SelectDataFileObjectRecord(r, &SelectDataFileObjectRecordInput{
		DirectoryId:  in.DirectoryId,
		FileType:     in.FileType,
		Hash:         in.Hash,
		Meta:         in.Meta,
		ParentFileId: in.ParentFileId,
	})
	if err != nil {
		return nil, errors.Forward(err, "20c496e5-7208-4cf0-9ac4-1750ff4960cb")
	}

	out := &SelectDataFileObjectOutput{
		Entity: selectDataFileObjectOut.Record,
	}

	return out, nil
}

type SelectDataFileObjectRecordInput struct {
	DirectoryId  *uuid.UUID
	FileType     entities.DataFileType
	Hash         []byte
	Meta         map[string]any
	ParentFileId *uuid.UUID
}

type SelectDataFileObjectRecordOutput struct {
	Record *entities.DataFileObject
}

func (rep *repository) SelectDataFileObjectRecord(r *arc.Request, in *SelectDataFileObjectRecordInput) (*SelectDataFileObjectRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"directory_id",
			"file_type",
			"hash",
			"status",
			"record_count",
			"file_dir",
			"file_name",
			"file_size",
			"priorities",
		).
		From("data_file_objects").
		Where("file_type = ?", in.FileType).
		Where("hash = ?", in.Hash)

	if in.DirectoryId != nil && *in.DirectoryId != uuid.Nil {
		builder = builder.Where("directory_id = ?", in.DirectoryId)
	}

	if len(in.Meta) > 0 {
		builder = builder.Where("meta @> ?", in.Meta)
	}

	if in.ParentFileId != nil && *in.ParentFileId != uuid.Nil {
		builder = builder.Where("parent_file_id = ?", in.ParentFileId)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "61b966bb-3f09-41a2-939b-db59239ee46d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "f4a83721-e5ac-4049-8e4b-ddc92c35378e")
	}

	out := new(SelectDataFileObjectRecordOutput)
	record := new(entities.DataFileObject)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.DirectoryId,
		&record.FileType,
		&record.Hash,
		&record.Status,
		&record.RecordCount,
		&record.FileDir,
		&record.FileName,
		&record.FileSize,
		&record.Priorities,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}

		return nil, &errors.Object{
			Id:     "3c02de9f-7015-4701-b230-57c694d17cf2",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	out.Record = record

	return out, nil
}

func (dom *domain) UpdateDataFileObject(r *arc.Request, in *entities.UpdateDataFileObjectInput) (*entities.UpdateDataFileObjectOutput, error) {
	updateDataFileObjectOut, err := dom.repository.UpdateDataFileObjectRecord(r, &UpdateDataFileObjectRecordInput{
		Id:          in.Id,
		UpdatedAt:   in.UpdatedAt,
		Meta:        in.Meta,
		Status:      in.Status,
		RecordCount: in.RecordCount,
		Priorities:  in.Priorities,
	})
	if err != nil {
		return nil, errors.Forward(err, "b9631a64-4539-4642-ad0c-dd3499c81824")
	}

	out := &entities.UpdateDataFileObjectOutput{
		Entity: updateDataFileObjectOut.Record,
	}

	return out, nil
}

type UpdateDataFileObjectRecordInput struct {
	Id          uuid.UUID
	UpdatedAt   time.Time
	Meta        map[string]any
	Status      int32
	RecordCount int32
	Priorities  []int32
}

type UpdateDataFileObjectRecordOutput struct {
	Record *entities.DataFileObject
}

func (rep *repository) UpdateDataFileObjectRecord(r *arc.Request, in *UpdateDataFileObjectRecordInput) (*UpdateDataFileObjectRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("data_file_objects").
		Where("id = ?", in.Id).
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at,
				meta,
				directory_id,
				parent_file_id,
				file_type,
				hash,
				status,
				record_count,
				file_dir,
				file_name,
				file_size,
				priorities
		`)

	if !in.UpdatedAt.IsZero() {
		builder = builder.Set("updated_at", in.UpdatedAt)
	}

	if len(in.Meta) > 0 {
		builder = builder.Set("meta", in.Meta)
	}

	if in.Status > 0 {
		builder = builder.Set("status", in.Status)
	}

	if in.RecordCount > 0 {
		builder = builder.Set("record_count", in.RecordCount)
	}

	if len(in.Priorities) > 0 {
		builder = builder.Set("priorities", in.Priorities)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "f7d42e9c-8b51-4d3c-a8f6-9e4f3b2c1a5d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "c2edb3b3-5763-4103-9ef9-de92bfe190f7")
	}

	out := new(UpdateDataFileObjectRecordOutput)
	record := new(entities.DataFileObject)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.DirectoryId,
		&record.ParentFileId,
		&record.FileType,
		&record.Hash,
		&record.Status,
		&record.RecordCount,
		&record.FileDir,
		&record.FileName,
		&record.FileSize,
		&record.Priorities,
	); err != nil {
		return nil, &errors.Object{
			Id:     "f4c4a05a-6898-4249-9cb5-01869d258cdd",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update row.",
			Cause:  err.Error(),
		}
	}

	out.Record = record

	return out, nil
}

type SelectUnparsedDataFileObjectRecordsInput struct {
	Limit         int32
	Meta          map[string]any
	PartnerId     *uuid.UUID
	PriorityGroup int32
	Statuses      []int32
	WorkerId      *uuid.UUID
}

type SelectUnparsedDataFileObjectRecordsOutput struct {
	Records []*entities.DataFileObject
}

func (repo *repository) SelectUnparsedDataFileObjectRecords(r *arc.Request, in *SelectUnparsedDataFileObjectRecordsInput) (*SelectUnparsedDataFileObjectRecordsOutput, error) {
	sql := `
		with unparsed_objects as (
			select
				data_file_objects.id,
				data_file_objects.meta,
				data_file_objects.directory_id,
				data_file_objects.priorities,
				data_file_objects.worker_id
			from data_file_objects
			where
				data_file_objects.status = any ($2)
				and coalesce(cardinality(data_file_objects.priorities), 0) > 0
				and data_file_objects.parent_file_id is null
		), ranked_objects as (
			select
				unparsed_objects.id,
				unparsed_objects.worker_id,
				rank() over (
					partition by unparsed_objects.priorities[1]
					order by unparsed_objects.priorities
				) as rn
			from unparsed_objects
			join data_file_directories
				on unparsed_objects.directory_id = data_file_directories.id
			where
				data_file_directories.partner_id = $3
				and ($4 or unparsed_objects.meta @> $5)
				and ($8 or unparsed_objects.priorities[1] = $9)
		), updated_objects as (
			update data_file_objects
			set
				status = $6,
				updated_at = now(),
				worker_id = $1
			where id in (
				select id
				from ranked_objects
				where
					rn = 1
					and (
						worker_id is null
						or worker_id <> $1
					)
				limit $7
			)
			returning *
		)
		select
			id,
			created_at,
			updated_at,
			meta,
			directory_id,
			parent_file_id,
			file_type,
			hash,
			status,
			record_count,
			file_dir,
			file_name,
			file_size,
			priorities,
			worker_id
		from updated_objects
		order by priorities
	`

	args := []any{
		in.WorkerId,
		in.Statuses,
		in.PartnerId,
		len(in.Meta) == 0,
		in.Meta,
		entities.DataFileObjectStatusInProgress,
		in.Limit,
		in.PriorityGroup == 0,
		in.PriorityGroup,
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "b47587d3-3aab-4dd0-86be-969404aeceff")
	}
	defer rows.Close()

	out := &SelectUnparsedDataFileObjectRecordsOutput{}

	for rows.Next() {
		record := new(entities.DataFileObject)

		if err := rows.Scan(
			&record.Id,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.Meta,
			&record.DirectoryId,
			&record.ParentFileId,
			&record.FileType,
			&record.Hash,
			&record.Status,
			&record.RecordCount,
			&record.FileDir,
			&record.FileName,
			&record.FileSize,
			&record.Priorities,
			&record.WorkerId,
		); err != nil {
			return nil, &errors.Object{
				Id:     "6d1496c7-7d07-49fd-9ce8-cb908e7dd518",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}
