package listings

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/geom"
	"abodemine/lib/ptr"
	"abodemine/models"
	"abodemine/projects/api/domains/auth"
)

const (
	// DefaultLimit is the default number of results to return
	DefaultLimit   = 10
	initialRadius  = 1.0     // initial radius in miles
	radiusLimit    = 100.    // maximum radius in miles
	maxPolygonArea = 30000.0 // maximum area of polygon in square miles
)

type Domain interface {
	SearchListings(r *arc.Request, filters *SearchListingsInput) (*SearchListingsOutput, error)

	SelectListing(r *arc.Request, in *SelectListingInput) (*SelectListingOutput, error)
}

type domain struct {
	arcDomain     arc.Domain
	addressDomain address.Domain
	authDomain    auth.Domain

	repository Repository
}

type NewDomainInput struct {
	ArcDomain     arc.Domain
	AddressDomain address.Domain
	AuthDomain    auth.Domain
	Repository    Repository
}

func NewDomain(in *NewDomainInput) *domain {
	rep := in.Repository

	if rep == nil {
		rep = &repository{}
	}

	return &domain{
		arcDomain:     in.ArcDomain,
		addressDomain: in.AddressDomain,
		authDomain:    in.AuthDomain,
		repository:    rep,
	}
}

func validateListingsSearchFilters(filter *SearchListingsInput) error {
	if filter == nil {
		return &errors.Object{
			Id:     "a68600b6-2850-4e45-8558-105c716d09a5",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Search filters cannot be nil",
			Label:  "invalid_search_filters",
		}
	}

	// Validate status change date format
	if filter.MinStatusChangeDate != nil {
		_, err := time.Parse("2006-01-02", *filter.MinStatusChangeDate)
		if err != nil {
			return &errors.Object{
				Id:     "7049e23f-25ab-4ae0-8257-79119aa14b7f",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: fmt.Sprintf("Invalid minimum status change date format: %s. Expected format: YYYY-MM-DD", *filter.MinStatusChangeDate),
				Label:  "invalid_status_change_date_format",
			}
		}
	}

	// Validate bedroom range
	if filter.MinBeds != nil && filter.MaxBeds != nil && *filter.MinBeds > *filter.MaxBeds {
		return &errors.Object{
			Id:     "34e1686e-b2ef-4321-b9d5-a1644dd8f0ed",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum bedrooms cannot be greater than maximum bedrooms",
			Label:  "invalid_bedroom_range",
		}
	}

	// Validate bathroom range
	if filter.MinBaths != nil && filter.MaxBaths != nil && *filter.MinBaths > *filter.MaxBaths {
		return &errors.Object{
			Id:     "16daa5bd-0a76-492e-9518-7dca679a1452",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum bathrooms cannot be greater than maximum bathrooms",
			Label:  "invalid_bathroom_range",
		}
	}

	// Validate price range
	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		return &errors.Object{
			Id:     "9ba5c7b4-465a-4f44-929d-b43b60ed8af3",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum price cannot be greater than maximum price",
			Label:  "invalid_price_range",
		}
	}

	// Validate square footage range
	if filter.MinSqFt != nil && filter.MaxSqFt != nil && *filter.MinSqFt > *filter.MaxSqFt {
		return &errors.Object{
			Id:     "0393744f-d002-45d9-b400-29f0d23d7b4a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum square footage cannot be greater than maximum square footage",
			Label:  "invalid_square_footage_range",
		}
	}

	// Validate year built range
	if filter.MinYearBuilt != nil && filter.MaxYearBuilt != nil && *filter.MinYearBuilt > *filter.MaxYearBuilt {
		return &errors.Object{
			Id:     "b6be2a89-62a8-4cb3-a525-eae0971c4243",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum year built cannot be greater than maximum year built",
			Label:  "invalid_year_built_range",
		}
	}

	// Validate days on market range
	if filter.MinDOM != nil && filter.MaxDOM != nil && *filter.MinDOM > *filter.MaxDOM {
		return &errors.Object{
			Id:     "e3e29f66-165d-4b99-ae96-0fee220ba45b",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum days on market cannot be greater than maximum days on market",
			Label:  "invalid_days_on_market_range",
		}
	}

	// Validate sale date range
	if filter.MinSaleDate != nil && filter.MaxSaleDate != nil {
		minDate, err := time.Parse("2006-01-02", *filter.MinSaleDate)
		if err != nil {
			return &errors.Object{
				Id:     "657b17ed-6569-4ad9-887c-b25500074f04",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid minimum sale date format",
				Label:  "invalid_sale_date_format",
			}
		}

		maxDate, err := time.Parse("2006-01-02", *filter.MaxSaleDate)
		if err != nil {
			return &errors.Object{
				Id:     "a2759977-fb6f-4e55-99e1-77b08b89e7b7",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid maximum sale date format",
				Label:  "invalid_sale_date_format",
			}
		}

		if minDate.After(maxDate) {
			return &errors.Object{
				Id:     "58d53fe4-3a80-42f7-987c-398ac9ed856d",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Minimum sale date cannot be after maximum sale date",
				Label:  "invalid_sale_date_range",
			}
		}
	}

	// Validate sale price range
	if filter.MinSalePrice != nil && filter.MaxSalePrice != nil && *filter.MinSalePrice > *filter.MaxSalePrice {
		return &errors.Object{
			Id:     "bdfde7f0-8b6b-4c1c-9b22-4dee5b94d4e4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum sale price cannot be greater than maximum sale price",
			Label:  "invalid_sale_price_range",
		}
	}

	// Validate lot size range
	if filter.MinLotSqFt != nil && filter.MaxLotSqFt != nil && *filter.MinLotSqFt > *filter.MaxLotSqFt {
		return &errors.Object{
			Id:     "52d4ee3a-2d83-4478-a185-65a6a080cc49",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum lot size cannot be greater than maximum lot size",
			Label:  "invalid_lot_size_range",
		}
	}

	// Validate lot acres range
	if filter.MinLotAcres != nil && filter.MaxLotAcres != nil && *filter.MinLotAcres > *filter.MaxLotAcres {
		return &errors.Object{
			Id:     "1ea795dd-0500-4359-95e2-1a855373ac48",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum lot acres cannot be greater than maximum lot acres",
			Label:  "invalid_lot_acres_range",
		}
	}

	// Validate garage spaces range
	if filter.MinGarageSpaces != nil && filter.MaxGarageSpaces != nil && *filter.MinGarageSpaces > *filter.MaxGarageSpaces {
		return &errors.Object{
			Id:     "c337f246-c820-4711-86e9-7b729de782d3",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum garage spaces cannot be greater than maximum garage spaces",
			Label:  "invalid_garage_spaces_range",
		}
	}

	// Validate stories range
	if filter.MinStories != nil && filter.MaxStories != nil && *filter.MinStories > *filter.MaxStories {
		return &errors.Object{
			Id:     "dee16348-0e42-405b-8723-6acf4af4e040",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Minimum stories cannot be greater than maximum stories",
			Label:  "invalid_stories_range",
		}
	}

	return nil
}

func (dom *domain) SearchListings(r *arc.Request, in *SearchListingsInput) (*SearchListingsOutput, error) {
	if !r.HasFlag(flags.ApiListingLayoutEnabled) {
		return nil, &errors.Object{
			Id:     "93a343a6-404b-499b-9f83-8e7e8f16581c",
			Code:   errors.Code_PERMISSION_DENIED,
			Detail: "The listing layout is not enabled for this organization.",
		}
	}

	if in.PageNumber < 1 {
		in.PageNumber = 1
	}

	if in.PageLimit != nil && *in.PageLimit > maxResultsPerPage {
		return nil, &errors.Object{
			Id:     "b0f3c4a1-5d8e-4c2b-8f6d-7e9a2f3b5c1e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Page limit must be less than or equal to 100",
			Label:  "invalid_page_limit",
		}
	}

	if in.PageLimit == nil {
		defaultValue := DefaultLimit
		in.PageLimit = &defaultValue
	} else if *in.PageLimit < 1 || *in.PageLimit > 100 {
		return nil, &errors.Object{
			Id:     "e792db31-cda4-45f7-8e75-8e54a0375e90",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Page limit must be between 1 and 100",
			Label:  "invalid_limit",
		}
	}

	isAupidOrAddressProvided := in.GeoFilter.Address != nil || in.GeoFilter.Aupid != ""
	if isAupidOrAddressProvided {
		var err error
		var lat, lon float64

		if in.GeoFilter.Aupid != "" {
			// get the lat long from the aupid
			lat, lon, err = dom.repository.GetLatLonFromAupid(r, in.GeoFilter.Aupid)
			if err != nil {
				return nil, errors.Forward(err,
					"25811001-2929-465a-bd15-ba65c6c02e04")
			}
		}

		if in.GeoFilter.Address != nil {

			in.AddressId, err = dom.getAddressId(r, in.GeoFilter.Address)
			if err != nil {
				return nil, errors.Forward(err, "71ad19d7-8c34-4fde-aa23-59bf014701d3")
			}

			// get the lat long from the address id
			lat, lon, err = dom.repository.GetLatLonFromAddressId(r, *in.AddressId)

			if err != nil {
				return nil, errors.Forward(err, "e2a16c32-b6ab-4504-8691-9393dcfac9f7")
			}
		}

		// check for radius
		radius := initialRadius
		if in.GeoFilter.GeoDistance != nil && in.GeoFilter.GeoDistance.Radius > 0 {
			radius = in.GeoFilter.GeoDistance.Radius
		}

		// build the geo filter
		in.GeoFilter = &GeoFilter{
			GeoDistance: &GeoDistance{
				Radius:   radius,
				Location: GeoPoint{Lat: lat, Lon: lon},
			},
		}
	}

	hasGeoFilter, err := validateGeoFilter(in.GeoFilter)
	if err != nil {
		return nil, errors.Forward(err, "ef9a7430-b66c-4ebc-b1f5-1201240a45d5")
	}

	if err := validateListingsSearchFilters(in); err != nil {
		return nil, errors.Forward(err, "d0a733ef-c7e9-489c-9b7a-326ebe968d2a")
	}

	if !hasGeoFilter && in.AddressId == nil && in.GeoFilter.Aupid == "" && len(in.Zip5Codes) == 0 && len(in.MlsNumbers) == 0 && in.Ouid == nil {
		return nil, &errors.Object{
			Id:     "fc6e6fc6-834d-4719-b963-1120e819770a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "At least one of the following search criteria is required: aupid, address, geoFilter, (ouid + mlsNumbers), zip5Codes.",
			Label:  "invalid_search_criteria",
		}
	}

	response, err := dom.searchListingsWithDynamicRadius(r, in)
	if err != nil {
		return nil, errors.Forward(err, "2d3fd605-7f1d-4b3c-83dc-31badd0e2706")
	}

	_, err = dom.authDomain.InsertApiQuotaTransaction(r, &auth.InsertApiQuotaTransactionInput{
		Entity: &arc.ApiQuotaTransaction{
			Description:         ptr.String("/api/v3/listing"),
			ListingLayoutAmount: int32(len(response.PropertyListing)),
		},
	})
	if err != nil {
		return nil, errors.Forward(err, "b4f73651-13a8-40b8-a630-c63f209a85d4")
	}

	return response, nil
}

func (dom *domain) searchListingsWithDynamicRadius(r *arc.Request, in *SearchListingsInput) (*SearchListingsOutput, error) {
	var out *SearchListingsOutput
	var err error

	maxRadius := 1.
	if in.GeoFilter.GeoDistance != nil && in.GeoFilter.GeoDistance.Radius > 0 {
		maxRadius = in.GeoFilter.GeoDistance.Radius
	}

	out, err = dom.repository.SearchMlsListings(r, in, maxRadius)
	if err != nil {
		return nil, errors.Forward(err, "d39ff000-6929-41ed-bd00-2d7327ff261a")
	}

	if out == nil || len(out.PropertyListing) == 0 {
		return nil, &errors.Object{
			Id:     "c4f2a1b3-5d8e-4c2b-8f6d-7e9a2f3b5c1e",
			Code:   errors.Code_NOT_FOUND,
			Detail: "No listings found",
			Label:  "no_listings_found",
		}

	}

	return &SearchListingsOutput{
		PropertyListing: out.PropertyListing,
		Pagination:      out.Pagination,
	}, nil
}

type SelectListingInput struct {
	Aupid *uuid.UUID
}

type SelectListingOutput struct {
	ListingEntities []*entities.Listing
}

func (dom *domain) SelectListing(r *arc.Request, in *SelectListingInput) (*SelectListingOutput, error) {
	selectListingRecordsOut, err := dom.repository.SelectListingRecord(r, &SelectListingRecordInput{
		Aupid: in.Aupid,
	})
	if err != nil {
		return nil, errors.Forward(err, "0e8b59dd-b497-4591-8087-6414f47cfdc3")
	}

	out := &SelectListingOutput{
		ListingEntities: selectListingRecordsOut.Records,
	}

	return out, nil
}

func (dom *domain) getAddressId(r *arc.Request, inAddress *models.ApiSearchAddress) (*uuid.UUID, error) {
	hasFullAddress := inAddress.FullStreetAddress != ""
	hasPartialAddress := inAddress.HouseNumber != "" && inAddress.StreetName != ""

	if !hasFullAddress && !hasPartialAddress {
		return nil, &errors.Object{
			Id:     "69b653cc-2a9f-4d74-b4d0-e03508035261",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "At least one address search criteria must be provided (Full Address or Zip5 + HouseNumber + StreetName + Others)",
		}
	}

	addressEntity, err := dom.addressDomain.SelectPropertyAddress(r, &address.SelectPropertyAddressInput{
		ApiSearchAddress: inAddress,
	})
	if err != nil {
		return nil, errors.Forward(err, "9ff9c1da-6928-4c15-bd04-af157610123c")
	}

	if addressEntity == nil || len(addressEntity.AddressEntities) == 0 {
		return nil, &errors.Object{
			Id:     "199801f2-4f1a-4846-85fc-fbf7a72adc80",
			Code:   errors.Code_NOT_FOUND,
			Detail: "No property address found",
		}
	}

	return addressEntity.AddressEntities[0].Id, nil
}

func validateGeoFilter(geo *GeoFilter) (bool, error) {
	if geo == nil {
		return false, nil
	}

	hasDistance := geo.GeoDistance != nil
	hasBoundingBox := geo.GeoBoundingBox != nil
	hasPolygon := geo.GeoPolygon != nil

	if (hasDistance && hasBoundingBox) || (hasDistance && hasPolygon) || (hasBoundingBox && hasPolygon) {
		return false, &errors.Object{
			Id:     "604935dd-65ba-41cd-bb73-1c99fe73c5f5",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Only one geo filter can be used at a time.",
		}
	}

	if hasDistance {
		if geo.GeoDistance.Radius > radiusLimit {
			return false, &errors.Object{
				Id:     "e47c2194-905b-468b-b233-7123cf72da9c",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: fmt.Sprintf("Radius must be less than or equal to %f.", radiusLimit),
			}
		}
	}

	if hasPolygon {
		libPoints := make([]geom.Point, 0, len(geo.GeoPolygon.Points))
		for _, v := range geo.GeoPolygon.Points {
			libPoints = append(libPoints, geom.Point{
				Lat: v.Lat,
				Lon: v.Lon,
			})
		}

		if len(libPoints) < 3 {
			return false, &errors.Object{
				Id:     "cb7a075a-4656-4d82-b7d4-8a9b8bb3c8f1",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Polygon must have at least 3 points",
			}
		}

		area := geom.CalculatePolygonArea(libPoints)
		if area > maxPolygonArea {
			return false, &errors.Object{
				Id:     "f64be2a2-9ded-4c3c-bedd-12db15326f89",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: fmt.Sprintf("Polygon area exceeds %f square miles", maxPolygonArea),
			}
		}
	}

	return hasDistance || hasBoundingBox || hasPolygon, nil
}
