package listings

import (
	"abodemine/domains/listings"
	"abodemine/models"
)

type SearchListingsInput struct {
	Statuses            []string `json:"statuses"`
	MinBeds             *int     `json:"minBeds"`
	MaxBeds             *int     `json:"maxBeds"`
	MinBaths            *float64 `json:"minBaths"`
	MaxBaths            *float64 `json:"maxBaths"`
	MinPrice            *int     `json:"minPrice"`
	MaxPrice            *int     `json:"maxPrice"`
	MinSqFt             *int     `json:"minSqFt"`
	MaxSqFt             *int     `json:"maxSqFt"`
	MinYearBuilt        *int     `json:"minYearBuilt"`
	MaxYearBuilt        *int     `json:"maxYearBuilt"`
	MinDom              *int     `json:"minDom"`
	MaxDom              *int     `json:"maxDom"`
	MinSaleDate         *string  `json:"minSaleDate"`
	MaxSaleDate         *string  `json:"maxSaleDate"`
	MinSalePrice        *int     `json:"minSalePrice"`
	MaxSalePrice        *int     `json:"maxSalePrice"`
	MinLotSqFt          *float64 `json:"minLotSqFt"`
	MaxLotSqFt          *float64 `json:"maxLotSqFt"`
	MinLotAcres         *float64 `json:"minLotAcres"`
	MinStatusChangeDate *string  `json:"minStatusChangeDate"`
	MaxLotAcres         *float64 `json:"maxLotAcres"`
	MinGarageSpaces     *int     `json:"minGarageSpaces"`
	MaxGarageSpaces     *int     `json:"maxGarageSpaces"`
	MinStories          *int     `json:"minStories"`
	MaxStories          *int     `json:"maxStories"`
	MlsPropertyType     *string  `json:"mlsPropertyType"`
	MlsPropertySubType  *string  `json:"mlsPropertySubType"`
	Zip5Codes           []string `json:"zip5Codes"`
	Ouid                *string  `json:"ouid"`
	MlsNumbers          []string `json:"mlsNumbers"`
	AsrPropertySubType  *string  `json:"asrPropertySubType"`

	PageLimit  *int `json:"pageLimit"`
	PageNumber int  `json:"pageNumber"`

	GeoFilter *GeoFilter `json:"geoFilter"`
}

// GeoFilter contains geo-specific filters
type GeoFilter struct {
	Aupid          string                   `json:"aupid"`
	Address        *models.ApiSearchAddress `json:"address"`
	GeoDistance    *GeoDistance             `json:"geoDistance"`
	GeoBoundingBox *GeoBoundingBox          `json:"geoBoundingBox"`
	GeoPolygon     *GeoPolygon              `json:"geoPolygon"`
}

type GeoDistance struct {
	Radius   float64  `json:"radius"`
	Location GeoPoint `json:"location"`
}

// GeoBoundingBoxFilter represents a rectangular area search
type GeoBoundingBox struct {
	Location *BoundingBoxLocation `json:"location"`
}

// BoundingBoxLocation contains the coordinates for the bounding box
type BoundingBoxLocation struct {
	TopLeft     GeoPoint `json:"topLeft"`
	BottomRight GeoPoint `json:"bottomRight"`
}

type GeoPolygon struct {
	Points []GeoPoint `json:"points"`
}

// GeoPoint represents a geographic coordinate
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (input *SearchListingsInput) ToDomainModel() listings.SearchListingsInput {
	domainInput := listings.SearchListingsInput{
		Statuses:            input.Statuses,
		MinBeds:             input.MinBeds,
		MaxBeds:             input.MaxBeds,
		MinBaths:            input.MinBaths,
		MaxBaths:            input.MaxBaths,
		MinPrice:            input.MinPrice,
		MaxPrice:            input.MaxPrice,
		MinSqFt:             input.MinSqFt,
		MaxSqFt:             input.MaxSqFt,
		MinYearBuilt:        input.MinYearBuilt,
		MaxYearBuilt:        input.MaxYearBuilt,
		MinDOM:              input.MinDom,
		MaxDOM:              input.MaxDom,
		MinSaleDate:         input.MinSaleDate,
		MaxSaleDate:         input.MaxSaleDate,
		MinSalePrice:        input.MinSalePrice,
		MaxSalePrice:        input.MaxSalePrice,
		MinLotSqFt:          input.MinLotSqFt,
		MaxLotSqFt:          input.MaxLotSqFt,
		MinLotAcres:         input.MinLotAcres,
		MaxLotAcres:         input.MaxLotAcres,
		MinGarageSpaces:     input.MinGarageSpaces,
		MaxGarageSpaces:     input.MaxGarageSpaces,
		MinStatusChangeDate: input.MinStatusChangeDate,
		PageLimit:           input.PageLimit,
		PageNumber:          input.PageNumber,
		MlsPropertyType:     input.MlsPropertyType,
		MlsPropertySubType:  input.MlsPropertySubType,
		Zip5Codes:           input.Zip5Codes,
		Ouid:                input.Ouid,
		MlsNumbers:          input.MlsNumbers,
		AsrPropertySubType:  input.AsrPropertySubType,
	}

	// explicitly check input.GeoFilter for nil to avoid panic
	if input.GeoFilter == nil {
		domainInput.GeoFilter = &listings.GeoFilter{}

		return domainInput
	}

	geoFilter := &listings.GeoFilter{
		Aupid: input.GeoFilter.Aupid,
	}

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

	if input.GeoFilter.GeoPolygon != nil {
		geoFilter.GeoPolygon = &listings.GeoPolygon{
			Points: make([]listings.GeoPoint, len(input.GeoFilter.GeoPolygon.Points)),
		}

		for i, point := range input.GeoFilter.GeoPolygon.Points {
			geoFilter.GeoPolygon.Points[i] = listings.GeoPoint{
				Lat: point.Lat,
				Lon: point.Lon,
			}
		}

	}

	domainInput.GeoFilter = geoFilter

	if input.GeoFilter.Address != nil {
		domainInput.GeoFilter.Address = input.GeoFilter.Address
	}

	return domainInput
}
