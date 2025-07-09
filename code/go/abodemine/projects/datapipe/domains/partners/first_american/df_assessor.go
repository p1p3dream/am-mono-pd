package first_american

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type Assessor struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any

	FIPS                        string
	PropertyID                  int64
	APN                         string
	APNSeqNbr                   *string
	OldAPN                      *string
	OldApnIndicator             *string
	TaxAccountNumber            *string
	SitusFullStreetAddress      *string
	SitusHouseNbr               *string
	SitusHouseNbrSuffix         *string
	SitusDirectionLeft          *string
	SitusStreet                 *string
	SitusMode                   *string
	SitusDirectionRight         *string
	SitusUnitType               *string
	SitusUnitNbr                *string
	SitusCity                   *string
	SitusState                  *string
	SitusZIP5                   *string
	SitusZIP4                   *string
	SitusCarrierCode            *string
	SitusLatitude               *float64
	SitusLongitude              *float64
	SitusGeoStatusCode          *string
	PropertyClassID             *string
	LandUseCode                 *string
	StateLandUseCode            *string
	CountyLandUseCode           *string
	Zoning                      *string
	SitusCensusTract            *string
	SitusCensusBlock            *int
	MobileHomeInd               *bool
	TimeshareCode               *bool
	SchoolDistrictName          *string
	LotSizeFrontageFeet         *string
	LotSizeDepthFeet            *int
	LotSizeAcres                *string
	LotSizeSqFt                 *int
	Owner1CorpInd               *bool
	Owner1LastName              *string
	Owner1FirstName             *string
	Owner1MiddleName            *string
	Owner1Suffix                *string
	Owner2CorpInd               *bool
	Owner2LastName              *string
	Owner2FirstName             *string
	Owner2MiddleName            *string
	Owner2Suffix                *string
	OwnerNAME1FULL              *string
	OwnerNAME2FULL              *string
	OwnerOccupied               *bool
	Owner1OwnershipRights       *string
	MailingFullStreetAddress    *string
	MailingHouseNbr             *string
	MailingHouseNbrSuffix       *string
	MailingDirectionLeft        *string
	MailingStreet               *string
	MailingMode                 *string
	MailingDirectionRight       *string
	MailingUnitType             *string
	MailingUnitNbr              *string
	MailingCity                 *string
	MailingState                *string
	MailingZIP5                 *string
	MailingZIP4                 *string
	MailingCarrierCode          *string
	MailingOptOut               *bool
	MailingCOName               *string
	MailingForeignAddressInd    *string
	AssdTotalValue              *int
	AssdLandValue               *int
	AssdImprovementValue        *int
	MarketTotalValue            *int
	MarketValueLand             *int
	MarketValueImprovement      *int
	TaxAmt                      *int64
	TaxYear                     *int
	TaxDeliquentYear            *int
	MarketYear                  *int
	AssdYear                    *int
	TaxRateCodeArea             *string
	SchoolTaxDistrict1Code      *string
	SchoolTaxDistrict2Code      *string
	SchoolTaxDistrict3Code      *string
	HomesteadInd                *bool
	VeteranInd                  *bool
	DisabledInd                 *bool
	WidowInd                    *bool
	SeniorInd                   *bool
	SchoolCollegeInd            *bool
	ReligiousInd                *bool
	WelfareInd                  *bool
	PublicUtilityInd            *bool
	CemeteryInd                 *bool
	HospitalInd                 *bool
	LibraryInd                  *bool
	BuildingArea                *decimal.Decimal
	BuildingAreaInd             *string
	SumBuildingSqFt             *decimal.Decimal
	SumLivingAreaSqFt           *int
	SumGroundFloorSqFt          *int
	SumGrossAreaSqFt            *int
	SumAdjAreaSqFt              *int
	AtticSqFt                   *int
	AtticUnfinishedSqFt         *int
	AtticFinishedSqFt           *int
	SumBasementSqFt             *int
	BasementUnfinishedSqFt      *int
	BasementFinishedSqFt        *int
	SumGarageSqFt               *int
	GarageUnFinishedSqFt        *int
	GarageFinishedSqFt          *int
	YearBuilt                   *int
	EffectiveYearBuilt          *int
	Bedrooms                    *decimal.Decimal
	TotalRooms                  *decimal.Decimal
	BathTotalCalc               *decimal.Decimal
	BathFull                    *decimal.Decimal
	BathsPartialNbr             *decimal.Decimal
	BathFixturesNbr             *decimal.Decimal
	Amenities                   *string
	AirConditioningCode         *int
	BasementCode                *int
	BuildingClassCode           *int
	BuildingConditionCode       *int
	ConstructionTypeCode        *int
	DeckInd                     *bool
	ExteriorWallsCode           *int
	InteriorWallsCode           *int
	FireplaceCode               *int
	FloorCoverCode              *string
	Garage                      *int
	HeatCode                    *int
	HeatingFuelTypeCode         *int
	SiteInfluenceCode           *string
	GarageParkingNbr            *int
	DrivewayCode                *string
	OtherRooms                  *string
	PatioCode                   *int
	PoolCode                    *int
	PorchCode                   *int
	BuildingQualityCode         *int
	RoofCoverCode               *int
	RoofTypeCode                *int
	SewerCode                   *int
	StoriesNbrCode              *int
	StyleCode                   *int
	SumResidentialUnits         *decimal.Decimal
	SumBuildingsNbr             *decimal.Decimal
	SumCommercialUnits          *decimal.Decimal
	TopographyCode              *string
	WaterCode                   *int
	LotCode                     *string
	LotNbr                      *string
	LandLot                     *string
	Block                       *string
	Section                     *string
	District                    *string
	LegalUnit                   *string
	Municipality                *string
	SubdivisionName             *string
	SubdivisionPhaseNbr         *string
	SubdivisionTractNbr         *string
	Meridian                    *string
	AssessorsMapRef             *string
	LegalDescription            *string
	CurrentSaleTransactionId    *int64
	CurrentSaleDocNbr           *string
	CurrentSaleBook             *string
	CurrentSalePage             *string
	CurrentSaleRecordingDate    *time.Time
	CurrentSaleContractDate     *time.Time
	CurrentSaleDocumentType     *string
	CurrentSalesPrice           *int
	CurrentSalesPriceCode       *int
	CurrentSaleBuyer1FullName   *string
	CurrentSaleBuyer2FullName   *string
	CurrentSaleSeller1FullName  *string
	CurrentSaleSeller2FullName  *string
	ConcurrentMtg1DocNbr        *string
	ConcurrentMtg1Book          *string
	ConcurrentMtg1Page          *string
	ConcurrentMtg1RecordingDate *time.Time
	ConcurrentMtg1LoanAmt       *int
	ConcurrentMtg1Lender        *string
	ConcurrentMtg1Term          *string
	ConcurrentMtg1InterestRate  *decimal.Decimal
	ConcurrentMtg1LoanDueDate   *time.Time
	ConcurrentMtg1LoanType      *int
	ConcurrentMtg1TypeFinancing *string
	ConcurrentMtg2DocNbr        *string
	ConcurrentMtg2Book          *string
	ConcurrentMtg2Page          *string
	ConcurrentMtg2RecordingDate *time.Time
	ConcurrentMtg2LoanAmt       *int
	ConcurrentMtg2Lender        *string
	ConcurrentMtg2Term          *string
	ConcurrentMtg2InterestRate  *decimal.Decimal
	ConcurrentMtg2LoanDueDate   *time.Time
	ConcurrentMtg2LoanType      *int
	ConcurrentMtg2Typefinancing *string
	PrevSaleTransactionId       *int64
	PrevSaleDocNbr              *string
	PrevSaleBook                *string
	PrevSalePage                *string
	PrevSaleRecordingDate       *time.Time
	PrevSaleContractDate        *time.Time
	PrevSaleDocumentType        *string
	PrevSalesPrice              *int
	PrevSalesPriceCode          *int
	PrevSaleBuyer1FullName      *string
	PrevSaleBuyer2FullName      *string
	PrevSaleSeller1FullName     *string
	PrevSaleSeller2FullName     *string
	PrevMtg1DocNbr              *string
	PrevMtg1Book                *string
	PrevMtg1Page                *string
	PrevMtg1RecordingDate       *time.Time
	PrevMtg1LoanAmt             *int
	PrevMtg1Lender              *string
	PrevMtg1Term                *int
	PrevMtg1InterestRate        *decimal.Decimal
	PrevMtg1LoanDueDate         *time.Time
	PrevMtg1LoanType            *int
	PrevMtg1TypeFinancing       *string
	TotalOpenLienNbr            *int
	TotalOpenLienAmt            *int
	Mtg1TransactionId           *int64
	Mtg1RecordingDate           *time.Time
	Mtg1LoanAmt                 *int
	Mtg1Lender                  *string
	Mtg1PrivateLender           *bool
	Mtg1Term                    *string
	Mtg1LoanDueDate             *time.Time
	Mtg1AdjRider                *bool
	Mtg1LoanType                *int
	Mtg1TypeFinancing           *string
	Mtg1LienPosition            *int
	Mtg2TransactionId           *int64
	Mtg2RecordingDate           *time.Time
	Mtg2LoanAmt                 *int
	Mtg2Lender                  *string
	Mtg2PrivateLender           *bool
	Mtg2Term                    *string
	Mtg2LoanDueDate             *time.Time
	Mtg2AdjRider                *bool
	Mtg2LoanType                *int
	Mtg2TypeFinancing           *string
	Mtg2LienPosition            *int
	Mtg3TransactionId           *int64
	Mtg3RecordingDate           *time.Time
	Mtg3LoanAmt                 *int
	Mtg3Lender                  *string
	Mtg3PrivateLender           *bool
	Mtg3Term                    *int
	Mtg3LoanDueDate             *time.Time
	Mtg3AdjRider                *bool
	Mtg3LoanType                *int
	Mtg3TypeFinancing           *string
	Mtg3LienPosition            *int
	Mtg4TransactionId           *int64
	Mtg4RecordingDate           *time.Time
	Mtg4LoanAmt                 *int
	Mtg4Lender                  *string
	Mtg4PrivateLender           *bool
	Mtg4Term                    *int
	Mtg4LoanDueDate             *time.Time
	Mtg4AdjRider                *bool
	Mtg4LoanType                *int
	Mtg4TypeFinancing           *string
	Mtg4LienPosition            *int
	FATimeStamp                 *time.Time
	FARecordType                *string
}

func (dr *Assessor) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(Assessor)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "FIPS":
			if field == "" {
				return nil, &errors.Object{
					Id:     "63bb09f1-ef2c-43ce-b1ab-728839876862",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "FIPS is required.",
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.FIPS = field
		case "PropertyID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "f9549361-7a1e-4078-9577-8716757ae23b",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.PropertyID = v
		case "APN":
			if field == "" {
				return nil, &errors.Object{
					Id:     "9e5af2a6-7704-479f-86ee-771a0f72ea00",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "APN is required.",
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.APN = field
		case "APNSeqNbr":
			record.APNSeqNbr = val.StringPtrIfNonZero(field)
		case "OldAPN":
			record.OldAPN = val.StringPtrIfNonZero(field)
		case "OldApnIndicator":
			record.OldApnIndicator = val.StringPtrIfNonZero(field)
		case "TaxAccountNumber":
			record.TaxAccountNumber = val.StringPtrIfNonZero(field)
		case "SitusFullStreetAddress":
			record.SitusFullStreetAddress = val.StringPtrIfNonZero(field)
		case "SitusHouseNbr":
			record.SitusHouseNbr = val.StringPtrIfNonZero(field)
		case "SitusHouseNbrSuffix":
			record.SitusHouseNbrSuffix = val.StringPtrIfNonZero(field)
		case "SitusDirectionLeft":
			record.SitusDirectionLeft = val.StringPtrIfNonZero(field)
		case "SitusStreet":
			record.SitusStreet = val.StringPtrIfNonZero(field)
		case "SitusMode":
			record.SitusMode = val.StringPtrIfNonZero(field)
		case "SitusDirectionRight":
			record.SitusDirectionRight = val.StringPtrIfNonZero(field)
		case "SitusUnitType":
			record.SitusUnitType = val.StringPtrIfNonZero(field)
		case "SitusUnitNbr":
			record.SitusUnitNbr = val.StringPtrIfNonZero(field)
		case "SitusCity":
			record.SitusCity = val.StringPtrIfNonZero(field)
		case "SitusState":
			record.SitusState = val.StringPtrIfNonZero(field)
		case "SitusZIP5":
			record.SitusZIP5 = val.StringPtrIfNonZero(field)
		case "SitusZIP4":
			record.SitusZIP4 = val.StringPtrIfNonZero(field)
		case "SitusCarrierCode":
			record.SitusCarrierCode = val.StringPtrIfNonZero(field)
		case "SitusLatitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b93bdf53-380b-4b4d-8e10-69eed8711985")
			}
			record.SitusLatitude = v
		case "SitusLongitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "31688cac-efe1-49df-8d96-53b36e6b9f48")
			}
			record.SitusLongitude = v
		case "SitusGeoStatusCode":
			record.SitusGeoStatusCode = val.StringPtrIfNonZero(field)
		case "PropertyClassID":
			record.PropertyClassID = val.StringPtrIfNonZero(field)
		case "LandUseCode":
			record.LandUseCode = val.StringPtrIfNonZero(field)
		case "StateLandUseCode":
			record.StateLandUseCode = val.StringPtrIfNonZero(field)
		case "CountyLandUseCode":
			record.CountyLandUseCode = val.StringPtrIfNonZero(field)
		case "Zoning":
			record.Zoning = val.StringPtrIfNonZero(field)
		case "SitusCensusTract":
			record.SitusCensusTract = val.StringPtrIfNonZero(field)
		case "SitusCensusBlock":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d89a076d-4d6f-4409-bd25-54dd7bb7f2c9")
			}
			record.SitusCensusBlock = v
		case "MobileHomeInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "307314dc-514f-464f-94d1-0f63f2a52c57")
			}
			record.MobileHomeInd = v
		case "TimeshareCode":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3d0cb91e-8ac0-4ef3-8898-7ecd4297179a")
			}
			record.TimeshareCode = v
		case "SchoolDistrictName":
			record.SchoolDistrictName = val.StringPtrIfNonZero(field)
		case "LotSizeFrontageFeet":
			record.LotSizeFrontageFeet = val.StringPtrIfNonZero(field)
		case "LotSizeDepthFeet":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "cbbfb6ca-69ac-49b6-acb2-c0f4e7211edd")
			}
			record.LotSizeDepthFeet = v
		case "LotSizeAcres":
			record.LotSizeAcres = val.StringPtrIfNonZero(field)
		case "LotSizeSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5649d7b7-2570-466e-85bf-88c6c7f5c6a1")
			}
			record.LotSizeSqFt = v
		case "Owner1CorpInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1b36657f-fb32-483d-9ff7-510192f36414")
			}
			record.Owner1CorpInd = v
		case "Owner1LastName":
			record.Owner1LastName = val.StringPtrIfNonZero(field)
		case "Owner1FirstName":
			record.Owner1FirstName = val.StringPtrIfNonZero(field)
		case "Owner1MiddleName":
			record.Owner1MiddleName = val.StringPtrIfNonZero(field)
		case "Owner1Suffix":
			record.Owner1Suffix = val.StringPtrIfNonZero(field)
		case "Owner2CorpInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "340eb422-5bb5-4174-8261-110f33fb472e")
			}
			record.Owner2CorpInd = v
		case "Owner2LastName":
			record.Owner2LastName = val.StringPtrIfNonZero(field)
		case "Owner2FirstName":
			record.Owner2FirstName = val.StringPtrIfNonZero(field)
		case "Owner2MiddleName":
			record.Owner2MiddleName = val.StringPtrIfNonZero(field)
		case "Owner2Suffix":
			record.Owner2Suffix = val.StringPtrIfNonZero(field)
		case "OwnerNAME1FULL":
			record.OwnerNAME1FULL = val.StringPtrIfNonZero(field)
		case "OwnerNAME2FULL":
			record.OwnerNAME2FULL = val.StringPtrIfNonZero(field)
		case "OwnerOccupied":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2f20eee5-b545-4a0a-b3e2-6f0c5dd04a69")
			}
			record.OwnerOccupied = v
		case "Owner1OwnershipRights":
			record.Owner1OwnershipRights = val.StringPtrIfNonZero(field)
		case "MailingFullStreetAddress":
			record.MailingFullStreetAddress = val.StringPtrIfNonZero(field)
		case "MailingHouseNbr":
			record.MailingHouseNbr = val.StringPtrIfNonZero(field)
		case "MailingHouseNbrSuffix":
			record.MailingHouseNbrSuffix = val.StringPtrIfNonZero(field)
		case "MailingDirectionLeft":
			record.MailingDirectionLeft = val.StringPtrIfNonZero(field)
		case "MailingStreet":
			record.MailingStreet = val.StringPtrIfNonZero(field)
		case "MailingMode":
			record.MailingMode = val.StringPtrIfNonZero(field)
		case "MailingDirectionRight":
			record.MailingDirectionRight = val.StringPtrIfNonZero(field)
		case "MailingUnitType":
			record.MailingUnitType = val.StringPtrIfNonZero(field)
		case "MailingUnitNbr":
			record.MailingUnitNbr = val.StringPtrIfNonZero(field)
		case "MailingCity":
			record.MailingCity = val.StringPtrIfNonZero(field)
		case "MailingState":
			record.MailingState = val.StringPtrIfNonZero(field)
		case "MailingZIP5":
			record.MailingZIP5 = val.StringPtrIfNonZero(field)
		case "MailingZIP4":
			record.MailingZIP4 = val.StringPtrIfNonZero(field)
		case "MailingCarrierCode":
			record.MailingCarrierCode = val.StringPtrIfNonZero(field)
		case "MailingOptOut":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2cf7057f-2d4d-4717-b62e-bca235f1b808")
			}
			record.MailingOptOut = v
		case "MailingCOName":
			record.MailingCOName = val.StringPtrIfNonZero(field)
		case "MailingForeignAddressInd":
			record.MailingForeignAddressInd = val.StringPtrIfNonZero(field)
		case "AssdTotalValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3c76d89a-808f-413c-9322-b9a87c9b9833")
			}
			record.AssdTotalValue = v
		case "AssdLandValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9e770927-244c-455b-99c9-14c97ac16aac")
			}
			record.AssdLandValue = v
		case "AssdImprovementValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8d4bbbee-635c-48a9-b74a-5db10f933521")
			}
			record.AssdImprovementValue = v
		case "MarketTotalValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5432644d-729c-4848-8026-2326d54ceaa0")
			}
			record.MarketTotalValue = v
		case "MarketValueLand":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "db37fcb9-b531-4acb-8496-d97d5d17fc0d")
			}
			record.MarketValueLand = v
		case "MarketValueImprovement":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "29e17fe0-7d3a-43df-b1b3-441da265e2ba")
			}
			record.MarketValueImprovement = v
		case "TaxAmt":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0edbf1bb-8b3f-4bca-889f-2a2f5ad82d2e")
			}
			record.TaxAmt = v
		case "TaxYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "cfddac2c-f21c-4f11-82f5-99fc73b65f49")
			}
			record.TaxYear = v
		case "TaxDeliquentYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "75ae9976-b5b0-46da-9ffb-1a29b906e2e2")
			}
			record.TaxDeliquentYear = v
		case "MarketYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3e7549fa-cef0-4e2f-80d6-38487c3dec16")
			}
			record.MarketYear = v
		case "AssdYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "95e4b566-7dab-4eb6-be07-0c22e29787e4")
			}
			record.AssdYear = v
		case "TaxRateCodeArea":
			record.TaxRateCodeArea = val.StringPtrIfNonZero(field)
		case "SchoolTaxDistrict1Code":
			record.SchoolTaxDistrict1Code = val.StringPtrIfNonZero(field)
		case "SchoolTaxDistrict2Code":
			record.SchoolTaxDistrict2Code = val.StringPtrIfNonZero(field)
		case "SchoolTaxDistrict3Code":
			record.SchoolTaxDistrict3Code = val.StringPtrIfNonZero(field)
		case "HomesteadInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "38ee1595-b385-41f2-ad17-ebacdb91593f")
			}
			record.HomesteadInd = v
		case "VeteranInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c17a4f5c-f68b-43e7-a884-13e0d9bf22cf")
			}
			record.VeteranInd = v
		case "DisabledInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "35b2e803-a124-452c-b0cf-30e2cd441b24")
			}
			record.DisabledInd = v
		case "WidowInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bdec0fa5-2da7-493c-b02a-dc7e2f7ec9ec")
			}
			record.WidowInd = v
		case "SeniorInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "295bfa32-a3c3-4efe-9037-b152068c0727")
			}
			record.SeniorInd = v
		case "SchoolCollegeInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "096744b5-d94b-408a-9df3-f77394931121")
			}
			record.SchoolCollegeInd = v
		case "ReligiousInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "72135ab3-1895-42f9-966f-e9d908069df9")
			}
			record.ReligiousInd = v
		case "WelfareInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "59b54c0b-c284-4e21-b6c9-c722edc43dca")
			}
			record.WelfareInd = v
		case "PublicUtilityInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "11ad7e9c-5f15-42b2-a861-d2e0869b95c4")
			}
			record.PublicUtilityInd = v
		case "CemeteryInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "14abf54c-07f1-4f3b-ba52-7f4c7e411f5b")
			}
			record.CemeteryInd = v
		case "HospitalInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "31d6e04f-da3b-439b-b7eb-73e78292685d")
			}
			record.HospitalInd = v
		case "LibraryInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a1a52289-f93d-446d-baef-c48be1016231")
			}
			record.LibraryInd = v
		case "BuildingArea":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "eccf6f75-c0c6-4b57-a039-85aba6dbb1ba")
			}
			record.BuildingArea = v
		case "BuildingAreaInd":
			record.BuildingAreaInd = val.StringPtrIfNonZero(field)
		case "SumBuildingSqFt":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "27679d0f-b741-4d04-b523-ce70f58935ef")
			}
			record.SumBuildingSqFt = v
		case "SumLivingAreaSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "094f836b-0511-465c-bc04-e73d9ead5a20")
			}
			record.SumLivingAreaSqFt = v
		case "SumGroundFloorSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "dac59a7b-3d06-4c7d-be0a-90bca7ae1ad7")
			}
			record.SumGroundFloorSqFt = v
		case "SumGrossAreaSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6bb4b218-a297-4cac-bf3e-e7ef2d45c8d7")
			}
			record.SumGrossAreaSqFt = v
		case "SumAdjAreaSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "04a14c0c-35e8-4642-96fc-34c704cd3946")
			}
			record.SumAdjAreaSqFt = v
		case "AtticSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6347238c-f9ae-4d66-beb5-6a52d569211c")
			}
			record.AtticSqFt = v
		case "AtticUnfinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "06b6eb80-95dd-4e21-ab43-1bf1db4ae149")
			}
			record.AtticUnfinishedSqFt = v
		case "AtticFinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1c196d2f-a5e8-40b9-b49a-877f59514637")
			}
			record.AtticFinishedSqFt = v
		case "SumBasementSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "860e7524-85b2-4457-92a5-0d5b25f80cf2")
			}
			record.SumBasementSqFt = v
		case "BasementUnfinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5c30a5ef-d125-4245-9da1-84104eef6112")
			}
			record.BasementUnfinishedSqFt = v
		case "BasementFinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "96779bca-e394-437d-adb2-8f7afac75cb0")
			}
			record.BasementFinishedSqFt = v
		case "SumGarageSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "548b7309-25a1-459b-894e-a3e3e4f6384a")
			}
			record.SumGarageSqFt = v
		case "GarageUnFinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ef34cfa9-977b-4648-9380-d6f5662cbadf")
			}
			record.GarageUnFinishedSqFt = v
		case "GarageFinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9f9bebf9-078d-446d-a24e-3e4f14a0f8cc")
			}
			record.GarageFinishedSqFt = v
		case "YearBuilt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3cb86c35-3362-4ed0-95ce-b4f3cbb94851")
			}
			record.YearBuilt = v
		case "EffectiveYearBuilt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c42b6e5c-23b4-4ef5-a537-7e852f058d9b")
			}
			record.EffectiveYearBuilt = v
		case "Bedrooms":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "096f5945-a9bb-4706-8e68-2fb1091670d6")
			}
			record.Bedrooms = v
		case "TotalRooms":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "dd75cedd-a762-4b9c-a9b3-3999c076a1b4")
			}
			record.TotalRooms = v
		case "BathTotalCalc":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "075cf1c9-99d5-4d35-a762-947658fe9040")
			}
			record.BathTotalCalc = v
		case "BathFull":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f7fb7761-6278-4019-8aaa-38bffbcdbca7")
			}
			record.BathFull = v
		case "BathsPartialNbr":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "669cfc00-0e00-4357-adcc-b6cdaa195158")
			}
			record.BathsPartialNbr = v
		case "BathFixturesNbr":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "63bafdfb-db56-4e57-8118-b43d63efd9e7")
			}
			record.BathFixturesNbr = v
		case "Amenities":
			record.Amenities = val.StringPtrIfNonZero(field)
		case "AirConditioningCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6e916baa-f4fa-4ea9-8dee-1b53491ccb6e")
			}
			record.AirConditioningCode = v
		case "BasementCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "139c3508-ec36-466a-894d-8c3e3be739f2")
			}
			record.BasementCode = v
		case "BuildingClassCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9074b8aa-a708-4549-a2d9-5d49457f857c")
			}
			record.BuildingClassCode = v
		case "BuildingConditionCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8a89fd8a-a175-42aa-a7f0-86ee3de7227b")
			}
			record.BuildingConditionCode = v
		case "ConstructionTypeCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9f907457-4a90-4eb1-ae1a-9293a606661b")
			}
			record.ConstructionTypeCode = v
		case "DeckInd":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3b799934-7efe-4943-8eb4-d26c20d88994")
			}
			record.DeckInd = v
		case "ExteriorWallsCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f43c3556-3cc0-4f59-9381-7ee67631733b")
			}
			record.ExteriorWallsCode = v
		case "InteriorWallsCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fa3d3758-f656-49ad-93d2-a2fad40e4904")
			}
			record.InteriorWallsCode = v
		case "FireplaceCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3c1a3efb-b667-42e1-829c-469c29f63bca")
			}
			record.FireplaceCode = v
		case "FloorCoverCode":
			record.FloorCoverCode = val.StringPtrIfNonZero(field)
		case "Garage":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f30e5a48-776b-490e-afb0-a957c6dcfaa9")
			}
			record.Garage = v
		case "HeatCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "775bcfdc-542b-4f63-8b70-3d16b42fd54c")
			}
			record.HeatCode = v
		case "HeatingFuelTypeCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "70c2acf8-1980-4fc8-9185-882b2869fa76")
			}
			record.HeatingFuelTypeCode = v
		case "SiteInfluenceCode":
			record.SiteInfluenceCode = val.StringPtrIfNonZero(field)
		case "GarageParkingNbr":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b1cc1885-9f0f-403f-9144-65ad49d7b9da")
			}
			record.GarageParkingNbr = v
		case "DrivewayCode":
			record.DrivewayCode = val.StringPtrIfNonZero(field)
		case "OtherRooms":
			record.OtherRooms = val.StringPtrIfNonZero(field)
		case "PatioCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b979fa32-24c0-46d9-ba03-b8a1618a05a1")
			}
			record.PatioCode = v
		case "PoolCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a4f8dc49-82cb-435d-ae37-18167d3e0186")
			}
			record.PoolCode = v
		case "PorchCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1f62392e-315d-4b71-b519-6dd2f946f209")
			}
			record.PorchCode = v
		case "BuildingQualityCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "59f01c94-4071-4c03-b23e-8969e30a0587")
			}
			record.BuildingQualityCode = v
		case "RoofCoverCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "861e85f9-1cce-4d66-9ec3-9d134cdc46eb")
			}
			record.RoofCoverCode = v
		case "RoofTypeCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "75feeba3-8010-47cf-84a1-c43161a1d97c")
			}
			record.RoofTypeCode = v
		case "SewerCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3f95b910-1ed0-4036-9ca2-a9e944977a00")
			}
			record.SewerCode = v
		case "StoriesNbrCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a744b9b7-f2d6-4770-931d-290fead660ae")
			}
			record.StoriesNbrCode = v
		case "StyleCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4ba4afec-1dac-4539-95a3-07b48cb8273b")
			}
			record.StyleCode = v
		case "SumResidentialUnits":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "81adb67c-5412-4028-8c80-543d63d75e25")
			}
			record.SumResidentialUnits = v
		case "SumBuildingsNbr":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "81f79df2-6d1e-4767-a070-6a4df0011cf1")
			}
			record.SumBuildingsNbr = v
		case "SumCommercialUnits":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "133566a7-aa1d-40e4-8c6e-5b6acd7593ba")
			}
			record.SumCommercialUnits = v
		case "TopographyCode":
			record.TopographyCode = val.StringPtrIfNonZero(field)
		case "WaterCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "99820da7-a4f2-4036-a7e6-37950a2a9b40")
			}
			record.WaterCode = v
		case "LotCode":
			record.LotCode = val.StringPtrIfNonZero(field)
		case "LotNbr":
			record.LotNbr = val.StringPtrIfNonZero(field)
		case "LandLot":
			record.LandLot = val.StringPtrIfNonZero(field)
		case "Block":
			record.Block = val.StringPtrIfNonZero(field)
		case "Section":
			record.Section = val.StringPtrIfNonZero(field)
		case "District":
			record.District = val.StringPtrIfNonZero(field)
		case "LegalUnit":
			record.LegalUnit = val.StringPtrIfNonZero(field)
		case "Municipality":
			record.Municipality = val.StringPtrIfNonZero(field)
		case "SubdivisionName":
			record.SubdivisionName = val.StringPtrIfNonZero(field)
		case "SubdivisionPhaseNbr":
			record.SubdivisionPhaseNbr = val.StringPtrIfNonZero(field)
		case "SubdivisionTractNbr":
			record.SubdivisionTractNbr = val.StringPtrIfNonZero(field)
		case "Meridian":
			record.Meridian = val.StringPtrIfNonZero(field)
		case "AssessorsMapRef":
			record.AssessorsMapRef = val.StringPtrIfNonZero(field)
		case "LegalDescription":
			record.LegalDescription = val.StringPtrIfNonZero(field)
		case "CurrentSaleTransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bc87ae7b-de4a-41a8-a5b4-9f029f7bda43")
			}
			record.CurrentSaleTransactionId = v
		case "CurrentSaleDocNbr":
			record.CurrentSaleDocNbr = val.StringPtrIfNonZero(field)
		case "CurrentSaleBook":
			record.CurrentSaleBook = val.StringPtrIfNonZero(field)
		case "CurrentSalePage":
			record.CurrentSalePage = val.StringPtrIfNonZero(field)
		case "CurrentSaleRecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "6604caea-a0f3-4855-b085-f54c8d612fc0")
			}
			record.CurrentSaleRecordingDate = v
		case "CurrentSaleContractDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "1af54d7a-ccc4-4c3e-bf60-b6cb8111237b")
			}
			record.CurrentSaleContractDate = v
		case "CurrentSaleDocumentType":
			record.CurrentSaleDocumentType = val.StringPtrIfNonZero(field)
		case "CurrentSalesPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "472ef2e4-4cd5-48be-9f1e-c5936e131836")
			}
			record.CurrentSalesPrice = v
		case "CurrentSalesPriceCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "66196705-0ab0-4b04-8ecd-c7fd6c822712")
			}
			record.CurrentSalesPriceCode = v
		case "CurrentSaleBuyer1FullName":
			record.CurrentSaleBuyer1FullName = val.StringPtrIfNonZero(field)
		case "CurrentSaleBuyer2FullName":
			record.CurrentSaleBuyer2FullName = val.StringPtrIfNonZero(field)
		case "CurrentSaleSeller1FullName":
			record.CurrentSaleSeller1FullName = val.StringPtrIfNonZero(field)
		case "CurrentSaleSeller2FullName":
			record.CurrentSaleSeller2FullName = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1DocNbr":
			record.ConcurrentMtg1DocNbr = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1Book":
			record.ConcurrentMtg1Book = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1Page":
			record.ConcurrentMtg1Page = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "b320ebf8-3c28-4d59-810e-f2e91d2839fd")
			}
			record.ConcurrentMtg1RecordingDate = v
		case "ConcurrentMtg1LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "489b45ec-a25e-447b-bca6-ef9dba62a6cf")
			}
			record.ConcurrentMtg1LoanAmt = v
		case "ConcurrentMtg1Lender":
			record.ConcurrentMtg1Lender = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1Term":
			record.ConcurrentMtg1Term = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg1InterestRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d5d21609-c2fa-4906-8d96-699b736557eb")
			}
			record.ConcurrentMtg1InterestRate = v
		case "ConcurrentMtg1LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "97b2fc43-70fa-4e57-ba23-77262eb3aaca")
			}
			record.ConcurrentMtg1LoanDueDate = v
		case "ConcurrentMtg1LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c0f82f17-7ea3-4f7a-9ec2-5f175f0368fe")
			}
			record.ConcurrentMtg1LoanType = v
		case "ConcurrentMtg1TypeFinancing":
			record.ConcurrentMtg1TypeFinancing = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2DocNbr":
			record.ConcurrentMtg2DocNbr = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2Book":
			record.ConcurrentMtg2Book = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2Page":
			record.ConcurrentMtg2Page = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "2d1cde80-27b7-450f-b7c3-d1e140a11ea9")
			}
			record.ConcurrentMtg2RecordingDate = v
		case "ConcurrentMtg2LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c15d0f8b-e48e-457b-b727-f448f7869c88")
			}
			record.ConcurrentMtg2LoanAmt = v
		case "ConcurrentMtg2Lender":
			record.ConcurrentMtg2Lender = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2Term":
			record.ConcurrentMtg2Term = val.StringPtrIfNonZero(field)
		case "ConcurrentMtg2InterestRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "99f9f339-b259-4a41-946c-6458d8d15e0b")
			}
			record.ConcurrentMtg2InterestRate = v
		case "ConcurrentMtg2LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "13551fce-3802-4213-b22a-c1fef416ef34")
			}
			record.ConcurrentMtg2LoanDueDate = v
		case "ConcurrentMtg2LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0b0c32df-2229-4e5b-a9f3-82a87fe240d5")
			}
			record.ConcurrentMtg2LoanType = v
		case "ConcurrentMtg2Typefinancing":
			record.ConcurrentMtg2Typefinancing = val.StringPtrIfNonZero(field)
		case "PrevSaleTransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "df81b042-0993-493e-a223-4c6932fef8ad")
			}
			record.PrevSaleTransactionId = v
		case "PrevSaleDocNbr":
			record.PrevSaleDocNbr = val.StringPtrIfNonZero(field)
		case "PrevSaleBook":
			record.PrevSaleBook = val.StringPtrIfNonZero(field)
		case "PrevSalePage":
			record.PrevSalePage = val.StringPtrIfNonZero(field)
		case "PrevSaleRecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "94657d53-107f-4945-b562-e960047e1254")
			}
			record.PrevSaleRecordingDate = v
		case "PrevSaleContractDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "d437b397-817b-4f92-824a-e1373a65781d")
			}
			record.PrevSaleContractDate = v
		case "PrevSaleDocumentType":
			record.PrevSaleDocumentType = val.StringPtrIfNonZero(field)
		case "PrevSalesPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e6606c8d-691c-404c-a48f-3aafb1e55d5f")
			}
			record.PrevSalesPrice = v
		case "PrevSalesPriceCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6982ca62-3d0c-4962-9364-d94d515cb7ec")
			}
			record.PrevSalesPriceCode = v
		case "PrevSaleBuyer1FullName":
			record.PrevSaleBuyer1FullName = val.StringPtrIfNonZero(field)
		case "PrevSaleBuyer2FullName":
			record.PrevSaleBuyer2FullName = val.StringPtrIfNonZero(field)
		case "PrevSaleSeller1FullName":
			record.PrevSaleSeller1FullName = val.StringPtrIfNonZero(field)
		case "PrevSaleSeller2FullName":
			record.PrevSaleSeller2FullName = val.StringPtrIfNonZero(field)
		case "PrevMtg1DocNbr":
			record.PrevMtg1DocNbr = val.StringPtrIfNonZero(field)
		case "PrevMtg1Book":
			record.PrevMtg1Book = val.StringPtrIfNonZero(field)
		case "PrevMtg1Page":
			record.PrevMtg1Page = val.StringPtrIfNonZero(field)
		case "PrevMtg1RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "4d7e60b7-3586-4ca0-8f41-afb41b772223")
			}
			record.PrevMtg1RecordingDate = v
		case "PrevMtg1LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9b92f8d6-b988-4322-84ba-4b45091042f2")
			}
			record.PrevMtg1LoanAmt = v
		case "PrevMtg1Lender":
			record.PrevMtg1Lender = val.StringPtrIfNonZero(field)
		case "PrevMtg1Term":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e807d25b-44b7-45ca-ab27-ad2fb728e12c")
			}
			record.PrevMtg1Term = v
		case "PrevMtg1InterestRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f6fce7c3-b3b6-46c1-a52a-5d50fb3b88f1")
			}
			record.PrevMtg1InterestRate = v
		case "PrevMtg1LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "30871afd-dcb2-46da-95f4-9538395ff538")
			}
			record.PrevMtg1LoanDueDate = v
		case "PrevMtg1LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "eda84bd2-d193-4029-9030-db73fa3f8580")
			}
			record.PrevMtg1LoanType = v
		case "PrevMtg1TypeFinancing":
			record.PrevMtg1TypeFinancing = val.StringPtrIfNonZero(field)
		case "TotalOpenLienNbr":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d09ac1c5-647b-40b2-9913-32769b5cc949")
			}
			record.TotalOpenLienNbr = v
		case "TotalOpenLienAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c125f6d5-6c3e-488c-bcee-0176102d488c")
			}
			record.TotalOpenLienAmt = v
		case "Mtg1TransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c99a27cf-3530-43b1-b75d-0d90705533f6")
			}
			record.Mtg1TransactionId = v
		case "Mtg1RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "bae5ac18-0cca-4032-a3e7-e49363370cb7")
			}
			record.Mtg1RecordingDate = v
		case "Mtg1LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4099dec0-fba2-498d-ba23-6fa4739d775f")
			}
			record.Mtg1LoanAmt = v
		case "Mtg1Lender":
			record.Mtg1Lender = val.StringPtrIfNonZero(field)
		case "Mtg1PrivateLender":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d3f21c50-a69b-4f4d-b123-fd6dd71ee88d")
			}
			record.Mtg1PrivateLender = v
		case "Mtg1Term":
			record.Mtg1Term = val.StringPtrIfNonZero(field)
		case "Mtg1LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "5b55c6a8-5f2e-4d04-9660-df188aff4e19")
			}
			record.Mtg1LoanDueDate = v
		case "Mtg1AdjRider":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fa1c600a-38f4-4d96-948f-9c849f367f80")
			}
			record.Mtg1AdjRider = v
		case "Mtg1LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c5c25303-555b-4d96-a46f-3fc9db6cb2c3")
			}
			record.Mtg1LoanType = v
		case "Mtg1TypeFinancing":
			record.Mtg1TypeFinancing = val.StringPtrIfNonZero(field)
		case "Mtg1LienPosition":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0956e4a4-a99d-4777-93c7-890546f6ec54")
			}
			record.Mtg1LienPosition = v
		case "Mtg2TransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9c53c5c0-91ce-4f41-b19c-cec956148c36")
			}
			record.Mtg2TransactionId = v
		case "Mtg2RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "12057f5a-e325-4de9-8c55-be677082d869")
			}
			record.Mtg2RecordingDate = v
		case "Mtg2LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ee2a666d-7293-46c8-8991-fdfe3c2c78bb")
			}
			record.Mtg2LoanAmt = v
		case "Mtg2Lender":
			record.Mtg2Lender = val.StringPtrIfNonZero(field)
		case "Mtg2PrivateLender":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e3b9ebb3-0a91-461e-83e6-072d0ade9d2a")
			}
			record.Mtg2PrivateLender = v
		case "Mtg2Term":
			record.Mtg2Term = val.StringPtrIfNonZero(field)
		case "Mtg2LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "b460e7cb-2c77-4b53-8b92-0fa9134aa514")
			}
			record.Mtg2LoanDueDate = v
		case "Mtg2AdjRider":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9d3730c6-465c-46b0-af7e-7b8968a59563")
			}
			record.Mtg2AdjRider = v
		case "Mtg2LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "52a14d14-ff02-42eb-adeb-f9bf320e7879")
			}
			record.Mtg2LoanType = v
		case "Mtg2TypeFinancing":
			record.Mtg2TypeFinancing = val.StringPtrIfNonZero(field)
		case "Mtg2LienPosition":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "70175a83-705a-4540-8456-b0d029e3b6ec")
			}
			record.Mtg2LienPosition = v
		case "Mtg3TransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a3b2f500-6f67-4499-8740-acab7d609ece")
			}
			record.Mtg3TransactionId = v
		case "Mtg3RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "0207f571-9be3-4925-a7a9-1a793792240c")
			}
			record.Mtg3RecordingDate = v
		case "Mtg3LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f6ff5e70-5a17-4280-a5a8-8dcc95661553")
			}
			record.Mtg3LoanAmt = v
		case "Mtg3Lender":
			record.Mtg3Lender = val.StringPtrIfNonZero(field)
		case "Mtg3PrivateLender":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2b051f50-272c-4a3b-97cd-bdd577e82eac")
			}
			record.Mtg3PrivateLender = v
		case "Mtg3Term":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6935bc39-79a8-4e20-b3eb-93d88ed4a011")
			}
			record.Mtg3Term = v
		case "Mtg3LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "736e7789-c84d-4602-b1ad-ac0708380268")
			}
			record.Mtg3LoanDueDate = v
		case "Mtg3AdjRider":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b8c1aedf-4a17-4037-91d9-ce4b6dfa8f9f")
			}
			record.Mtg3AdjRider = v
		case "Mtg3LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "014a4a4d-67f2-4bbb-b0e3-b9afa6b83d7e")
			}
			record.Mtg3LoanType = v
		case "Mtg3TypeFinancing":
			record.Mtg3TypeFinancing = val.StringPtrIfNonZero(field)
		case "Mtg3LienPosition":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "34ea0383-281b-40d7-a11e-c09c2cd40db3")
			}
			record.Mtg3LienPosition = v
		case "Mtg4TransactionId":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fb9db43b-7532-4e24-868e-05539be121f3")
			}
			record.Mtg4TransactionId = v
		case "Mtg4RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "87d69065-91d2-477c-912e-34ae3e0e14fb")
			}
			record.Mtg4RecordingDate = v
		case "Mtg4LoanAmt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b3132948-d7f6-4a85-b8dc-a12f3faa5093")
			}
			record.Mtg4LoanAmt = v
		case "Mtg4Lender":
			record.Mtg4Lender = val.StringPtrIfNonZero(field)
		case "Mtg4PrivateLender":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b282e2a6-b7c9-4480-ac0f-849ffaf9932a")
			}
			record.Mtg4PrivateLender = v
		case "Mtg4Term":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "02a0c224-cdf1-4476-b3df-e83f6a5471f5")
			}
			record.Mtg4Term = v
		case "Mtg4LoanDueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "4f694d0e-d7a3-4867-ab82-4f7af7481019")
			}
			record.Mtg4LoanDueDate = v
		case "Mtg4AdjRider":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "68d68fe5-6c42-4729-a41d-5145c695b674")
			}
			record.Mtg4AdjRider = v
		case "Mtg4LoanType":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "35451dba-f052-4027-a9f6-a30471d50d9e")
			}
			record.Mtg4LoanType = v
		case "Mtg4TypeFinancing":
			record.Mtg4TypeFinancing = val.StringPtrIfNonZero(field)
		case "Mtg4LienPosition":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "72e1e093-744d-4381-bd2e-64f779178d48")
			}
			record.Mtg4LienPosition = v
		case "FATimeStamp":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "0778258e-ae3a-4e78-bcfb-ee22ee43c4d8")
			}
			record.FATimeStamp = v
		case "FARecordType":
			record.FARecordType = val.StringPtrIfNonZero(field)
		default:
			return nil, &errors.Object{
				Id:     "41aa7625-d7e4-4136-95ce-02c40afaf5a3",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Unknown header.",
				Meta: map[string]any{
					"field_index": k,
					"field_value": field,
					"header":      header,
				},
			}
		}
	}

	return record, nil
}

func (dr *Assessor) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"fips",
		"property_id",
		"apn",
		"apn_seq_nbr",
		"old_apn",
		"old_apn_indicator",
		"tax_account_number",
		"situs_full_street_address",
		"situs_house_nbr",
		"situs_house_nbr_suffix",
		"situs_direction_left",
		"situs_street",
		"situs_mode",
		"situs_direction_right",
		"situs_unit_type",
		"situs_unit_nbr",
		"situs_city",
		"situs_state",
		"situs_zip5",
		"situs_zip4",
		"situs_carrier_code",
		"situs_latitude",
		"situs_longitude",
		"situs_geo_status_code",
		"property_class_id",
		"land_use_code",
		"state_land_use_code",
		"county_land_use_code",
		"zoning",
		"situs_census_tract",
		"situs_census_block",
		"mobile_home_ind",
		"timeshare_code",
		"school_district_name",
		"lot_size_frontage_feet",
		"lot_size_depth_feet",
		"lot_size_acres",
		"lot_size_sq_ft",
		"owner1corp_ind",
		"owner1last_name",
		"owner1first_name",
		"owner1middle_name",
		"owner1suffix",
		"owner2corp_ind",
		"owner2last_name",
		"owner2first_name",
		"owner2middle_name",
		"owner2suffix",
		"owner_name1full",
		"owner_name2full",
		"owner_occupied",
		"owner1ownership_rights",
		"mailing_full_street_address",
		"mailing_house_nbr",
		"mailing_house_nbr_suffix",
		"mailing_direction_left",
		"mailing_street",
		"mailing_mode",
		"mailing_direction_right",
		"mailing_unit_type",
		"mailing_unit_nbr",
		"mailing_city",
		"mailing_state",
		"mailing_zip5",
		"mailing_zip4",
		"mailing_carrier_code",
		"mailing_opt_out",
		"mailing_co_name",
		"mailing_foreign_address_ind",
		"assd_total_value",
		"assd_land_value",
		"assd_improvement_value",
		"market_total_value",
		"market_value_land",
		"market_value_improvement",
		"tax_amt",
		"tax_year",
		"tax_deliquent_year",
		"market_year",
		"assd_year",
		"tax_rate_code_area",
		"school_tax_district1code",
		"school_tax_district2code",
		"school_tax_district3code",
		"homestead_ind",
		"veteran_ind",
		"disabled_ind",
		"widow_ind",
		"senior_ind",
		"school_college_ind",
		"religious_ind",
		"welfare_ind",
		"public_utility_ind",
		"cemetery_ind",
		"hospital_ind",
		"library_ind",
		"building_area",
		"building_area_ind",
		"sum_building_sq_ft",
		"sum_living_area_sq_ft",
		"sum_ground_floor_sq_ft",
		"sum_gross_area_sq_ft",
		"sum_adj_area_sq_ft",
		"attic_sq_ft",
		"attic_unfinished_sq_ft",
		"attic_finished_sq_ft",
		"sum_basement_sq_ft",
		"basement_unfinished_sq_ft",
		"basement_finished_sq_ft",
		"sum_garage_sq_ft",
		"garage_un_finished_sq_ft",
		"garage_finished_sq_ft",
		"year_built",
		"effective_year_built",
		"bedrooms",
		"total_rooms",
		"bath_total_calc",
		"bath_full",
		"baths_partial_nbr",
		"bath_fixtures_nbr",
		"amenities",
		"air_conditioning_code",
		"basement_code",
		"building_class_code",
		"building_condition_code",
		"construction_type_code",
		"deck_ind",
		"exterior_walls_code",
		"interior_walls_code",
		"fireplace_code",
		"floor_cover_code",
		"garage",
		"heat_code",
		"heating_fuel_type_code",
		"site_influence_code",
		"garage_parking_nbr",
		"driveway_code",
		"other_rooms",
		"patio_code",
		"pool_code",
		"porch_code",
		"building_quality_code",
		"roof_cover_code",
		"roof_type_code",
		"sewer_code",
		"stories_nbr_code",
		"style_code",
		"sum_residential_units",
		"sum_buildings_nbr",
		"sum_commercial_units",
		"topography_code",
		"water_code",
		"lot_code",
		"lot_nbr",
		"land_lot",
		"block",
		"section",
		"district",
		"legal_unit",
		"municipality",
		"subdivision_name",
		"subdivision_phase_nbr",
		"subdivision_tract_nbr",
		"meridian",
		"assessors_map_ref",
		"legal_description",
		"current_sale_transaction_id",
		"current_sale_doc_nbr",
		"current_sale_book",
		"current_sale_page",
		"current_sale_recording_date",
		"current_sale_contract_date",
		"current_sale_document_type",
		"current_sales_price",
		"current_sales_price_code",
		"current_sale_buyer1full_name",
		"current_sale_buyer2full_name",
		"current_sale_seller1full_name",
		"current_sale_seller2full_name",
		"concurrent_mtg1doc_nbr",
		"concurrent_mtg1book",
		"concurrent_mtg1page",
		"concurrent_mtg1recording_date",
		"concurrent_mtg1loan_amt",
		"concurrent_mtg1lender",
		"concurrent_mtg1term",
		"concurrent_mtg1interest_rate",
		"concurrent_mtg1loan_due_date",
		"concurrent_mtg1loan_type",
		"concurrent_mtg1type_financing",
		"concurrent_mtg2doc_nbr",
		"concurrent_mtg2book",
		"concurrent_mtg2page",
		"concurrent_mtg2recording_date",
		"concurrent_mtg2loan_amt",
		"concurrent_mtg2lender",
		"concurrent_mtg2term",
		"concurrent_mtg2interest_rate",
		"concurrent_mtg2loan_due_date",
		"concurrent_mtg2loan_type",
		"concurrent_mtg2typefinancing",
		"prev_sale_transaction_id",
		"prev_sale_doc_nbr",
		"prev_sale_book",
		"prev_sale_page",
		"prev_sale_recording_date",
		"prev_sale_contract_date",
		"prev_sale_document_type",
		"prev_sales_price",
		"prev_sales_price_code",
		"prev_sale_buyer1full_name",
		"prev_sale_buyer2full_name",
		"prev_sale_seller1full_name",
		"prev_sale_seller2full_name",
		"prev_mtg1doc_nbr",
		"prev_mtg1book",
		"prev_mtg1page",
		"prev_mtg1recording_date",
		"prev_mtg1loan_amt",
		"prev_mtg1lender",
		"prev_mtg1term",
		"prev_mtg1interest_rate",
		"prev_mtg1loan_due_date",
		"prev_mtg1loan_type",
		"prev_mtg1type_financing",
		"total_open_lien_nbr",
		"total_open_lien_amt",
		"mtg1transaction_id",
		"mtg1recording_date",
		"mtg1loan_amt",
		"mtg1lender",
		"mtg1private_lender",
		"mtg1term",
		"mtg1loan_due_date",
		"mtg1adj_rider",
		"mtg1loan_type",
		"mtg1type_financing",
		"mtg1lien_position",
		"mtg2transaction_id",
		"mtg2recording_date",
		"mtg2loan_amt",
		"mtg2lender",
		"mtg2private_lender",
		"mtg2term",
		"mtg2loan_due_date",
		"mtg2adj_rider",
		"mtg2loan_type",
		"mtg2type_financing",
		"mtg2lien_position",
		"mtg3transaction_id",
		"mtg3recording_date",
		"mtg3loan_amt",
		"mtg3lender",
		"mtg3private_lender",
		"mtg3term",
		"mtg3loan_due_date",
		"mtg3adj_rider",
		"mtg3loan_type",
		"mtg3type_financing",
		"mtg3lien_position",
		"mtg4transaction_id",
		"mtg4recording_date",
		"mtg4loan_amt",
		"mtg4lender",
		"mtg4private_lender",
		"mtg4term",
		"mtg4loan_due_date",
		"mtg4adj_rider",
		"mtg4loan_type",
		"mtg4type_financing",
		"mtg4lien_position",
		"fa_time_stamp",
		"fa_record_type",
	}
}

func (dr *Assessor) SQLTable() string {
	return "fa_df_assessor"
}

func (dr *Assessor) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "d9d51222-adc4-4b0e-bf2b-ad14b772f901",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to generate UUID.",
				Cause:  err.Error(),
			}
		}
		dr.AMId = u
	}

	now := time.Now()

	if dr.AMCreatedAt.IsZero() {
		dr.AMCreatedAt = now
	}

	values := []any{
		dr.AMId,
		dr.AMCreatedAt,
		now,
		dr.AMMeta,
		dr.FIPS,
		dr.PropertyID,
		dr.APN,
		dr.APNSeqNbr,
		dr.OldAPN,
		dr.OldApnIndicator,
		dr.TaxAccountNumber,
		dr.SitusFullStreetAddress,
		dr.SitusHouseNbr,
		dr.SitusHouseNbrSuffix,
		dr.SitusDirectionLeft,
		dr.SitusStreet,
		dr.SitusMode,
		dr.SitusDirectionRight,
		dr.SitusUnitType,
		dr.SitusUnitNbr,
		dr.SitusCity,
		dr.SitusState,
		dr.SitusZIP5,
		dr.SitusZIP4,
		dr.SitusCarrierCode,
		dr.SitusLatitude,
		dr.SitusLongitude,
		dr.SitusGeoStatusCode,
		dr.PropertyClassID,
		dr.LandUseCode,
		dr.StateLandUseCode,
		dr.CountyLandUseCode,
		dr.Zoning,
		dr.SitusCensusTract,
		dr.SitusCensusBlock,
		dr.MobileHomeInd,
		dr.TimeshareCode,
		dr.SchoolDistrictName,
		dr.LotSizeFrontageFeet,
		dr.LotSizeDepthFeet,
		dr.LotSizeAcres,
		dr.LotSizeSqFt,
		dr.Owner1CorpInd,
		dr.Owner1LastName,
		dr.Owner1FirstName,
		dr.Owner1MiddleName,
		dr.Owner1Suffix,
		dr.Owner2CorpInd,
		dr.Owner2LastName,
		dr.Owner2FirstName,
		dr.Owner2MiddleName,
		dr.Owner2Suffix,
		dr.OwnerNAME1FULL,
		dr.OwnerNAME2FULL,
		dr.OwnerOccupied,
		dr.Owner1OwnershipRights,
		dr.MailingFullStreetAddress,
		dr.MailingHouseNbr,
		dr.MailingHouseNbrSuffix,
		dr.MailingDirectionLeft,
		dr.MailingStreet,
		dr.MailingMode,
		dr.MailingDirectionRight,
		dr.MailingUnitType,
		dr.MailingUnitNbr,
		dr.MailingCity,
		dr.MailingState,
		dr.MailingZIP5,
		dr.MailingZIP4,
		dr.MailingCarrierCode,
		dr.MailingOptOut,
		dr.MailingCOName,
		dr.MailingForeignAddressInd,
		dr.AssdTotalValue,
		dr.AssdLandValue,
		dr.AssdImprovementValue,
		dr.MarketTotalValue,
		dr.MarketValueLand,
		dr.MarketValueImprovement,
		dr.TaxAmt,
		dr.TaxYear,
		dr.TaxDeliquentYear,
		dr.MarketYear,
		dr.AssdYear,
		dr.TaxRateCodeArea,
		dr.SchoolTaxDistrict1Code,
		dr.SchoolTaxDistrict2Code,
		dr.SchoolTaxDistrict3Code,
		dr.HomesteadInd,
		dr.VeteranInd,
		dr.DisabledInd,
		dr.WidowInd,
		dr.SeniorInd,
		dr.SchoolCollegeInd,
		dr.ReligiousInd,
		dr.WelfareInd,
		dr.PublicUtilityInd,
		dr.CemeteryInd,
		dr.HospitalInd,
		dr.LibraryInd,
		dr.BuildingArea,
		dr.BuildingAreaInd,
		dr.SumBuildingSqFt,
		dr.SumLivingAreaSqFt,
		dr.SumGroundFloorSqFt,
		dr.SumGrossAreaSqFt,
		dr.SumAdjAreaSqFt,
		dr.AtticSqFt,
		dr.AtticUnfinishedSqFt,
		dr.AtticFinishedSqFt,
		dr.SumBasementSqFt,
		dr.BasementUnfinishedSqFt,
		dr.BasementFinishedSqFt,
		dr.SumGarageSqFt,
		dr.GarageUnFinishedSqFt,
		dr.GarageFinishedSqFt,
		dr.YearBuilt,
		dr.EffectiveYearBuilt,
		dr.Bedrooms,
		dr.TotalRooms,
		dr.BathTotalCalc,
		dr.BathFull,
		dr.BathsPartialNbr,
		dr.BathFixturesNbr,
		dr.Amenities,
		dr.AirConditioningCode,
		dr.BasementCode,
		dr.BuildingClassCode,
		dr.BuildingConditionCode,
		dr.ConstructionTypeCode,
		dr.DeckInd,
		dr.ExteriorWallsCode,
		dr.InteriorWallsCode,
		dr.FireplaceCode,
		dr.FloorCoverCode,
		dr.Garage,
		dr.HeatCode,
		dr.HeatingFuelTypeCode,
		dr.SiteInfluenceCode,
		dr.GarageParkingNbr,
		dr.DrivewayCode,
		dr.OtherRooms,
		dr.PatioCode,
		dr.PoolCode,
		dr.PorchCode,
		dr.BuildingQualityCode,
		dr.RoofCoverCode,
		dr.RoofTypeCode,
		dr.SewerCode,
		dr.StoriesNbrCode,
		dr.StyleCode,
		dr.SumResidentialUnits,
		dr.SumBuildingsNbr,
		dr.SumCommercialUnits,
		dr.TopographyCode,
		dr.WaterCode,
		dr.LotCode,
		dr.LotNbr,
		dr.LandLot,
		dr.Block,
		dr.Section,
		dr.District,
		dr.LegalUnit,
		dr.Municipality,
		dr.SubdivisionName,
		dr.SubdivisionPhaseNbr,
		dr.SubdivisionTractNbr,
		dr.Meridian,
		dr.AssessorsMapRef,
		dr.LegalDescription,
		dr.CurrentSaleTransactionId,
		dr.CurrentSaleDocNbr,
		dr.CurrentSaleBook,
		dr.CurrentSalePage,
		dr.CurrentSaleRecordingDate,
		dr.CurrentSaleContractDate,
		dr.CurrentSaleDocumentType,
		dr.CurrentSalesPrice,
		dr.CurrentSalesPriceCode,
		dr.CurrentSaleBuyer1FullName,
		dr.CurrentSaleBuyer2FullName,
		dr.CurrentSaleSeller1FullName,
		dr.CurrentSaleSeller2FullName,
		dr.ConcurrentMtg1DocNbr,
		dr.ConcurrentMtg1Book,
		dr.ConcurrentMtg1Page,
		dr.ConcurrentMtg1RecordingDate,
		dr.ConcurrentMtg1LoanAmt,
		dr.ConcurrentMtg1Lender,
		dr.ConcurrentMtg1Term,
		dr.ConcurrentMtg1InterestRate,
		dr.ConcurrentMtg1LoanDueDate,
		dr.ConcurrentMtg1LoanType,
		dr.ConcurrentMtg1TypeFinancing,
		dr.ConcurrentMtg2DocNbr,
		dr.ConcurrentMtg2Book,
		dr.ConcurrentMtg2Page,
		dr.ConcurrentMtg2RecordingDate,
		dr.ConcurrentMtg2LoanAmt,
		dr.ConcurrentMtg2Lender,
		dr.ConcurrentMtg2Term,
		dr.ConcurrentMtg2InterestRate,
		dr.ConcurrentMtg2LoanDueDate,
		dr.ConcurrentMtg2LoanType,
		dr.ConcurrentMtg2Typefinancing,
		dr.PrevSaleTransactionId,
		dr.PrevSaleDocNbr,
		dr.PrevSaleBook,
		dr.PrevSalePage,
		dr.PrevSaleRecordingDate,
		dr.PrevSaleContractDate,
		dr.PrevSaleDocumentType,
		dr.PrevSalesPrice,
		dr.PrevSalesPriceCode,
		dr.PrevSaleBuyer1FullName,
		dr.PrevSaleBuyer2FullName,
		dr.PrevSaleSeller1FullName,
		dr.PrevSaleSeller2FullName,
		dr.PrevMtg1DocNbr,
		dr.PrevMtg1Book,
		dr.PrevMtg1Page,
		dr.PrevMtg1RecordingDate,
		dr.PrevMtg1LoanAmt,
		dr.PrevMtg1Lender,
		dr.PrevMtg1Term,
		dr.PrevMtg1InterestRate,
		dr.PrevMtg1LoanDueDate,
		dr.PrevMtg1LoanType,
		dr.PrevMtg1TypeFinancing,
		dr.TotalOpenLienNbr,
		dr.TotalOpenLienAmt,
		dr.Mtg1TransactionId,
		dr.Mtg1RecordingDate,
		dr.Mtg1LoanAmt,
		dr.Mtg1Lender,
		dr.Mtg1PrivateLender,
		dr.Mtg1Term,
		dr.Mtg1LoanDueDate,
		dr.Mtg1AdjRider,
		dr.Mtg1LoanType,
		dr.Mtg1TypeFinancing,
		dr.Mtg1LienPosition,
		dr.Mtg2TransactionId,
		dr.Mtg2RecordingDate,
		dr.Mtg2LoanAmt,
		dr.Mtg2Lender,
		dr.Mtg2PrivateLender,
		dr.Mtg2Term,
		dr.Mtg2LoanDueDate,
		dr.Mtg2AdjRider,
		dr.Mtg2LoanType,
		dr.Mtg2TypeFinancing,
		dr.Mtg2LienPosition,
		dr.Mtg3TransactionId,
		dr.Mtg3RecordingDate,
		dr.Mtg3LoanAmt,
		dr.Mtg3Lender,
		dr.Mtg3PrivateLender,
		dr.Mtg3Term,
		dr.Mtg3LoanDueDate,
		dr.Mtg3AdjRider,
		dr.Mtg3LoanType,
		dr.Mtg3TypeFinancing,
		dr.Mtg3LienPosition,
		dr.Mtg4TransactionId,
		dr.Mtg4RecordingDate,
		dr.Mtg4LoanAmt,
		dr.Mtg4Lender,
		dr.Mtg4PrivateLender,
		dr.Mtg4Term,
		dr.Mtg4LoanDueDate,
		dr.Mtg4AdjRider,
		dr.Mtg4LoanType,
		dr.Mtg4TypeFinancing,
		dr.Mtg4LienPosition,
		dr.FATimeStamp,
		dr.FARecordType,
	}

	return values, nil
}

func (dr *Assessor) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		LoadFunc: dr.LoadFunc,
		Mode:     entities.DataRecordModeLoadFunc,
	}
}

func (dr *Assessor) LoadFunc(r *arc.Request, in *entities.LoadDataRecordInput) (*entities.LoadDataRecordOutput, error) {
	newBuilder := func() squirrel.InsertBuilder {
		return squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Insert(in.DataRecord.SQLTable())
	}

	// Adjust batch size to account for the number of ids passed as params.
	batchSize := in.BatchSize - int32(math.Ceil(float64(in.BatchSize)/float64(len(in.Columns))))
	builder := newBuilder()
	deletedCount := int64(0)
	dfObject := in.DataFileObject
	recordCount := int32(0)
	processedRecords := int64(0)
	scanner := in.Scanner
	deleteIds := []int64{}
	insertIds := []int64{}

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "614361a6-0635-4b26-a77a-fc4e1d41c09a")
	}

	loadRecords := func() error {
		// Skip if no records to process.
		if recordCount == 0 {
			return nil
		}

		var markedExpr any

		if len(deleteIds) > 0 {
			markedExpr = squirrel.Eq{"property_id": deleteIds}
		} else {
			markedExpr = false
		}

		var deletedExpr any

		if len(insertIds) > 0 {
			deletedExpr = squirrel.Eq{"property_id": insertIds}
		} else {
			deletedExpr = false
		}

		prefixExpr := squirrel.Expr(
			`
			with marked_as_deleted_records as (
				update fa_df_assessor
				set am_deleted_at = now()
				where ?
			), deleted_records as (
				delete from fa_df_assessor
				where ?
				returning *
			), archived_records as (
				insert into fa_assessor_history
				select
					*,
					now() as am_archived_at
				from deleted_records
			)
			`,
			markedExpr,
			deletedExpr,
		)

		var sql string
		var args []any

		if len(insertIds) == 0 {
			// Create a select-only query when there are no insertions.
			sql, args, err = squirrel.StatementBuilder.
				PlaceholderFormat(squirrel.Dollar).
				Select("count(*)").
				PrefixExpr(prefixExpr).
				From("deleted_records").
				ToSql()
		} else {
			builder = builder.
				PrefixExpr(prefixExpr).
				Suffix(`returning (select count(*) from deleted_records)`)
			sql, args, err = builder.ToSql()
		}

		if err != nil {
			return &errors.Object{
				Id:     "0a5982e2-db38-4f37-9b94-d1868b2588f8",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to build SQL.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"deleteCount": len(deleteIds),
					"insertCount": len(insertIds),
				},
			}
		}

		tx, err := pgxPool.Begin(r.Context())
		if err != nil {
			return &errors.Object{
				Id:     "8826e710-2a99-4366-9559-9ad9f1fef1e3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to begin transaction.",
				Cause:  err.Error(),
			}
		}

		defer extutils.RollbackPgxTx(r.Context(), tx, "98250594-9322-4236-b2ab-ef02a38c1c60")

		// This clone won't replace the original arc.
		r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

		row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
		if err != nil {
			return errors.Forward(err, "5967e8cd-b6cd-4b32-af51-7973302f1935")
		}

		var batchDeletes int64

		if err := row.Scan(&batchDeletes); err != nil {
			return &errors.Object{
				Id:     "a3d44407-330d-46f1-b568-c9c4284b4ff1",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to scan record.",
				Cause:  err.Error(),
			}
		}

		deletedCount += batchDeletes

		updateObjectOut, err := in.UpdateDataFileObjectFunc(r, &entities.UpdateDataFileObjectInput{
			Id:          dfObject.Id,
			UpdatedAt:   time.Now(),
			RecordCount: dfObject.RecordCount + recordCount,
			Status:      entities.DataFileObjectStatusInProgress,
		})
		if err != nil {
			return errors.Forward(err, "c3dcd270-fc70-4dfc-8b51-fe3d40bc5457")
		}

		if err := tx.Commit(r.Context()); err != nil {
			return &errors.Object{
				Id:     "38892a84-1343-498a-8827-b1a0e1ef338d",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to commit transaction.",
				Cause:  err.Error(),
			}
		}

		builder = newBuilder()
		deleteIds = []int64{}
		dfObject = updateObjectOut.Entity
		insertIds = []int64{}
		recordCount = 0

		return nil
	}

	for scanner.Scan() {
		if recordCount == batchSize {
			if err := loadRecords(); err != nil {
				return nil, err
			}
		}

		fields := strings.Split(scanner.Text(), in.FieldSeparator)

		record, err := in.DataRecord.New(in.Headers, fields)
		if err != nil {
			return nil, errors.Forward(err, "91c65fcc-17da-4cf7-ab94-d1fd27935f65")
		}

		recordValues, err := record.SQLValues()
		if err != nil {
			return nil, errors.Forward(err, "440acdd5-d104-4a70-ad0b-7fc42915022c")
		}

		assessorRecord := record.(*Assessor)

		switch strings.ToUpper(val.PtrDeref(assessorRecord.FARecordType)) {
		case "D":
			deleteIds = append(deleteIds, assessorRecord.PropertyID)
		default:
			builder = builder.Values(recordValues...)
			insertIds = append(insertIds, assessorRecord.PropertyID)
		}

		recordCount++
		processedRecords++
	}

	if err := loadRecords(); err != nil {
		return nil, err
	}

	out := &entities.LoadDataRecordOutput{
		DeletedRecords:   deletedCount,
		ProcessedRecords: processedRecords,
	}

	return out, nil
}
