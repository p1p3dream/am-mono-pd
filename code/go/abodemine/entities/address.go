package entities

import (
	"time"

	"github.com/google/uuid"
)

type PropertyAddress struct {
	Id        *uuid.UUID     `json:"id,omitempty"`
	CreatedAt *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`

	FullStreetAddress   *string `json:"fullStreetAddress,omitempty"`
	HouseNumber         *string `json:"houseNumber,omitempty"`
	StreetPreDirection  *string `json:"streetPreDirection,omitempty"`
	StreetName          *string `json:"streetName,omitempty"`
	StreetPostDirection *string `json:"streetPostDirection,omitempty"`
	StreetSuffix        *string `json:"streetSuffix,omitempty"`
	UnitType            *string `json:"unitType,omitempty"`
	UnitNumber          *string `json:"unitNumber,omitempty"`
	City                *string `json:"city,omitempty"`
	State               *string `json:"state,omitempty"`
	Zip5                *string `json:"zip5,omitempty"`
	County              *string `json:"county,omitempty"`

	// OpenSearch specific fields.

	// The property id related to this address.
	// Currently, this is stored in meta->>'property_id',
	// and its counter relation is on properties.address_id.
	Aupid         *uuid.UUID `json:"aupid,omitempty"`
	Fips          *string    `json:"fips,omitempty"`
	StateFullName *string    `json:"stateFullName,omitempty"`
}

func (e *PropertyAddress) OpenSearchId() string {
	if e.Id == nil {
		return ""
	}
	return e.Id.String()
}
