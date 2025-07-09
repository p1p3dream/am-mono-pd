package attom_data

import (
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/projects/datapipe/entities"
)

type PropertyDelete struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any

	ATTOMID int64
}

func (dr *PropertyDelete) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(PropertyDelete)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "[ATTOM ID]":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "66c4eef2-fa56-43e5-8eda-1471bb8bd1cd",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.ATTOMID = v
		default:
			return nil, &errors.Object{
				Id:     "0b2ac439-8072-4149-a364-1ecd39288139",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Unknown header.",
				Meta: map[string]any{
					"field_index": k,
					"field_value": field,
					"header":      header,
				},
			}
		}
	}

	return record, nil
}

func (dr *PropertyDelete) SQLColumns() []string {
	return []string{
		"attomid",
	}
}

func (dr *PropertyDelete) SQLTable() string {
	return "ad_df_assessor"
}

func (dr *PropertyDelete) SQLValues() ([]any, error) {
	values := []any{
		dr.ATTOMID,
	}

	return values, nil
}

func (dr *PropertyDelete) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		LoadFunc: dr.LoadFunc,
		Mode:     entities.DataRecordModeLoadFunc,
	}
}

func (dr *PropertyDelete) LoadFunc(r *arc.Request, in *entities.LoadDataRecordInput) (*entities.LoadDataRecordOutput, error) {
	newBuilder := func() squirrel.UpdateBuilder {
		return squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Update(in.DataRecord.SQLTable()).
			Set("am_deleted_at", time.Now())
	}

	builder := newBuilder()
	column := in.Columns[0]
	dfObject := in.DataFileObject
	recordCount := int32(0)
	processedRecords := int64(0)
	scanner := in.Scanner
	values := make([]any, in.BatchSize)

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "0f6e8991-9330-4b70-90ed-8a7e7abc6f7f")
	}

	loadRecords := func() error {
		// Skip if no records to delete.
		if recordCount == 0 {
			return nil
		}

		builder = builder.Where(squirrel.Eq{column: values[:recordCount]})

		sql, args, err := builder.ToSql()
		if err != nil {
			return &errors.Object{
				Id:     "505b6a3f-25db-48e6-ab40-6e54a98c9973",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to build SQL.",
				Cause:  err.Error(),
			}
		}

		tx, err := pgxPool.Begin(r.Context())
		if err != nil {
			return &errors.Object{
				Id:     "ce37747f-fb6c-499a-b14b-13c9fb95c7dc",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to begin transaction.",
				Cause:  err.Error(),
			}
		}

		defer extutils.RollbackPgxTx(r.Context(), tx, "4ca0b94b-f3e4-44a6-a9e3-835f5e07eac0")

		// This clone won't replace the original arc.
		r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

		rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
		if err != nil {
			return errors.Forward(err, "6383fc1c-042b-432e-b7dc-765d1001c387")
		}

		rows.Close()

		if rows.Err() != nil {
			return &errors.Object{
				Id:     "b5d2d2d9-b576-4658-837c-ac9e9e1cb242",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query rows.",
				Cause:  rows.Err().Error(),
			}
		}

		updateObjectOut, err := in.UpdateDataFileObjectFunc(r, &entities.UpdateDataFileObjectInput{
			Id:          dfObject.Id,
			UpdatedAt:   time.Now(),
			RecordCount: dfObject.RecordCount + recordCount,
			Status:      entities.DataFileObjectStatusInProgress,
		})
		if err != nil {
			return errors.Forward(err, "c985dd82-d1aa-4007-b705-ac0d5375b7b4")
		}

		if err := tx.Commit(r.Context()); err != nil {
			return &errors.Object{
				Id:     "8f405162-772b-43bd-b2df-b235b0c2969d",
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
			if err := loadRecords(); err != nil {
				return nil, err
			}
		}

		fields := strings.Split(scanner.Text(), in.FieldSeparator)

		record, err := in.DataRecord.New(in.Headers, fields)
		if err != nil {
			return nil, errors.Forward(err, "36b871d5-3e7b-4100-a76a-f38df79d40dd")
		}

		recordValues, err := record.SQLValues()
		if err != nil {
			return nil, errors.Forward(err, "689878cf-255d-4e05-9176-bb7d7c1b247d")
		}

		if len(recordValues) != 1 {
			return nil, &errors.Object{
				Id:     "7704d81e-957b-4827-9d0f-4a252912656e",
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

	if err := loadRecords(); err != nil {
		return nil, err
	}

	out := &entities.LoadDataRecordOutput{
		ProcessedRecords: processedRecords,
	}

	return out, nil
}
