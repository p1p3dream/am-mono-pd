package user

import (
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type Domain interface {
	InsertUser(r *arc.Request, in *InsertUserInput) (*InsertUserOutput, error)
	SelectUser(r *arc.Request, in *SelectUserInput) (*SelectUserOutput, error)
}

type domain struct {
	repository Repository
}

type NewDomainInput struct {
	Repository Repository
}

func NewDomain(in *NewDomainInput) *domain {
	rep := in.Repository

	if rep == nil {
		rep = &repository{}
	}

	return &domain{
		repository: rep,
	}
}

type InsertUserInput struct {
	OrganizationId uuid.UUID
	User           *User
}

type InsertUserOutput struct {
	User *User
}

func (dom *domain) InsertUser(r *arc.Request, in *InsertUserInput) (*InsertUserOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "3c8767a1-0c66-4f08-a6d4-a97b856812f7",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	record := in.User

	if record == nil {
		return nil, &errors.Object{
			Id:     "602b0f63-a2f1-4d20-8545-fdc942a0aee4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing user.",
		}
	}

	if record.OrganizationId == uuid.Nil {
		record.OrganizationId = in.OrganizationId
	}

	if record.Id == uuid.Nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "c4305f9b-a2d3-47a8-8ff2-d4d4584e726c")
		}
		record.Id = id
	}

	insertUserRecordOut, err := dom.repository.InsertUserRecord(r, &InsertUserRecordInput{
		Record: record,
	})
	if err != nil {
		return nil, errors.Forward(err, "f7bbdfac-9a55-4ec4-9a2b-99ec60ac843a")
	}

	out := &InsertUserOutput{
		User: insertUserRecordOut.Record,
	}

	return out, nil
}

type SelectUserInput struct {
	OrganizationId uuid.UUID
	Id             uuid.UUID
	ExternalId     string
}

type SelectUserOutput struct {
	User *User
}

func (dom *domain) SelectUser(r *arc.Request, in *SelectUserInput) (*SelectUserOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "34d52b60-9d01-4748-88f5-03191cce3614",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	selectUserRecordOut, err := dom.repository.SelectUserRecord(r, &SelectUserRecordInput{
		OrganizationId: in.OrganizationId,
		Id:             in.Id,
		ExternalId:     in.ExternalId,
	})
	if err != nil {
		lastErr := errors.Last(err)

		if lastErr.Code == errors.Code_NOT_FOUND {
			return nil, errors.Wrap(err, &errors.Object{
				Id:     "e4bf0f9d-ee15-4729-a853-51c3c6bd4744",
				Code:   errors.Code_NOT_FOUND,
				Detail: "User not found.",
			})
		}

		return nil, errors.Forward(err, "53d9a332-e52c-45c5-b7d0-4902f990c45c")
	}

	out := &SelectUserOutput{
		User: selectUserRecordOut.Record,
	}

	return out, nil
}
