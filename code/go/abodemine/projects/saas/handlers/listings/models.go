package listings

import (
	"abodemine/domains/listings"
)

type SearchListingsInput struct {
	Aupid string `json:"aupid"`

	Statuses     []string `json:"statuses"`
	MinBeds      int      `json:"min_beds"`
	MaxBeds      int      `json:"max_beds"`
	MinBaths     *float64 `json:"min_baths"`
	MaxBaths     *float64 `json:"max_baths"`
	MinPrice     *int     `json:"min_price"`
	MaxPrice     *int     `json:"max_price"`
	MinSqFt      int      `json:"min_sq_ft"`
	MaxSqFt      int      `json:"max_sq_ft"`
	MinYearBuilt int      `json:"min_year_built"`
	MaxYearBuilt int      `json:"max_year_built"`
	MinDOM       int      `json:"min_dom"`
	MaxDOM       int      `json:"max_dom"`

	PageLimit  *int `json:"pageLimit"`
	PageNumber int  `json:"pageNumber"`

	Address   *AddressFilter `json:"address"`
	GeoFilter *GeoFilter     `json:"geo_filter"`
}

// AddressFilter contains address-specific filters
type AddressFilter struct {
	FullStreetAddress string   `json:"full_street_address,omitempty"`
	HouseNbr          string   `json:"house_nbr,omitempty"`
	Street            string   `json:"street,omitempty"`
	City              string   `json:"city,omitempty"`
	State             string   `json:"state,omitempty"`
	ZIP5              []string `json:"zip5,omitempty"`
	ZIP4              string   `json:"zip4,omitempty"`
	UnitType          string   `json:"unit_type,omitempty"`
	UnitNbr           string   `json:"unit_nbr,omitempty"`
	DirectionLeft     string   `json:"direction_left,omitempty"`
	DirectionRight    string   `json:"direction_right,omitempty"`
}

// GeoFilter contains geo-specific filters
type GeoFilter struct {
	GeoDistance    *GeoDistanceFilter    `json:"geo_distance"`
	GeoBoundingBox *GeoBoundingBoxFilter `json:"geo_bounding_box"`
}

type GeoDistanceFilter struct {
	Radius   float64  `json:"radius"` // e.g. "1km"
	Location GeoPoint `json:"location"`
}

// GeoBoundingBoxFilter represents a rectangular area search
type GeoBoundingBoxFilter struct {
	Location *BoundingBoxLocation `json:"location"`
}

// BoundingBoxLocation contains the coordinates for the bounding box
type BoundingBoxLocation struct {
	TopLeft     GeoPoint `json:"top_left"`
	BottomRight GeoPoint `json:"bottom_right"`
}

// GeoPoint represents a geographic coordinate
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (input *SearchListingsInput) ToDomainModel() listings.SearchListingsInput {
	domainInput := listings.SearchListingsInput{
		// Aupid:        input.Aupid,
		Statuses:     input.Statuses,
		MinBeds:      &input.MinBeds,
		MaxBeds:      &input.MaxBeds,
		MinBaths:     input.MinBaths,
		MaxBaths:     input.MaxBaths,
		MinPrice:     input.MinPrice,
		MaxPrice:     input.MaxPrice,
		MinSqFt:      &input.MinSqFt,
		MaxSqFt:      &input.MaxSqFt,
		MinYearBuilt: &input.MinYearBuilt,
		MaxYearBuilt: &input.MaxYearBuilt,
		MinDOM:       &input.MinDOM,
		MaxDOM:       &input.MaxDOM,
		PageLimit:    input.PageLimit,
		PageNumber:   input.PageNumber,
	}

	if input.Address != nil {
		// addressFilter := &listings.AddressFilter{
		// 	City:              input.Address.City,
		// 	State:             input.Address.State,
		// 	ZIP5:              input.Address.ZIP5,
		// 	UnitType:          input.Address.UnitType,
		// 	UnitNbr:           input.Address.UnitNbr,
		// 	DirectionLeft:     input.Address.DirectionLeft,
		// 	DirectionRight:    input.Address.DirectionRight,
		// 	FullStreetAddress: input.Address.FullStreetAddress,
		// }

		// domainInput.Address = addressFilter
	}

	if input.GeoFilter != nil {
		geoFilter := &listings.GeoFilter{}

		if input.GeoFilter.GeoDistance != nil {
			geoFilter.GeoDistance = &listings.GeoDistance{
				Radius: input.GeoFilter.GeoDistance.Radius,
			}

			geoFilter.GeoDistance.Location = listings.GeoPoint{
				Lat: input.GeoFilter.GeoDistance.Location.Lat,
				Lon: input.GeoFilter.GeoDistance.Location.Lon,
			}
		}

		if input.GeoFilter.GeoBoundingBox != nil {
			geoFilter.GeoBoundingBox = &listings.GeoBoundingBox{}

			if input.GeoFilter.GeoBoundingBox.Location != nil {
				geoFilter.GeoBoundingBox.Location = &listings.BoundingBoxLocation{}

				geoFilter.GeoBoundingBox.Location.TopLeft = listings.GeoPoint{
					Lat: input.GeoFilter.GeoBoundingBox.Location.TopLeft.Lat,
					Lon: input.GeoFilter.GeoBoundingBox.Location.TopLeft.Lon,
				}

				geoFilter.GeoBoundingBox.Location.BottomRight = listings.GeoPoint{
					Lat: input.GeoFilter.GeoBoundingBox.Location.BottomRight.Lat,
					Lon: input.GeoFilter.GeoBoundingBox.Location.BottomRight.Lon,
				}
			}
		}

		domainInput.GeoFilter = geoFilter
	}

	return domainInput
}
