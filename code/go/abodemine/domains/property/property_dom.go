package property

import (
	"sync/atomic"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/domains/assessor"
	"abodemine/domains/avm"
	"abodemine/domains/listings"
	"abodemine/domains/recorder"
	"abodemine/entities"
	"abodemine/lib/errors"
)

type Domain interface {
	SelectProperty(r *arc.Request, in *SelectPropertyInput) (*SelectPropertyOutput, error)
}

type domain struct {
	addressDomain  address.Domain
	assessorDomain assessor.Domain
	avmDomain      avm.Domain
	listingDomain  listings.Domain
	recorderDomain recorder.Domain
}

type NewDomainInput struct {
	AddressDomain  address.Domain
	AssessorDomain assessor.Domain
	AvmDomain      avm.Domain
	ListingDomain  listings.Domain
	RecorderDomain recorder.Domain
}

func NewDomain(in *NewDomainInput) Domain {
	return &domain{
		addressDomain:  in.AddressDomain,
		assessorDomain: in.AssessorDomain,
		avmDomain:      in.AvmDomain,
		listingDomain:  in.ListingDomain,
		recorderDomain: in.RecorderDomain,
	}
}

type SelectPropertyInput struct {
	Aupids []*uuid.UUID

	IncludeAddress      bool
	IncludeAssessor     bool
	IncludeComps        bool
	IncludeListing      bool
	IncludeRecorder     bool
	IncludeSaleEstimate bool
	IncludeRentEstimate bool
}

type SelectPropertyOutput struct {
	PropertyEntities []*entities.Property

	AddressLayoutSum   int32
	AssessorLayoutSum  int32
	CompsLayoutSum     int32
	ListingLayoutSum   int32
	RecorderLayoutSum  int32
	RentalAvmLayoutSum int32
	SaleAvmLayoutSum   int32
}

func (dom *domain) SelectProperty(r *arc.Request, in *SelectPropertyInput) (*SelectPropertyOutput, error) {
	if len(in.Aupids) == 0 {
		return nil, &errors.Object{
			Id:     "a0f2b1c4-5d3e-4b8e-9f7c-5a6d3f8b1c2e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Aupid is required.",
		}
	}

	aupid := in.Aupids[0]

	if aupid == nil {
		return nil, &errors.Object{
			Id:     "1d420a37-8372-4e12-9f72-104e8a51432e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid aupid.",
		}
	}

	var g errgroup.Group
	var propertyAddressEnt *entities.PropertyAddress
	var assessorEnt *entities.Assessor
	var saleAvmEnt *entities.SaleAvm
	var compsEnts []*entities.Property
	var listingEnts []*entities.Listing
	var recorderEnts []*entities.Recorder
	var rentalAvmEnt *entities.RentalAvm

	var addressLayoutSum atomic.Int32
	var assessorLayoutSum atomic.Int32
	var compsLayoutSum atomic.Int32
	var listingLayoutSum atomic.Int32
	var recorderLayoutSum atomic.Int32
	var rentalAvmLayoutSum atomic.Int32
	var saleAvmLayoutSum atomic.Int32

	// Run 4 goroutines at most.
	g.SetLimit(4)

	if in.IncludeAddress {
		g.Go(func() error {
			selectPropertyAddressOut, err := dom.addressDomain.SelectPropertyAddress(r, &address.SelectPropertyAddressInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "32128858-ef02-447f-9fca-189c315f44ee")
			}

			if len(selectPropertyAddressOut.AddressEntities) > 0 {
				propertyAddressEnt = selectPropertyAddressOut.AddressEntities[0]
				addressLayoutSum.Add(1)
			}

			return nil
		})
	}

	if in.IncludeAssessor {
		g.Go(func() error {
			selectAssessorOut, err := dom.assessorDomain.SelectAssessor(r, &assessor.SelectAssessorInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "58dbb05f-e8ac-4f7e-8088-bea354ede0cf")
			}

			if len(selectAssessorOut.AssessorEntities) > 0 {
				assessorEnt = selectAssessorOut.AssessorEntities[0]
				assessorLayoutSum.Add(1)
			}

			return nil
		})
	}

	if in.IncludeListing {
		g.Go(func() error {
			selectListingOut, err := dom.listingDomain.SelectListing(r, &listings.SelectListingInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "f837bba5-dfd2-41f0-9567-0e197911fa56")
			}

			listingEnts = selectListingOut.ListingEntities
			listingLayoutSum.Add(int32(len(listingEnts)))

			return nil
		})
	}

	if in.IncludeRecorder {
		g.Go(func() error {
			selectRecorderOut, err := dom.recorderDomain.SelectRecorder(r, &recorder.SelectRecorderInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "98dbb429-5532-4790-8310-9379577689cf")
			}

			recorderEnts = selectRecorderOut.RecorderEntities
			recorderLayoutSum.Add(int32(len(recorderEnts)))

			return nil
		})
	}

	if in.IncludeComps || in.IncludeSaleEstimate {
		g.Go(func() error {
			selectSaleAvmOut, err := dom.avmDomain.SelectSaleAvm(r, &avm.SelectSaleAvmInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "f71d43aa-8e03-4a51-be56-3363c567e8d5")
			}

			if len(selectSaleAvmOut.SaleAvmEntities) == 0 {
				return nil
			}

			saleAvmEntity := selectSaleAvmOut.SaleAvmEntities[0]
			saleAvmLayoutSum.Add(1)

			if in.IncludeComps {
				for _, comp := range saleAvmEntity.Comps {
					g.Go(func() error {
						selectPropertyOut, err := dom.SelectProperty(r, &SelectPropertyInput{
							Aupids:              []*uuid.UUID{comp},
							IncludeSaleEstimate: true,
						})
						if err != nil {
							return errors.Forward(err, "ab41d5bf-5f32-427f-8f18-ed3680f7aa52")
						}

						if len(selectPropertyOut.PropertyEntities) > 0 {
							compsEnts = append(compsEnts, selectPropertyOut.PropertyEntities[0])
							compsLayoutSum.Add(1)
						}

						return nil
					})
				}
			}

			if in.IncludeSaleEstimate {
				saleAvmEnt = saleAvmEntity
			}

			return nil
		})
	}

	if in.IncludeRentEstimate {
		g.Go(func() error {
			selectRentalAvmOut, err := dom.avmDomain.SelectRentalAvm(r, &avm.SelectRentalAvmInput{
				Aupid: aupid,
			})
			if err != nil {
				return errors.Forward(err, "28905ed8-0ef6-4033-a407-3a2ef7eb6268")
			}

			if len(selectRentalAvmOut.RentalAvmEntities) > 0 {
				rentalAvmEnt = selectRentalAvmOut.RentalAvmEntities[0]
				rentalAvmLayoutSum.Add(1)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Forward(err, "a39e8145-5644-46bf-a59b-6a8dc0b5a5ce")
	}

	var properties []*entities.Property

	anyLayoutFound := propertyAddressEnt != nil ||
		assessorEnt != nil ||
		len(compsEnts) > 0 ||
		len(listingEnts) > 0 ||
		len(recorderEnts) > 0 ||
		rentalAvmEnt != nil ||
		saleAvmEnt != nil

	if anyLayoutFound {
		properties = []*entities.Property{{
			Aupid:    aupid,
			Address:  propertyAddressEnt,
			Assessor: assessorEnt,
			Comps:    compsEnts,
			Listing:  listingEnts,
			Recorder: recorderEnts,
			Rental:   rentalAvmEnt,
			Sale:     saleAvmEnt,
		}}
	}

	out := &SelectPropertyOutput{
		PropertyEntities: properties,

		AddressLayoutSum:   addressLayoutSum.Load(),
		AssessorLayoutSum:  assessorLayoutSum.Load(),
		CompsLayoutSum:     compsLayoutSum.Load(),
		ListingLayoutSum:   listingLayoutSum.Load(),
		RecorderLayoutSum:  recorderLayoutSum.Load(),
		RentalAvmLayoutSum: rentalAvmLayoutSum.Load(),
		SaleAvmLayoutSum:   saleAvmLayoutSum.Load(),
	}

	return out, nil
}
