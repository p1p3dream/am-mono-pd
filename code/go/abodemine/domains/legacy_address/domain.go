package legacy_address

import (
	"abodemine/domains/arc"
	"abodemine/lib/errors"
)

type Domain interface {
	SearchByAddress(r *arc.Request, in *SearchRecordInput, index []string) (*SearchRecordOutput, error)
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

type SearchRecordInput struct {
	FullStreetAddress string `json:"full_street_address,omitempty"`
	HouseNbr          string `json:"house_nbr,omitempty"`
	Street            string `json:"street,omitempty"`
	City              string `json:"city,omitempty"`
	State             string `json:"state,omitempty"`
	ZIP5              string `json:"zip5,omitempty"`
	ZIP4              string `json:"zip4,omitempty"`
	UnitType          string `json:"unit_type,omitempty"`
	UnitNbr           string `json:"unit_nbr,omitempty"`
	DirectionLeft     string `json:"direction_left,omitempty"`
	DirectionRight    string `json:"direction_right,omitempty"`

	Limit         *int     `json:"limit,omitempty"`
	DesiredFields []string `json:"desired_fields,omitempty"`
}

type SearchRecordOutput struct {
	APN               string `json:"APN"`
	Fips              string `json:"Fips"`
	FinalValue        string `json:"FinalValue"`
	LowValue          string `json:"LowValue"`
	HighValue         string `json:"HighValue"`
	ConfidenceScore   int64  `json:"ConfidenceScore"`
	StandardDeviation int64  `json:"StandardDeviation"`
	ValuationDate     string `json:"ValuationDate"`

	// Comparable properties
	Comp1PropertyID string `json:"Comp1PropertyID"`
	Comp2PropertyID string `json:"Comp2PropertyID"`
	Comp3PropertyID string `json:"Comp3PropertyID"`
	Comp4PropertyID string `json:"Comp4PropertyID"`
	Comp5PropertyID string `json:"Comp5PropertyID"`
	Comp6PropertyID string `json:"Comp6PropertyID"`
	Comp7PropertyID string `json:"Comp7PropertyID"`

	// Address fields
	SitusFullStreetAddress string `json:"SitusFullStreetAddress"`
	SitusHouseNbr          string `json:"SitusHouseNbr"`
	SitusHouseNbrSuffix    string `json:"SitusHouseNbrSuffix"`
	SitusStreet            string `json:"SitusStreet"`
	SitusCity              string `json:"SitusCity"`
	SitusState             string `json:"SitusState"`
	SitusZIP5              string `json:"SitusZIP5"`
	SitusZIP4              string `json:"SitusZIP4"`
	SitusUnitType          string `json:"SitusUnitType"`
	SitusUnitNbr           string `json:"SitusUnitNbr"`
	SitusDirectionLeft     string `json:"SitusDirectionLeft"`
	SitusDirectionRight    string `json:"SitusDirectionRight"`
	SitusMode              string `json:"SitusMode"`
	SitusCarrierCode       string `json:"SitusCarrierCode"`

	// IDs
	Aupid      string `json:"aupid"`
	PropertyID string `json:"PropertyID"`
}

func (dom *domain) SearchByAddress(r *arc.Request, in *SearchRecordInput, index []string) (*SearchRecordOutput, error) {
	if in.Limit == nil {
		defaultValue := int(1)
		in.Limit = &defaultValue
	}

	out, err := dom.repository.SearchByAddress(r, in, index)
	if err != nil {
		return nil, &errors.Object{
			Id:     "04fb238b-9be4-4a6d-b229-059b8515ce92",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to search by address",
			Cause:  err.Error(),
		}
	}

	return out, nil
}
