package attom_data

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

	ATTOMID                                    int64
	SitusStateCode                             *string
	SitusCounty                                *string
	PropertyJurisdictionName                   *string
	SitusStateCountyFIPS                       *string
	CombinedStatisticalArea                    *string
	CBSAName                                   *string
	CBSACode                                   *int
	MSAName                                    *string
	MSACode                                    *int
	MetropolitanDivision                       *string
	MinorCivilDivisionName                     *string
	MinorCivilDivisionCode                     *int
	NeighborhoodCode                           *string
	CensusFIPSPlaceCode                        *string
	CensusTract                                *int
	CensusBlockGroup                           *int
	CensusBlock                                *int
	ParcelNumberRaw                            *string
	ParcelNumberFormatted                      *string
	ParcelNumberYearAdded                      *int
	ParcelNumberAlternate                      *string
	ParcelMapBook                              *string
	ParcelMapPage                              *string
	ParcelNumberYearChange                     *int
	ParcelNumberPrevious                       *string
	ParcelAccountNumber                        *string
	PropertyAddressFull                        *string
	PropertyAddressHouseNumber                 *string
	PropertyAddressStreetDirection             *string
	PropertyAddressStreetName                  *string
	PropertyAddressStreetSuffix                *string
	PropertyAddressStreetPostDirection         *string
	PropertyAddressUnitPrefix                  *string
	PropertyAddressUnitValue                   *string
	PropertyAddressCity                        *string
	PropertyAddressState                       *string
	PropertyAddressZIP                         *string
	PropertyAddressZIP4                        *string
	PropertyAddressCRRT                        *string
	PropertyAddressInfoPrivacy                 *string
	CongressionalDistrictHouse                 *int
	PropertyLatitude                           *float64
	PropertyLongitude                          *float64
	GeoQuality                                 *int
	LegalDescription                           *string
	LegalRange                                 *string
	LegalTownship                              *string
	LegalSection                               *string
	LegalQuarter                               *string
	LegalQuarterQuarter                        *string
	LegalSubdivision                           *string
	LegalPhase                                 *string
	LegalTractNumber                           *string
	LegalBlock1                                *string
	LegalBlock2                                *string
	LegalLotNumber1                            *string
	LegalLotNumber2                            *string
	LegalLotNumber3                            *string
	LegalUnit                                  *string
	PartyOwner1NameFull                        *string
	PartyOwner1NameFirst                       *string
	PartyOwner1NameMiddle                      *string
	PartyOwner1NameLast                        *string
	PartyOwner1NameSuffix                      *string
	TrustDescription                           *string
	CompanyFlag                                *bool
	PartyOwner2NameFull                        *string
	PartyOwner2NameFirst                       *string
	PartyOwner2NameMiddle                      *string
	PartyOwner2NameLast                        *string
	PartyOwner2NameSuffix                      *string
	OwnerTypeDescription1                      *string
	OwnershipVestingRelationCode               *int
	PartyOwner3NameFull                        *string
	PartyOwner3NameFirst                       *string
	PartyOwner3NameMiddle                      *string
	PartyOwner3NameLast                        *string
	PartyOwner3NameSuffix                      *string
	PartyOwner4NameFull                        *string
	PartyOwner4NameFirst                       *string
	PartyOwner4NameMiddle                      *string
	PartyOwner4NameLast                        *string
	PartyOwner4NameSuffix                      *string
	OwnerTypeDescription2                      *string
	ContactOwnerMailingCounty                  *string
	ContactOwnerMailingFIPS                    *string
	ContactOwnerMailAddressFull                *string
	ContactOwnerMailAddressHouseNumber         *string
	ContactOwnerMailAddressStreetDirection     *string
	ContactOwnerMailAddressStreetName          *string
	ContactOwnerMailAddressStreetSuffix        *string
	ContactOwnerMailAddressStreetPostDirection *string
	ContactOwnerMailAddressUnitPrefix          *string
	ContactOwnerMailAddressUnit                *string
	ContactOwnerMailAddressCity                *string
	ContactOwnerMailAddressState               *string
	ContactOwnerMailAddressZIP                 *string
	ContactOwnerMailAddressZIP4                *string
	ContactOwnerMailAddressCRRT                *string
	ContactOwnerMailAddressInfoFormat          *string
	ContactOwnerMailInfoPrivacy                *string
	StatusOwnerOccupiedFlag                    *bool
	DeedOwner1NameFull                         *string
	DeedOwner1NameFirst                        *string
	DeedOwner1NameMiddle                       *string
	DeedOwner1NameLast                         *string
	DeedOwner1NameSuffix                       *string
	DeedOwner2NameFull                         *string
	DeedOwner2NameFirst                        *string
	DeedOwner2NameMiddle                       *string
	DeedOwner2NameLast                         *string
	DeedOwner2NameSuffix                       *string
	DeedOwner3NameFull                         *string
	DeedOwner3NameFirst                        *string
	DeedOwner3NameMiddle                       *string
	DeedOwner3NameLast                         *string
	DeedOwner3NameSuffix                       *string
	DeedOwner4NameFull                         *string
	DeedOwner4NameFirst                        *string
	DeedOwner4NameMiddle                       *string
	DeedOwner4NameLast                         *string
	DeedOwner4NameSuffix                       *string
	TaxYearAssessed                            *int
	TaxAssessedValueTotal                      *int
	TaxAssessedValueImprovements               *int
	TaxAssessedValueLand                       *int
	TaxAssessedImprovementsPerc                *decimal.Decimal
	PreviousAssessedValue                      *int
	TaxMarketValueYear                         *int
	TaxMarketValueTotal                        *int
	TaxMarketValueImprovements                 *int
	TaxMarketValueLand                         *int
	TaxMarketImprovementsPerc                  *decimal.Decimal
	TaxFiscalYear                              *int
	TaxRateArea                                *string
	TaxBilledAmount                            *decimal.Decimal
	TaxDelinquentYear                          *int
	LastAssessorTaxRollUpdate                  *time.Time
	AssrLastUpdated                            *time.Time
	TaxExemptionHomeownerFlag                  *bool
	TaxExemptionDisabledFlag                   *bool
	TaxExemptionSeniorFlag                     *bool
	TaxExemptionVeteranFlag                    *bool
	TaxExemptionWidowFlag                      *bool
	TaxExemptionAdditional                     *bool
	YearBuilt                                  *int
	YearBuiltEffective                         *int
	ZonedCodeLocal                             *string
	PropertyUseMuni                            *string
	PropertyUseGroup                           *string
	PropertyUseStandardized                    *int
	AssessorLastSaleDate                       *time.Time
	AssessorLastSaleAmount                     *int
	AssessorPriorSaleDate                      *time.Time
	AssessorPriorSaleAmount                    *int
	LastOwnershipTransferDate                  *time.Time
	LastOwnershipTransferDocumentNumber        *string
	LastOwnershipTransferTransactionID         *int64
	DeedLastSaleDocumentBook                   *string
	DeedLastSaleDocumentPage                   *string
	DeedLastDocumentNumber                     *string
	DeedLastSaleDate                           *time.Time
	DeedLastSalePrice                          *int
	DeedLastSaleTransactionID                  *int64
	AreaBuilding                               *int
	AreaBuildingDefinitionCode                 *int
	AreaGross                                  *int
	Area1stFloor                               *int
	Area2ndFloor                               *int
	AreaUpperFloors                            *int
	AreaLotAcres                               *float64
	AreaLotSF                                  *decimal.Decimal
	AreaLotDepth                               *decimal.Decimal
	AreaLotWidth                               *decimal.Decimal
	RoomsAtticArea                             *int
	RoomsAtticFlag                             *bool
	RoomsBasementArea                          *int
	RoomsBasementAreaFinished                  *int
	RoomsBasementAreaUnfinished                *int
	ParkingGarage                              *int
	ParkingGarageArea                          *int
	ParkingCarport                             *bool
	ParkingCarportArea                         *int
	HVACCoolingDetail                          *int
	HVACHeatingDetail                          *int
	HVACHeatingFuel                            *int
	UtilitiesSewageUsage                       *int
	UtilitiesWaterSource                       *int
	UtilitiesMobileHomeHookupFlag              *bool
	Foundation                                 *int
	Construction                               *int
	InteriorStructure                          *int
	PlumbingFixturesCount                      *int
	ConstructionFireResistanceClass            *int
	SafetyFireSprinklersFlag                   *bool
	FlooringMaterialPrimary                    *int
	BathCount                                  *decimal.Decimal
	BathPartialCount                           *int
	BedroomsCount                              *int
	RoomsCount                                 *int
	StoriesCount                               *int
	UnitsCount                                 *int
	RoomsBonusRoomFlag                         *bool
	RoomsBreakfastNookFlag                     *bool
	RoomsCellarFlag                            *bool
	RoomsCellarWineFlag                        *bool
	RoomsExerciseFlag                          *bool
	RoomsFamilyCode                            *bool
	RoomsGameFlag                              *bool
	RoomsGreatFlag                             *bool
	RoomsHobbyFlag                             *bool
	RoomsLaundryFlag                           *bool
	RoomsMediaFlag                             *bool
	RoomsMudFlag                               *bool
	RoomsOfficeArea                            *int
	RoomsOfficeFlag                            *bool
	RoomsSafeRoomFlag                          *bool
	RoomsSittingFlag                           *bool
	RoomsStormShelter                          *bool
	RoomsStudyFlag                             *bool
	RoomsSunroomFlag                           *bool
	RoomsUtilityArea                           *int
	RoomsUtilityCode                           *bool
	Fireplace                                  *int
	FireplaceCount                             *int
	AccessabilityElevatorFlag                  *bool
	AccessabilityHandicapFlag                  *bool
	EscalatorFlag                              *bool
	CentralVacuumFlag                          *bool
	ContentIntercomFlag                        *bool
	ContentSoundSystemFlag                     *bool
	WetBarFlag                                 *bool
	SecurityAlarmFlag                          *bool
	StructureStyle                             *int
	Exterior1Code                              *string
	RoofMaterial                               *int
	RoofConstruction                           *int
	ContentStormShutterFlag                    *bool
	ContentOverheadDoorFlag                    *bool
	ViewDescription                            *string
	PorchCode                                  *string
	PorchArea                                  *int
	PatioArea                                  *int
	DeckFlag                                   *bool
	DeckArea                                   *int
	FeatureBalconyFlag                         *bool
	BalconyArea                                *int
	BreezewayFlag                              *bool
	ParkingRVParkingFlag                       *bool
	ParkingSpaceCount                          *int
	DrivewayArea                               *int
	DrivewayMaterial                           *string
	Pool                                       *int
	PoolArea                                   *int
	ContentSaunaFlag                           *bool
	TopographyCode                             *int
	FenceCode                                  *bool
	FenceArea                                  *int
	CourtyardFlag                              *bool
	CourtyardArea                              *int
	ArborPergolaFlag                           *bool
	SprinklersFlag                             *bool
	GolfCourseGreenFlag                        *bool
	TennisCourtFlag                            *bool
	SportsCourtFlag                            *bool
	ArenaFlag                                  *bool
	WaterFeatureFlag                           *bool
	PondFlag                                   *bool
	BoatLiftFlag                               *bool
	BuildingsCount                             *int
	BathHouseArea                              *int
	BathHouseFlag                              *bool
	BoatAccessFlag                             *bool
	BoatHouseArea                              *int
	BoatHouseFlag                              *bool
	CabinArea                                  *int
	CabinFlag                                  *bool
	CanopyArea                                 *int
	CanopyFlag                                 *bool
	GazeboArea                                 *int
	GazeboFlag                                 *bool
	GraineryArea                               *int
	GraineryFlag                               *bool
	GreenHouseArea                             *int
	GreenHouseFlag                             *bool
	GuestHouseArea                             *int
	GuestHouseFlag                             *bool
	KennelArea                                 *int
	KennelFlag                                 *bool
	LeanToArea                                 *int
	LeanToFlag                                 *bool
	LoadingPlatformArea                        *int
	LoadingPlatformFlag                        *bool
	MilkHouseArea                              *int
	MilkHouseFlag                              *bool
	OutdoorKitchenFireplaceFlag                *bool
	PoolHouseArea                              *int
	PoolHouseFlag                              *bool
	PoultryHouseArea                           *int
	PoultryHouseFlag                           *bool
	QuonsetArea                                *int
	QuonsetFlag                                *bool
	ShedArea                                   *int
	ShedCode                                   *bool
	SiloArea                                   *int
	SiloFlag                                   *bool
	StableArea                                 *int
	StableFlag                                 *bool
	StorageBuildingArea                        *int
	StorageBuildingFlag                        *bool
	UtilityBuildingArea                        *int
	UtilityBuildingFlag                        *bool
	PoleStructureArea                          *int
	PoleStructureFlag                          *bool
	CommunityRecRoomFlag                       *bool
	PublicationDate                            *time.Time
	ParcelShellRecord                          *bool
}

func (dr *Assessor) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(Assessor)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "[ATTOM ID]":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:     "aa8f0c82-7924-44a1-9b60-a265cd59b8fa",
					Code:   errors.Code_UNKNOWN,
					Detail: "ATTOMID is required.",
					Cause:  err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.ATTOMID = v
		case "SitusStateCode":
			record.SitusStateCode = val.StringPtrIfNonZero(field)
		case "SitusCounty":
			record.SitusCounty = val.StringPtrIfNonZero(field)
		case "PropertyJurisdictionName":
			record.PropertyJurisdictionName = val.StringPtrIfNonZero(field)
		case "SitusStateCountyFIPS":
			record.SitusStateCountyFIPS = val.StringPtrIfNonZero(field)
		case "CombinedStatisticalArea":
			record.CombinedStatisticalArea = val.StringPtrIfNonZero(field)
		case "CBSAName":
			record.CBSAName = val.StringPtrIfNonZero(field)
		case "CBSACode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2daa0bef-50a9-4cfd-811f-36d4f45392e9")
			}
			record.CBSACode = v
		case "MSAName":
			record.MSAName = val.StringPtrIfNonZero(field)
		case "MSACode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "cc1b64c3-bfce-4399-bab2-3a237d726aba")
			}
			record.MSACode = v
		case "MetropolitanDivision":
			record.MetropolitanDivision = val.StringPtrIfNonZero(field)
		case "MinorCivilDivisionName":
			record.MinorCivilDivisionName = val.StringPtrIfNonZero(field)
		case "MinorCivilDivisionCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "43bd276d-d8f3-4b41-984a-96789c6f4c15")
			}
			record.MinorCivilDivisionCode = v
		case "NeighborhoodCode":
			record.NeighborhoodCode = val.StringPtrIfNonZero(field)
		case "CensusFIPSPlaceCode":
			record.CensusFIPSPlaceCode = val.StringPtrIfNonZero(field)
		case "CensusTract":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ed4137ed-be3a-4915-a17c-1857d532b63b")
			}
			record.CensusTract = v
		case "CensusBlockGroup":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4f6b655b-8899-472c-b741-ec5672a90320")
			}
			record.CensusBlockGroup = v
		case "CensusBlock":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a300e49f-832f-4d6b-a191-7bc58fe46e60")
			}
			record.CensusBlock = v
		case "ParcelNumberRaw":
			record.ParcelNumberRaw = val.StringPtrIfNonZero(field)
		case "ParcelNumberFormatted":
			record.ParcelNumberFormatted = val.StringPtrIfNonZero(field)
		case "ParcelNumberYearAdded":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f3f0094e-0ae0-4fbe-933e-96bbd4eabd8d")
			}
			record.ParcelNumberYearAdded = v
		case "ParcelNumberAlternate":
			record.ParcelNumberAlternate = val.StringPtrIfNonZero(field)
		case "ParcelMapBook":
			record.ParcelMapBook = val.StringPtrIfNonZero(field)
		case "ParcelMapPage":
			record.ParcelMapPage = val.StringPtrIfNonZero(field)
		case "ParcelNumberYearChange":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a83d7a3e-bf64-4c01-b550-9c4359eca121")
			}
			record.ParcelNumberYearChange = v
		case "ParcelNumberPrevious":
			record.ParcelNumberPrevious = val.StringPtrIfNonZero(field)
		case "ParcelAccountNumber":
			record.ParcelAccountNumber = val.StringPtrIfNonZero(field)
		case "PropertyAddressFull":
			record.PropertyAddressFull = val.StringPtrIfNonZero(field)
		case "PropertyAddressHouseNumber":
			record.PropertyAddressHouseNumber = val.StringPtrIfNonZero(field)
		case "PropertyAddressStreetDirection":
			record.PropertyAddressStreetDirection = val.StringPtrIfNonZero(field)
		case "PropertyAddressStreetName":
			record.PropertyAddressStreetName = val.StringPtrIfNonZero(field)
		case "PropertyAddressStreetSuffix":
			record.PropertyAddressStreetSuffix = val.StringPtrIfNonZero(field)
		case "PropertyAddressStreetPostDirection":
			record.PropertyAddressStreetPostDirection = val.StringPtrIfNonZero(field)
		case "PropertyAddressUnitPrefix":
			record.PropertyAddressUnitPrefix = val.StringPtrIfNonZero(field)
		case "PropertyAddressUnitValue":
			record.PropertyAddressUnitValue = val.StringPtrIfNonZero(field)
		case "PropertyAddressCity":
			record.PropertyAddressCity = val.StringPtrIfNonZero(field)
		case "PropertyAddressState":
			record.PropertyAddressState = val.StringPtrIfNonZero(field)
		case "PropertyAddressZIP":
			record.PropertyAddressZIP = val.StringPtrIfNonZero(field)
		case "PropertyAddressZIP4":
			record.PropertyAddressZIP4 = val.StringPtrIfNonZero(field)
		case "PropertyAddressCRRT":
			record.PropertyAddressCRRT = val.StringPtrIfNonZero(field)
		case "PropertyAddressInfoPrivacy":
			record.PropertyAddressInfoPrivacy = val.StringPtrIfNonZero(field)
		case "CongressionalDistrictHouse":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "10294589-a0f7-46f5-86f0-2b9e4bd14f10")
			}
			record.CongressionalDistrictHouse = v
		case "PropertyLatitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "782de475-f53b-40e2-bdcb-5be0a6c05318")
			}
			record.PropertyLatitude = v
		case "PropertyLongitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4cb2aa2f-d67c-4e64-899a-35f9f553f514")
			}
			record.PropertyLongitude = v
		case "GeoQuality":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2e158b96-79a3-42ca-bd02-0c3cfe72a199")
			}
			record.GeoQuality = v
		case "LegalDescription":
			record.LegalDescription = val.StringPtrIfNonZero(field)
		case "LegalRange":
			record.LegalRange = val.StringPtrIfNonZero(field)
		case "LegalTownship":
			record.LegalTownship = val.StringPtrIfNonZero(field)
		case "LegalSection":
			record.LegalSection = val.StringPtrIfNonZero(field)
		case "LegalQuarter":
			record.LegalQuarter = val.StringPtrIfNonZero(field)
		case "LegalQuarterQuarter":
			record.LegalQuarterQuarter = val.StringPtrIfNonZero(field)
		case "LegalSubdivision":
			record.LegalSubdivision = val.StringPtrIfNonZero(field)
		case "LegalPhase":
			record.LegalPhase = val.StringPtrIfNonZero(field)
		case "LegalTractNumber":
			record.LegalTractNumber = val.StringPtrIfNonZero(field)
		case "LegalBlock1":
			record.LegalBlock1 = val.StringPtrIfNonZero(field)
		case "LegalBlock2":
			record.LegalBlock2 = val.StringPtrIfNonZero(field)
		case "LegalLotNumber1":
			record.LegalLotNumber1 = val.StringPtrIfNonZero(field)
		case "LegalLotNumber2":
			record.LegalLotNumber2 = val.StringPtrIfNonZero(field)
		case "LegalLotNumber3":
			record.LegalLotNumber3 = val.StringPtrIfNonZero(field)
		case "LegalUnit":
			record.LegalUnit = val.StringPtrIfNonZero(field)
		case "PartyOwner1NameFull":
			record.PartyOwner1NameFull = val.StringPtrIfNonZero(field)
		case "PartyOwner1NameFirst":
			record.PartyOwner1NameFirst = val.StringPtrIfNonZero(field)
		case "PartyOwner1NameMiddle":
			record.PartyOwner1NameMiddle = val.StringPtrIfNonZero(field)
		case "PartyOwner1NameLast":
			record.PartyOwner1NameLast = val.StringPtrIfNonZero(field)
		case "PartyOwner1NameSuffix":
			record.PartyOwner1NameSuffix = val.StringPtrIfNonZero(field)
		case "TrustDescription":
			record.TrustDescription = val.StringPtrIfNonZero(field)
		case "CompanyFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4b6f5df9-1e8e-4f68-a21c-4deefe26da38")
			}
			record.CompanyFlag = v
		case "PartyOwner2NameFull":
			record.PartyOwner2NameFull = val.StringPtrIfNonZero(field)
		case "PartyOwner2NameFirst":
			record.PartyOwner2NameFirst = val.StringPtrIfNonZero(field)
		case "PartyOwner2NameMiddle":
			record.PartyOwner2NameMiddle = val.StringPtrIfNonZero(field)
		case "PartyOwner2NameLast":
			record.PartyOwner2NameLast = val.StringPtrIfNonZero(field)
		case "PartyOwner2NameSuffix":
			record.PartyOwner2NameSuffix = val.StringPtrIfNonZero(field)
		case "OwnerTypeDescription1":
			record.OwnerTypeDescription1 = val.StringPtrIfNonZero(field)
		case "OwnershipVestingRelationCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fb58ede6-8859-49be-be65-d095ea5135cd")
			}
			record.OwnershipVestingRelationCode = v
		case "PartyOwner3NameFull":
			record.PartyOwner3NameFull = val.StringPtrIfNonZero(field)
		case "PartyOwner3NameFirst":
			record.PartyOwner3NameFirst = val.StringPtrIfNonZero(field)
		case "PartyOwner3NameMiddle":
			record.PartyOwner3NameMiddle = val.StringPtrIfNonZero(field)
		case "PartyOwner3NameLast":
			record.PartyOwner3NameLast = val.StringPtrIfNonZero(field)
		case "PartyOwner3NameSuffix":
			record.PartyOwner3NameSuffix = val.StringPtrIfNonZero(field)
		case "PartyOwner4NameFull":
			record.PartyOwner4NameFull = val.StringPtrIfNonZero(field)
		case "PartyOwner4NameFirst":
			record.PartyOwner4NameFirst = val.StringPtrIfNonZero(field)
		case "PartyOwner4NameMiddle":
			record.PartyOwner4NameMiddle = val.StringPtrIfNonZero(field)
		case "PartyOwner4NameLast":
			record.PartyOwner4NameLast = val.StringPtrIfNonZero(field)
		case "PartyOwner4NameSuffix":
			record.PartyOwner4NameSuffix = val.StringPtrIfNonZero(field)
		case "OwnerTypeDescription2":
			record.OwnerTypeDescription2 = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailingCounty":
			record.ContactOwnerMailingCounty = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailingFIPS":
			record.ContactOwnerMailingFIPS = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressFull":
			record.ContactOwnerMailAddressFull = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressHouseNumber":
			record.ContactOwnerMailAddressHouseNumber = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressStreetDirection":
			record.ContactOwnerMailAddressStreetDirection = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressStreetName":
			record.ContactOwnerMailAddressStreetName = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressStreetSuffix":
			record.ContactOwnerMailAddressStreetSuffix = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressStreetPostDirection":
			record.ContactOwnerMailAddressStreetPostDirection = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressUnitPrefix":
			record.ContactOwnerMailAddressUnitPrefix = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressUnit":
			record.ContactOwnerMailAddressUnit = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressCity":
			record.ContactOwnerMailAddressCity = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressState":
			record.ContactOwnerMailAddressState = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressZIP":
			record.ContactOwnerMailAddressZIP = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressZIP4":
			record.ContactOwnerMailAddressZIP4 = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressCRRT":
			record.ContactOwnerMailAddressCRRT = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailAddressInfoFormat":
			record.ContactOwnerMailAddressInfoFormat = val.StringPtrIfNonZero(field)
		case "ContactOwnerMailInfoPrivacy":
			record.ContactOwnerMailInfoPrivacy = val.StringPtrIfNonZero(field)
		case "StatusOwnerOccupiedFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "61cb260a-ee85-44be-9500-96a3e2c67779")
			}
			record.StatusOwnerOccupiedFlag = v
		case "DeedOwner1NameFull":
			record.DeedOwner1NameFull = val.StringPtrIfNonZero(field)
		case "DeedOwner1NameFirst":
			record.DeedOwner1NameFirst = val.StringPtrIfNonZero(field)
		case "DeedOwner1NameMiddle":
			record.DeedOwner1NameMiddle = val.StringPtrIfNonZero(field)
		case "DeedOwner1NameLast":
			record.DeedOwner1NameLast = val.StringPtrIfNonZero(field)
		case "DeedOwner1NameSuffix":
			record.DeedOwner1NameSuffix = val.StringPtrIfNonZero(field)
		case "DeedOwner2NameFull":
			record.DeedOwner2NameFull = val.StringPtrIfNonZero(field)
		case "DeedOwner2NameFirst":
			record.DeedOwner2NameFirst = val.StringPtrIfNonZero(field)
		case "DeedOwner2NameMiddle":
			record.DeedOwner2NameMiddle = val.StringPtrIfNonZero(field)
		case "DeedOwner2NameLast":
			record.DeedOwner2NameLast = val.StringPtrIfNonZero(field)
		case "DeedOwner2NameSuffix":
			record.DeedOwner2NameSuffix = val.StringPtrIfNonZero(field)
		case "DeedOwner3NameFull":
			record.DeedOwner3NameFull = val.StringPtrIfNonZero(field)
		case "DeedOwner3NameFirst":
			record.DeedOwner3NameFirst = val.StringPtrIfNonZero(field)
		case "DeedOwner3NameMiddle":
			record.DeedOwner3NameMiddle = val.StringPtrIfNonZero(field)
		case "DeedOwner3NameLast":
			record.DeedOwner3NameLast = val.StringPtrIfNonZero(field)
		case "DeedOwner3NameSuffix":
			record.DeedOwner3NameSuffix = val.StringPtrIfNonZero(field)
		case "DeedOwner4NameFull":
			record.DeedOwner4NameFull = val.StringPtrIfNonZero(field)
		case "DeedOwner4NameFirst":
			record.DeedOwner4NameFirst = val.StringPtrIfNonZero(field)
		case "DeedOwner4NameMiddle":
			record.DeedOwner4NameMiddle = val.StringPtrIfNonZero(field)
		case "DeedOwner4NameLast":
			record.DeedOwner4NameLast = val.StringPtrIfNonZero(field)
		case "DeedOwner4NameSuffix":
			record.DeedOwner4NameSuffix = val.StringPtrIfNonZero(field)
		case "TaxYearAssessed":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b2e67606-1ed3-40d7-9837-3b3c3a2f9aa4")
			}
			record.TaxYearAssessed = v
		case "TaxAssessedValueTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "68acdb42-55c3-44e5-aec9-baa6fe25d795")
			}
			record.TaxAssessedValueTotal = v
		case "TaxAssessedValueImprovements":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ddd2c143-5499-46a5-a538-151efd16117f")
			}
			record.TaxAssessedValueImprovements = v
		case "TaxAssessedValueLand":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0c873f6a-fed0-449a-b80a-9b2912abe739")
			}
			record.TaxAssessedValueLand = v
		case "TaxAssessedImprovementsPerc":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "58627a99-9166-4db4-8301-cc60d955776b")
			}
			record.TaxAssessedImprovementsPerc = v
		case "PreviousAssessedValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "073164e2-25d7-4bb1-a83f-d011f52876af")
			}
			record.PreviousAssessedValue = v
		case "TaxMarketValueYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a221b9e0-8574-44ae-a697-a651ab5fe29c")
			}
			record.TaxMarketValueYear = v
		case "TaxMarketValueTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4be12ec7-d45c-44ac-9a28-f8c8b7358a6d")
			}
			record.TaxMarketValueTotal = v
		case "TaxMarketValueImprovements":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7e3fff59-1acd-48af-b73f-5a4a784fc24c")
			}
			record.TaxMarketValueImprovements = v
		case "TaxMarketValueLand":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4cd11c3e-6c76-440f-af85-d27a3ca5ed29")
			}
			record.TaxMarketValueLand = v
		case "TaxMarketImprovementsPerc":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a92b6d5b-f0d9-436e-bb5a-40ddc36cf46e")
			}
			record.TaxMarketImprovementsPerc = v
		case "TaxFiscalYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "66d84f25-66c5-4146-a359-7cd6537ed24b")
			}
			record.TaxFiscalYear = v
		case "TaxRateArea":
			record.TaxRateArea = val.StringPtrIfNonZero(field)
		case "TaxBilledAmount":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c53dac5e-3d62-46d2-8c2c-77e8ea5f8e26")
			}
			record.TaxBilledAmount = v
		case "TaxDelinquentYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fca1141f-e91d-45f5-abf3-cf511bc09879")
			}
			record.TaxDelinquentYear = v
		case "LastAssessorTaxRollUpdate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "14b64f5b-5d85-46ad-826e-d016878bd3e9")
			}
			record.LastAssessorTaxRollUpdate = v
		case "AssrLastUpdated":
			if field == "" {
				field = "1800-01-01"
			}

			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "5309adcf-bc31-42d0-993e-a467ff2fcf5a")
			}
			record.AssrLastUpdated = v
		case "TaxExemptionHomeownerFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "18e6a435-44f7-4096-9a66-357b9b992c08")
			}
			record.TaxExemptionHomeownerFlag = v
		case "TaxExemptionDisabledFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8f00dad5-568c-48eb-84a7-9a8c741db24a")
			}
			record.TaxExemptionDisabledFlag = v
		case "TaxExemptionSeniorFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "78405c70-9b58-4625-8645-2c2dfd0159d3")
			}
			record.TaxExemptionSeniorFlag = v
		case "TaxExemptionVeteranFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1ab64f25-c124-413d-9512-940385429f38")
			}
			record.TaxExemptionVeteranFlag = v
		case "TaxExemptionWidowFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8da759e8-2325-4196-baf3-fa470e837751")
			}
			record.TaxExemptionWidowFlag = v
		case "TaxExemptionAdditional":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8c59689f-1edc-4e03-8c52-d05f65ddae92")
			}
			record.TaxExemptionAdditional = v
		case "YearBuilt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fb3587c5-4956-4b49-8778-e362f167d079")
			}
			record.YearBuilt = v
		case "YearBuiltEffective":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a09fcf99-ebe1-4a7f-b048-5276f50b2d47")
			}
			record.YearBuiltEffective = v
		case "ZonedCodeLocal":
			record.ZonedCodeLocal = val.StringPtrIfNonZero(field)
		case "PropertyUseMuni":
			record.PropertyUseMuni = val.StringPtrIfNonZero(field)
		case "PropertyUseGroup":
			record.PropertyUseGroup = val.StringPtrIfNonZero(field)
		case "PropertyUseStandardized":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0d0a501c-3769-4fda-8e2a-ea9b9b4f0f86")
			}
			record.PropertyUseStandardized = v
		case "AssessorLastSaleDate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "b250cffd-645f-4299-9d95-b7b7ad0a9f98")
			}
			record.AssessorLastSaleDate = v
		case "AssessorLastSaleAmount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9f3ae355-dfa0-456b-9273-4e288a4f669a")
			}
			record.AssessorLastSaleAmount = v
		case "AssessorPriorSaleDate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "88d96fcb-aab6-4d4e-9c6e-231e91b0de00")
			}
			record.AssessorPriorSaleDate = v
		case "AssessorPriorSaleAmount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "50520801-0b00-46b3-8c43-223f5e44a1ae")
			}
			record.AssessorPriorSaleAmount = v
		case "LastOwnershipTransferDate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "216b1bf1-ca93-4927-ab83-0cee7a6e81ca")
			}
			record.LastOwnershipTransferDate = v
		case "LastOwnershipTransferDocumentNumber":
			record.LastOwnershipTransferDocumentNumber = val.StringPtrIfNonZero(field)
		case "LastOwnershipTransferTransactionID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ca2fec54-3c8c-4510-b13d-62453da03824")
			}
			record.LastOwnershipTransferTransactionID = v
		case "DeedLastSaleDocumentBook":
			record.DeedLastSaleDocumentBook = val.StringPtrIfNonZero(field)
		case "DeedLastSaleDocumentPage":
			record.DeedLastSaleDocumentPage = val.StringPtrIfNonZero(field)
		case "DeedLastDocumentNumber":
			record.DeedLastDocumentNumber = val.StringPtrIfNonZero(field)
		case "DeedLastSaleDate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "72ee4cc2-49ac-43de-a4f6-9d9ea1800bda")
			}
			record.DeedLastSaleDate = v
		case "DeedLastSalePrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4da386c3-61e9-4a02-a654-ed3c4713cefb")
			}
			record.DeedLastSalePrice = v
		case "DeedLastSaleTransactionID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e47e5612-df48-4df5-9183-703331792ca9")
			}
			record.DeedLastSaleTransactionID = v
		case "AreaBuilding":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "52fbe39b-be6f-4651-b50e-ed08b1dc6082")
			}
			record.AreaBuilding = v
		case "AreaBuildingDefinitionCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f30db5c2-2faa-42d9-bf39-7cd992cdc5cd")
			}
			record.AreaBuildingDefinitionCode = v
		case "AreaGross":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b216425e-e4f6-4a86-80cf-f870042171f6")
			}
			record.AreaGross = v
		case "Area1stFloor":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9cfb217c-e1df-4eb8-87d2-704789dd632c")
			}
			record.Area1stFloor = v
		case "Area2ndFloor":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "53a170d1-679f-4049-81e8-b7e3e00f7e09")
			}
			record.Area2ndFloor = v
		case "AreaUpperFloors":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e7908e87-8d71-4316-8de4-b462079097b7")
			}
			record.AreaUpperFloors = v
		case "AreaLotAcres":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3a833377-36f2-40f2-ac29-3b507f19ce57")
			}
			record.AreaLotAcres = v
		case "AreaLotSF":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3ad72928-0094-4e22-9f53-6956636137a7")
			}
			record.AreaLotSF = v
		case "AreaLotDepth":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1fc7810a-78db-460d-ae8b-81054b7456d1")
			}
			record.AreaLotDepth = v
		case "AreaLotWidth":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "cec95914-8ae9-40d8-9d1d-6fb126bb4a3b")
			}
			record.AreaLotWidth = v
		case "RoomsAtticArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "42fe0c2b-8d3e-41ba-8a29-2dfc462ca14e")
			}
			record.RoomsAtticArea = v
		case "RoomsAtticFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8541f5fa-3d9a-485f-98f6-c2a54b8a7cec")
			}
			record.RoomsAtticFlag = v
		case "RoomsBasementArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8b5dad71-5d75-4128-9aa3-cc85591429d0")
			}
			record.RoomsBasementArea = v
		case "RoomsBasementAreaFinished":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d9188b9f-624b-489b-a697-626618ffc674")
			}
			record.RoomsBasementAreaFinished = v
		case "RoomsBasementAreaUnfinished":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ce55ee2e-6cd0-442c-babc-970686738e03")
			}
			record.RoomsBasementAreaUnfinished = v
		case "ParkingGarage":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6d687a06-83a1-4b26-9084-ed5b2f75b75d")
			}
			record.ParkingGarage = v
		case "ParkingGarageArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "30f65d31-0410-41f0-bfed-d1b346567178")
			}
			record.ParkingGarageArea = v
		case "ParkingCarport":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d26d4710-c189-49f4-bbac-a1a54130170a")
			}
			record.ParkingCarport = v
		case "ParkingCarportArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a760c57e-f6dd-42cf-abe6-64192ca11546")
			}
			record.ParkingCarportArea = v
		case "HVACCoolingDetail":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ef88aee5-fd25-42e5-ac11-95e4b36407e2")
			}
			record.HVACCoolingDetail = v
		case "HVACHeatingDetail":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c146b2a9-8016-4cb4-98d1-9baaf81a1251")
			}
			record.HVACHeatingDetail = v
		case "HVACHeatingFuel":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c807155e-8971-49df-962e-0d231c5a3894")
			}
			record.HVACHeatingFuel = v
		case "UtilitiesSewageUsage":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1ab51a99-4eff-4535-a510-04b0696d3d29")
			}
			record.UtilitiesSewageUsage = v
		case "UtilitiesWaterSource":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "590149ef-72cb-4576-8fe8-c836afe1d757")
			}
			record.UtilitiesWaterSource = v
		case "UtilitiesMobileHomeHookupFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2549adde-cd4b-4d29-9ece-85b92a33a1e0")
			}
			record.UtilitiesMobileHomeHookupFlag = v
		case "Foundation":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4974c02f-5b47-4660-b966-8e29dbdcb36b")
			}
			record.Foundation = v
		case "Construction":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b78b1999-17f3-43a6-b55a-209975ac41c2")
			}
			record.Construction = v
		case "InteriorStructure":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9e7b7c71-3187-4c31-905d-5cdc889bfbce")
			}
			record.InteriorStructure = v
		case "PlumbingFixturesCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "24d0d2a6-29d6-4c5b-b93d-03cee816272b")
			}
			record.PlumbingFixturesCount = v
		case "ConstructionFireResistanceClass":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4b396868-44dd-4d81-bf48-984db0fa9b19")
			}
			record.ConstructionFireResistanceClass = v
		case "SafetyFireSprinklersFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "44b56891-a403-4b00-bd0a-760f2533c8b4")
			}
			record.SafetyFireSprinklersFlag = v
		case "FlooringMaterialPrimary":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "83dbfecf-17ec-4509-b6b5-48e60fba1d7d")
			}
			record.FlooringMaterialPrimary = v
		case "BathCount":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e2c9bfa9-10fe-4cc1-b4b0-7d9f545a3fca")
			}
			record.BathCount = v
		case "BathPartialCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "21979957-3ffe-4acc-bef4-29e3290ddbee")
			}
			record.BathPartialCount = v
		case "BedroomsCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a05ba604-99be-4506-aa65-88b719440330")
			}
			record.BedroomsCount = v
		case "RoomsCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "85502820-acd4-4b62-8c8c-95a2190aa416")
			}
			record.RoomsCount = v
		case "StoriesCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b55b70dd-c86f-4623-943e-fd6bf38ab7cc")
			}
			record.StoriesCount = v
		case "UnitsCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b7342b51-d044-4729-89bd-958b7df5c68c")
			}
			record.UnitsCount = v
		case "RoomsBonusRoomFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9bbc9f41-c91d-4811-a624-92ae3817c1c1")
			}
			record.RoomsBonusRoomFlag = v
		case "RoomsBreakfastNookFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9b5033ce-b7e3-4e02-bd0b-6ea4cd7e44b7")
			}
			record.RoomsBreakfastNookFlag = v
		case "RoomsCellarFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3179a1e2-d494-4a4b-bd61-06560f9260b0")
			}
			record.RoomsCellarFlag = v
		case "RoomsCellarWineFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8d4633e0-b541-4146-aebc-725d822ff0ff")
			}
			record.RoomsCellarWineFlag = v
		case "RoomsExerciseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4b5abbad-dcc9-4714-8e0a-86c376cf274c")
			}
			record.RoomsExerciseFlag = v
		case "RoomsFamilyCode":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6588bd28-3fd1-40d1-b029-00ed55a02675")
			}
			record.RoomsFamilyCode = v
		case "RoomsGameFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6146b93c-cf83-4d9b-b5c9-2801a41e5e72")
			}
			record.RoomsGameFlag = v
		case "RoomsGreatFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3cad3935-1778-4253-949f-d26296a0c305")
			}
			record.RoomsGreatFlag = v
		case "RoomsHobbyFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1f7011c3-2f2a-42a6-bcdd-5da1f094bf84")
			}
			record.RoomsHobbyFlag = v
		case "RoomsLaundryFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "00389d7c-9c29-48cc-ab1b-42bf2cf21c24")
			}
			record.RoomsLaundryFlag = v
		case "RoomsMediaFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8f780838-dbc8-4e73-bf10-f0782c02be67")
			}
			record.RoomsMediaFlag = v
		case "RoomsMudFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e6f20a4e-9f3d-4148-9136-b33c3a13fd00")
			}
			record.RoomsMudFlag = v
		case "RoomsOfficeArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "913a9293-91e6-416e-9e4a-72c571ef370a")
			}
			record.RoomsOfficeArea = v
		case "RoomsOfficeFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ba3d4f12-a6ac-4b57-9335-30cfd2fce1f1")
			}
			record.RoomsOfficeFlag = v
		case "RoomsSafeRoomFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "324e5930-f533-40e6-8e0d-eb2de9f809a3")
			}
			record.RoomsSafeRoomFlag = v
		case "RoomsSittingFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "19e754a1-42bb-4fa3-b630-72e6c18419fa")
			}
			record.RoomsSittingFlag = v
		case "RoomsStormShelter":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b4f15cef-77a8-4bcd-9fdb-049cca3fa5b4")
			}
			record.RoomsStormShelter = v
		case "RoomsStudyFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a31aacb0-b7f1-43ac-96a5-ac102db63dbb")
			}
			record.RoomsStudyFlag = v
		case "RoomsSunroomFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5e3faa49-6041-40ed-a7d5-0c9d0cc3692f")
			}
			record.RoomsSunroomFlag = v
		case "RoomsUtilityArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5bd08183-5e1f-40dc-b620-90cb47b35cfa")
			}
			record.RoomsUtilityArea = v
		case "RoomsUtilityCode":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6c36aee2-68fd-4363-858f-4d652de34b09")
			}
			record.RoomsUtilityCode = v
		case "Fireplace":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c1109592-2f37-4feb-bfde-89d985abaae0")
			}
			record.Fireplace = v
		case "FireplaceCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fd460fe3-2515-41e8-96bf-8214dacd54c8")
			}
			record.FireplaceCount = v
		case "AccessabilityElevatorFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e34e9b4e-c909-4478-902c-63b45ac11893")
			}
			record.AccessabilityElevatorFlag = v
		case "AccessabilityHandicapFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "47d56ad6-5f85-4278-9e15-a75111251ceb")
			}
			record.AccessabilityHandicapFlag = v
		case "EscalatorFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3bf0f099-2344-481b-99c1-c787f9a6c994")
			}
			record.EscalatorFlag = v
		case "CentralVacuumFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c29c01f3-91e6-454d-854b-5b420ff3a43d")
			}
			record.CentralVacuumFlag = v
		case "ContentIntercomFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fd044d19-ed0f-48b6-a4d8-11f5287e0720")
			}
			record.ContentIntercomFlag = v
		case "ContentSoundSystemFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "459ede05-54df-4dea-bb69-d90872d73d0e")
			}
			record.ContentSoundSystemFlag = v
		case "WetBarFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3ea6c4b8-d6f2-4049-a809-9a375d2074e9")
			}
			record.WetBarFlag = v
		case "SecurityAlarmFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f24e78d9-0726-460d-9c2d-df75813c9f6f")
			}
			record.SecurityAlarmFlag = v
		case "StructureStyle":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bbaf2cf4-0ce4-406e-9122-d3b43b98531e")
			}
			record.StructureStyle = v
		case "Exterior1Code":
			record.Exterior1Code = val.StringPtrIfNonZero(field)
		case "RoofMaterial":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ee63c66b-63a7-4c57-930b-418c1fa78451")
			}
			record.RoofMaterial = v
		case "RoofConstruction":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f7157958-80e5-4a9d-8126-008114005156")
			}
			record.RoofConstruction = v
		case "ContentStormShutterFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "876da3a9-0547-462f-b12f-8f49e6a9b30e")
			}
			record.ContentStormShutterFlag = v
		case "ContentOverheadDoorFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1782e6a9-b39a-49a6-9a93-94ac3f099f70")
			}
			record.ContentOverheadDoorFlag = v
		case "ViewDescription":
			record.ViewDescription = val.StringPtrIfNonZero(field)
		case "PorchCode":
			record.PorchCode = val.StringPtrIfNonZero(field)
		case "PorchArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "751b8c9f-9535-4582-bac5-0858eccac517")
			}
			record.PorchArea = v
		case "PatioArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1b718bad-b63c-4217-9068-c608471b1aae")
			}
			record.PatioArea = v
		case "DeckFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "62fa0f50-0ed2-4073-9fdd-882435cb7091")
			}
			record.DeckFlag = v
		case "DeckArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ea3d8837-0b0d-465b-910e-044238c2d4c9")
			}
			record.DeckArea = v
		case "FeatureBalconyFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0d239edd-df63-4cbe-9e45-158e59124987")
			}
			record.FeatureBalconyFlag = v
		case "BalconyArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2a979eb9-61e6-4b63-ba4b-93993834f0d9")
			}
			record.BalconyArea = v
		case "BreezewayFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c5fa5226-2a0c-43c3-bbb6-1a5157c21613")
			}
			record.BreezewayFlag = v
		case "ParkingRVParkingFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "dbebaf5d-9d35-4b7b-b3c6-63e9e691cc04")
			}
			record.ParkingRVParkingFlag = v
		case "ParkingSpaceCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "956d689c-3917-4e26-8301-b0a8a89cebb2")
			}
			record.ParkingSpaceCount = v
		case "DrivewayArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ac3ac42b-09e8-4e6e-b6de-bb65fb931b41")
			}
			record.DrivewayArea = v
		case "DrivewayMaterial":
			record.DrivewayMaterial = val.StringPtrIfNonZero(field)
		case "Pool":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "19542ba0-b669-4a22-8904-132fc1030237")
			}
			record.Pool = v
		case "PoolArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "179ed3ca-c07c-4a34-bcd7-d41624cb73f3")
			}
			record.PoolArea = v
		case "ContentSaunaFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "04850b2a-2e3a-420e-bc77-fdbb398e6e48")
			}
			record.ContentSaunaFlag = v
		case "TopographyCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0d22be34-eb44-40ad-a4f6-0859dd9fbfdf")
			}
			record.TopographyCode = v
		case "FenceCode":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7b79c626-ab1a-4cba-96ff-a804053e2e63")
			}
			record.FenceCode = v
		case "FenceArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b0f47eba-d596-4f21-b4ee-96121e42ee1e")
			}
			record.FenceArea = v
		case "CourtyardFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "eea32df5-b3ef-4768-adc1-5b8fcca6025c")
			}
			record.CourtyardFlag = v
		case "CourtyardArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "cb09c052-d3ec-4e96-a17e-be108fa21301")
			}
			record.CourtyardArea = v
		case "ArborPergolaFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "151ea071-e65c-421a-9f86-f10fbc371715")
			}
			record.ArborPergolaFlag = v
		case "SprinklersFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "edf9e7ca-a388-4f21-ad82-53097bc19248")
			}
			record.SprinklersFlag = v
		case "GolfCourseGreenFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2167b9b4-e301-4ea8-84e9-eb90d9ef1e2b")
			}
			record.GolfCourseGreenFlag = v
		case "TennisCourtFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b328d31c-36c3-4a6e-beb7-14444e923e07")
			}
			record.TennisCourtFlag = v
		case "SportsCourtFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6b05c669-2266-45a9-82ef-be37c27fcf45")
			}
			record.SportsCourtFlag = v
		case "ArenaFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e3192328-ae35-49f6-8db5-161e488b9457")
			}
			record.ArenaFlag = v
		case "WaterFeatureFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ee0b260b-cfcc-4b56-925a-1095baa6f874")
			}
			record.WaterFeatureFlag = v
		case "PondFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "07d5e188-8727-4ad0-a578-03e79a57f812")
			}
			record.PondFlag = v
		case "BoatLiftFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "df7cacbf-367d-44d1-aee4-ae796d5a1e8c")
			}
			record.BoatLiftFlag = v
		case "BuildingsCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4b502d8d-27fc-4fff-a22f-82ac1f25f30b")
			}
			record.BuildingsCount = v
		case "BathHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c0d6bd0c-672d-4581-a0d3-350d862f48b2")
			}
			record.BathHouseArea = v
		case "BathHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "51f5a191-b0cc-435f-8c53-42620fe7199b")
			}
			record.BathHouseFlag = v
		case "BoatAccessFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5c44321f-a816-4b10-b7b7-70fa295f21aa")
			}
			record.BoatAccessFlag = v
		case "BoatHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "878b6ab8-4110-4a3f-b80d-f3d0f9afcf29")
			}
			record.BoatHouseArea = v
		case "BoatHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "adf709a2-de16-4baa-82e9-ed4e4d737810")
			}
			record.BoatHouseFlag = v
		case "CabinArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ef72c35b-eef0-41a0-895f-5f17609dd3bb")
			}
			record.CabinArea = v
		case "CabinFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a6e8e315-a540-43ea-b26d-ebae4c25cf4f")
			}
			record.CabinFlag = v
		case "CanopyArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "dfe6b0f7-12cb-4526-a911-471f61ce8da0")
			}
			record.CanopyArea = v
		case "CanopyFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "40894658-f7a1-411c-a02d-f5b31739eddf")
			}
			record.CanopyFlag = v
		case "GazeboArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "268bd178-6e61-465f-a1f1-b948c03dfd09")
			}
			record.GazeboArea = v
		case "GazeboFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "63170744-ad12-4b2e-8c16-776f00def5c0")
			}
			record.GazeboFlag = v
		case "GraineryArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "592e3b76-2888-4286-8f8e-938ff1075107")
			}
			record.GraineryArea = v
		case "GraineryFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "47a87396-8d01-43fe-bb28-c2a67c2f8e6e")
			}
			record.GraineryFlag = v
		case "GreenHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3547d05c-29b9-4812-ae35-66d80d6c7a2e")
			}
			record.GreenHouseArea = v
		case "GreenHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6c613e05-a10c-486e-ad42-248ae8a5ffaf")
			}
			record.GreenHouseFlag = v
		case "GuestHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "35670f41-468b-4d82-8adb-3001eb8ed331")
			}
			record.GuestHouseArea = v
		case "GuestHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e550179d-61e0-402e-a527-0c8632366e35")
			}
			record.GuestHouseFlag = v
		case "KennelArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "dc948b56-c730-4583-b432-ac7cc66c7e8a")
			}
			record.KennelArea = v
		case "KennelFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a415f111-2bdd-4ea1-bfef-63e2750ec3db")
			}
			record.KennelFlag = v
		case "LeanToArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bd869700-7ecc-471b-8b1b-07afcb33138f")
			}
			record.LeanToArea = v
		case "LeanToFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a28a553d-a117-4c31-86e6-9a5bfc1ffe1a")
			}
			record.LeanToFlag = v
		case "LoadingPlatformArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "40da5dd4-acd2-41d3-8895-013500eb2c44")
			}
			record.LoadingPlatformArea = v
		case "LoadingPlatformFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9fc22e31-5206-45d0-8bc2-c1c4816d4c1d")
			}
			record.LoadingPlatformFlag = v
		case "MilkHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bc4cccd3-f316-4fe7-888a-57011227e174")
			}
			record.MilkHouseArea = v
		case "MilkHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "49a81ff0-f42f-4f5e-9460-2acaefa599e0")
			}
			record.MilkHouseFlag = v
		case "OutdoorKitchenFireplaceFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4e238dbf-409a-4414-bc89-34e5a3876a1f")
			}
			record.OutdoorKitchenFireplaceFlag = v
		case "PoolHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6e1cd79a-378b-4423-b25b-94ba236da4bd")
			}
			record.PoolHouseArea = v
		case "PoolHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8aa592a1-6b82-4571-8f97-c37e98a7016d")
			}
			record.PoolHouseFlag = v
		case "PoultryHouseArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2aaab82d-8692-4dde-9695-9328b4bd9748")
			}
			record.PoultryHouseArea = v
		case "PoultryHouseFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0e9de34a-eaaf-44e6-9bc4-c6cef5978b74")
			}
			record.PoultryHouseFlag = v
		case "QuonsetArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d5847262-394f-4370-915f-3ad407ab07bc")
			}
			record.QuonsetArea = v
		case "QuonsetFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "085b67a9-e975-4356-87b9-d0069fb905d3")
			}
			record.QuonsetFlag = v
		case "ShedArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3f41eb14-68f4-4c61-aef6-68780d0257ce")
			}
			record.ShedArea = v
		case "ShedCode":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5631fcc6-50cf-40c9-b8cf-8200943abaa6")
			}
			record.ShedCode = v
		case "SiloArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3f7ec33f-6034-448b-a257-352d2524e133")
			}
			record.SiloArea = v
		case "SiloFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "415f6da6-d2d8-4492-84ea-b76408f09b7f")
			}
			record.SiloFlag = v
		case "StableArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7bbcca5c-73c5-48aa-8c14-89584802858c")
			}
			record.StableArea = v
		case "StableFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "fcf63526-1201-427f-ba5a-9f2702783589")
			}
			record.StableFlag = v
		case "StorageBuildingArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1526a8b0-9597-463a-b49b-9c9e26c51435")
			}
			record.StorageBuildingArea = v
		case "StorageBuildingFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9f5716cf-c2dd-4763-b20e-4fe6078b5e9b")
			}
			record.StorageBuildingFlag = v
		case "UtilityBuildingArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "87325a33-162f-4383-b98d-261a64200c93")
			}
			record.UtilityBuildingArea = v
		case "UtilityBuildingFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d3f50088-834e-4562-9493-3ce05331f353")
			}
			record.UtilityBuildingFlag = v
		case "PoleStructureArea":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0e698306-31a5-4c88-b3ab-cabd628c93d1")
			}
			record.PoleStructureArea = v
		case "PoleStructureFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "536f088b-d3a8-4acf-8684-d7093a04c73d")
			}
			record.PoleStructureFlag = v
		case "CommunityRecRoomFlag":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3d86d0f7-805c-4bcb-8b1f-8b6b01800998")
			}
			record.CommunityRecRoomFlag = v
		case "PublicationDate":
			v, err := val.TimePtrFromStringIfNonZero("2006-01-02", field)
			if err != nil {
				return nil, errors.Forward(err, "b44da502-1d92-4d03-9dee-c5b7ac7582d2")
			}
			record.PublicationDate = v
		case "ParcelShellRecord":
			// Special case for AttomData.
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1b7cff31-bfc2-46d3-a9fc-37a7e5048a8d")
			}
			record.ParcelShellRecord = v
		default:
			return nil, &errors.Object{
				Id:     "62a07d9f-1516-4b26-b96b-42e40627e74f",
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
		"attomid",
		"situs_state_code",
		"situs_county",
		"property_jurisdiction_name",
		"situs_state_county_fips",
		"combined_statistical_area",
		"cbsa_name",
		"cbsa_code",
		"msa_name",
		"msa_code",
		"metropolitan_division",
		"minor_civil_division_name",
		"minor_civil_division_code",
		"neighborhood_code",
		"census_fips_place_code",
		"census_tract",
		"census_block_group",
		"census_block",
		"parcel_number_raw",
		"parcel_number_formatted",
		"parcel_number_year_added",
		"parcel_number_alternate",
		"parcel_map_book",
		"parcel_map_page",
		"parcel_number_year_change",
		"parcel_number_previous",
		"parcel_account_number",
		"property_address_full",
		"property_address_house_number",
		"property_address_street_direction",
		"property_address_street_name",
		"property_address_street_suffix",
		"property_address_street_post_direction",
		"property_address_unit_prefix",
		"property_address_unit_value",
		"property_address_city",
		"property_address_state",
		"property_address_zip",
		"property_address_zip4",
		"property_address_crrt",
		"property_address_info_privacy",
		"congressional_district_house",
		"property_latitude",
		"property_longitude",
		"geo_quality",
		"legal_description",
		"legal_range",
		"legal_township",
		"legal_section",
		"legal_quarter",
		"legal_quarter_quarter",
		"legal_subdivision",
		"legal_phase",
		"legal_tract_number",
		"legal_block1",
		"legal_block2",
		"legal_lot_number1",
		"legal_lot_number2",
		"legal_lot_number3",
		"legal_unit",
		"party_owner1name_full",
		"party_owner1name_first",
		"party_owner1name_middle",
		"party_owner1name_last",
		"party_owner1name_suffix",
		"trust_description",
		"company_flag",
		"party_owner2name_full",
		"party_owner2name_first",
		"party_owner2name_middle",
		"party_owner2name_last",
		"party_owner2name_suffix",
		"owner_type_description1",
		"ownership_vesting_relation_code",
		"party_owner3name_full",
		"party_owner3name_first",
		"party_owner3name_middle",
		"party_owner3name_last",
		"party_owner3name_suffix",
		"party_owner4name_full",
		"party_owner4name_first",
		"party_owner4name_middle",
		"party_owner4name_last",
		"party_owner4name_suffix",
		"owner_type_description2",
		"contact_owner_mailing_county",
		"contact_owner_mailing_fips",
		"contact_owner_mail_address_full",
		"contact_owner_mail_address_house_number",
		"contact_owner_mail_address_street_direction",
		"contact_owner_mail_address_street_name",
		"contact_owner_mail_address_street_suffix",
		"contact_owner_mail_address_street_post_direction",
		"contact_owner_mail_address_unit_prefix",
		"contact_owner_mail_address_unit",
		"contact_owner_mail_address_city",
		"contact_owner_mail_address_state",
		"contact_owner_mail_address_zip",
		"contact_owner_mail_address_zip4",
		"contact_owner_mail_address_crrt",
		"contact_owner_mail_address_info_format",
		"contact_owner_mail_info_privacy",
		"status_owner_occupied_flag",
		"deed_owner1name_full",
		"deed_owner1name_first",
		"deed_owner1name_middle",
		"deed_owner1name_last",
		"deed_owner1name_suffix",
		"deed_owner2name_full",
		"deed_owner2name_first",
		"deed_owner2name_middle",
		"deed_owner2name_last",
		"deed_owner2name_suffix",
		"deed_owner3name_full",
		"deed_owner3name_first",
		"deed_owner3name_middle",
		"deed_owner3name_last",
		"deed_owner3name_suffix",
		"deed_owner4name_full",
		"deed_owner4name_first",
		"deed_owner4name_middle",
		"deed_owner4name_last",
		"deed_owner4name_suffix",
		"tax_year_assessed",
		"tax_assessed_value_total",
		"tax_assessed_value_improvements",
		"tax_assessed_value_land",
		"tax_assessed_improvements_perc",
		"previous_assessed_value",
		"tax_market_value_year",
		"tax_market_value_total",
		"tax_market_value_improvements",
		"tax_market_value_land",
		"tax_market_improvements_perc",
		"tax_fiscal_year",
		"tax_rate_area",
		"tax_billed_amount",
		"tax_delinquent_year",
		"last_assessor_tax_roll_update",
		"assr_last_updated",
		"tax_exemption_homeowner_flag",
		"tax_exemption_disabled_flag",
		"tax_exemption_senior_flag",
		"tax_exemption_veteran_flag",
		"tax_exemption_widow_flag",
		"tax_exemption_additional",
		"year_built",
		"year_built_effective",
		"zoned_code_local",
		"property_use_muni",
		"property_use_group",
		"property_use_standardized",
		"assessor_last_sale_date",
		"assessor_last_sale_amount",
		"assessor_prior_sale_date",
		"assessor_prior_sale_amount",
		"last_ownership_transfer_date",
		"last_ownership_transfer_document_number",
		"last_ownership_transfer_transaction_id",
		"deed_last_sale_document_book",
		"deed_last_sale_document_page",
		"deed_last_document_number",
		"deed_last_sale_date",
		"deed_last_sale_price",
		"deed_last_sale_transaction_id",
		"area_building",
		"area_building_definition_code",
		"area_gross",
		"area1st_floor",
		"area2nd_floor",
		"area_upper_floors",
		"area_lot_acres",
		"area_lot_sf",
		"area_lot_depth",
		"area_lot_width",
		"rooms_attic_area",
		"rooms_attic_flag",
		"rooms_basement_area",
		"rooms_basement_area_finished",
		"rooms_basement_area_unfinished",
		"parking_garage",
		"parking_garage_area",
		"parking_carport",
		"parking_carport_area",
		"hvac_cooling_detail",
		"hvac_heating_detail",
		"hvac_heating_fuel",
		"utilities_sewage_usage",
		"utilities_water_source",
		"utilities_mobile_home_hookup_flag",
		"foundation",
		"construction",
		"interior_structure",
		"plumbing_fixtures_count",
		"construction_fire_resistance_class",
		"safety_fire_sprinklers_flag",
		"flooring_material_primary",
		"bath_count",
		"bath_partial_count",
		"bedrooms_count",
		"rooms_count",
		"stories_count",
		"units_count",
		"rooms_bonus_room_flag",
		"rooms_breakfast_nook_flag",
		"rooms_cellar_flag",
		"rooms_cellar_wine_flag",
		"rooms_exercise_flag",
		"rooms_family_code",
		"rooms_game_flag",
		"rooms_great_flag",
		"rooms_hobby_flag",
		"rooms_laundry_flag",
		"rooms_media_flag",
		"rooms_mud_flag",
		"rooms_office_area",
		"rooms_office_flag",
		"rooms_safe_room_flag",
		"rooms_sitting_flag",
		"rooms_storm_shelter",
		"rooms_study_flag",
		"rooms_sunroom_flag",
		"rooms_utility_area",
		"rooms_utility_code",
		"fireplace",
		"fireplace_count",
		"accessability_elevator_flag",
		"accessability_handicap_flag",
		"escalator_flag",
		"central_vacuum_flag",
		"content_intercom_flag",
		"content_sound_system_flag",
		"wet_bar_flag",
		"security_alarm_flag",
		"structure_style",
		"exterior1code",
		"roof_material",
		"roof_construction",
		"content_storm_shutter_flag",
		"content_overhead_door_flag",
		"view_description",
		"porch_code",
		"porch_area",
		"patio_area",
		"deck_flag",
		"deck_area",
		"feature_balcony_flag",
		"balcony_area",
		"breezeway_flag",
		"parking_rv_parking_flag",
		"parking_space_count",
		"driveway_area",
		"driveway_material",
		"pool",
		"pool_area",
		"content_sauna_flag",
		"topography_code",
		"fence_code",
		"fence_area",
		"courtyard_flag",
		"courtyard_area",
		"arbor_pergola_flag",
		"sprinklers_flag",
		"golf_course_green_flag",
		"tennis_court_flag",
		"sports_court_flag",
		"arena_flag",
		"water_feature_flag",
		"pond_flag",
		"boat_lift_flag",
		"buildings_count",
		"bath_house_area",
		"bath_house_flag",
		"boat_access_flag",
		"boat_house_area",
		"boat_house_flag",
		"cabin_area",
		"cabin_flag",
		"canopy_area",
		"canopy_flag",
		"gazebo_area",
		"gazebo_flag",
		"grainery_area",
		"grainery_flag",
		"green_house_area",
		"green_house_flag",
		"guest_house_area",
		"guest_house_flag",
		"kennel_area",
		"kennel_flag",
		"lean_to_area",
		"lean_to_flag",
		"loading_platform_area",
		"loading_platform_flag",
		"milk_house_area",
		"milk_house_flag",
		"outdoor_kitchen_fireplace_flag",
		"pool_house_area",
		"pool_house_flag",
		"poultry_house_area",
		"poultry_house_flag",
		"quonset_area",
		"quonset_flag",
		"shed_area",
		"shed_code",
		"silo_area",
		"silo_flag",
		"stable_area",
		"stable_flag",
		"storage_building_area",
		"storage_building_flag",
		"utility_building_area",
		"utility_building_flag",
		"pole_structure_area",
		"pole_structure_flag",
		"community_rec_room_flag",
		"publication_date",
		"parcel_shell_record",
	}
}

func (dr *Assessor) SQLTable() string {
	return "ad_df_assessor"
}

func (dr *Assessor) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "cec7f313-05c2-420a-b94e-b5a9e00c41dc",
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
		dr.ATTOMID,
		dr.SitusStateCode,
		dr.SitusCounty,
		dr.PropertyJurisdictionName,
		dr.SitusStateCountyFIPS,
		dr.CombinedStatisticalArea,
		dr.CBSAName,
		dr.CBSACode,
		dr.MSAName,
		dr.MSACode,
		dr.MetropolitanDivision,
		dr.MinorCivilDivisionName,
		dr.MinorCivilDivisionCode,
		dr.NeighborhoodCode,
		dr.CensusFIPSPlaceCode,
		dr.CensusTract,
		dr.CensusBlockGroup,
		dr.CensusBlock,
		dr.ParcelNumberRaw,
		dr.ParcelNumberFormatted,
		dr.ParcelNumberYearAdded,
		dr.ParcelNumberAlternate,
		dr.ParcelMapBook,
		dr.ParcelMapPage,
		dr.ParcelNumberYearChange,
		dr.ParcelNumberPrevious,
		dr.ParcelAccountNumber,
		dr.PropertyAddressFull,
		dr.PropertyAddressHouseNumber,
		dr.PropertyAddressStreetDirection,
		dr.PropertyAddressStreetName,
		dr.PropertyAddressStreetSuffix,
		dr.PropertyAddressStreetPostDirection,
		dr.PropertyAddressUnitPrefix,
		dr.PropertyAddressUnitValue,
		dr.PropertyAddressCity,
		dr.PropertyAddressState,
		dr.PropertyAddressZIP,
		dr.PropertyAddressZIP4,
		dr.PropertyAddressCRRT,
		dr.PropertyAddressInfoPrivacy,
		dr.CongressionalDistrictHouse,
		dr.PropertyLatitude,
		dr.PropertyLongitude,
		dr.GeoQuality,
		dr.LegalDescription,
		dr.LegalRange,
		dr.LegalTownship,
		dr.LegalSection,
		dr.LegalQuarter,
		dr.LegalQuarterQuarter,
		dr.LegalSubdivision,
		dr.LegalPhase,
		dr.LegalTractNumber,
		dr.LegalBlock1,
		dr.LegalBlock2,
		dr.LegalLotNumber1,
		dr.LegalLotNumber2,
		dr.LegalLotNumber3,
		dr.LegalUnit,
		dr.PartyOwner1NameFull,
		dr.PartyOwner1NameFirst,
		dr.PartyOwner1NameMiddle,
		dr.PartyOwner1NameLast,
		dr.PartyOwner1NameSuffix,
		dr.TrustDescription,
		dr.CompanyFlag,
		dr.PartyOwner2NameFull,
		dr.PartyOwner2NameFirst,
		dr.PartyOwner2NameMiddle,
		dr.PartyOwner2NameLast,
		dr.PartyOwner2NameSuffix,
		dr.OwnerTypeDescription1,
		dr.OwnershipVestingRelationCode,
		dr.PartyOwner3NameFull,
		dr.PartyOwner3NameFirst,
		dr.PartyOwner3NameMiddle,
		dr.PartyOwner3NameLast,
		dr.PartyOwner3NameSuffix,
		dr.PartyOwner4NameFull,
		dr.PartyOwner4NameFirst,
		dr.PartyOwner4NameMiddle,
		dr.PartyOwner4NameLast,
		dr.PartyOwner4NameSuffix,
		dr.OwnerTypeDescription2,
		dr.ContactOwnerMailingCounty,
		dr.ContactOwnerMailingFIPS,
		dr.ContactOwnerMailAddressFull,
		dr.ContactOwnerMailAddressHouseNumber,
		dr.ContactOwnerMailAddressStreetDirection,
		dr.ContactOwnerMailAddressStreetName,
		dr.ContactOwnerMailAddressStreetSuffix,
		dr.ContactOwnerMailAddressStreetPostDirection,
		dr.ContactOwnerMailAddressUnitPrefix,
		dr.ContactOwnerMailAddressUnit,
		dr.ContactOwnerMailAddressCity,
		dr.ContactOwnerMailAddressState,
		dr.ContactOwnerMailAddressZIP,
		dr.ContactOwnerMailAddressZIP4,
		dr.ContactOwnerMailAddressCRRT,
		dr.ContactOwnerMailAddressInfoFormat,
		dr.ContactOwnerMailInfoPrivacy,
		dr.StatusOwnerOccupiedFlag,
		dr.DeedOwner1NameFull,
		dr.DeedOwner1NameFirst,
		dr.DeedOwner1NameMiddle,
		dr.DeedOwner1NameLast,
		dr.DeedOwner1NameSuffix,
		dr.DeedOwner2NameFull,
		dr.DeedOwner2NameFirst,
		dr.DeedOwner2NameMiddle,
		dr.DeedOwner2NameLast,
		dr.DeedOwner2NameSuffix,
		dr.DeedOwner3NameFull,
		dr.DeedOwner3NameFirst,
		dr.DeedOwner3NameMiddle,
		dr.DeedOwner3NameLast,
		dr.DeedOwner3NameSuffix,
		dr.DeedOwner4NameFull,
		dr.DeedOwner4NameFirst,
		dr.DeedOwner4NameMiddle,
		dr.DeedOwner4NameLast,
		dr.DeedOwner4NameSuffix,
		dr.TaxYearAssessed,
		dr.TaxAssessedValueTotal,
		dr.TaxAssessedValueImprovements,
		dr.TaxAssessedValueLand,
		dr.TaxAssessedImprovementsPerc,
		dr.PreviousAssessedValue,
		dr.TaxMarketValueYear,
		dr.TaxMarketValueTotal,
		dr.TaxMarketValueImprovements,
		dr.TaxMarketValueLand,
		dr.TaxMarketImprovementsPerc,
		dr.TaxFiscalYear,
		dr.TaxRateArea,
		dr.TaxBilledAmount,
		dr.TaxDelinquentYear,
		dr.LastAssessorTaxRollUpdate,
		dr.AssrLastUpdated,
		dr.TaxExemptionHomeownerFlag,
		dr.TaxExemptionDisabledFlag,
		dr.TaxExemptionSeniorFlag,
		dr.TaxExemptionVeteranFlag,
		dr.TaxExemptionWidowFlag,
		dr.TaxExemptionAdditional,
		dr.YearBuilt,
		dr.YearBuiltEffective,
		dr.ZonedCodeLocal,
		dr.PropertyUseMuni,
		dr.PropertyUseGroup,
		dr.PropertyUseStandardized,
		dr.AssessorLastSaleDate,
		dr.AssessorLastSaleAmount,
		dr.AssessorPriorSaleDate,
		dr.AssessorPriorSaleAmount,
		dr.LastOwnershipTransferDate,
		dr.LastOwnershipTransferDocumentNumber,
		dr.LastOwnershipTransferTransactionID,
		dr.DeedLastSaleDocumentBook,
		dr.DeedLastSaleDocumentPage,
		dr.DeedLastDocumentNumber,
		dr.DeedLastSaleDate,
		dr.DeedLastSalePrice,
		dr.DeedLastSaleTransactionID,
		dr.AreaBuilding,
		dr.AreaBuildingDefinitionCode,
		dr.AreaGross,
		dr.Area1stFloor,
		dr.Area2ndFloor,
		dr.AreaUpperFloors,
		dr.AreaLotAcres,
		dr.AreaLotSF,
		dr.AreaLotDepth,
		dr.AreaLotWidth,
		dr.RoomsAtticArea,
		dr.RoomsAtticFlag,
		dr.RoomsBasementArea,
		dr.RoomsBasementAreaFinished,
		dr.RoomsBasementAreaUnfinished,
		dr.ParkingGarage,
		dr.ParkingGarageArea,
		dr.ParkingCarport,
		dr.ParkingCarportArea,
		dr.HVACCoolingDetail,
		dr.HVACHeatingDetail,
		dr.HVACHeatingFuel,
		dr.UtilitiesSewageUsage,
		dr.UtilitiesWaterSource,
		dr.UtilitiesMobileHomeHookupFlag,
		dr.Foundation,
		dr.Construction,
		dr.InteriorStructure,
		dr.PlumbingFixturesCount,
		dr.ConstructionFireResistanceClass,
		dr.SafetyFireSprinklersFlag,
		dr.FlooringMaterialPrimary,
		dr.BathCount,
		dr.BathPartialCount,
		dr.BedroomsCount,
		dr.RoomsCount,
		dr.StoriesCount,
		dr.UnitsCount,
		dr.RoomsBonusRoomFlag,
		dr.RoomsBreakfastNookFlag,
		dr.RoomsCellarFlag,
		dr.RoomsCellarWineFlag,
		dr.RoomsExerciseFlag,
		dr.RoomsFamilyCode,
		dr.RoomsGameFlag,
		dr.RoomsGreatFlag,
		dr.RoomsHobbyFlag,
		dr.RoomsLaundryFlag,
		dr.RoomsMediaFlag,
		dr.RoomsMudFlag,
		dr.RoomsOfficeArea,
		dr.RoomsOfficeFlag,
		dr.RoomsSafeRoomFlag,
		dr.RoomsSittingFlag,
		dr.RoomsStormShelter,
		dr.RoomsStudyFlag,
		dr.RoomsSunroomFlag,
		dr.RoomsUtilityArea,
		dr.RoomsUtilityCode,
		dr.Fireplace,
		dr.FireplaceCount,
		dr.AccessabilityElevatorFlag,
		dr.AccessabilityHandicapFlag,
		dr.EscalatorFlag,
		dr.CentralVacuumFlag,
		dr.ContentIntercomFlag,
		dr.ContentSoundSystemFlag,
		dr.WetBarFlag,
		dr.SecurityAlarmFlag,
		dr.StructureStyle,
		dr.Exterior1Code,
		dr.RoofMaterial,
		dr.RoofConstruction,
		dr.ContentStormShutterFlag,
		dr.ContentOverheadDoorFlag,
		dr.ViewDescription,
		dr.PorchCode,
		dr.PorchArea,
		dr.PatioArea,
		dr.DeckFlag,
		dr.DeckArea,
		dr.FeatureBalconyFlag,
		dr.BalconyArea,
		dr.BreezewayFlag,
		dr.ParkingRVParkingFlag,
		dr.ParkingSpaceCount,
		dr.DrivewayArea,
		dr.DrivewayMaterial,
		dr.Pool,
		dr.PoolArea,
		dr.ContentSaunaFlag,
		dr.TopographyCode,
		dr.FenceCode,
		dr.FenceArea,
		dr.CourtyardFlag,
		dr.CourtyardArea,
		dr.ArborPergolaFlag,
		dr.SprinklersFlag,
		dr.GolfCourseGreenFlag,
		dr.TennisCourtFlag,
		dr.SportsCourtFlag,
		dr.ArenaFlag,
		dr.WaterFeatureFlag,
		dr.PondFlag,
		dr.BoatLiftFlag,
		dr.BuildingsCount,
		dr.BathHouseArea,
		dr.BathHouseFlag,
		dr.BoatAccessFlag,
		dr.BoatHouseArea,
		dr.BoatHouseFlag,
		dr.CabinArea,
		dr.CabinFlag,
		dr.CanopyArea,
		dr.CanopyFlag,
		dr.GazeboArea,
		dr.GazeboFlag,
		dr.GraineryArea,
		dr.GraineryFlag,
		dr.GreenHouseArea,
		dr.GreenHouseFlag,
		dr.GuestHouseArea,
		dr.GuestHouseFlag,
		dr.KennelArea,
		dr.KennelFlag,
		dr.LeanToArea,
		dr.LeanToFlag,
		dr.LoadingPlatformArea,
		dr.LoadingPlatformFlag,
		dr.MilkHouseArea,
		dr.MilkHouseFlag,
		dr.OutdoorKitchenFireplaceFlag,
		dr.PoolHouseArea,
		dr.PoolHouseFlag,
		dr.PoultryHouseArea,
		dr.PoultryHouseFlag,
		dr.QuonsetArea,
		dr.QuonsetFlag,
		dr.ShedArea,
		dr.ShedCode,
		dr.SiloArea,
		dr.SiloFlag,
		dr.StableArea,
		dr.StableFlag,
		dr.StorageBuildingArea,
		dr.StorageBuildingFlag,
		dr.UtilityBuildingArea,
		dr.UtilityBuildingFlag,
		dr.PoleStructureArea,
		dr.PoleStructureFlag,
		dr.CommunityRecRoomFlag,
		dr.PublicationDate,
		dr.ParcelShellRecord,
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
	ids := make([]int64, batchSize)

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "28c5299f-0cf7-4021-9775-b7afa73284d5")
	}

	loadRecords := func() error {
		// Skip if no records to process.
		if recordCount == 0 {
			return nil
		}

		builder = builder.
			PrefixExpr(squirrel.Expr(
				`
				with deleted_records as (
					delete from ad_df_assessor
					where ?
					returning *
				), archived_records as (
					insert into ad_assessor_history
					select
						*,
						now() as am_archived_at
					from deleted_records
				)
				`,
				squirrel.Eq{"attomid": ids[:recordCount]},
			)).Suffix(`returning (select count(*) from deleted_records)`)

		sql, args, err := builder.ToSql()
		if err != nil {
			return &errors.Object{
				Id:     "7003d14a-35bf-47ca-85ab-e0b0b8c848c3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to build SQL.",
				Cause:  err.Error(),
			}
		}

		tx, err := pgxPool.Begin(r.Context())
		if err != nil {
			return &errors.Object{
				Id:     "b2b5f524-219e-4eb3-bb21-68be9922e8ee",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to begin transaction.",
				Cause:  err.Error(),
			}
		}

		defer extutils.RollbackPgxTx(r.Context(), tx, "d42a57c0-7a37-4287-92fa-96310fd0fb93")

		// This clone won't replace the original arc.
		r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

		row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
		if err != nil {
			return errors.Forward(err, "96044a70-e0ad-4c0f-8ed4-759dc21669b2")
		}

		var batchDeletes int64

		if err := row.Scan(&batchDeletes); err != nil {
			return &errors.Object{
				Id:     "c0f1d8b2-3a4e-4f5c-9b6d-7e0f1c2b3d4e",
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
			return errors.Forward(err, "0b7fdc1c-1308-407c-bf3f-4076d4fd39d8")
		}

		if err := tx.Commit(r.Context()); err != nil {
			return &errors.Object{
				Id:     "e5504d88-f065-4221-89b1-f2bfdf41653a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to commit transaction.",
				Cause:  err.Error(),
			}
		}

		builder = newBuilder()
		dfObject = updateObjectOut.Entity
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
			return nil, errors.Forward(err, "42591b19-06d9-4bd1-a3b7-4db5b33e18ec")
		}

		recordValues, err := record.SQLValues()
		if err != nil {
			return nil, errors.Forward(err, "ea1836a7-c3fa-49c0-801d-1aa9f19b5383")
		}

		builder = builder.Values(recordValues...)
		ids[recordCount] = record.(*Assessor).ATTOMID

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
