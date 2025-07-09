package worker

import (
	"archive/zip"
	"bufio"
	"io"
	"path"
	"runtime/debug"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/storage"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type ProcessDataSourceDirInput struct {
	Backend           storage.Backend
	FileBufferSize    int
	DataFileType      entities.DataFileType
	DataSource        entities.DataSource
	IgnoreSubDirs     bool
	IsRootDir         bool
	Meta              map[string]any
	ParentDirectoryId *uuid.UUID
	PartnerId         uuid.UUID
	Path              string
	Priorities        []int32
	PriorityGroup     int32
	WorkerId          *uuid.UUID
}

type ProcessDataSourceDirOutput struct {
	TotalObjectsToLoad int32
}

func (dom *domain) ProcessDataSourceDirV2(r *arc.Request, in *ProcessDataSourceDirInput) (*ProcessDataSourceDirOutput, error) {
	directoryPath := path.Join(
		"/",
		// We don't know if the backend is local fs or remote, so we replace
		// the path separator with "/".
		path.Clean(strings.ReplaceAll(in.Path, in.Backend.PathSeparator(), "/")),
	)

	log.Info().
		Str("path", directoryPath).
		Msg("Processing data source directory.")

	selectDirectoryOut, err := dom.SelectDataFileDirectory(r, &SelectDataFileDirectoryInput{
		Meta:      in.Meta,
		PartnerId: in.PartnerId,
		Path:      directoryPath,
	})
	if err != nil {
		return nil, errors.Forward(err, "11900908-8c03-441e-9ab4-291578a5cdb4")
	}

	directory := selectDirectoryOut.Entity

	if directory == nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "df583c4a-1e3d-44c0-bcc0-b7fd823df923")
		}

		now := time.Now()

		insertDirectoryOut, err := dom.InsertDataFileDirectory(r, &InsertDataFileDirectoryInput{
			Entity: &entities.DataFileDirectory{
				Id:                id,
				CreatedAt:         now,
				UpdatedAt:         now,
				Meta:              in.Meta,
				ParentDirectoryId: in.ParentDirectoryId,
				PartnerId:         in.PartnerId,
				Status:            entities.DataFileDirectoryStatusToDo,
				Path:              directoryPath,
				Name:              path.Base(directoryPath),
			},
		})
		if err != nil {
			return nil, errors.Forward(err, "df6caebd-dcd3-43bb-943c-f213481fd00e")
		}

		directory = insertDirectoryOut.Entity
	}

	out := &ProcessDataSourceDirOutput{}

	switch directory.Status {
	case entities.DataFileDirectoryStatusDone:
		log.Info().
			Str("path", in.Path).
			Msg("Directory already processed.")
		return out, nil
	case entities.DataFileDirectoryStatusIgnored:
		log.Info().
			Str("path", in.Path).
			Msg("Directory ignored.")
		return out, nil
	}

	stgObjects, err := in.Backend.List(
		r.Context(),
		directoryPath,
		&storage.ListOptions{
			WithSize: true,
		},
	)
	if err != nil {
		return nil, errors.Forward(err, "9514ede7-036b-4d62-84d9-bb5e3fe5139e")
	}

	var totalObjectsToLoad int32

	for _, stgObject := range stgObjects {
		processDataSourceObjectOut, err := dom.ProcessDataSourceObjectV2(r, &ProcessDataSourceObjectInput{
			Backend:       in.Backend,
			DataFileType:  in.DataFileType,
			DataSource:    in.DataSource,
			DirectoryId:   &directory.Id,
			IgnoreSubDirs: in.IgnoreSubDirs,
			Meta:          in.Meta,
			StorageObject: stgObject,
			PartnerId:     in.PartnerId,
			Path:          path.Join(directoryPath, stgObject.Name),
			Priorities:    in.Priorities,
			WorkerId:      in.WorkerId,
		})
		if err != nil {
			return nil, errors.Forward(err, "267834cd-d86e-4198-adfe-bc79afb25314")
		}

		totalObjectsToLoad += processDataSourceObjectOut.TotalObjectsToLoad
	}

	if !in.IsRootDir {
		out.TotalObjectsToLoad = totalObjectsToLoad
		return out, nil
	}

	log.Info().
		Int32("totalObjectsToLoad", totalObjectsToLoad).
		Msg("Processing data sources.")

	g, gctx := errgroup.WithContext(r.Context())
	gr := r.Clone(arc.CloneRequestWithContext(gctx))

	remainingObjectsToLoad := atomic.Int32{}
	remainingObjectsToLoad.Store(totalObjectsToLoad)
	stopProcessing := atomic.Bool{}

	g.SetLimit(in.FileBufferSize)

	for batchNumber := 1; ; batchNumber++ {
		if stopProcessing.Load() {
			break
		}

		unparsedDataFileObjectsOut, err := dom.repository.SelectUnparsedDataFileObjectRecords(gr, &SelectUnparsedDataFileObjectRecordsInput{
			Limit:         int32(in.FileBufferSize),
			Meta:          in.Meta,
			PartnerId:     &in.PartnerId,
			PriorityGroup: in.PriorityGroup,
			Statuses: []int32{
				entities.DataFileObjectStatusToDo,
				entities.DataFileObjectStatusInProgress,
			},
			WorkerId: in.WorkerId,
		})
		if err != nil {
			// Prevent this error from shadowing the errgroup error.
			log.Error().
				Err(errors.Forward(err, "1f275025-b748-4637-bfd4-be7778c70f93")).
				Send()
			break
		}

		dfObjects := unparsedDataFileObjectsOut.Records

		if len(dfObjects) == 0 {
			// The returned list may be empty if the remaining objects
			// depend on currently running loaders. Therefore, we check if
			// there are any remaining objects to load.

			remaining := remainingObjectsToLoad.Load()

			if remaining == 0 {
				// No more objects to process.
				break
			}

		WAIT_LOADER:
			for remaining == remainingObjectsToLoad.Load() {
				// Wait for a loader to finish before checking again.
				select {
				case <-gctx.Done():
					break WAIT_LOADER
				case <-time.After(250 * time.Millisecond):
					// Check again.
				}
			}

			continue
		}

		for _, dfObject := range dfObjects {
			g.Go(func() error {
				defer func() {
					remainingObjectsToLoad.Add(-1)

					if r := recover(); r != nil {
						stopProcessing.Store(true)
						log.Error().
							Err(&errors.Object{
								Id:     "18017cd6-3be2-4032-9b82-8b9c7e51ac36",
								Code:   errors.Code_INTERNAL,
								Detail: "Panic in goroutine.",
							}).
							Str("stack", string(debug.Stack())).
							Send()
					}
				}()

				loadDataSourceObjectOut, err := dom.LoadDataSourceObject(gr, &LoadDataSourceObjectInput{
					Backend:        in.Backend,
					BatchNumber:    batchNumber,
					DataFileObject: dfObject,
					DataSource:     in.DataSource,
					Meta:           in.Meta,
					PartnerId:      in.PartnerId,
				})
				if err != nil {
					return errors.Forward(err, "31229781-ec19-4096-ae57-c8d554332591")
				}

				_ = loadDataSourceObjectOut

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Forward(err, "8aca4c54-54fa-4477-b697-963572bf0de1")
	}

	return out, nil
}

type ProcessDataSourceObjectInput struct {
	Backend       storage.Backend
	DataFileType  entities.DataFileType
	DataSource    entities.DataSource
	DirectoryId   *uuid.UUID
	IgnoreSubDirs bool
	Meta          map[string]any
	ParentFileId  *uuid.UUID
	PartnerId     uuid.UUID
	Path          string
	Priorities    []int32
	WorkerId      *uuid.UUID

	StorageObject *storage.Object
	ZipFile       *zip.File
}

type ProcessDataSourceObjectOutput struct {
	TotalObjectsToLoad int32
}

func (dom *domain) ProcessDataSourceObjectV2(r *arc.Request, in *ProcessDataSourceObjectInput) (*ProcessDataSourceObjectOutput, error) {
	obj := in.StorageObject

	entry, err := in.DataSource.CreateDataFileEntry(r, &entities.CreateDataFileEntryInput{
		Path:          in.Path,
		StorageObject: obj,
	})
	if err != nil {
		return nil, errors.Forward(err, "39e450ce-b9a3-4012-93fa-12fccc5001b0")
	}

	out := &ProcessDataSourceObjectOutput{}

	if entry.Ignore {
		return out, nil
	}

	fileType := entry.FileType

	switch {
	case fileType == 0:
		if in.DataFileType == 0 {
			log.Info().
				Str("path", in.Path).
				Msg("Skipping unknown data file type.")
			return out, nil
		}

		// Use the parent directory's DataFileType.
		fileType = in.DataFileType
	case fileType == entities.DataFileTypeSelectedDirectory:
		fileType = 0
	}

	if obj.IsDirectory() {
		if entry.EnterDirectory || !in.IgnoreSubDirs {
			// Recurse into subdirectories.
			processDataSourceDirOut, err := dom.ProcessDataSourceDirV2(r, &ProcessDataSourceDirInput{
				Backend:           in.Backend,
				DataFileType:      fileType,
				DataSource:        in.DataSource,
				IgnoreSubDirs:     entry.IgnoreSubDirs,
				Meta:              in.Meta,
				PartnerId:         in.PartnerId,
				ParentDirectoryId: in.DirectoryId,
				Path:              in.Path,
				Priorities:        entry.Priorities,
				WorkerId:          in.WorkerId,
			})
			if err != nil {
				return nil, errors.Forward(err, "c5b919f7-b3f6-48d5-83f0-842f24b4ff74")
			}

			out.TotalObjectsToLoad = processDataSourceDirOut.TotalObjectsToLoad
		}

		// No further processing necessary for directories.
		return out, nil
	}

	switch strings.ToLower(path.Ext(in.Path)) {
	case ".txt", ".zip":
		// This file format is supported.
	default:
		log.Info().
			Str("path", in.Path).
			Msg("Ignoring unsupported file format.")
		return out, nil
	}

	if fileType == 0 {
		if in.DataFileType == 0 {
			// Skip unknown files.
			return out, nil
		}

		// Use the parent directory's DataFileType.
		fileType = in.DataFileType
	}

	priorities := val.Ternary(
		len(entry.Priorities) > 0,
		entry.Priorities,
		in.Priorities,
	)

	ensureDataFileObjectOut, err := dom.EnsureDataFileObject(r, &EnsureDataFileObjectInput{
		DataFileType: fileType,
		DirectoryId:  in.DirectoryId,
		FileSize:     obj.Size,
		Meta:         in.Meta,
		ParentFileId: in.ParentFileId,
		Path:         in.Path,
		// If entry has priorities, use them.
		// Otherwise, use the parent directory's.
		Priorities: priorities,
	})
	if err != nil {
		return nil, errors.Forward(err, "4a78e463-e110-474e-8ed7-9ecefeb4bac6")
	}

	dfObject := ensureDataFileObjectOut.Entity

	switch dfObject.Status {
	case entities.DataFileObjectStatusToDo:
		log.Info().
			Str("path", in.Path).
			Msg("Found new data source object.")
		out.TotalObjectsToLoad++
	case entities.DataFileObjectStatusInProgress:
		log.Info().
			Str("path", in.Path).
			Msg("Found pending data source object.")
		out.TotalObjectsToLoad++
	case entities.DataFileObjectStatusDone,
		entities.DataFileObjectStatusIgnored:
		// Skip.
		return out, nil
	}

	// Check if returning object has the same priorities as before.
	// If not, update it.
	if slices.Compare(dfObject.Priorities, priorities) != 0 {
		_, err := dom.UpdateDataFileObject(r, &entities.UpdateDataFileObjectInput{
			Id:         dfObject.Id,
			UpdatedAt:  time.Now(),
			Priorities: priorities,
		})
		if err != nil {
			return nil, errors.Forward(err, "a2c31572-7343-4672-b1ac-a267a8490980")
		}
	}

	return out, nil
}

type LoadDataSourceObjectInput struct {
	Backend        storage.Backend
	BatchNumber    int
	DataFileObject *entities.DataFileObject
	DataSource     entities.DataSource
	Meta           map[string]any
	PartnerId      uuid.UUID
}

type LoadDataSourceObjectOutput struct {
	ObjectFound bool
	RecordCount int32
}

func (dom *domain) LoadDataSourceObject(r *arc.Request, in *LoadDataSourceObjectInput) (*LoadDataSourceObjectOutput, error) {
	out := &LoadDataSourceObjectOutput{}

	if r.Context().Err() != nil {
		// Prevent further processing if the context is cancelled.
		return out, nil
	}

	dfObject := in.DataFileObject

	log.Info().
		Int("batchNumber", in.BatchNumber).
		Str("fileDir", dfObject.FileDir).
		Str("fileName", dfObject.FileName).
		Msg("Loading data source object.")

	fileExt := strings.ToLower(path.Ext(dfObject.FileName))

	if fileExt != ".zip" {
		log.Warn().
			Int("batchNumber", in.BatchNumber).
			Str("fileExt", fileExt).
			Str("fileDir", dfObject.FileDir).
			Str("fileName", dfObject.FileName).
			Msg("Unsupported file format.")
		return out, nil
	}

	_, err := dom.LoadZipDataSourceObject(r, &LoadZipDataSourceObjectInput{
		Backend:        in.Backend,
		DataFileObject: dfObject,
		DataSource:     in.DataSource,
	})
	if err != nil {
		return nil, errors.Forward(err, "8f7b872e-a1e7-4ba3-812f-949794716483")
	}

	updateDataFileObjectOut, err := dom.UpdateDataFileObject(r, &entities.UpdateDataFileObjectInput{
		Id:        dfObject.Id,
		UpdatedAt: time.Now(),
		Status:    entities.DataFileObjectStatusDone,
	})
	if err != nil {
		return nil, errors.Forward(err, "0c77f461-47f1-4f22-9d23-5916ba9e5572")
	}

	_ = updateDataFileObjectOut

	return out, nil
}

type LoadZipDataSourceObjectInput struct {
	Backend        storage.Backend
	DataFileObject *entities.DataFileObject
	DataSource     entities.DataSource
}

type LoadZipDataSourceObjectOutput struct{}

func (dom *domain) LoadZipDataSourceObject(r *arc.Request, in *LoadZipDataSourceObjectInput) (*LoadZipDataSourceObjectOutput, error) {
	dfObject := in.DataFileObject
	stgObject := &storage.Object{
		// Since we store only the relative path on dfObject,
		// ensure we also use the backend's path prefix.
		Dir:  in.Backend.PathJoin(in.Backend.Path(), dfObject.FileDir),
		Name: dfObject.FileName,
		Size: val.PtrDeref(dfObject.FileSize),
	}

	zipArchive, err := in.Backend.Get(r.Context(), stgObject)
	if err != nil {
		return nil, &errors.Object{
			Id:     "73920644-ce4e-4433-9f21-cffd6ea36485",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to open zip file.",
			Cause:  err.Error(),
		}
	}
	defer zipArchive.Close()

	zipReader, err := zip.NewReader(zipArchive, stgObject.Size)
	if err != nil {
		return nil, &errors.Object{
			Id:     "7f5082d5-1018-46bb-8c23-68047a281850",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create zip reader.",
			Cause:  err.Error(),
		}
	}

	for _, zipFile := range zipReader.File {
		fileExt := path.Ext(zipFile.Name)

		if strings.ToLower(fileExt) != ".txt" {
			continue
		}

		_, err = dom.LoadTxtDataSourceObject(r, &LoadTxtDataSourceObjectInput{
			Backend:        in.Backend,
			DataFileObject: in.DataFileObject,
			DataSource:     in.DataSource,
			Path:           zipFile.Name,
			ZipFile:        zipFile,
		})
		if err != nil {
			return nil, errors.Forward(err, "bf8d9503-6f86-49f0-82ea-b269f16cbabd")
		}
	}

	out := &LoadZipDataSourceObjectOutput{}

	return out, nil
}

type LoadTxtDataSourceObjectInput struct {
	Backend        storage.Backend
	DataFileObject *entities.DataFileObject
	DataSource     entities.DataSource
	Path           string

	StorageObject *storage.Object
	ZipFile       *zip.File
}

type LoadTxtDataSourceObjectOutput struct{}

func (dom *domain) LoadTxtDataSourceObject(r *arc.Request, in *LoadTxtDataSourceObjectInput) (*LoadTxtDataSourceObjectOutput, error) {
	log.Info().
		Str("path", in.Path).
		Msg("Processing txt data file.")

	var readCloser io.ReadCloser
	var err error

	switch {
	case in.ZipFile != nil:
		readCloser, err = in.ZipFile.Open()
	case in.StorageObject != nil:
		readCloser, err = in.Backend.Get(r.Context(), in.StorageObject)
	}

	if err != nil {
		return nil, &errors.Object{
			Id:     "fe051d40-0a2c-4ed6-9d1a-0d9a9468a844",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to open data file.",
			Cause:  err.Error(),
		}
	}

	defer func() {
		if err := readCloser.Close(); err != nil {
			log.Error().
				Err(&errors.Object{
					Id:     "0711fa9f-0864-4adc-9ac0-bef5f9356f85",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to close data file.",
					Cause:  err.Error(),
				}).
				Send()
		}
	}()

	parentDfObject := in.DataFileObject

	ensureDataFileObjectOut, err := dom.EnsureDataFileObject(r, &EnsureDataFileObjectInput{
		DataFileType: parentDfObject.FileType,
		Meta:         parentDfObject.Meta,
		Path:         in.Path,
		ParentFileId: &parentDfObject.Id,
		WorkerId:     parentDfObject.WorkerId,
	})
	if err != nil {
		return nil, errors.Forward(err, "04feb016-7f45-4507-a7f7-59e3df4ee212")
	}

	dfObject := ensureDataFileObjectOut.Entity

	dataRecord, err := in.DataSource.DataRecordByFileType(dfObject.FileType)
	if err != nil {
		return nil, errors.Forward(err, "b2b1a76a-3387-4689-b14d-e1dbd62ac315")
	}

	loadParams := dataRecord.LoadParams()

	if loadParams == nil {
		return nil, &errors.Object{
			Id:     "7e3a9703-9152-4e78-b95f-928fb7756328",
			Code:   errors.Code_INTERNAL,
			Detail: "Undefined load params.",
		}
	}

	columns := dataRecord.SQLColumns()
	// Ensure we don't exceed the maximum number of params for the operation.
	batchSize := int32(min(65535/len(columns), 1000))
	headers := make(map[int]string)
	scanner := bufio.NewScanner(transform.NewReader(readCloser, charmap.ISO8859_1.NewDecoder()))

	// Make headers first so we don't have to check if it's
	// the first line for every record.
	scanner.Scan()

	fieldSeparator := in.DataSource.FieldSeparatorByFileType(dfObject.FileType)

	for i, field := range strings.Split(
		scanner.Text(),
		fieldSeparator,
	) {
		headers[i] = field
	}

	// Skip previously processed records, if any.
	for range dfObject.RecordCount {
		scanner.Scan()
	}

	loadRecordIn := &entities.LoadDataRecordInput{
		BatchSize:                batchSize,
		Columns:                  columns,
		DataFileObject:           dfObject,
		DataRecord:               dataRecord,
		FieldSeparator:           fieldSeparator,
		Headers:                  headers,
		Scanner:                  scanner,
		UpdateDataFileObjectFunc: dom.UpdateDataFileObject,
	}

	var loadRecordOut *entities.LoadDataRecordOutput

	switch loadParams.Mode {
	case entities.DataRecordModeBatchDelete:
		loadRecordOut, err = dom.BatchDeleteDataRecord(r, loadRecordIn)
		if err != nil {
			return nil, errors.Forward(err, "c203e3b0-8dc3-4d95-b17c-5fb846890999")
		}
	case entities.DataRecordModeBatchInsert:
		loadRecordOut, err = dom.BatchInsertDataRecord(r, loadRecordIn)
		if err != nil {
			return nil, errors.Forward(err, "0e5b6a8a-7a66-41b4-aec1-ac143b844fef")
		}
	case entities.DataRecordModeLoadFunc:
		if loadParams.LoadFunc == nil {
			return nil, &errors.Object{
				Id:     "de7b7a03-a0da-4193-bf4c-07d854ed6d35",
				Code:   errors.Code_INTERNAL,
				Detail: "Undefined load function.",
			}
		}

		loadRecordOut, err = loadParams.LoadFunc(r, loadRecordIn)
		if err != nil {
			return nil, errors.Forward(err, "408db9df-4e10-44ab-89f2-921ca2441324")
		}
	default:
		return nil, &errors.Object{
			Id:     "1cc5ed49-0094-493a-a98e-8a829bdf0643",
			Code:   errors.Code_INTERNAL,
			Detail: "Unsupported processing mode.",
			Meta: map[string]any{
				"mode": loadParams.Mode,
			},
		}
	}

	log.Info().
		Str("path", in.Path).
		Int64("deletedRecords", loadRecordOut.DeletedRecords).
		Int64("processedRecords", loadRecordOut.ProcessedRecords).
		Send()

	if err := scanner.Err(); err != nil {
		return nil, &errors.Object{
			Id:     "6ded960f-e6f1-4fd5-94e4-db21e38d6e29",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to scan file.",
			Cause:  err.Error(),
		}
	}

	_, err = dom.UpdateDataFileObject(r, &entities.UpdateDataFileObjectInput{
		Id:        dfObject.Id,
		UpdatedAt: time.Now(),
		Status:    entities.DataFileObjectStatusDone,
	})
	if err != nil {
		return nil, errors.Forward(err, "797368fb-b0e3-4790-8765-b731f3630179")
	}

	out := &LoadTxtDataSourceObjectOutput{}

	return out, nil
}

func (dom *domain) BatchDeleteDataRecord(r *arc.Request, in *entities.LoadDataRecordInput) (*entities.LoadDataRecordOutput, error) {
	if len(in.Columns) != 1 {
		return nil, &errors.Object{
			Id:     "2340d40f-a214-46c2-930b-e8a42bc67ee7",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Batch delete requires exactly one column.",
			Meta: map[string]any{
				"columns": in.Columns,
			},
		}
	}

	newBuilder := func() squirrel.DeleteBuilder {
		return squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Delete(in.DataRecord.SQLTable())
	}

	var recordCount int32

	builder := newBuilder()
	column := in.Columns[0]
	dfObject := in.DataFileObject
	scanner := in.Scanner
	processedRecords := int64(0)
	values := make([]any, in.BatchSize)

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "9a42e151-cf85-4241-9267-2491e6700ac4")
	}

	deleteRecords := func() error {
		// Skip if no records to delete.
		if recordCount == 0 {
			return nil
		}

		builder = builder.Where(squirrel.Eq{column: values[:recordCount]})

		sql, args, err := builder.ToSql()
		if err != nil {
			return &errors.Object{
				Id:     "842747a2-158e-4f59-a720-63ed2ec726d1",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to build SQL.",
				Cause:  err.Error(),
			}
		}

		tx, err := pgxPool.Begin(r.Context())
		if err != nil {
			return &errors.Object{
				Id:     "a5f0b506-7c74-48a9-b126-bac3d8c59306",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to begin transaction.",
				Cause:  err.Error(),
			}
		}

		// This clone won't replace the original arc.
		r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

		_, err = dom.repository.RemoveDataRecords(r, &RemoveDataRecordsInput{
			SQL:  sql,
			Args: args,
		})
		if err != nil {
			extutils.RollbackPgxTx(r.Context(), tx, "3735a982-17c7-475f-9af3-b37c1f201c25")
			return errors.Forward(err, "af55e930-8dc4-45a9-8ca7-29c9443602e7")
		}

		updateObjectOut, err := dom.UpdateDataFileObject(r, &entities.UpdateDataFileObjectInput{
			Id:          dfObject.Id,
			UpdatedAt:   time.Now(),
			RecordCount: dfObject.RecordCount + recordCount,
		})
		if err != nil {
			extutils.RollbackPgxTx(r.Context(), tx, "86d98971-aba6-45a7-a655-582ddf188879")
			return errors.Forward(err, "1e5135c1-49ab-4af6-8004-30b884f37f34")
		}

		if err := tx.Commit(r.Context()); err != nil {
			return &errors.Object{
				Id:     "5889e22b-4b49-4cc2-a542-ef67c92b1a02",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to commit transaction.",
				Cause:  err.Error(),
			}
		}

		builder = newBuilder()
		dfObject = updateObjectOut.Entity
		recordCount = 0

		return nil
	}

	for scanner.Scan() {
		if recordCount == in.BatchSize {
			if err := deleteRecords(); err != nil {
				return nil, err
			}
		}

		fields := strings.Split(scanner.Text(), in.FieldSeparator)

		record, err := in.DataRecord.New(in.Headers, fields)
		if err != nil {
			return nil, errors.Forward(err, "5192164d-dfbe-4584-a615-4051923f072d")
		}

		recordValues, err := record.SQLValues()
		if err != nil {
			return nil, errors.Forward(err, "98f38053-5fcf-4f97-8eec-907c94df703e")
		}

		if len(recordValues) != 1 {
			return nil, &errors.Object{
				Id:     "0fc5d005-5059-4af6-8756-ad815996871a",
				Code:   errors.Code_INTERNAL,
				Detail: "Invalid number of values to delete.",
				Meta: map[string]any{
					"expected": 1,
					"actual":   len(recordValues),
				},
			}
		}

		values[recordCount] = recordValues[0]

		recordCount++
		processedRecords++
	}

	if err := deleteRecords(); err != nil {
		return nil, err
	}

	out := &entities.LoadDataRecordOutput{
		ProcessedRecords: processedRecords,
	}

	return out, nil
}

func (dom *domain) BatchInsertDataRecord(r *arc.Request, in *entities.LoadDataRecordInput) (*entities.LoadDataRecordOutput, error) {
	newBuilder := func() squirrel.InsertBuilder {
		return squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Insert(in.DataRecord.SQLTable()).
			Columns(in.Columns...)
	}

	var recordCount int32

	builder := newBuilder()
	dfObject := in.DataFileObject
	scanner := in.Scanner
	processedRecords := int64(0)

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "8256c560-0397-486a-a441-94b52141e907")
	}

	insertRecords := func() error {
		// Skip if no records to insert.
		if recordCount == 0 {
			return nil
		}

		sql, args, err := builder.ToSql()
		if err != nil {
			return &errors.Object{
				Id:     "581d50e2-5cfa-4143-a53a-e4871d5f14b7",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to build SQL.",
				Cause:  err.Error(),
			}
		}

		tx, err := pgxPool.Begin(r.Context())
		if err != nil {
			return &errors.Object{
				Id:     "cce383bf-e68f-463c-a242-2711da7465a3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to begin transaction.",
				Cause:  err.Error(),
			}
		}

		// This clone won't replace the original arc.
		r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

		_, err = dom.repository.CreateDataRecords(r, &CreateDataRecordsInput{
			SQL:  sql,
			Args: args,
		})
		if err != nil {
			extutils.RollbackPgxTx(r.Context(), tx, "75b305db-45b0-4b09-96db-e7f45a0609fa")
			return errors.Forward(err, "caa2b1ce-8be2-4b95-9531-249c9c54bc22")
		}

		updateObjectOut, err := dom.UpdateDataFileObject(r, &entities.UpdateDataFileObjectInput{
			Id:          dfObject.Id,
			UpdatedAt:   time.Now(),
			RecordCount: dfObject.RecordCount + recordCount,
		})
		if err != nil {
			extutils.RollbackPgxTx(r.Context(), tx, "127b114f-4af5-49a2-9ecd-3f1f885b0c36")
			return errors.Forward(err, "418e4dcd-3632-48f5-b69f-5d962c3657a0")
		}

		if err := tx.Commit(r.Context()); err != nil {
			return &errors.Object{
				Id:     "b2516dfc-9598-418a-bd8e-7d5d3019f76a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to commit transaction.",
				Cause:  err.Error(),
			}
		}

		builder = newBuilder()
		dfObject = updateObjectOut.Entity
		recordCount = 0

		return nil
	}

	for scanner.Scan() {
		if recordCount == in.BatchSize {
			if err := insertRecords(); err != nil {
				return nil, err
			}
		}

		fields := strings.Split(scanner.Text(), in.FieldSeparator)

		record, err := in.DataRecord.New(in.Headers, fields)
		if err != nil {
			return nil, errors.Forward(err, "ab1ce0c2-d94c-42a3-9563-0aca3962d4d5")
		}

		recordValues, err := record.SQLValues()
		if err != nil {
			return nil, errors.Forward(err, "2c07ea36-d0d3-4c53-9256-2f9da15416cb")
		}

		builder = builder.Values(recordValues...)

		recordCount++
		processedRecords++
	}

	if err := insertRecords(); err != nil {
		return nil, err
	}

	out := &entities.LoadDataRecordOutput{
		ProcessedRecords: processedRecords,
	}

	return out, nil
}
