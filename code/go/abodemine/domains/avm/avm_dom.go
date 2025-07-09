package avm

import (
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type Domain interface {
	SelectRentalAvm(r *arc.Request, in *SelectRentalAvmInput) (*SelectRentalAvmOutput, error)
	SelectSaleAvm(r *arc.Request, in *SelectSaleAvmInput) (*SelectSaleAvmOutput, error)
}

type domain struct {
	repository Repository
}

type NewDomainInput struct {
	Repository Repository
}

func NewDomain(in *NewDomainInput) Domain {
	return &domain{
		repository: val.Ternary(
			in.Repository == nil,
			NewRepository(),
			in.Repository,
		),
	}
}

type SelectSaleAvmInput struct {
	Aupid *uuid.UUID
}

type SelectSaleAvmOutput struct {
	SaleAvmEntities []*entities.SaleAvm
}

func (dom *domain) SelectSaleAvm(r *arc.Request, in *SelectSaleAvmInput) (*SelectSaleAvmOutput, error) {
	selectSaleAvmRecordsOut, err := dom.repository.SelectSaleAvmRecord(r, &SelectSaleAvmRecordInput{
		Aupid: in.Aupid,
	})
	if err != nil {
		return nil, errors.Forward(err, "42950e2b-0a01-41ce-ac5c-9418336cf2df")
	}

	out := &SelectSaleAvmOutput{
		SaleAvmEntities: selectSaleAvmRecordsOut.Records,
	}

	return out, nil
}
