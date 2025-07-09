package address

import (
	"strings"

	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/stringutil"
	"abodemine/lib/val"
	"abodemine/models"
	"abodemine/repositories/opensearch"
)

type Domain interface {
	SelectPropertyAddress(r *arc.Request, in *SelectPropertyAddressInput) (*SelectPropertyAddressOutput, error)

	SelectFips(r *arc.Request, in *SelectFipsInput) (*SelectFipsOutput, error)

	SelectZip5(r *arc.Request, in *SelectZip5Input) (*SelectZip5Output, error)
	UpdateZip5Table(r *arc.Request) error
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

type SearchRecordInput struct {
	FullStreetAddress string   `json:"full_street_address,omitempty"`
	HouseNbr          string   `json:"house_nbr,omitempty"`
	Street            string   `json:"street,omitempty"`
	City              string   `json:"city,omitempty"`
	State             string   `json:"state,omitempty"`
	ZIP5              string   `json:"zip5,omitempty"`
	ZIP4              string   `json:"zip4,omitempty"`
	UnitType          string   `json:"unit_type,omitempty"`
	UnitNbr           string   `json:"unit_nbr,omitempty"`
	DirectionLeft     string   `json:"direction_left,omitempty"`
	DirectionRight    string   `json:"direction_right,omitempty"`
	Limit             *int     `json:"limit,omitempty"`
	DesiredFields     []string `json:"desired_fields,omitempty"`
}

type SearchRecordOutput struct {
	PropertyID string `json:"property_id"`
}

type SelectPropertyAddressInput struct {
	Aupid            *uuid.UUID
	ApiSearchAddress *models.ApiSearchAddress
	Fips             string
	IdGt             *uuid.UUID

	Columns []string
	Limit   uint64
	OrderBy []string

	IncludePropertyRefs      bool
	IncludeStateFullName     bool
	ReturnOpenSearchDocument bool
}

type SelectPropertyAddressOutput struct {
	AddressDocuments    []opensearch.Document
	AddressEntities     []*entities.PropertyAddress
	PropertyRefEntities []*entities.PropertyRef
}

func (dom *domain) SelectPropertyAddress(r *arc.Request, in *SelectPropertyAddressInput) (*SelectPropertyAddressOutput, error) {
	if in.ApiSearchAddress != nil {
		addr := in.ApiSearchAddress
		addr.FullStreetAddress = strings.ToUpper(strings.TrimSpace(addr.FullStreetAddress))
		addr.State = strings.ToUpper(strings.TrimSpace(addr.State))
		addr.City = strings.ToUpper(strings.TrimSpace(addr.City))
		addr.Zip5 = strings.ToUpper(strings.TrimSpace(addr.Zip5))
		addr.HouseNumber = strings.ToUpper(strings.TrimSpace(addr.HouseNumber))
		addr.StreetName = strings.ToUpper(strings.TrimSpace(addr.StreetName))
		addr.StreetSuffix = strings.ToUpper(strings.TrimSpace(addr.StreetSuffix))
		addr.UnitType = strings.ToUpper(strings.TrimSpace(addr.UnitType))
		addr.UnitNumber = strings.ToUpper(strings.TrimSpace(addr.UnitNumber))
		addr.StreetPreDirection = strings.ToUpper(strings.TrimSpace(addr.StreetPreDirection))
		addr.StreetPostDirection = strings.ToUpper(strings.TrimSpace(addr.StreetPostDirection))

		if r.HasFlag(flags.SearchUsingAddressesTable) {
			if addr.FullStreetAddress == "" {
				addr.FullStreetAddress = stringutil.JoinNonEmpty(
					" ",
					addr.HouseNumber,
					addr.StreetPreDirection,
					addr.StreetName,
					addr.StreetSuffix,
					addr.StreetPostDirection,
					addr.UnitType,
					addr.UnitNumber,
				)
			}

			addr.HouseNumber = ""
			addr.StreetPreDirection = ""
			addr.StreetName = ""
			addr.StreetSuffix = ""
			addr.StreetPostDirection = ""
			addr.UnitType = ""
			addr.UnitNumber = ""
		} else {
			searchPropertyRecordOut, err := dom.repository.SearchPropertyAddressRecord(r, &SearchPropertyAddressRecordInput{
				IncludePropertyRefs: in.IncludePropertyRefs,
				ApiSearchAddress:    in.ApiSearchAddress,
				Limit:               1,
			})
			if err != nil {
				return nil, errors.Forward(err, "81356b30-a4b6-44d9-8298-a23ee951a51e")
			}

			return &SelectPropertyAddressOutput{
				AddressEntities:     searchPropertyRecordOut.Records,
				PropertyRefEntities: searchPropertyRecordOut.PropertyRefs,
			}, nil
		}
	}

	selectPropertyAddressRecordOut, err := dom.repository.SelectPropertyAddressRecord(r, &SelectPropertyAddressRecordInput{
		Aupid:                    in.Aupid,
		ApiSearchAddress:         in.ApiSearchAddress,
		Fips:                     in.Fips,
		IdGt:                     in.IdGt,
		Columns:                  in.Columns,
		OrderBy:                  in.OrderBy,
		Limit:                    in.Limit,
		IncludePropertyRefs:      in.IncludePropertyRefs,
		IncludeStateFullName:     in.IncludeStateFullName,
		ReturnOpenSearchDocument: in.ReturnOpenSearchDocument,
	})
	if err != nil {
		return nil, errors.Forward(err, "b09d4173-3046-4aad-8656-ed796b6c3f62")
	}

	out := &SelectPropertyAddressOutput{
		AddressDocuments:    selectPropertyAddressRecordOut.Documents,
		AddressEntities:     selectPropertyAddressRecordOut.Records,
		PropertyRefEntities: selectPropertyAddressRecordOut.PropertyRefs,
	}

	return out, nil
}

type SelectFipsInput struct {
	OrderBy string
}

type SelectFipsOutput struct {
	Models []*Fips
}

func (dom *domain) SelectFips(r *arc.Request, in *SelectFipsInput) (*SelectFipsOutput, error) {
	selectFipsOut, err := dom.repository.SelectFipsRecord(r, &SelectFipsRecordInput{
		OrderBy: in.OrderBy,
	})
	if err != nil {
		return nil, errors.Forward(err, "40ba15f9-4c9d-4a59-8478-908b7dddaccb")
	}

	out := &SelectFipsOutput{
		Models: selectFipsOut.Records,
	}

	return out, nil
}

type SelectZip5Input struct {
	OrderBy string
}

type SelectZip5Output struct {
	Models []*Zip5
}

func (dom *domain) SelectZip5(r *arc.Request, in *SelectZip5Input) (*SelectZip5Output, error) {
	selectZip5Out, err := dom.repository.SelectZip5Record(r, &SelectZip5RecordInput{
		OrderBy: in.OrderBy,
	})
	if err != nil {
		return nil, errors.Forward(err, "cb0d2854-8bbb-428c-b536-628907de570d")
	}

	out := &SelectZip5Output{
		Models: selectZip5Out.Records,
	}

	return out, nil
}

func (dom *domain) UpdateZip5Table(r *arc.Request) error {
	if err := dom.repository.UpdateZip5Table(r); err != nil {
		return errors.Forward(err, "9542dc9e-dbc0-4a38-b348-753ac140ee8d")
	}

	return nil
}
