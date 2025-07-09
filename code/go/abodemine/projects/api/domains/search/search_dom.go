package search

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/domains/property"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/ptr"
	"abodemine/models"
	"abodemine/projects/api/domains/auth"
)

const (
	// DefaultLimit is the default number of results to return
	DefaultLimit int = 10
)

type Domain interface {
	SearchProperty(r *arc.Request, in *SearchPropertyInput) (*SearchPropertyOutput, error)
}

type domain struct {
	addressDomain  address.Domain
	authDomain     auth.Domain
	propertyDomain property.Domain
}

type NewDomainInput struct {
	AddressDomain  address.Domain
	AuthDomain     auth.Domain
	PropertyDomain property.Domain
}

func NewDomain(in *NewDomainInput) *domain {
	return &domain{
		addressDomain:  in.AddressDomain,
		authDomain:     in.AuthDomain,
		propertyDomain: in.PropertyDomain,
	}
}

type SearchPropertyInput struct {
	Aupid   *uuid.UUID
	Layouts []string

	ApiSearchAddress *models.ApiSearchAddress
}

type SearchPropertyOutput struct {
	PropertyEntities []*entities.Property
}

func (dom *domain) SearchProperty(r *arc.Request, in *SearchPropertyInput) (*SearchPropertyOutput, error) {
	searchAddress := in.ApiSearchAddress

	if in.Aupid == nil {
		if searchAddress == nil {
			return nil, &errors.Object{
				Id:     "376ed856-623d-48e4-944f-23dd75577436",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "At least one search criteria must be provided (Aupid or Address).",
			}
		}

		hasFullAddress := searchAddress.FullStreetAddress != ""
		hasHouseAndStreet := searchAddress.HouseNumber != "" && searchAddress.StreetName != ""
		hasCityAndState := searchAddress.City != "" && searchAddress.State != ""
		hasZip5 := searchAddress.Zip5 != ""

		hasValidSearchAddress := (hasFullAddress || hasHouseAndStreet) && (hasCityAndState || hasZip5)
		if !hasValidSearchAddress {
			return nil, &errors.Object{
				Id:     "9f19a3a3-10ff-4039-a819-f044ea46c89b",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "A valid search address requires a street component (Full Address or House Number with Street Name) and a location component (City and State, or Zip5).",
			}
		}
	}

	selectPropertyInput := &property.SelectPropertyInput{}

	for _, layout := range in.Layouts {
		switch strings.ToUpper(strings.TrimSpace(layout)) {
		case "ADDRESS":
			if !r.HasFlag(flags.ApiAddressLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "51f1e6e4-82f2-448c-9338-a367b86eb417",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The address layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeAddress = true
		case "ASSESSOR":
			if !r.HasFlag(flags.ApiAssessorLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "ea1c2af5-46ef-4f99-b444-06da417d6a64",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The assessor layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeAssessor = true
		case "COMPS":
			if !r.HasFlag(flags.ApiCompsLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "0812546d-6159-4cff-a400-ee3e2a1a4b72",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The comps layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeComps = true
		case "LISTING":
			if !r.HasFlag(flags.ApiListingLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "7c10e640-c00b-4881-9832-059c8a786fcb",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The listing layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeListing = true
		case "RECORDER":
			if !r.HasFlag(flags.ApiRecorderLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "f9c12746-0671-46e5-9256-f688ec7d5ba8",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The recorder layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeRecorder = true
		case "RENTESTIMATE":
			if !r.HasFlag(flags.ApiRentEstimateLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "5ef02783-993d-44c1-aa2b-ab337c7e4bf5",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The rentEstimate layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeRentEstimate = true
		case "SALEESTIMATE":
			if !r.HasFlag(flags.ApiSaleEstimateLayoutEnabled) {
				return nil, &errors.Object{
					Id:     "d1c10dd7-5bd7-4ec2-82df-40b47bca5e0a",
					Code:   errors.Code_PERMISSION_DENIED,
					Detail: "The saleEstimate layout is not enabled for this organization.",
				}
			}
			selectPropertyInput.IncludeSaleEstimate = true
		default:
			return nil, &errors.Object{
				Id:     "004b8aeb-27e1-4a48-a8cb-4a1e96c1f8f5",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: fmt.Sprintf("Invalid layout type: %s.", layout),
			}
		}
	}

	var aupids []*uuid.UUID

	out := &SearchPropertyOutput{}

	if in.Aupid != nil {
		aupids = []*uuid.UUID{in.Aupid}
	} else {
		selectPropertyAddressOut, err := dom.addressDomain.SelectPropertyAddress(r, &address.SelectPropertyAddressInput{
			IncludePropertyRefs: true,
			ApiSearchAddress:    in.ApiSearchAddress,
		})
		if err != nil {
			return nil, errors.Forward(err, "7f8829a4-d047-4566-9ecf-9b07205e03fb")
		}

		if len(selectPropertyAddressOut.PropertyRefEntities) == 0 {
			return out, nil
		}

		aupids = make([]*uuid.UUID, len(selectPropertyAddressOut.PropertyRefEntities))

		for i, propertyRef := range selectPropertyAddressOut.PropertyRefEntities {
			aupids[i] = propertyRef.Aupid
		}
	}

	selectPropertyInput.Aupids = aupids

	selectPropertyOut, err := dom.propertyDomain.SelectProperty(r, selectPropertyInput)
	if err != nil {
		return nil, errors.Forward(err, "67664b7c-7f53-4028-bf70-3120af907ce5")
	}

	_, err = dom.authDomain.InsertApiQuotaTransaction(r, &auth.InsertApiQuotaTransactionInput{
		Entity: &arc.ApiQuotaTransaction{
			Description:              ptr.String("/api/v3/search"),
			AddressLayoutAmount:      selectPropertyOut.AddressLayoutSum,
			AssessorLayoutAmount:     selectPropertyOut.AssessorLayoutSum,
			CompsLayoutAmount:        selectPropertyOut.CompsLayoutSum,
			ListingLayoutAmount:      selectPropertyOut.ListingLayoutSum,
			RecorderLayoutAmount:     selectPropertyOut.RecorderLayoutSum,
			RentEstimateLayoutAmount: selectPropertyOut.RentalAvmLayoutSum,
			SaleEstimateLayoutAmount: selectPropertyOut.SaleAvmLayoutSum,
		},
	})
	if err != nil {
		return nil, errors.Forward(err, "7c5e8f3c-b982-4429-ae19-a42b5a182935")
	}

	out.PropertyEntities = selectPropertyOut.PropertyEntities

	return out, nil
}
