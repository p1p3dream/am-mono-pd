package assessor

import (
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type Domain interface {
	SelectAssessor(r *arc.Request, in *SelectAssessorInput) (*SelectAssessorOutput, error)
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

type SelectAssessorInput struct {
	Aupid *uuid.UUID
}

type SelectAssessorOutput struct {
	AssessorEntities []*entities.Assessor
}

func (dom *domain) SelectAssessor(r *arc.Request, in *SelectAssessorInput) (*SelectAssessorOutput, error) {
	selectAssessorRecordOut, err := dom.repository.SelectAssessorRecord(r, &SelectAssessorRecordInput{
		Aupid: in.Aupid,
	})
	if err != nil {
		return nil, errors.Forward(err, "b20dae41-232d-40a8-a13a-d28ae7c33037")
	}

	out := &SelectAssessorOutput{
		AssessorEntities: selectAssessorRecordOut.Records,
	}

	return out, nil
}
