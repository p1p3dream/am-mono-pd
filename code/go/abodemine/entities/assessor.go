package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Assessor struct {
	// Id        uuid.UUID      `json:"id,omitempty"`
	// CreatedAt time.Time      `json:"createdAt,omitempty"`
	// UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	// Meta      map[string]any `json:"meta,omitempty"`

	Fips                                *string          `json:"fips,omitempty"`
	Apn                                 *string          `json:"apn,omitempty"`
	ApnSeqNbr                           *string          `json:"apnSeqNbr,omitempty"`
	TaxAccountNumber                    *string          `json:"taxAccountNumber,omitempty"`
	FullStreetAddress                   *string          `json:"fullStreetAddress,omitempty"`
	HouseNumber                         *string          `json:"houseNumber,omitempty"`
	StreetPreDirection                  *string          `json:"streetPreDirection,omitempty"`
	StreetName                          *string          `json:"streetName,omitempty"`
	StreetPostDirection                 *string          `json:"streetPostDirection,omitempty"`
	StreetSuffix                        *string          `json:"streetSuffix,omitempty"`
	UnitType                            *string          `json:"unitType,omitempty"`
	UnitNumber                          *string          `json:"unitNumber,omitempty"`
	City                                *string          `json:"city,omitempty"`
	State                               *string          `json:"state,omitempty"`
	Zip5                                *string          `json:"zip5,omitempty"`
	Zip4                                *string          `json:"zip4,omitempty"`
	County                              *string          `json:"county,omitempty"`
	HouseNumberSuffix                   *string          `json:"houseNumberSuffix,omitempty"`
	CarrierCode                         *string          `json:"carrierCode,omitempty"`
	Latitude                            *float64         `json:"latitude,string,omitempty"`
	Longitude                           *float64         `json:"longitude,string,omitempty"`
	GeoStatusCode                       *string          `json:"geoStatusCode,omitempty"`
	PropertyJurisdictionName            *string          `json:"propertyJurisdictionName,omitempty"`
	CensusFipsPlaceCode                 *string          `json:"censusFipsPlaceCode,omitempty"`
	PropertyClassId                     *string          `json:"propertyClassId,omitempty"`
	LandUseCode                         *string          `json:"landUseCode,omitempty"`
	StateLandUseCode                    *string          `json:"stateLandUseCode,omitempty"`
	CountyLandUseCode                   *string          `json:"countyLandUseCode,omitempty"`
	Zoning                              *string          `json:"zoning,omitempty"`
	CensusTract                         *string          `json:"censusTract,omitempty"`
	CensusBlock                         *int             `json:"censusBlock,omitempty"`
	HasMobileHome                       *bool            `json:"hasMobileHome,omitempty"`
	IsTimeShare                         *bool            `json:"isTimeShare,omitempty"`
	SchoolDistrictName                  *string          `json:"schoolDistrictName,omitempty"`
	CombinedStatisticalArea             *string          `json:"combinedStatisticalArea,omitempty"`
	CbsaName                            *string          `json:"cbsaName,omitempty"`
	CbsaCode                            *int             `json:"cbsaCode,omitempty"`
	MsaName                             *string          `json:"msaName,omitempty"`
	MsaCode                             *int             `json:"msaCode,omitempty"`
	MetropolitanDivision                *string          `json:"metropolitanDivision,omitempty"`
	MinorCivilDivisionName              *string          `json:"minorCivilDivisionName,omitempty"`
	MinorCivilDivisionCode              *int             `json:"minorCivilDivisionCode,omitempty"`
	NeighborhoodCode                    *string          `json:"neighborhoodCode,omitempty"`
	CensusBlockGroup                    *int             `json:"censusBlockGroup,omitempty"`
	LotSizeFrontageFeet                 *string          `json:"lotSizeFrontageFeet,omitempty"`
	LotSizeDepthFeet                    *int             `json:"lotSizeDepthFeet,omitempty"`
	LotSizeAcres                        *string          `json:"lotSizeAcres,omitempty"`
	LotSizeSqFt                         *int             `json:"lotSizeSqFt,omitempty"`
	IsOwner1Corp                        *bool            `json:"isOwner1Corp,omitempty"`
	IsOwner1Trust                       *bool            `json:"isOwner1Trust,omitempty"`
	Owner1TypeDescription               *string          `json:"owner1TypeDescription,omitempty"`
	Owner1FullName                      *string          `json:"owner1FullName,omitempty"`
	Owner1LastName                      *string          `json:"owner1LastName,omitempty"`
	Owner1FirstName                     *string          `json:"owner1FirstName,omitempty"`
	Owner1MiddleName                    *string          `json:"owner1MiddleName,omitempty"`
	Owner1Suffix                        *string          `json:"owner1Suffix,omitempty"`
	IsOwner2Corp                        *bool            `json:"isOwner2Corp,omitempty"`
	Owner2FullName                      *string          `json:"owner2FullName,omitempty"`
	Owner2LastName                      *string          `json:"owner2LastName,omitempty"`
	Owner2FirstName                     *string          `json:"owner2FirstName,omitempty"`
	Owner2MiddleName                    *string          `json:"owner2MiddleName,omitempty"`
	Owner2Suffix                        *string          `json:"owner2Suffix,omitempty"`
	OwnerOccupied                       *bool            `json:"ownerOccupied,omitempty"`
	Owner1OwnershipRights               *string          `json:"owner1OwnershipRights,omitempty"`
	Owner3FullName                      *string          `json:"owner3FullName,omitempty"`
	Owner3LastName                      *string          `json:"owner3LastName,omitempty"`
	Owner3FirstName                     *string          `json:"owner3FirstName,omitempty"`
	Owner3MiddleName                    *string          `json:"owner3MiddleName,omitempty"`
	Owner3Suffix                        *string          `json:"owner3Suffix,omitempty"`
	Owner4FullName                      *string          `json:"owner4FullName,omitempty"`
	Owner4LastName                      *string          `json:"owner4LastName,omitempty"`
	Owner4FirstName                     *string          `json:"owner4FirstName,omitempty"`
	Owner4MiddleName                    *string          `json:"owner4MiddleName,omitempty"`
	Owner4Suffix                        *string          `json:"owner4Suffix,omitempty"`
	MailingFullStreetAddress            *string          `json:"mailingFullStreetAddress,omitempty"`
	MailingHouseNumber                  *string          `json:"mailingHouseNumber,omitempty"`
	MailingStreetPreDirection           *string          `json:"mailingStreetPreDirection,omitempty"`
	MailingStreet                       *string          `json:"mailingStreet,omitempty"`
	MailingStreetPostDirection          *string          `json:"mailingStreetPostDirection,omitempty"`
	MailingStreetSuffix                 *string          `json:"mailingMode,omitempty"`
	MailingUnitType                     *string          `json:"mailingUnitType,omitempty"`
	MailingUnitNumber                   *string          `json:"mailingUnitNumber,omitempty"`
	MailingCity                         *string          `json:"mailingCity,omitempty"`
	MailingState                        *string          `json:"mailingState,omitempty"`
	MailingZip5                         *string          `json:"mailingZip5,omitempty"`
	MailingZip4                         *string          `json:"mailingZip4,omitempty"`
	MailingCounty                       *string          `json:"mailingCounty,omitempty"`
	MailingHouseNumberSuffix            *string          `json:"mailingHouseNumberSuffix,omitempty"`
	MailingCarrierCode                  *string          `json:"mailingCarrierCode,omitempty"`
	MailingOptOut                       *bool            `json:"mailingOptOut,omitempty"`
	MailingCoName                       *string          `json:"mailingCoName,omitempty"`
	IsMailingForeignAddress             *string          `json:"isMailingForeignAddress,omitempty"`
	MailingFips                         *string          `json:"mailingFips,omitempty"`
	MailAddressInfoFormat               *string          `json:"mailAddressInfoFormat,omitempty"`
	AssessedTotalValue                  *int             `json:"assessedTotalValue,omitempty"`
	AssessedLandValue                   *int             `json:"assessedLandValue,omitempty"`
	AssessedImprovementValue            *int             `json:"assessedImprovementValue,omitempty"`
	TaxAssessedImprovementsPerc         *decimal.Decimal `json:"taxAssessedImprovementsPerc,omitempty"`
	MarketTotalValue                    *int             `json:"marketTotalValue,omitempty"`
	MarketValueLand                     *int             `json:"marketValueLand,omitempty"`
	MarketValueImprovement              *int             `json:"marketValueImprovement,omitempty"`
	TaxMarketImprovementsPerc           *decimal.Decimal `json:"taxMarketImprovementsPerc,omitempty"`
	TaxAmount                           *int64           `json:"taxAmount,string,omitempty"`
	TaxYear                             *int             `json:"taxYear,omitempty"`
	TaxDeliquentYear                    *int             `json:"taxDeliquentYear,omitempty"`
	MarketYear                          *int             `json:"marketYear,omitempty"`
	AssessedYear                        *int             `json:"assessedYear,omitempty"`
	PreviousAssessedValue               *int             `json:"previousAssessedValue,omitempty"`
	TaxRateCodeArea                     *string          `json:"taxRateCodeArea,omitempty"`
	SchoolTaxDistrict1Code              *string          `json:"schoolTaxDistrict1Code,omitempty"`
	SchoolTaxDistrict2Code              *string          `json:"schoolTaxDistrict2Code,omitempty"`
	SchoolTaxDistrict3Code              *string          `json:"schoolTaxDistrict3Code,omitempty"`
	HasTaxExemptionHomeowner            *bool            `json:"hasTaxExemptionHomeowner,omitempty"`
	HasTaxExemptionVeteran              *bool            `json:"hasTaxExemptionVeteran,omitempty"`
	HasTaxExemptionDisabled             *bool            `json:"hasTaxExemptionDisabled,omitempty"`
	HasTaxExemptionWidow                *bool            `json:"hasTaxExemptionWidow,omitempty"`
	HasTaxExemptionSenior               *bool            `json:"hasTaxExemptionSenior,omitempty"`
	HasTaxExemptionSchoolCollege        *bool            `json:"hasTaxExemptionSchoolCollege,omitempty"`
	HasTaxExemptionReligious            *bool            `json:"hasTaxExemptionReligious,omitempty"`
	HasTaxExemptionWelfare              *bool            `json:"hasTaxExemptionWelfare,omitempty"`
	HasTaxExemptionPublicUtility        *bool            `json:"hasTaxExemptionPublicUtility,omitempty"`
	HasTaxExemptionCemetery             *bool            `json:"hasTaxExemptionCemetery,omitempty"`
	HasTaxExemptionHospital             *bool            `json:"hasTaxExemptionHospital,omitempty"`
	HasTaxExemptionLibrary              *bool            `json:"hasTaxExemptionLibrary,omitempty"`
	BuildingArea                        *decimal.Decimal `json:"buildingArea,omitempty"`
	BuildingAreaType                    *string          `json:"buildingAreaType,omitempty"`
	SumBuildingSqFt                     *decimal.Decimal `json:"sumBuildingSqFt,omitempty"`
	SumLivingAreaSqFt                   *int             `json:"sumLivingAreaSqFt,omitempty"`
	SumGroundFloorSqFt                  *int             `json:"sumGroundFloorSqFt,omitempty"`
	Area2ndFloor                        *int             `json:"area2ndFloor,omitempty"`
	AreaUpperFloors                     *int             `json:"areaUpperFloors,omitempty"`
	BuildingsCount                      *int             `json:"buildingsCount,omitempty"`
	SumGrossAreaSqFt                    *int             `json:"sumGrossAreaSqFt,omitempty"`
	SumAdjAreaSqFt                      *int             `json:"sumAdjAreaSqFt,omitempty"`
	AtticSqFt                           *int             `json:"atticSqFt,omitempty"`
	AtticUnfinishedSqFt                 *int             `json:"atticUnfinishedSqFt,omitempty"`
	AtticFinishedSqFt                   *int             `json:"atticFinishedSqFt,omitempty"`
	HasRoomsAttic                       *bool            `json:"hasRoomsAttic,omitempty"`
	SumBasementSqFt                     *int             `json:"sumBasementSqFt,omitempty"`
	BasementUnfinishedSqFt              *int             `json:"basementUnfinishedSqFt,omitempty"`
	BasementFinishedSqFt                *int             `json:"basementFinishedSqFt,omitempty"`
	SumGarageSqFt                       *int             `json:"sumGarageSqFt,omitempty"`
	ParkingCarportArea                  *int             `json:"parkingCarportArea,omitempty"`
	GarageUnfinishedSqFt                *int             `json:"garageUnfinishedSqFt,omitempty"`
	GarageFinishedSqFt                  *int             `json:"garageFinishedSqFt,omitempty"`
	YearBuilt                           *int             `json:"yearBuilt,omitempty"`
	EffectiveYearBuilt                  *int             `json:"effectiveYearBuilt,omitempty"`
	Bedrooms                            *decimal.Decimal `json:"bedrooms,omitempty"`
	TotalRooms                          *decimal.Decimal `json:"totalRooms,omitempty"`
	BathTotalCalc                       *decimal.Decimal `json:"bathTotalCalc,omitempty"`
	BathFull                            *decimal.Decimal `json:"bathFull,omitempty"`
	BathsPartialNbr                     *decimal.Decimal `json:"bathsPartialNbr,omitempty"`
	BathFixturesNbr                     *decimal.Decimal `json:"bathFixturesNbr,omitempty"`
	Amenities                           *string          `json:"amenities,omitempty"`
	AirConditioningCode                 *int             `json:"airConditioningCode,omitempty"`
	BasementCode                        *int             `json:"basementCode,omitempty"`
	BuildingClassCode                   *int             `json:"buildingClassCode,omitempty"`
	BuildingConditionCode               *int             `json:"buildingConditionCode,omitempty"`
	ConstructionTypeCode                *int             `json:"constructionTypeCode,omitempty"`
	Foundation                          *int             `json:"foundation,omitempty"`
	HasDeck                             *bool            `json:"hasDeck,omitempty"`
	DeckArea                            *int             `json:"deckArea,omitempty"`
	ExteriorWallsCode                   *int             `json:"exteriorWallsCode,omitempty"`
	InteriorWallsCode                   *int             `json:"interiorWallsCode,omitempty"`
	FireplaceCode                       *int             `json:"fireplaceCode,omitempty"`
	HasFireplace                        *int             `json:"hasFireplace,omitempty"`
	FloorCoverCode                      *string          `json:"floorCoverCode,omitempty"`
	Garage                              *int             `json:"garage,omitempty"`
	ParkingCarport                      *bool            `json:"parkingCarport,omitempty"`
	HeatCode                            *int             `json:"heatCode,omitempty"`
	HeatingFuelTypeCode                 *int             `json:"heatingFuelTypeCode,omitempty"`
	SiteInfluenceCode                   *string          `json:"siteInfluenceCode,omitempty"`
	GarageParkingNbr                    *int             `json:"garageParkingNbr,omitempty"`
	DrivewayCode                        *string          `json:"drivewayCode,omitempty"`
	DrivewayArea                        *int             `json:"drivewayArea,omitempty"`
	OtherRooms                          *string          `json:"otherRooms,omitempty"`
	PatioCode                           *int             `json:"patioCode,omitempty"`
	PatioArea                           *int             `json:"patioArea,omitempty"`
	BalconyArea                         *int             `json:"balconyArea,omitempty"`
	CourtyardArea                       *int             `json:"courtyardArea,omitempty"`
	CanopyArea                          *int             `json:"canopyArea,omitempty"`
	GazeboArea                          *int             `json:"gazeboArea,omitempty"`
	PoolCode                            *int             `json:"poolCode,omitempty"`
	PoolArea                            *int             `json:"poolArea,omitempty"`
	PorchCode                           *int             `json:"porchCode,omitempty"`
	PorchArea                           *int             `json:"porchArea,omitempty"`
	BuildingQualityCode                 *int             `json:"buildingQualityCode,omitempty"`
	RoofCoverCode                       *int             `json:"roofCoverCode,omitempty"`
	RoofTypeCode                        *int             `json:"roofTypeCode,omitempty"`
	SewerCode                           *int             `json:"sewerCode,omitempty"`
	StoriesNbrCode                      *int             `json:"storiesNbrCode,omitempty"`
	StyleCode                           *int             `json:"styleCode,omitempty"`
	SumResidentialUnits                 *decimal.Decimal `json:"sumResidentialUnits,omitempty"`
	SumBuildingsNbr                     *decimal.Decimal `json:"sumBuildingsNbr,omitempty"`
	SumCommercialUnits                  *decimal.Decimal `json:"sumCommercialUnits,omitempty"`
	TopographyCode                      *string          `json:"topographyCode,omitempty"`
	WaterCode                           *int             `json:"waterCode,omitempty"`
	ViewDescription                     *string          `json:"viewDescription,omitempty"`
	GuestHouseArea                      *int             `json:"guestHouseArea,omitempty"`
	ShedArea                            *int             `json:"shedArea,omitempty"`
	PoleStructureArea                   *int             `json:"poleStructureArea,omitempty"`
	LotCode                             *string          `json:"lotCode,omitempty"`
	LotNbr                              *string          `json:"lotNbr,omitempty"`
	LegalLotNumber2                     *string          `json:"legalLotNumber2,omitempty"`
	LegalLotNumber3                     *string          `json:"legalLotNumber3,omitempty"`
	LandLot                             *string          `json:"landLot,omitempty"`
	Block                               *string          `json:"block,omitempty"`
	Block2                              *string          `json:"block2,omitempty"`
	Section                             *string          `json:"section,omitempty"`
	District                            *string          `json:"district,omitempty"`
	LegalUnit                           *string          `json:"legalUnit,omitempty"`
	Municipality                        *string          `json:"municipality,omitempty"`
	SubdivisionName                     *string          `json:"subdivisionName,omitempty"`
	SubdivisionPhaseNbr                 *string          `json:"subdivisionPhaseNbr,omitempty"`
	SubdivisionTractNbr                 *string          `json:"subdivisionTractNbr,omitempty"`
	Meridian                            *string          `json:"meridian,omitempty"`
	AssessorsMapRef                     *string          `json:"assessorsMapRef,omitempty"`
	LegalDescription                    *string          `json:"legalDescription,omitempty"`
	CurrentSaleTransactionId            *int64           `json:"-"`
	DeedLastSaleTransactionId           *int64           `json:"deedLastSaleTransactionId,string,omitempty"`
	CurrentSaleDocNbr                   *string          `json:"currentSaleDocNbr,omitempty"`
	CurrentSaleBook                     *string          `json:"currentSaleBook,omitempty"`
	CurrentSalePage                     *string          `json:"currentSalePage,omitempty"`
	CurrentSaleRecordingDate            *time.Time       `json:"currentSaleRecordingDate,omitempty"`
	CurrentSaleContractDate             *time.Time       `json:"currentSaleContractDate,omitempty"`
	LastOwnershipTransferDate           *time.Time       `json:"lastOwnershipTransferDate,omitempty"`
	LastOwnershipTransferDocumentNumber *string          `json:"lastOwnershipTransferDocumentNumber,omitempty"`
	LastOwnershipTransferTransactionId  *int             `json:"lastOwnershipTransferTransactionId,omitempty"`
	CurrentSaleDocumentType             *string          `json:"currentSaleDocumentType,omitempty"`
	CurrentSalesPrice                   *int             `json:"currentSalesPrice,omitempty"`
	AssessorLastSaleAmount              *int             `json:"assessorLastSaleAmount,omitempty"`
	CurrentSalesPriceCode               *int             `json:"currentSalesPriceCode,omitempty"`
	CurrentSaleBuyer1FullName           *string          `json:"currentSaleBuyer1FullName,omitempty"`
	CurrentSaleBuyer2FullName           *string          `json:"currentSaleBuyer2FullName,omitempty"`
	CurrentSaleSeller1FullName          *string          `json:"currentSaleSeller1FullName,omitempty"`
	CurrentSaleSeller2FullName          *string          `json:"currentSaleSeller2FullName,omitempty"`
	PrevSaleTransactionId               *int64           `json:"-"`
	PrevSaleDocNbr                      *string          `json:"prevSaleDocNbr,omitempty"`
	PrevSaleBook                        *string          `json:"prevSaleBook,omitempty"`
	PrevSalePage                        *string          `json:"prevSalePage,omitempty"`
	PrevSaleRecordingDate               *time.Time       `json:"prevSaleRecordingDate,omitempty"`
	PrevSaleContractDate                *time.Time       `json:"prevSaleContractDate,omitempty"`
	PrevSaleDocumentType                *string          `json:"prevSaleDocumentType,omitempty"`
	PrevSalesPrice                      *int             `json:"prevSalesPrice,omitempty"`
	PrevSalesPriceCode                  *int             `json:"prevSalesPriceCode,omitempty"`
	PrevSaleBuyer1FullName              *string          `json:"prevSaleBuyer1FullName,omitempty"`
	PrevSaleBuyer2FullName              *string          `json:"prevSaleBuyer2FullName,omitempty"`
	PrevSaleSeller1FullName             *string          `json:"prevSaleSeller1FullName,omitempty"`
	PrevSaleSeller2FullName             *string          `json:"prevSaleSeller2FullName,omitempty"`
}
