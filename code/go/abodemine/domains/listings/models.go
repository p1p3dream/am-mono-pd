package listings

import (
	"abodemine/entities"
	"abodemine/models"

	"github.com/google/uuid"
)

type SearchListingsInput struct {
	Statuses            []string
	MinBeds             *int
	MaxBeds             *int
	MinBaths            *float64
	MaxBaths            *float64
	MinPrice            *int
	MaxPrice            *int
	MinSqFt             *int
	MaxSqFt             *int
	MinYearBuilt        *int
	MaxYearBuilt        *int
	MinDOM              *int
	MaxDOM              *int
	MinSaleDate         *string
	MaxSaleDate         *string
	MinSalePrice        *int
	MaxSalePrice        *int
	MinLotSqFt          *float64
	MaxLotSqFt          *float64
	MinLotAcres         *float64
	MaxLotAcres         *float64
	MinGarageSpaces     *int
	MaxGarageSpaces     *int
	MinStories          *int
	MaxStories          *int
	MinStatusChangeDate *string
	MlsPropertyType     *string
	MlsPropertySubType  *string
	AsrPropertySubType  *string
	Zip5Codes           []string
	Ouid                *string
	MlsNumbers          []string

	PageLimit  *int
	PageNumber int

	GeoFilter *GeoFilter

	AddressId *uuid.UUID
}

// GeoFilter contains geo-specific filters
type GeoFilter struct {
	Aupid   string
	Address *models.ApiSearchAddress

	GeoDistance    *GeoDistance
	GeoBoundingBox *GeoBoundingBox
	GeoPolygon     *GeoPolygon
}

type GeoDistance struct {
	Radius   float64
	Location GeoPoint
}

// GeoBoundingBoxFilter represents a rectangular area search
type GeoBoundingBox struct {
	Location *BoundingBoxLocation
}

// BoundingBoxLocation contains the coordinates for the bounding box
type BoundingBoxLocation struct {
	TopLeft     GeoPoint
	BottomRight GeoPoint
}

type GeoPolygon struct {
	Points []GeoPoint
}

// GeoPoint represents a geographic coordinate
type GeoPoint struct {
	Lat float64
	Lon float64
}

// PaginationInfo contains pagination metadata for the API response
type PaginationInfo struct {
	PageNumber   int `json:"pageNumber,omitempty"`
	PageLimit    int `json:"pageLimit,omitempty"`
	TotalResults int `json:"totalResults,omitempty"`
}

// APIResponse structures the API response as specified
type SearchListingsOutput struct {
	PropertyListing []ListingComparison `json:"propertyListing"`
	Pagination      *PaginationInfo     `json:"pagination,omitempty"`
}

type ListingComparison struct {
	entities.Listing

	DistanceMeters float64 `json:"distanceMeters,omitempty"`
	DistanceMiles  float64 `json:"distanceMiles,omitempty"`
	Aupid          *string `json:"aupid,omitempty"`
}
