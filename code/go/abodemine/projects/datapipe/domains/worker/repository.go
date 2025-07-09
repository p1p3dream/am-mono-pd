package worker

import (
	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
)

type Repository interface {
	InsertDataFileDirectoryRecord(r *arc.Request, in *InsertDataFileDirectoryRecordInput) (*InsertDataFileDirectoryRecordOutput, error)
	SelectDataFileDirectoryRecord(r *arc.Request, in *SelectDataFileDirectoryRecordInput) (*SelectDataFileDirectoryRecordOutput, error)
	UpdateDataFileDirectoryRecord(r *arc.Request, in *UpdateDataFileDirectoryRecordInput) (*UpdateDataFileDirectoryRecordOutput, error)

	InsertDataFileObjectRecord(r *arc.Request, in *InsertDataFileObjectRecordInput) (*InsertDataFileObjectRecordOutput, error)
	SelectDataFileObjectRecord(r *arc.Request, in *SelectDataFileObjectRecordInput) (*SelectDataFileObjectRecordOutput, error)
	UpdateDataFileObjectRecord(r *arc.Request, in *UpdateDataFileObjectRecordInput) (*UpdateDataFileObjectRecordOutput, error)

	SelectUnparsedDataFileObjectRecords(r *arc.Request, in *SelectUnparsedDataFileObjectRecordsInput) (*SelectUnparsedDataFileObjectRecordsOutput, error)

	CreateDataRecords(r *arc.Request, in *CreateDataRecordsInput) (*CreateDataRecordsOutput, error)
	RemoveDataRecords(r *arc.Request, in *RemoveDataRecordsInput) (*RemoveDataRecordsOutput, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

type CreateDataRecordsInput struct {
	SQL  string
	Args []any
}

type CreateDataRecordsOutput struct{}

func (repo *repository) CreateDataRecords(r *arc.Request, in *CreateDataRecordsInput) (*CreateDataRecordsOutput, error) {
	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, in.SQL, in.Args)
	if err != nil {
		return nil, errors.Forward(err, "281fdbc4-2d19-4ec6-8406-a0a7a756a6e9")
	}

	rows.Close()

	if rows.Err() != nil {
		return nil, &errors.Object{
			Id:     "63186515-cad8-453e-9d6f-0b3ef885d224",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to query rows.",
			Cause:  rows.Err().Error(),
		}
	}

	out := &CreateDataRecordsOutput{}

	return out, nil
}

type RemoveDataRecordsInput struct {
	SQL  string
	Args []any
}

type RemoveDataRecordsOutput struct{}

func (repo *repository) RemoveDataRecords(r *arc.Request, in *RemoveDataRecordsInput) (*RemoveDataRecordsOutput, error) {
	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, in.SQL, in.Args)
	if err != nil {
		return nil, errors.Forward(err, "bec9cf91-ab90-4837-ace5-2d2cba19d850")
	}

	rows.Close()

	if rows.Err() != nil {
		return nil, &errors.Object{
			Id:     "e0782fe1-744b-4d9a-8f4f-a43194825f8a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to query rows.",
			Cause:  rows.Err().Error(),
		}
	}

	out := &RemoveDataRecordsOutput{}

	return out, nil
}
