package avm

import (
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
)

type SelectRentalAvmInput struct {
	Aupid *uuid.UUID
}

type SelectRentalAvmOutput struct {
	RentalAvmEntities []*entities.RentalAvm
}

func (dom *domain) SelectRentalAvm(r *arc.Request, in *SelectRentalAvmInput) (*SelectRentalAvmOutput, error) {
	selectRentalAvmRecordsOut, err := dom.repository.SelectRentalAvmRecord(r, &SelectRentalAvmRecordInput{
		Aupid: in.Aupid,
	})
	if err != nil {
		return nil, errors.Forward(err, "7f17c5e8-30d5-447b-a15a-60858a72fe54")
	}

	out := &SelectRentalAvmOutput{
		RentalAvmEntities: selectRentalAvmRecordsOut.Records,
	}

	return out, nil
}
