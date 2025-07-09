package entities

import "github.com/google/uuid"

type Property struct {
	Aupid *uuid.UUID `json:"aupid,omitempty"`

	Address  *PropertyAddress `json:"address,omitempty"`
	Assessor *Assessor        `json:"assessor,omitempty"`
	Comps    []*Property      `json:"comps,omitempty"`
	Listing  []*Listing       `json:"listing,omitempty"`
	Recorder []*Recorder      `json:"recorder,omitempty"`
	Rental   *RentalAvm       `json:"rentEstimate,omitempty"`
	Sale     *SaleAvm         `json:"saleEstimate,omitempty"`
}

// PropertyRef holds references to a given
// Abodemine Unique Property Id (AUPID).
type PropertyRef struct {
	Aupid *uuid.UUID

	ADAttomId    *int64
	FAPropertyId *int64
}
