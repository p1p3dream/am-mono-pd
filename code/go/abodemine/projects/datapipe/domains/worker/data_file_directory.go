package worker

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/projects/datapipe/entities"
)

type InsertDataFileDirectoryInput struct {
	Entity *entities.DataFileDirectory
}

type InsertDataFileDirectoryOutput struct {
	Entity *entities.DataFileDirectory
}

func (dom *domain) InsertDataFileDirectory(r *arc.Request, in *InsertDataFileDirectoryInput) (*InsertDataFileDirectoryOutput, error) {
	insertDataFileDirectoryOut, err := dom.repository.InsertDataFileDirectoryRecord(r, &InsertDataFileDirectoryRecordInput{
		Record: in.Entity,
	})
	if err != nil {
		return nil, errors.Forward(err, "86a3bc6f-c481-4691-965e-f996c3427a97")
	}

	out := &InsertDataFileDirectoryOutput{
		Entity: insertDataFileDirectoryOut.Record,
	}

	return out, nil
}

type InsertDataFileDirectoryRecordInput struct {
	Record *entities.DataFileDirectory
}

type InsertDataFileDirectoryRecordOutput struct {
	Record *entities.DataFileDirectory
}

func (rep *repository) InsertDataFileDirectoryRecord(r *arc.Request, in *InsertDataFileDirectoryRecordInput) (*InsertDataFileDirectoryRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("data_file_directories").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"partner_id",
			"parent_directory_id",
			"status",
			"path",
			"name",
		).
		Values(
			in.Record.Id,
			in.Record.CreatedAt,
			in.Record.UpdatedAt,
			in.Record.Meta,
			in.Record.PartnerId,
			in.Record.ParentDirectoryId,
			in.Record.Status,
			in.Record.Path,
			in.Record.Name,
		).
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at,
				meta,
				partner_id,
				parent_directory_id,
				status,
				path,
				name
		`)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "6a1b2c3d-4e5f-6g7h-8i9j-0k1l2m3n4o5p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "c5b0eeb1-9e0d-4d08-b098-d84a9e052b75")
	}

	out := new(InsertDataFileDirectoryRecordOutput)
	record := new(entities.DataFileDirectory)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.PartnerId,
		&record.ParentDirectoryId,
		&record.Status,
		&record.Path,
		&record.Name,
	); err != nil {
		return nil, &errors.Object{
			Id:     "2a3b4c5d-6e7f-8g9h-0i1j-2k3l4m5n6o7p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to insert row.",
			Cause:  err.Error(),
		}
	}

	out.Record = record

	return out, nil
}

type SelectDataFileDirectoryInput struct {
	Meta      map[string]any
	PartnerId uuid.UUID
	Path      string
}

type SelectDataFileDirectoryOutput struct {
	Entity *entities.DataFileDirectory
}

func (dom *domain) SelectDataFileDirectory(r *arc.Request, in *SelectDataFileDirectoryInput) (*SelectDataFileDirectoryOutput, error) {
	selectDataFileDirectoryOut, err := dom.repository.SelectDataFileDirectoryRecord(r, &SelectDataFileDirectoryRecordInput{
		Meta:      in.Meta,
		PartnerId: in.PartnerId,
		Path:      in.Path,
	})
	if err != nil {
		return nil, errors.Forward(err, "46eb7f88-cf00-47ae-9a9a-77b7977460f5")
	}

	out := &SelectDataFileDirectoryOutput{
		Entity: selectDataFileDirectoryOut.Record,
	}

	return out, nil
}

type SelectDataFileDirectoryRecordInput struct {
	Meta      map[string]any
	PartnerId uuid.UUID
	Path      string
}

type SelectDataFileDirectoryRecordOutput struct {
	Record *entities.DataFileDirectory
}

func (rep *repository) SelectDataFileDirectoryRecord(r *arc.Request, in *SelectDataFileDirectoryRecordInput) (*SelectDataFileDirectoryRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"partner_id",
			"parent_directory_id",
			"status",
			"path",
			"name",
		).
		From("data_file_directories").
		Where("partner_id = ?", in.PartnerId).
		Where("path = ?", in.Path)

	if len(in.Meta) > 0 {
		builder = builder.Where("meta @> ?", in.Meta)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "3a4b5c6d-7e8f-9g0h-1i2j-3k4l5m6n7o8p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "c7144610-e7c3-4f66-a161-d9bd2d8053c7")
	}

	out := new(SelectDataFileDirectoryRecordOutput)
	record := new(entities.DataFileDirectory)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.PartnerId,
		&record.ParentDirectoryId,
		&record.Status,
		&record.Path,
		&record.Name,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}
		return nil, &errors.Object{
			Id:     "5a6b7c8d-9e0f-1g2h-3i4j-5k6l7m8n9o0p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	out.Record = record

	return out, nil
}

type UpdateDataFileDirectoryInput struct {
	Id        uuid.UUID
	UpdatedAt time.Time
	Meta      map[string]any
	Status    int
}

type UpdateDataFileDirectoryOutput struct {
	Entity *entities.DataFileDirectory
}

func (dom *domain) UpdateDataFileDirectory(r *arc.Request, in *UpdateDataFileDirectoryInput) (*UpdateDataFileDirectoryOutput, error) {
	UpdateDataFileDirectoryOut, err := dom.repository.UpdateDataFileDirectoryRecord(r, &UpdateDataFileDirectoryRecordInput{
		Id:        in.Id,
		UpdatedAt: in.UpdatedAt,
		Meta:      in.Meta,
		Status:    in.Status,
	})
	if err != nil {
		return nil, errors.Forward(err, "c1ee28e1-801f-4aa7-adce-4ad82df32a2e")
	}

	out := &UpdateDataFileDirectoryOutput{
		Entity: UpdateDataFileDirectoryOut.Record,
	}

	return out, nil
}

type UpdateDataFileDirectoryRecordInput struct {
	Id        uuid.UUID
	UpdatedAt time.Time
	Meta      map[string]any
	Status    int
}

type UpdateDataFileDirectoryRecordOutput struct {
	Record *entities.DataFileDirectory
}

func (rep *repository) UpdateDataFileDirectoryRecord(r *arc.Request, in *UpdateDataFileDirectoryRecordInput) (*UpdateDataFileDirectoryRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("data_file_directories").
		Where("id = ?", in.Id).
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at,
				meta,
				partner_id,
				parent_directory_id,
				status,
				path,
				name,
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

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "6a7b8c9d-0e1f-2g3h-4i5j-6k7l8m9n0o1p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "56679202-476c-4442-916f-1e638c25c966")
	}

	out := new(UpdateDataFileDirectoryRecordOutput)
	record := new(entities.DataFileDirectory)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.PartnerId,
		&record.ParentDirectoryId,
		&record.Status,
		&record.Path,
		&record.Name,
		&record.Priorities,
	); err != nil {
		return nil, &errors.Object{
			Id:     "8a9b0c1d-2e3f-4g5h-6i7j-8k9l0m1n2o3p",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update row.",
			Cause:  err.Error(),
		}
	}

	out.Record = record

	return out, nil
}
