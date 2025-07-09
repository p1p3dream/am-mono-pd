package attom_data

import (
	"strconv"

	"abodemine/lib/errors"
	"abodemine/projects/datapipe/entities"
)

type RecorderDelete struct {
	TransactionID int64
}

func (dr *RecorderDelete) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(RecorderDelete)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "TransactionID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "f93c97fa-709a-4f04-ba0a-678144f4d5c7",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.TransactionID = v
		default:
			return nil, &errors.Object{
				Id:     "7e6aa757-687e-4d46-b365-5c2612803313",
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

func (dr *RecorderDelete) SQLColumns() []string {
	return []string{
		"transaction_id",
	}
}

func (dr *RecorderDelete) SQLTable() string {
	return "ad_df_recorder"
}

func (dr *RecorderDelete) SQLValues() ([]any, error) {
	values := []any{
		dr.TransactionID,
	}

	return values, nil
}

func (dr *RecorderDelete) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		Mode: entities.DataRecordModeBatchDelete,
	}
}
