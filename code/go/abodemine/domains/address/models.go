package address

import (
	"time"

	"abodemine/entities"
)

// AddressDocument represents an address document for OpenSearch
type AddressDocument struct {
	// Aupid       string    `json:"aupid" mapping_type:"keyword"`
	AMId        string    `json:"am_id" mapping_type:"keyword"`
	AMUpdatedAt time.Time `json:"am_updated_at" mapping_type:"date"`

	PropertyId   *string `json:"property_id" mapping_type:"keyword"`
	ADAttomId    *int64  `json:"ad_attom_id" mapping_type:"keyword"`
	FAPropertyId *int64  `json:"fa_property_id" mapping_type:"keyword"`

	FullAddress   *string `json:"full_address" mapping_type:"text"`
	FIPS          string  `json:"fips" mapping_type:"keyword"`
	State         *string `json:"state" mapping_type:"keyword"`
	StateFullName *string `json:"state_full_name" mapping_type:"text"`
	County        *string `json:"county" mapping_type:"text"`
	ZIP5          string  `json:"zip5" mapping_type:"keyword"`
	// ZIP4                    *string    `json:"zip4" mapping_type:"keyword"`
	PreDirectional  *string `json:"pre_directional" mapping_type:"keyword"`
	StreetNumber    *string `json:"street_number" mapping_type:"keyword"`
	Street          *string `json:"street" mapping_type:"text"`
	PostDirectional *string `json:"post_directional" mapping_type:"keyword"`
	StreetType      *string `json:"street_type" mapping_type:"keyword"`
	UnitType        *string `json:"unit_type" mapping_type:"keyword"`
	UnitNbr         *string `json:"unit_nbr" mapping_type:"keyword"`
	// VacantIndicator         *string    `json:"vacant_indicator" mapping_type:"keyword"`
	// NonUSPSAddressIndicator *string    `json:"non_usps_address_indicator" mapping_type:"keyword"`
	// NotCurrentlyDeliverable *string    `json:"not_currently_deliverable" mapping_type:"keyword"`
	// CommunityName           *string    `json:"community_name" mapping_type:"keyword"`
	City *string `json:"city" mapping_type:"keyword"`
	// PostalCommunity         *string    `json:"postal_community" mapping_type:"text"`
	// PlaceName               *string    `json:"place_name" mapping_type:"keyword"`
	// SubdivisionName         *string    `json:"subdivision_name" mapping_type:"keyword"`
	Latitude  float64 `json:"-"`
	Longitude float64 `json:"-"`
	// PropertyClassID         *string    `json:"property_class_id" mapping_type:"keyword"`
	// PropertyID              *string    `json:"-"`
	// AddressType             *string    `json:"address_type" mapping_type:"keyword"`
	// FAPropertyID            *string    `json:"fa_property_id" mapping_type:"keyword"`
	// AddressMasterID         string     `json:"address_master_id" mapping_type:"keyword"`
	// LastUpdate              *time.Time `json:"last_update" mapping_type:"date"`
	// DPVFootnotes            *string    `json:"dpv_footnotes" mapping_type:"keyword"`
	// DeliveryPointCheckDigit *string    `json:"delivery_point_check_digit" mapping_type:"keyword"`
	// DeliveryPointCode       *string    `json:"delivery_point_code" mapping_type:"keyword"`
	// DPVCount                *string    `json:"dpv_count" mapping_type:"keyword"`
	Source *string `json:"source" mapping_type:"keyword"`

	// Additional field for geo queries
	Location string `json:"location,omitempty" mapping_type:"geo_point"`
}

// SearchAddressInput represents the input parameters for an address search
type SearchAddressInput struct {
	// Address fields
	FullAddress         string
	HouseNumber         string
	StreetDirection     string
	StreetName          string
	StreetSuffix        string
	StreetPostDirection string
	City                string
	State               string
	Zip5                string
	Unit                string

	// Property identifiers
	Apn      string
	Aupid    string
	OldAupid string

	// Pagination
	Limit *int
	Page  int

	// Field filtering
	IncludeFields  []string
	OmitFields     []string
	IncludeLayouts []string
	OmitLayouts    []string
	Custom         map[string]any

	// Additional filters
	RecorderDocTypes []string
}

// SearchAddressOutput represents the response from an address search
type SearchAddressOutput struct {
	// Basic property information
	HouseNumber         string
	StreetDirection     string
	StreetName          string
	StreetSuffix        string
	StreetPostDirection string
	Zip5                string
	City                string
	State               string
	Aupid               string

	// Additional data
	Assessor     map[string]any
	Hits         int64
	Recorder     [][]map[string]any
	RentEstimate map[string]any
	SaleEstimate map[string]any

	// Pagination info
	Pagination AddressPaginationInfo
}

// AddressPaginationInfo contains pagination metadata for the API response
type AddressPaginationInfo struct {
	Page         int `json:"page,omitempty"`
	Limit        int `json:"limit,omitempty"`
	TotalResults int `json:"totalResults,omitempty"`
}

type PropertyAddressSearchQuery struct {
	Size  int32 `json:"size,omitempty"`
	Query struct {
		Bool struct {
			Must []map[string]any `json:"must"`
		} `json:"bool"`
	} `json:"query"`
}

type PropertyAddressSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source entities.PropertyAddress `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Fips struct {
	Fips string
}

type Zip5 struct {
	Zip5 string
}
