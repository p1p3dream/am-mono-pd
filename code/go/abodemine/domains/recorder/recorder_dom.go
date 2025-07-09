package recorder

import (
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type Domain interface {
	SelectRecorder(r *arc.Request, in *SelectRecorderInput) (*SelectRecorderOutput, error)
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

type SelectRecorderInput struct {
	Aupid *uuid.UUID
}

type SelectRecorderOutput struct {
	RecorderEntities []*entities.Recorder
}

func (dom *domain) SelectRecorder(r *arc.Request, in *SelectRecorderInput) (*SelectRecorderOutput, error) {
	selectRecorderRecordsOut, err := dom.repository.SelectRecorderRecord(r, &SelectRecorderRecordInput{
		Aupid: in.Aupid,
	})
	if err != nil {
		return nil, errors.Forward(err, "68a0433d-1a00-4550-96dc-aa34f208f881")
	}

	out := &SelectRecorderOutput{
		RecorderEntities: selectRecorderRecordsOut.Records,
	}

	return out, nil
}
