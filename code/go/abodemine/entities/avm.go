package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RentalAvm struct {
	EstimatedRentalValue    *int       `json:"estimatedRentalValue,omitempty"`
	EstimatedMinRentalValue *int       `json:"estimatedMinRentalValue,omitempty"`
	EstimatedMaxRentalValue *int       `json:"estimatedMaxRentalValue,omitempty"`
	ValuationDate           *time.Time `json:"valuationDate,omitempty"`
}

type SaleAvm struct {
	// Id        uuid.UUID      `json:"id,omitempty"`
	// CreatedAt time.Time      `json:"createdAt,omitempty"`
	// UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	// Meta      map[string]any `json:"meta,omitempty"`

	Fips                *string          `json:"fips,omitempty"`
	Apn                 *string          `json:"apn,omitempty"`
	FullStreetAddress   *string          `json:"fullStreetAddress,omitempty"`
	HouseNumber         *string          `json:"houseNumber,omitempty"`
	StreetPreDirection  *string          `json:"streetPreDirection,omitempty"`
	StreetName          *string          `json:"streetName,omitempty"`
	StreetPostDirection *string          `json:"streetPostDirection,omitempty"`
	StreetSuffix        *string          `json:"streetSuffix,omitempty"`
	UnitType            *string          `json:"unitType,omitempty"`
	UnitNumber          *string          `json:"unitNumber,omitempty"`
	City                *string          `json:"city,omitempty"`
	State               *string          `json:"state,omitempty"`
	Zip5                *string          `json:"zip5,omitempty"`
	Zip4                *string          `json:"zip4,omitempty"`
	HouseNumberSuffix   *string          `json:"houseNumberSuffix,omitempty"`
	CarrierCode         *string          `json:"carrierCode,omitempty"`
	FinalValue          *decimal.Decimal `json:"finalValue,omitempty"`
	HighValue           *decimal.Decimal `json:"highValue,omitempty"`
	LowValue            *decimal.Decimal `json:"lowValue,omitempty"`
	ConfidenceScore     *float64         `json:"confidenceScore,string,omitempty"`
	StandardDeviation   *float64         `json:"standardDeviation,string,omitempty"`
	ValuationDate       *time.Time       `json:"valuationDate,omitempty"`
	Comps               []*uuid.UUID     `json:"comps,omitempty"`
}
