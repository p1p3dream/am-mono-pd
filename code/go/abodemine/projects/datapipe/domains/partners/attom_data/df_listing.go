package attom_data

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type Listing struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any

	ATTOMID                            int64
	MLSRecordID                        int64
	MLSListingID                       int64
	StatusChangeDate                   time.Time
	PropertyAddressFull                *string
	PropertyAddressHouseNumber         *string
	PropertyAddressStreetDirection     *string
	PropertyAddressStreetName          *string
	PropertyAddressStreetSuffix        *string
	PropertyAddressStreetPostDirection *string

	PropertyAddressUnitPrefix *string
	PropertyAddressUnitValue  *string
	PropertyAddressCity       *string
	PropertyAddressState      *string
	PropertyAddressZIP        *string
	PropertyAddressZIP4       *string
	SitusCounty               *string
	Township                  *string
	MLSListingAddress         *string
	MLSListingCity            *string

	MLSListingState        *string
	MLSListingZip          *string
	MLSListingCountyFIPS   *string
	MLSNumber              *string
	MLSSource              *string
	ListingStatus          *string
	MLSSoldDate            *time.Time
	MLSSoldPrice           *int
	AssessorLastSaleDate   *time.Time
	AssessorLastSaleAmount *string

	MarketValue              *int
	MarketValueDate          *time.Time
	AvgMarketPricePerSqFt    *int
	ListingDate              *time.Time
	LatestListingPrice       *int
	PreviousListingPrice     *int
	LatestPriceChangeDate    *time.Time
	PendingDate              *time.Time
	SpecialListingConditions *string
	OriginalListingDate      *time.Time

	OriginalListingPrice   *int
	LeaseOption            *string
	LeaseTerm              *string
	LeaseIncludes          *string
	Concessions            *string
	ConcessionsAmount      *int
	ConcessionsComments    *string
	ContingencyDate        *time.Time
	ContingencyDescription *string
	MLSPropertyType        *string

	MLSPropertySubType   *string
	ATTOMPropertyType    *string
	ATTOMPropertySubType *string
	OwnershipDescription *string
	Latitude             *float64
	Longitude            *float64
	APNFormatted         *string
	LegalDescription     *string
	LegalSubdivision     *string
	DaysOnMarket         *int

	CumulativeDaysOnMarket     *int
	ListingAgentFullName       *string
	ListingAgentMLSID          *string
	ListingAgentStateLicense   *string
	ListingAgentAOR            *string
	ListingAgentPreferredPhone *string
	ListingAgentEmail          *string
	ListingOfficeName          *string
	ListingOfficeMlsId         *string
	ListingOfficeAOR           *string

	ListingOfficePhone           *string
	ListingOfficeEmail           *string
	ListingCoAgentFullName       *string
	ListingCoAgentMLSID          *string
	ListingCoAgentStateLicense   *string
	ListingCoAgentAOR            *string
	ListingCoAgentPreferredPhone *string
	ListingCoAgentEmail          *string
	ListingCoAgentOfficeName     *string
	ListingCoAgentOfficeMlsId    *string

	ListingCoAgentOfficeAOR   *string
	ListingCoAgentOfficePhone *string
	ListingCoAgentOfficeEmail *string
	BuyerAgentFullName        *string
	BuyerAgentMLSID           *string
	BuyerAgentStateLicense    *string
	BuyerAgentAOR             *string
	BuyerAgentPreferredPhone  *string
	BuyerAgentEmail           *string
	BuyerOfficeName           *string

	BuyerOfficeMlsId           *string
	BuyerOfficeAOR             *string
	BuyerOfficePhone           *string
	BuyerOfficeEmail           *string
	BuyerCoAgentFullName       *string
	BuyerCoAgentMLSID          *string
	BuyerCoAgentStateLicense   *string
	BuyerCoAgentAOR            *string
	BuyerCoAgentPreferredPhone *string
	BuyerCoAgentEmail          *string

	BuyerCoAgentOfficeName  *string
	BuyerCoAgentOfficeMlsId *string
	BuyerCoAgentOfficeAOR   *string
	BuyerCoAgentOfficePhone *string
	BuyerCoAgentOfficeEmail *string
	PublicListingRemarks    *string
	HomeWarrantyYN          *bool
	TaxYearAssessed         *int
	TaxAssessedValueTotal   *int
	TaxAmount               *int

	TaxAnnualOther      *int
	OwnerName           *string
	OwnerVesting        *string
	YearBuilt           *int
	YearBuiltEffective  *int
	YearBuiltSource     *string
	NewConstructionYN   *bool
	BuilderName         *string
	AdditionalParcelsYN *bool
	NumberOfLots        *int

	LotSizeSquareFeet    *float64
	LotSizeAcres         *float64
	LotSizeSource        *string
	LotDimensions        *string
	LotFeatureList       *string
	FrontageLength       *string
	FrontageType         *string
	FrontageRoadType     *string
	LivingAreaSquareFeet *int
	LivingAreaSource     *string

	Levels               *string
	Stories              *decimal.Decimal
	BuildingStoriesTotal *decimal.Decimal
	BuildingKeywords     *string
	BuildingAreaTotal    *int
	NumberOfUnitsTotal   *int
	NumberOfBuildings    *int
	PropertyAttachedYN   *bool
	OtherStructures      *string
	RoomsTotal           *int

	BedroomsTotal          *int
	BathroomsFull          *decimal.Decimal
	BathroomsHalf          *int
	BathroomsQuarter       *int
	BathroomsThreeQuarters *int
	BasementFeatures       *string
	BelowGradeSquareFeet   *int
	BasementTotalSqFt      *int
	BasementFinishedSqFt   *int
	BasementUnfinishedSqFt *int

	PropertyCondition     *string
	RepairsYN             *bool
	RepairsDescription    *string
	Disclosures           *string
	ConstructionMaterials *string
	GarageYN              *bool
	AttachedGarageYN      *bool
	GarageSpaces          *decimal.Decimal
	CarportYN             *bool
	CarportSpaces         *float64

	ParkingFeatures   *string
	ParkingOther      *string
	OpenParkingSpaces *float64
	ParkingTotal      *float64
	PoolPrivateYN     *bool
	PoolFeatures      *string
	Occupancy         *string
	ViewYN            *bool
	View              *string
	Topography        *string

	HeatingYN                  *bool
	HeatingFeatures            *string
	CoolingYN                  *bool
	Cooling                    *string
	FireplaceYN                *bool
	Fireplace                  *string
	FireplaceNumber            *float64
	FoundationFeatures         *string
	Roof                       *string
	ArchitecturalStyleFeatures *string

	PatioAndPorchFeatures  *string
	Utilities              *string
	ElectricIncluded       *bool
	ElectricDescription    *string
	WaterIncluded          *bool
	WaterSource            *string
	Sewer                  *string
	GasDescription         *string
	OtherEquipmentIncluded *string
	LaundryFeatures        *string

	Appliances         *string
	InteriorFeatures   *string
	ExteriorFeatures   *string
	FencingFeatures    *string
	PetsAllowed        *string
	HorseZoningYN      *bool
	SeniorCommunityYN  *bool
	WaterbodyName      *string
	WaterfrontYN       *bool
	WaterfrontFeatures *string

	ZoningCode               *string
	ZoningDescription        *string
	CurrentUse               *string
	PossibleUse              *string
	AssociationYN            *bool
	Association1Name         *string
	Association1Phone        *string
	Association1Fee          *int
	Association1FeeFrequency *string
	Association2Name         *string

	Association2Phone        *string
	Association2Fee          *int
	Association2FeeFrequency *string
	AssociationFeeIncludes   *string
	AssociationAmenities     *string
	SchoolElementary         *string
	SchoolElementaryDistrict *string
	SchoolMiddle             *string
	SchoolMiddleDistrict     *string
	SchoolHigh               *string

	SchoolHighDistrict            *string
	GreenVerificationYN           *bool
	GreenBuildingVerificationType *string
	GreenEnergyEfficient          *string
	GreenEnergyGeneration         *string
	GreenIndoorAirQuality         *string
	GreenLocation                 *string
	GreenSustainability           *string
	GreenWaterConservation        *string
	LandLeaseYN                   *bool

	LandLeaseAmount          *decimal.Decimal
	LandLeaseAmountFrequency *string
	LandLeaseExpirationDate  *time.Time
	CapRate                  *decimal.Decimal
	GrossIncome              *decimal.Decimal
	IncomeIncludes           *string
	GrossScheduledIncome     *decimal.Decimal
	NetOperatingIncome       *decimal.Decimal
	TotalActualRent          *decimal.Decimal
	ExistingLeaseType        *string

	FinancialDataSource  *string
	RentControlYN        *bool
	UnitTypeDescription  *string
	UnitTypeFurnished    *string
	NumberOfUnitsLeased  *float64
	NumberOfUnitsMoMo    *float64
	NumberOfUnitsVacant  *float64
	VacancyAllowance     *float64
	VacancyAllowanceRate *float64
	OperatingExpense     *decimal.Decimal

	CableTvExpense              *decimal.Decimal
	ElectricExpense             *decimal.Decimal
	FuelExpense                 *decimal.Decimal
	FurnitureReplacementExpense *decimal.Decimal
	GardenerExpense             *decimal.Decimal
	InsuranceExpense            *decimal.Decimal
	OperatingExpenseIncludes    *string
	LicensesExpense             *decimal.Decimal
	MaintenanceExpense          *decimal.Decimal
	ManagerExpense              *decimal.Decimal

	NewTaxesExpense               *decimal.Decimal
	OtherExpense                  *decimal.Decimal
	PestControlExpense            *decimal.Decimal
	PoolExpense                   *decimal.Decimal
	ProfessionalManagementExpense *decimal.Decimal
	SuppliesExpense               *decimal.Decimal
	TrashExpense                  *decimal.Decimal
	WaterSewerExpense             *decimal.Decimal
	WorkmansCompensationExpense   *decimal.Decimal
	OwnerPays                     *string

	TenantPays          *string
	ListingMarketingURL *string
	PhotosCount         *int
	PhotoKey            *string
	PhotoURLPrefix      *string

	// Added in DataFileTypeListingV20250417.
	CurrentStatus *bool

	fileType entities.DataFileType
}

func (dr *Listing) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := &Listing{
		fileType: dr.fileType,
	}

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "ATTOM ID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "b5440ef0-0458-4f75-9069-3b3c53505042",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": v,
					},
				}
			}
			record.ATTOMID = v
		case "MLSRecordID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "2edee4c8-0423-4f8c-a93e-74e9805db9f3",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": v,
					},
				}
			}
			record.MLSRecordID = v
		case "MLSListingID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "fb61c5d4-e618-4860-baae-4bcacab2296e",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": v,
					},
				}
			}
			record.MLSListingID = v
		case "StatusChangeDate":
			v, err := time.Parse(consts.USSlashDate, field)
			if err != nil {
				return nil, &errors.Object{
					Id:    "1e44df38-b834-4d7a-885c-e0bfdf19d460",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": v,
					},
				}
			}
			record.StatusChangeDate = v
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
		case "SitusCounty":
			record.SitusCounty = val.StringPtrIfNonZero(field)
		case "Township":
			record.Township = val.StringPtrIfNonZero(field)
		case "MLSListingAddress":
			record.MLSListingAddress = val.StringPtrIfNonZero(field)
		case "MLSListingCity":
			record.MLSListingCity = val.StringPtrIfNonZero(field)
		case "MLSListingState":
			record.MLSListingState = val.StringPtrIfNonZero(field)
		case "MLSListingZip":
			record.MLSListingZip = val.StringPtrIfNonZero(field)
		case "MLSListingCountyFIPS":
			record.MLSListingCountyFIPS = val.StringPtrIfNonZero(field)
		case "MLSNumber":
			record.MLSNumber = val.StringPtrIfNonZero(field)
		case "MLSSource":
			record.MLSSource = val.StringPtrIfNonZero(field)
		case "ListingStatus":
			record.ListingStatus = val.StringPtrIfNonZero(field)
		case "MLSSoldDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "d2cf716b-d444-4758-9ebf-85a4e8f23dcb")
			}
			record.MLSSoldDate = v
		case "MLSSoldPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "aff904a0-5554-432a-abb1-fb2374113318")
			}
			record.MLSSoldPrice = v
		case "AssessorLastSaleDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "78e58fa0-ce95-44c3-8997-7bf3659f7d46")
			}
			record.AssessorLastSaleDate = v
		case "AssessorLastSaleAmount":
			record.AssessorLastSaleAmount = val.StringPtrIfNonZero(field)
		case "MarketValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2ac4f610-907f-49c9-8f4b-e0d0347171f0")
			}
			record.MarketValue = v
		case "MarketValueDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "12ab7c2b-a41b-47b3-80aa-f92b7875b766")
			}
			record.MarketValueDate = v
		case "AvgMarketPricePerSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9deb5a99-ff19-4749-869d-d0b5f0ea3f9a")
			}
			record.AvgMarketPricePerSqFt = v
		case "ListingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "7c49e43b-d548-44a1-adce-0a6b01184318")
			}
			record.ListingDate = v
		case "LatestListingPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f63654df-b4c2-4576-a67b-9e2f019167a9")
			}
			record.LatestListingPrice = v
		case "PreviousListingPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "daa8b210-bc53-45fb-bb32-f549ed375503")
			}
			record.PreviousListingPrice = v
		case "LatestPriceChangeDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "1bf8971d-d68b-4e16-ac58-63b57947095a")
			}
			record.LatestPriceChangeDate = v
		case "PendingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "d4978a49-3282-4154-988d-a3a89b327f21")
			}
			record.PendingDate = v
		case "SpecialListingConditions":
			record.SpecialListingConditions = val.StringPtrIfNonZero(field)
		case "OriginalListingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "bbc7acac-bbff-4235-8cc8-4e5eff35075c")
			}
			record.OriginalListingDate = v
		case "OriginalListingPrice":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a4989da0-95ec-481c-afac-d50d2c06828a")
			}
			record.OriginalListingPrice = v
		case "LeaseOption":
			record.LeaseOption = val.StringPtrIfNonZero(field)
		case "LeaseTerm":
			record.LeaseTerm = val.StringPtrIfNonZero(field)
		case "LeaseIncludes":
			record.LeaseIncludes = val.StringPtrIfNonZero(field)
		case "Concessions":
			record.Concessions = val.StringPtrIfNonZero(field)
		case "ConcessionsAmount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6c86814a-4382-4517-8c47-51ee7813f819")
			}
			record.ConcessionsAmount = v
		case "ConcessionsComments":
			record.ConcessionsComments = val.StringPtrIfNonZero(field)
		case "ContingencyDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "7a90bf9a-a056-4255-a7b9-dc270a56f622")
			}
			record.ContingencyDate = v
		case "ContingencyDescription":
			record.ContingencyDescription = val.StringPtrIfNonZero(field)
		case "MLSPropertyType":
			record.MLSPropertyType = val.StringPtrIfNonZero(field)
		case "MLSPropertySubType":
			record.MLSPropertySubType = val.StringPtrIfNonZero(field)
		case "ATTOMPropertyType":
			record.ATTOMPropertyType = val.StringPtrIfNonZero(field)
		case "ATTOMPropertySubType":
			record.ATTOMPropertySubType = val.StringPtrIfNonZero(field)
		case "OwnershipDescription":
			record.OwnershipDescription = val.StringPtrIfNonZero(field)
		case "Latitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "313383d3-9186-4f28-97fc-74735ceb9fb7")
			}
			record.Latitude = v
		case "Longitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b093c406-07fe-453c-b299-5afb12e4745f")
			}
			record.Longitude = v
		case "APNFormatted":
			record.APNFormatted = val.StringPtrIfNonZero(field)
		case "LegalDescription":
			record.LegalDescription = val.StringPtrIfNonZero(field)
		case "LegalSubdivision":
			record.LegalSubdivision = val.StringPtrIfNonZero(field)
		case "DaysOnMarket":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3c098528-dfce-4c55-a626-92ba2bf8e697")
			}
			record.DaysOnMarket = v
		case "CumulativeDaysOnMarket":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "226e2e06-2e73-4b68-9c98-62c9e18c5630")
			}
			record.CumulativeDaysOnMarket = v
		case "ListingAgentFullName":
			record.ListingAgentFullName = val.StringPtrIfNonZero(field)
		case "ListingAgentMLSID":
			record.ListingAgentMLSID = val.StringPtrIfNonZero(field)
		case "ListingAgentStateLicense":
			record.ListingAgentStateLicense = val.StringPtrIfNonZero(field)
		case "ListingAgentAOR":
			record.ListingAgentAOR = val.StringPtrIfNonZero(field)
		case "ListingAgentPreferredPhone":
			record.ListingAgentPreferredPhone = val.StringPtrIfNonZero(field)
		case "ListingAgentEmail":
			record.ListingAgentEmail = val.StringPtrIfNonZero(field)
		case "ListingOfficeName":
			record.ListingOfficeName = val.StringPtrIfNonZero(field)
		case "ListingOfficeMlsId":
			record.ListingOfficeMlsId = val.StringPtrIfNonZero(field)
		case "ListingOfficeAOR":
			record.ListingOfficeAOR = val.StringPtrIfNonZero(field)
		case "ListingOfficePhone":
			record.ListingOfficePhone = val.StringPtrIfNonZero(field)
		case "ListingOfficeEmail":
			record.ListingOfficeEmail = val.StringPtrIfNonZero(field)
		case "ListingCoAgentFullName":
			record.ListingCoAgentFullName = val.StringPtrIfNonZero(field)
		case "ListingCoAgentMLSID":
			record.ListingCoAgentMLSID = val.StringPtrIfNonZero(field)
		case "ListingCoAgentStateLicense":
			record.ListingCoAgentStateLicense = val.StringPtrIfNonZero(field)
		case "ListingCoAgentAOR":
			record.ListingCoAgentAOR = val.StringPtrIfNonZero(field)
		case "ListingCoAgentPreferredPhone":
			record.ListingCoAgentPreferredPhone = val.StringPtrIfNonZero(field)
		case "ListingCoAgentEmail":
			record.ListingCoAgentEmail = val.StringPtrIfNonZero(field)
		case "ListingCoAgentOfficeName":
			record.ListingCoAgentOfficeName = val.StringPtrIfNonZero(field)
		case "ListingCoAgentOfficeMlsId":
			record.ListingCoAgentOfficeMlsId = val.StringPtrIfNonZero(field)
		case "ListingCoAgentOfficeAOR":
			record.ListingCoAgentOfficeAOR = val.StringPtrIfNonZero(field)
		case "ListingCoAgentOfficePhone":
			record.ListingCoAgentOfficePhone = val.StringPtrIfNonZero(field)
		case "ListingCoAgentOfficeEmail":
			record.ListingCoAgentOfficeEmail = val.StringPtrIfNonZero(field)
		case "BuyerAgentFullName":
			record.BuyerAgentFullName = val.StringPtrIfNonZero(field)
		case "BuyerAgentMLSID":
			record.BuyerAgentMLSID = val.StringPtrIfNonZero(field)
		case "BuyerAgentStateLicense":
			record.BuyerAgentStateLicense = val.StringPtrIfNonZero(field)
		case "BuyerAgentAOR":
			record.BuyerAgentAOR = val.StringPtrIfNonZero(field)
		case "BuyerAgentPreferredPhone":
			record.BuyerAgentPreferredPhone = val.StringPtrIfNonZero(field)
		case "BuyerAgentEmail":
			record.BuyerAgentEmail = val.StringPtrIfNonZero(field)
		case "BuyerOfficeName":
			record.BuyerOfficeName = val.StringPtrIfNonZero(field)
		case "BuyerOfficeMlsId":
			record.BuyerOfficeMlsId = val.StringPtrIfNonZero(field)
		case "BuyerOfficeAOR":
			record.BuyerOfficeAOR = val.StringPtrIfNonZero(field)
		case "BuyerOfficePhone":
			record.BuyerOfficePhone = val.StringPtrIfNonZero(field)
		case "BuyerOfficeEmail":
			record.BuyerOfficeEmail = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentFullName":
			record.BuyerCoAgentFullName = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentMLSID":
			record.BuyerCoAgentMLSID = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentStateLicense":
			record.BuyerCoAgentStateLicense = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentAOR":
			record.BuyerCoAgentAOR = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentPreferredPhone":
			record.BuyerCoAgentPreferredPhone = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentEmail":
			record.BuyerCoAgentEmail = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentOfficeName":
			record.BuyerCoAgentOfficeName = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentOfficeMlsId":
			record.BuyerCoAgentOfficeMlsId = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentOfficeAOR":
			record.BuyerCoAgentOfficeAOR = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentOfficePhone":
			record.BuyerCoAgentOfficePhone = val.StringPtrIfNonZero(field)
		case "BuyerCoAgentOfficeEmail":
			record.BuyerCoAgentOfficeEmail = val.StringPtrIfNonZero(field)
		case "PublicListingRemarks":
			record.PublicListingRemarks = val.StringPtrIfNonZero(field)
		case "HomeWarrantyYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "e53e6008-6e81-4a4a-b6a6-1629ad798a60")
			}
			record.HomeWarrantyYN = v
		case "TaxYearAssessed":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f6b2d7e6-ffb8-4be0-933d-c57a76efd048")
			}
			record.TaxYearAssessed = v
		case "TaxAssessedValueTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8911fe77-a699-483e-b73a-9c37c1d21b3f")
			}
			record.TaxAssessedValueTotal = v
		case "TaxAmount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9d1f6765-8a97-4178-96ce-96705e1fdc77")
			}
			record.TaxAmount = v
		case "TaxAnnualOther":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a664b73f-5bb4-4d8f-922e-110670acae23")
			}
			record.TaxAnnualOther = v
		case "OwnerName":
			record.OwnerName = val.StringPtrIfNonZero(field)
		case "OwnerVesting":
			record.OwnerVesting = val.StringPtrIfNonZero(field)
		case "YearBuilt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a457d7e4-99db-4540-8ebd-cbf9f9e4322b")
			}
			record.YearBuilt = v
		case "YearBuiltEffective":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "efe62b7d-43c8-4cd9-874d-5b79551275b6")
			}
			record.YearBuiltEffective = v
		case "YearBuiltSource":
			record.YearBuiltSource = val.StringPtrIfNonZero(field)
		case "NewConstructionYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "544c0a8b-10a4-4450-b344-f23eda9f3d73")
			}
			record.NewConstructionYN = v
		case "BuilderName":
			record.BuilderName = val.StringPtrIfNonZero(field)
		case "AdditionalParcelsYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "8deb0f68-1dd6-4a09-8c51-79116d5dad74")
			}
			record.AdditionalParcelsYN = v
		case "NumberOfLots":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6f939dd2-0ffe-4a78-b2ef-346c7e233a36")
			}
			record.NumberOfLots = v
		case "LotSizeSquareFeet":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7c6955de-829f-47fa-9d18-36178ddc863f")
			}
			record.LotSizeSquareFeet = v
		case "LotSizeAcres":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "30fffd4d-6d15-409d-ba65-9d4e233e06cd")
			}
			record.LotSizeAcres = v
		case "LotSizeSource":
			record.LotSizeSource = val.StringPtrIfNonZero(field)
		case "LotDimensions":
			record.LotDimensions = val.StringPtrIfNonZero(field)
		case "LotFeatureList":
			record.LotFeatureList = val.StringPtrIfNonZero(field)
		case "FrontageLength":
			record.FrontageLength = val.StringPtrIfNonZero(field)
		case "FrontageType":
			record.FrontageType = val.StringPtrIfNonZero(field)
		case "FrontageRoadType":
			record.FrontageRoadType = val.StringPtrIfNonZero(field)
		case "LivingAreaSquareFeet":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "80179cfe-a01a-4faa-a61e-cc102630e227")
			}
			record.LivingAreaSquareFeet = v
		case "LivingAreaSource":
			record.LivingAreaSource = val.StringPtrIfNonZero(field)
		case "Levels":
			record.Levels = val.StringPtrIfNonZero(field)
		case "Stories":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "16e21ebb-2fe7-4045-bfdf-007cbe49fa47")
			}
			record.Stories = v
		case "BuildingStoriesTotal":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "061e27fe-67c3-4e8c-b0dd-ac84849ac1b0")
			}
			record.BuildingStoriesTotal = v
		case "BuildingKeywords":
			record.BuildingKeywords = val.StringPtrIfNonZero(field)
		case "BuildingAreaTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9d8d19bd-61bb-48ed-a77d-ad7ca0461e3b")
			}
			record.BuildingAreaTotal = v
		case "NumberOfUnitsTotal":
			// Although the type is int they use float in the text files.
			v, err := val.IntPtrFromFloat64StringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "82480f1a-5900-4bc4-978b-10e6ad5af2ca")
			}
			record.NumberOfUnitsTotal = v
		case "NumberOfBuildings":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "df590b1f-c61a-4d0e-b2cb-4b7ec575a536")
			}
			record.NumberOfBuildings = v
		case "PropertyAttachedYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "a99309b4-43cc-4631-a136-d48b2d87be71")
			}
			record.PropertyAttachedYN = v
		case "OtherStructures":
			record.OtherStructures = val.StringPtrIfNonZero(field)
		case "RoomsTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3528f20b-f685-4f46-9288-183d4824d1d1")
			}
			record.RoomsTotal = v
		case "BedroomsTotal":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0e45a626-0f28-46ea-9015-313a7069588c")
			}
			record.BedroomsTotal = v
		case "BathroomsFull":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "418ec435-ffac-463f-a76b-dbe111de9a5e")
			}
			record.BathroomsFull = v
		case "BathroomsHalf":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "89ce2424-29ab-427b-84bd-e280c15e87ff")
			}
			record.BathroomsHalf = v
		case "BathroomsQuarter":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bc689edb-1153-4ad3-a47a-37b9f572097f")
			}
			record.BathroomsQuarter = v
		case "BathroomsThreeQuarters":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "24ffed12-4142-44ef-af21-d04620ea28b4")
			}
			record.BathroomsThreeQuarters = v
		case "BasementFeatures":
			record.BasementFeatures = val.StringPtrIfNonZero(field)
		case "BelowGradeSquareFeet":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f58a77ee-817d-4432-a901-81da7664778c")
			}
			record.BelowGradeSquareFeet = v
		case "BasementTotalSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d4911f92-6526-476c-bc28-438450466992")
			}
			record.BasementTotalSqFt = v
		case "BasementFinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8337734c-8a10-4ce8-bb25-cc954dca0a5b")
			}
			record.BasementFinishedSqFt = v
		case "BasementUnfinishedSqFt":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "52a830a1-7ed3-488d-bffe-46ca07b20ada")
			}
			record.BasementUnfinishedSqFt = v
		case "PropertyCondition":
			record.PropertyCondition = val.StringPtrIfNonZero(field)
		case "RepairsYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "8f7e8e97-ebde-465a-ab3c-2cb925bd01f6")
			}
			record.RepairsYN = v
		case "RepairsDescription":
			record.RepairsDescription = val.StringPtrIfNonZero(field)
		case "Disclosures":
			record.Disclosures = val.StringPtrIfNonZero(field)
		case "ConstructionMaterials":
			record.ConstructionMaterials = val.StringPtrIfNonZero(field)
		case "GarageYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "690fd3c0-850e-41a8-9bf7-0901d66b0c00")
			}
			record.GarageYN = v
		case "AttachedGarageYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "354e7ac2-25db-4a5f-9537-c9acba29e225")
			}
			record.AttachedGarageYN = v
		case "GarageSpaces":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "50049289-badc-4b29-a704-0d78dba3cffe")
			}
			record.GarageSpaces = v
		case "CarportYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "c78e5937-cc7b-4a37-b72d-c06507738774")
			}
			record.CarportYN = v
		case "CarportSpaces":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e4763a19-6e3a-4aa0-9bc8-ac9b3c756a97")
			}
			record.CarportSpaces = v
		case "ParkingFeatures":
			record.ParkingFeatures = val.StringPtrIfNonZero(field)
		case "ParkingOther":
			record.ParkingOther = val.StringPtrIfNonZero(field)
		case "OpenParkingSpaces":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2fe534ea-3255-49f8-aad0-f3aaa706be25")
			}
			record.OpenParkingSpaces = v
		case "ParkingTotal":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f14d51ed-119e-4672-9cf5-748415506971")
			}
			record.ParkingTotal = v
		case "PoolPrivateYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "f843e6c4-a155-491b-a35d-6302dc565a29")
			}
			record.PoolPrivateYN = v
		case "PoolFeatures":
			record.PoolFeatures = val.StringPtrIfNonZero(field)
		case "Occupancy":
			record.Occupancy = val.StringPtrIfNonZero(field)
		case "ViewYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "1b43667c-a961-4aae-ad6f-d484a8cc399e")
			}
			record.ViewYN = v
		case "View":
			record.View = val.StringPtrIfNonZero(field)
		case "Topography":
			record.Topography = val.StringPtrIfNonZero(field)
		case "HeatingYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "5618c15a-6018-48c8-9e2c-29d4beea8f6e")
			}
			record.HeatingYN = v
		case "HeatingFeatures":
			record.HeatingFeatures = val.StringPtrIfNonZero(field)
		case "CoolingYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "c8a298be-ce09-4cca-9918-267d044c36a4")
			}
			record.CoolingYN = v
		case "Cooling":
			record.Cooling = val.StringPtrIfNonZero(field)
		case "FireplaceYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "b560d878-250b-4c93-b178-e5b8b0c01f91")
			}
			record.FireplaceYN = v
		case "Fireplace":
			record.Fireplace = val.StringPtrIfNonZero(field)
		case "FireplaceNumber":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3af45654-887d-43f4-a2b4-84da5c29eaed")
			}
			record.FireplaceNumber = v
		case "FoundationFeatures":
			record.FoundationFeatures = val.StringPtrIfNonZero(field)
		case "Roof":
			record.Roof = val.StringPtrIfNonZero(field)
		case "ArchitecturalStyleFeatures":
			record.ArchitecturalStyleFeatures = val.StringPtrIfNonZero(field)
		case "PatioAndPorchFeatures":
			record.PatioAndPorchFeatures = val.StringPtrIfNonZero(field)
		case "Utilities":
			record.Utilities = val.StringPtrIfNonZero(field)
		case "ElectricIncluded":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "88d16f26-0d83-4514-883b-8b260d9745ca")
			}
			record.ElectricIncluded = v
		case "ElectricDescription":
			record.ElectricDescription = val.StringPtrIfNonZero(field)
		case "WaterIncluded":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "4465075f-0c24-4f49-88e6-23ee5d84858f")
			}
			record.WaterIncluded = v
		case "WaterSource":
			record.WaterSource = val.StringPtrIfNonZero(field)
		case "Sewer":
			record.Sewer = val.StringPtrIfNonZero(field)
		case "GasDescription":
			record.GasDescription = val.StringPtrIfNonZero(field)
		case "OtherEquipmentIncluded":
			record.OtherEquipmentIncluded = val.StringPtrIfNonZero(field)
		case "LaundryFeatures":
			record.LaundryFeatures = val.StringPtrIfNonZero(field)
		case "Appliances":
			record.Appliances = val.StringPtrIfNonZero(field)
		case "InteriorFeatures":
			record.InteriorFeatures = val.StringPtrIfNonZero(field)
		case "ExteriorFeatures":
			record.ExteriorFeatures = val.StringPtrIfNonZero(field)
		case "FencingFeatures":
			record.FencingFeatures = val.StringPtrIfNonZero(field)
		case "PetsAllowed":
			record.PetsAllowed = val.StringPtrIfNonZero(field)
		case "HorseZoningYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "3099b46c-b3a6-432d-ac65-f2d17d599203")
			}
			record.HorseZoningYN = v
		case "SeniorCommunityYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "972b8b74-65a4-4ba8-ae5a-3f510e5bebec")
			}
			record.SeniorCommunityYN = v
		case "WaterbodyName":
			record.WaterbodyName = val.StringPtrIfNonZero(field)
		case "WaterfrontYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "2a7beff0-c3fe-469c-b3b5-ba6d81a8b71d")
			}
			record.WaterfrontYN = v
		case "WaterfrontFeatures":
			record.WaterfrontFeatures = val.StringPtrIfNonZero(field)
		case "ZoningCode":
			record.ZoningCode = val.StringPtrIfNonZero(field)
		case "ZoningDescription":
			record.ZoningDescription = val.StringPtrIfNonZero(field)
		case "CurrentUse":
			record.CurrentUse = val.StringPtrIfNonZero(field)
		case "PossibleUse":
			record.PossibleUse = val.StringPtrIfNonZero(field)
		case "AssociationYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "8be503c1-b8d0-4910-bb30-295093847225")
			}
			record.AssociationYN = v
		case "Association1Name":
			record.Association1Name = val.StringPtrIfNonZero(field)
		case "Association1Phone":
			record.Association1Phone = val.StringPtrIfNonZero(field)
		case "Association1Fee":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2e5fdf4c-8ef2-413c-92b6-ea84cf1f8610")
			}
			record.Association1Fee = v
		case "Association1FeeFrequency":
			record.Association1FeeFrequency = val.StringPtrIfNonZero(field)
		case "Association2Name":
			record.Association2Name = val.StringPtrIfNonZero(field)
		case "Association2Phone":
			record.Association2Phone = val.StringPtrIfNonZero(field)
		case "Association2Fee":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d76a423b-0142-41aa-9cc7-84b21c66cf0d")
			}
			record.Association2Fee = v
		case "Association2FeeFrequency":
			record.Association2FeeFrequency = val.StringPtrIfNonZero(field)
		case "AssociationFeeIncludes":
			record.AssociationFeeIncludes = val.StringPtrIfNonZero(field)
		case "AssociationAmenities":
			record.AssociationAmenities = val.StringPtrIfNonZero(field)
		case "SchoolElementary":
			record.SchoolElementary = val.StringPtrIfNonZero(field)
		case "SchoolElementaryDistrict":
			record.SchoolElementaryDistrict = val.StringPtrIfNonZero(field)
		case "SchoolMiddle":
			record.SchoolMiddle = val.StringPtrIfNonZero(field)
		case "SchoolMiddleDistrict":
			record.SchoolMiddleDistrict = val.StringPtrIfNonZero(field)
		case "SchoolHigh":
			record.SchoolHigh = val.StringPtrIfNonZero(field)
		case "SchoolHighDistrict":
			record.SchoolHighDistrict = val.StringPtrIfNonZero(field)
		case "GreenVerificationYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "28361dba-9e34-4dc7-8abb-f53f7ae85ca3")
			}
			record.GreenVerificationYN = v
		case "GreenBuildingVerificationType":
			record.GreenBuildingVerificationType = val.StringPtrIfNonZero(field)
		case "GreenEnergyEfficient":
			record.GreenEnergyEfficient = val.StringPtrIfNonZero(field)
		case "GreenEnergyGeneration":
			record.GreenEnergyGeneration = val.StringPtrIfNonZero(field)
		case "GreenIndoorAirQuality":
			record.GreenIndoorAirQuality = val.StringPtrIfNonZero(field)
		case "GreenLocation":
			record.GreenLocation = val.StringPtrIfNonZero(field)
		case "GreenSustainability":
			record.GreenSustainability = val.StringPtrIfNonZero(field)
		case "GreenWaterConservation":
			record.GreenWaterConservation = val.StringPtrIfNonZero(field)
		case "LandLeaseYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "d5cc4410-89d1-4399-a9aa-1982669781f4")
			}
			record.LandLeaseYN = v
		case "LandLeaseAmount":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "50cac833-e1f2-43b8-994a-1dcad63c896e")
			}
			record.LandLeaseAmount = v
		case "LandLeaseAmountFrequency":
			record.LandLeaseAmountFrequency = val.StringPtrIfNonZero(field)
		case "LandLeaseExpirationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.USSlashDate, field)
			if err != nil {
				return nil, errors.Forward(err, "14f21808-1399-4225-adff-568d07fb4f55")
			}
			record.LandLeaseExpirationDate = v
		case "CapRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ac1d8563-99ca-4f88-bb85-3b4edc6fe2ce")
			}
			record.CapRate = v
		case "GrossIncome":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7684bda3-1fd6-4395-b138-4bc4bcce44cc")
			}
			record.GrossIncome = v
		case "IncomeIncludes":
			record.IncomeIncludes = val.StringPtrIfNonZero(field)
		case "GrossScheduledIncome":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "796f7df6-c0d4-43aa-8364-db95510d438d")
			}
			record.GrossScheduledIncome = v
		case "NetOperatingIncome":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "204a934e-ecf4-40df-998c-1caa35fe0e88")
			}
			record.NetOperatingIncome = v
		case "TotalActualRent":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1e0f4c12-99cb-4b14-a16c-0413c80956e3")
			}
			record.TotalActualRent = v
		case "ExistingLeaseType":
			record.ExistingLeaseType = val.StringPtrIfNonZero(field)
		case "FinancialDataSource":
			record.FinancialDataSource = val.StringPtrIfNonZero(field)
		case "RentControlYN":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "af4bc4aa-d0c8-4bc1-9ec3-8490d806da10")
			}
			record.RentControlYN = v
		case "UnitTypeDescription":
			record.UnitTypeDescription = val.StringPtrIfNonZero(field)
		case "UnitTypeFurnished":
			record.UnitTypeFurnished = val.StringPtrIfNonZero(field)
		case "NumberOfUnitsLeased":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f1c5451e-bed0-490e-ae94-b518e1688338")
			}
			record.NumberOfUnitsLeased = v
		case "NumberOfUnitsMoMo":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9b19a0cc-c04c-46fb-acbe-d583191550ae")
			}
			record.NumberOfUnitsMoMo = v
		case "NumberOfUnitsVacant":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "35132b04-e583-44aa-a83c-4abd26dc249b")
			}
			record.NumberOfUnitsVacant = v
		case "VacancyAllowance":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "8d17f237-d021-40fc-9263-1b23e683d102")
			}
			record.VacancyAllowance = v
		case "VacancyAllowanceRate":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d5fc2d3a-7e2b-47ec-9511-1912cab7d841")
			}
			record.VacancyAllowanceRate = v
		case "OperatingExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "44cffc03-610c-49f7-a4a8-15f2f45281ca")
			}
			record.OperatingExpense = v
		case "CableTvExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f2025832-36e8-463c-bb02-11f0b13770ba")
			}
			record.CableTvExpense = v
		case "ElectricExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e6d503d9-9e5e-409c-8d49-942cee15562a")
			}
			record.ElectricExpense = v
		case "FuelExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f77a54f6-86c0-4dfa-aa41-e9536e516871")
			}
			record.FuelExpense = v
		case "FurnitureReplacementExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "70e99c0e-1a0d-4a6e-8a1f-7c8908394cd2")
			}
			record.FurnitureReplacementExpense = v
		case "GardenerExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "aba63f89-6a9e-4317-a110-3e4184332278")
			}
			record.GardenerExpense = v
		case "InsuranceExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0d70b4d4-0d21-4c96-a3fe-008a5d79e54b")
			}
			record.InsuranceExpense = v
		case "OperatingExpenseIncludes":
			record.OperatingExpenseIncludes = val.StringPtrIfNonZero(field)
		case "LicensesExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bd8a052e-77ee-4417-bdc8-6927adf77fae")
			}
			record.LicensesExpense = v
		case "MaintenanceExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "194aa0c1-e383-47d0-a2e0-bd9d23202b22")
			}
			record.MaintenanceExpense = v
		case "ManagerExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "37d4ceea-f6e4-4bdb-a7fb-86f29e599d21")
			}
			record.ManagerExpense = v
		case "NewTaxesExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "94bc167d-9433-498c-84ac-eb6a178b0e81")
			}
			record.NewTaxesExpense = v
		case "OtherExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e4ca1fba-74d9-4ccd-91df-911bf2b613b5")
			}
			record.OtherExpense = v
		case "PestControlExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d2b434ed-8905-4e76-816a-ba3e36f223bc")
			}
			record.PestControlExpense = v
		case "PoolExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "90d8c8cc-2a63-4038-b247-8aaef1377a2a")
			}
			record.PoolExpense = v
		case "ProfessionalManagementExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ef518129-a29a-4dcd-ba25-18ba21430490")
			}
			record.ProfessionalManagementExpense = v
		case "SuppliesExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3f2a1529-1b3f-4d2d-8c6e-8896950f08d1")
			}
			record.SuppliesExpense = v
		case "TrashExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "69e71346-41e4-4741-ae95-0ae647559e2c")
			}
			record.TrashExpense = v
		case "WaterSewerExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "01c7d6f4-b424-4420-b5ab-47a0474c3a58")
			}
			record.WaterSewerExpense = v
		case "WorkmansCompensationExpense":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7e0af996-f943-454e-a09c-3c7bfaa842e8")
			}
			record.WorkmansCompensationExpense = v
		case "OwnerPays":
			record.OwnerPays = val.StringPtrIfNonZero(field)
		case "TenantPays":
			record.TenantPays = val.StringPtrIfNonZero(field)
		case "ListingMarketingURL":
			record.ListingMarketingURL = val.StringPtrIfNonZero(field)
		case "PhotosCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6f189498-18c3-4537-bd04-1b9b9577ba50")
			}
			record.PhotosCount = v
		case "PhotoKey":
			record.PhotoKey = val.StringPtrIfNonZero(field)
		case "PhotoURLPrefix":
			record.PhotoURLPrefix = val.StringPtrIfNonZero(field)
		case "CurrentStatus":
			v, err := val.BoolPtrFromYNString(field)
			if err != nil {
				return nil, errors.Forward(err, "7ecabb89-9a47-4f93-920a-94b26b92a8ba")
			}
			record.CurrentStatus = v
		default:
			return nil, &errors.Object{
				Id:     "c5e47c6e-4c69-4b2b-8593-5cd60573035e",
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

func (dr *Listing) SQLColumns() []string {
	columns := []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"attom_id",
		"mls_record_id",
		"mls_listing_id",
		"status_change_date",
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
		"situs_county",
		"township",
		"mls_listing_address",
		"mls_listing_city",
		"mls_listing_state",
		"mls_listing_zip",
		"mls_listing_county_fips",
		"mls_number",
		"mls_source",
		"listing_status",
		"mls_sold_date",
		"mls_sold_price",
		"assessor_last_sale_date",
		"assessor_last_sale_amount",
		"market_value",
		"market_value_date",
		"avg_market_price_per_sq_ft",
		"listing_date",
		"latest_listing_price",
		"previous_listing_price",
		"latest_price_change_date",
		"pending_date",
		"special_listing_conditions",
		"original_listing_date",
		"original_listing_price",
		"lease_option",
		"lease_term",
		"lease_includes",
		"concessions",
		"concessions_amount",
		"concessions_comments",
		"contingency_date",
		"contingency_description",
		"mls_property_type",
		"mls_property_sub_type",
		"attom_property_type",
		"attom_property_sub_type",
		"ownership_description",
		"latitude",
		"longitude",
		"apn_formatted",
		"legal_description",
		"legal_subdivision",
		"days_on_market",
		"cumulative_days_on_market",
		"listing_agent_full_name",
		"listing_agent_mls_id",
		"listing_agent_state_license",
		"listing_agent_aor",
		"listing_agent_preferred_phone",
		"listing_agent_email",
		"listing_office_name",
		"listing_office_mls_id",
		"listing_office_aor",
		"listing_office_phone",
		"listing_office_email",
		"listing_co_agent_full_name",
		"listing_co_agent_mls_id",
		"listing_co_agent_state_license",
		"listing_co_agent_aor",
		"listing_co_agent_preferred_phone",
		"listing_co_agent_email",
		"listing_co_agent_office_name",
		"listing_co_agent_office_mls_id",
		"listing_co_agent_office_aor",
		"listing_co_agent_office_phone",
		"listing_co_agent_office_email",
		"buyer_agent_full_name",
		"buyer_agent_mls_id",
		"buyer_agent_state_license",
		"buyer_agent_aor",
		"buyer_agent_preferred_phone",
		"buyer_agent_email",
		"buyer_office_name",
		"buyer_office_mls_id",
		"buyer_office_aor",
		"buyer_office_phone",
		"buyer_office_email",
		"buyer_co_agent_full_name",
		"buyer_co_agent_mls_id",
		"buyer_co_agent_state_license",
		"buyer_co_agent_aor",
		"buyer_co_agent_preferred_phone",
		"buyer_co_agent_email",
		"buyer_co_agent_office_name",
		"buyer_co_agent_office_mls_id",
		"buyer_co_agent_office_aor",
		"buyer_co_agent_office_phone",
		"buyer_co_agent_office_email",
		"public_listing_remarks",
		"home_warranty_yn",
		"tax_year_assessed",
		"tax_assessed_value_total",
		"tax_amount",
		"tax_annual_other",
		"owner_name",
		"owner_vesting",
		"year_built",
		"year_built_effective",
		"year_built_source",
		"new_construction_yn",
		"builder_name",
		"additional_parcels_yn",
		"number_of_lots",
		"lot_size_square_feet",
		"lot_size_acres",
		"lot_size_source",
		"lot_dimensions",
		"lot_feature_list",
		"frontage_length",
		"frontage_type",
		"frontage_road_type",
		"living_area_square_feet",
		"living_area_source",
		"levels",
		"stories",
		"building_stories_total",
		"building_keywords",
		"building_area_total",
		"number_of_units_total",
		"number_of_buildings",
		"property_attached_yn",
		"other_structures",
		"rooms_total",
		"bedrooms_total",
		"bathrooms_full",
		"bathrooms_half",
		"bathrooms_quarter",
		"bathrooms_three_quarters",
		"basement_features",
		"below_grade_square_feet",
		"basement_total_sq_ft",
		"basement_finished_sq_ft",
		"basement_unfinished_sq_ft",
		"property_condition",
		"repairs_yn",
		"repairs_description",
		"disclosures",
		"construction_materials",
		"garage_yn",
		"attached_garage_yn",
		"garage_spaces",
		"carport_yn",
		"carport_spaces",
		"parking_features",
		"parking_other",
		"open_parking_spaces",
		"parking_total",
		"pool_private_yn",
		"pool_features",
		"occupancy",
		"view_yn",
		"view_col",
		"topography",
		"heating_yn",
		"heating_features",
		"cooling_yn",
		"cooling",
		"fireplace_yn",
		"fireplace",
		"fireplace_number",
		"foundation_features",
		"roof",
		"architectural_style_features",
		"patio_and_porch_features",
		"utilities",
		"electric_included",
		"electric_description",
		"water_included",
		"water_source",
		"sewer",
		"gas_description",
		"other_equipment_included",
		"laundry_features",
		"appliances",
		"interior_features",
		"exterior_features",
		"fencing_features",
		"pets_allowed",
		"horse_zoning_yn",
		"senior_community_yn",
		"waterbody_name",
		"waterfront_yn",
		"waterfront_features",
		"zoning_code",
		"zoning_description",
		"current_use",
		"possible_use",
		"association_yn",
		"association1_name",
		"association1_phone",
		"association1_fee",
		"association1_fee_frequency",
		"association2_name",
		"association2_phone",
		"association2_fee",
		"association2_fee_frequency",
		"association_fee_includes",
		"association_amenities",
		"school_elementary",
		"school_elementary_district",
		"school_middle",
		"school_middle_district",
		"school_high",
		"school_high_district",
		"green_verification_yn",
		"green_building_verification_type",
		"green_energy_efficient",
		"green_energy_generation",
		"green_indoor_air_quality",
		"green_location",
		"green_sustainability",
		"green_water_conservation",
		"land_lease_yn",
		"land_lease_amount",
		"land_lease_amount_frequency",
		"land_lease_expiration_date",
		"cap_rate",
		"gross_income",
		"income_includes",
		"gross_scheduled_income",
		"net_operating_income",
		"total_actual_rent",
		"existing_lease_type",
		"financial_data_source",
		"rent_control_yn",
		"unit_type_description",
		"unit_type_furnished",
		"number_of_units_leased",
		"number_of_units_mo_mo",
		"number_of_units_vacant",
		"vacancy_allowance",
		"vacancy_allowance_rate",
		"operating_expense",
		"cable_tv_expense",
		"electric_expense",
		"fuel_expense",
		"furniture_replacement_expense",
		"gardener_expense",
		"insurance_expense",
		"operating_expense_includes",
		"licenses_expense",
		"maintenance_expense",
		"manager_expense",
		"new_taxes_expense",
		"other_expense",
		"pest_control_expense",
		"pool_expense",
		"professional_management_expense",
		"supplies_expense",
		"trash_expense",
		"water_sewer_expense",
		"workmans_compensation_expense",
		"owner_pays",
		"tenant_pays",
		"listing_marketing_url",
		"photos_count",
		"photo_key",
		"photo_url_prefix",
	}

	if dr.fileType >= DataFileTypeListingV20250417 {
		columns = append(columns, "current_status")
	}

	return columns
}

func (dr *Listing) SQLTable() string {
	return "ad_df_listing"
}

func (dr *Listing) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "23115865-839a-4372-8718-61995de9439b",
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
		dr.MLSRecordID,
		dr.MLSListingID,
		dr.StatusChangeDate,
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
		dr.SitusCounty,
		dr.Township,
		dr.MLSListingAddress,
		dr.MLSListingCity,
		dr.MLSListingState,
		dr.MLSListingZip,
		dr.MLSListingCountyFIPS,
		dr.MLSNumber,
		dr.MLSSource,
		dr.ListingStatus,
		dr.MLSSoldDate,
		dr.MLSSoldPrice,
		dr.AssessorLastSaleDate,
		dr.AssessorLastSaleAmount,
		dr.MarketValue,
		dr.MarketValueDate,
		dr.AvgMarketPricePerSqFt,
		dr.ListingDate,
		dr.LatestListingPrice,
		dr.PreviousListingPrice,
		dr.LatestPriceChangeDate,
		dr.PendingDate,
		dr.SpecialListingConditions,
		dr.OriginalListingDate,
		dr.OriginalListingPrice,
		dr.LeaseOption,
		dr.LeaseTerm,
		dr.LeaseIncludes,
		dr.Concessions,
		dr.ConcessionsAmount,
		dr.ConcessionsComments,
		dr.ContingencyDate,
		dr.ContingencyDescription,
		dr.MLSPropertyType,
		dr.MLSPropertySubType,
		dr.ATTOMPropertyType,
		dr.ATTOMPropertySubType,
		dr.OwnershipDescription,
		dr.Latitude,
		dr.Longitude,
		dr.APNFormatted,
		dr.LegalDescription,
		dr.LegalSubdivision,
		dr.DaysOnMarket,
		dr.CumulativeDaysOnMarket,
		dr.ListingAgentFullName,
		dr.ListingAgentMLSID,
		dr.ListingAgentStateLicense,
		dr.ListingAgentAOR,
		dr.ListingAgentPreferredPhone,
		dr.ListingAgentEmail,
		dr.ListingOfficeName,
		dr.ListingOfficeMlsId,
		dr.ListingOfficeAOR,
		dr.ListingOfficePhone,
		dr.ListingOfficeEmail,
		dr.ListingCoAgentFullName,
		dr.ListingCoAgentMLSID,
		dr.ListingCoAgentStateLicense,
		dr.ListingCoAgentAOR,
		dr.ListingCoAgentPreferredPhone,
		dr.ListingCoAgentEmail,
		dr.ListingCoAgentOfficeName,
		dr.ListingCoAgentOfficeMlsId,
		dr.ListingCoAgentOfficeAOR,
		dr.ListingCoAgentOfficePhone,
		dr.ListingCoAgentOfficeEmail,
		dr.BuyerAgentFullName,
		dr.BuyerAgentMLSID,
		dr.BuyerAgentStateLicense,
		dr.BuyerAgentAOR,
		dr.BuyerAgentPreferredPhone,
		dr.BuyerAgentEmail,
		dr.BuyerOfficeName,
		dr.BuyerOfficeMlsId,
		dr.BuyerOfficeAOR,
		dr.BuyerOfficePhone,
		dr.BuyerOfficeEmail,
		dr.BuyerCoAgentFullName,
		dr.BuyerCoAgentMLSID,
		dr.BuyerCoAgentStateLicense,
		dr.BuyerCoAgentAOR,
		dr.BuyerCoAgentPreferredPhone,
		dr.BuyerCoAgentEmail,
		dr.BuyerCoAgentOfficeName,
		dr.BuyerCoAgentOfficeMlsId,
		dr.BuyerCoAgentOfficeAOR,
		dr.BuyerCoAgentOfficePhone,
		dr.BuyerCoAgentOfficeEmail,
		dr.PublicListingRemarks,
		dr.HomeWarrantyYN,
		dr.TaxYearAssessed,
		dr.TaxAssessedValueTotal,
		dr.TaxAmount,
		dr.TaxAnnualOther,
		dr.OwnerName,
		dr.OwnerVesting,
		dr.YearBuilt,
		dr.YearBuiltEffective,
		dr.YearBuiltSource,
		dr.NewConstructionYN,
		dr.BuilderName,
		dr.AdditionalParcelsYN,
		dr.NumberOfLots,
		dr.LotSizeSquareFeet,
		dr.LotSizeAcres,
		dr.LotSizeSource,
		dr.LotDimensions,
		dr.LotFeatureList,
		dr.FrontageLength,
		dr.FrontageType,
		dr.FrontageRoadType,
		dr.LivingAreaSquareFeet,
		dr.LivingAreaSource,
		dr.Levels,
		dr.Stories,
		dr.BuildingStoriesTotal,
		dr.BuildingKeywords,
		dr.BuildingAreaTotal,
		dr.NumberOfUnitsTotal,
		dr.NumberOfBuildings,
		dr.PropertyAttachedYN,
		dr.OtherStructures,
		dr.RoomsTotal,
		dr.BedroomsTotal,
		dr.BathroomsFull,
		dr.BathroomsHalf,
		dr.BathroomsQuarter,
		dr.BathroomsThreeQuarters,
		dr.BasementFeatures,
		dr.BelowGradeSquareFeet,
		dr.BasementTotalSqFt,
		dr.BasementFinishedSqFt,
		dr.BasementUnfinishedSqFt,
		dr.PropertyCondition,
		dr.RepairsYN,
		dr.RepairsDescription,
		dr.Disclosures,
		dr.ConstructionMaterials,
		dr.GarageYN,
		dr.AttachedGarageYN,
		dr.GarageSpaces,
		dr.CarportYN,
		dr.CarportSpaces,
		dr.ParkingFeatures,
		dr.ParkingOther,
		dr.OpenParkingSpaces,
		dr.ParkingTotal,
		dr.PoolPrivateYN,
		dr.PoolFeatures,
		dr.Occupancy,
		dr.ViewYN,
		dr.View,
		dr.Topography,
		dr.HeatingYN,
		dr.HeatingFeatures,
		dr.CoolingYN,
		dr.Cooling,
		dr.FireplaceYN,
		dr.Fireplace,
		dr.FireplaceNumber,
		dr.FoundationFeatures,
		dr.Roof,
		dr.ArchitecturalStyleFeatures,
		dr.PatioAndPorchFeatures,
		dr.Utilities,
		dr.ElectricIncluded,
		dr.ElectricDescription,
		dr.WaterIncluded,
		dr.WaterSource,
		dr.Sewer,
		dr.GasDescription,
		dr.OtherEquipmentIncluded,
		dr.LaundryFeatures,
		dr.Appliances,
		dr.InteriorFeatures,
		dr.ExteriorFeatures,
		dr.FencingFeatures,
		dr.PetsAllowed,
		dr.HorseZoningYN,
		dr.SeniorCommunityYN,
		dr.WaterbodyName,
		dr.WaterfrontYN,
		dr.WaterfrontFeatures,
		dr.ZoningCode,
		dr.ZoningDescription,
		dr.CurrentUse,
		dr.PossibleUse,
		dr.AssociationYN,
		dr.Association1Name,
		dr.Association1Phone,
		dr.Association1Fee,
		dr.Association1FeeFrequency,
		dr.Association2Name,
		dr.Association2Phone,
		dr.Association2Fee,
		dr.Association2FeeFrequency,
		dr.AssociationFeeIncludes,
		dr.AssociationAmenities,
		dr.SchoolElementary,
		dr.SchoolElementaryDistrict,
		dr.SchoolMiddle,
		dr.SchoolMiddleDistrict,
		dr.SchoolHigh,
		dr.SchoolHighDistrict,
		dr.GreenVerificationYN,
		dr.GreenBuildingVerificationType,
		dr.GreenEnergyEfficient,
		dr.GreenEnergyGeneration,
		dr.GreenIndoorAirQuality,
		dr.GreenLocation,
		dr.GreenSustainability,
		dr.GreenWaterConservation,
		dr.LandLeaseYN,
		dr.LandLeaseAmount,
		dr.LandLeaseAmountFrequency,
		dr.LandLeaseExpirationDate,
		dr.CapRate,
		dr.GrossIncome,
		dr.IncomeIncludes,
		dr.GrossScheduledIncome,
		dr.NetOperatingIncome,
		dr.TotalActualRent,
		dr.ExistingLeaseType,
		dr.FinancialDataSource,
		dr.RentControlYN,
		dr.UnitTypeDescription,
		dr.UnitTypeFurnished,
		dr.NumberOfUnitsLeased,
		dr.NumberOfUnitsMoMo,
		dr.NumberOfUnitsVacant,
		dr.VacancyAllowance,
		dr.VacancyAllowanceRate,
		dr.OperatingExpense,
		dr.CableTvExpense,
		dr.ElectricExpense,
		dr.FuelExpense,
		dr.FurnitureReplacementExpense,
		dr.GardenerExpense,
		dr.InsuranceExpense,
		dr.OperatingExpenseIncludes,
		dr.LicensesExpense,
		dr.MaintenanceExpense,
		dr.ManagerExpense,
		dr.NewTaxesExpense,
		dr.OtherExpense,
		dr.PestControlExpense,
		dr.PoolExpense,
		dr.ProfessionalManagementExpense,
		dr.SuppliesExpense,
		dr.TrashExpense,
		dr.WaterSewerExpense,
		dr.WorkmansCompensationExpense,
		dr.OwnerPays,
		dr.TenantPays,
		dr.ListingMarketingURL,
		dr.PhotosCount,
		dr.PhotoKey,
		dr.PhotoURLPrefix,
	}

	if dr.fileType >= DataFileTypeListingV20250417 {
		values = append(values, dr.CurrentStatus)
	}

	return values, nil
}

func (dr *Listing) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		LoadFunc: dr.LoadFunc,
		Mode:     entities.DataRecordModeLoadFunc,
	}
}

func (dr *Listing) LoadFunc(r *arc.Request, in *entities.LoadDataRecordInput) (*entities.LoadDataRecordOutput, error) {
	// newListingsBuilder := func() squirrel.InsertBuilder {
	// 	return squirrel.StatementBuilder.
	// 		PlaceholderFormat(squirrel.Dollar).
	// 		Insert(in.DataRecord.SQLTable()).
	// 		Columns(in.Columns...).
	// 		Suffix("returning am_id, attom_id")
	// }

	// currentListingsDeleteCount := int64(0)
	// currentListingsInsertCount := int64(0)
	// dfObject := in.DataFileObject
	// listingsBuilder := newListingsBuilder()
	// listingsAmIds := make([]any, in.BatchSize)
	// listingsAttomIds := make([]any, in.BatchSize)
	// recordCount := int32(0)
	processedRecords := int64(0)
	// scanner := in.Scanner

	// pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	// if err != nil {
	// 	return nil, errors.Forward(err, "3a0f2e33-0a07-4296-9142-97c82e230b66")
	// }

	// loadRecords := func() error {
	// 	// Skip if no records to insert.
	// 	if recordCount == 0 {
	// 		return nil
	// 	}

	// 	listingsSql, listingsArgs, err := listingsBuilder.ToSql()
	// 	if err != nil {
	// 		return &errors.Object{
	// 			Id:     "c80241ca-ca60-4ab4-959e-29607444a5f7",
	// 			Code:   errors.Code_UNKNOWN,
	// 			Detail: "Failed to build SQL.",
	// 			Cause:  err.Error(),
	// 		}
	// 	}

	// 	tx, err := pgxPool.Begin(r.Context())
	// 	if err != nil {
	// 		return &errors.Object{
	// 			Id:     "5a05ca66-f366-45ce-b124-d0af38d1e738",
	// 			Code:   errors.Code_UNKNOWN,
	// 			Detail: "Failed to begin transaction.",
	// 			Cause:  err.Error(),
	// 		}
	// 	}

	// 	defer extutils.RollbackPgxTx(r.Context(), tx, "c0c22d86-226f-4584-89f2-9a6c48990a35")

	// 	// This clone won't replace the original arc.
	// 	r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

	// 	// Insert the Listing records.
	// 	listingsRows, err := extutils.PgxQuery(
	// 		r,
	// 		consts.ConfigKeyPostgresDatapipe,
	// 		listingsSql,
	// 		listingsArgs,
	// 	)
	// 	if err != nil {
	// 		return errors.Forward(err, "9872eff0-d2d7-42c7-9f99-19de83008f42")
	// 	}

	// 	for i := 0; listingsRows.Next(); i++ {
	// 		var amId uuid.UUID
	// 		var attomId int64

	// 		if err := listingsRows.Scan(
	// 			&amId,
	// 			&attomId,
	// 		); err != nil {
	// 			listingsRows.Close()
	// 			return &errors.Object{
	// 				Id:     "14378854-0c23-4926-be3b-31da8476e714",
	// 				Code:   errors.Code_UNKNOWN,
	// 				Detail: "Failed to scan rows.",
	// 				Cause:  err.Error(),
	// 			}
	// 		}

	// 		listingsAmIds[i] = amId
	// 		listingsAttomIds[i] = attomId
	// 	}

	// 	listingsRows.Close()

	// 	if listingsRows.Err() != nil {
	// 		return &errors.Object{
	// 			Id:     "678ad78b-98fc-48e7-b05c-47a6b20e4a25",
	// 			Code:   errors.Code_UNKNOWN,
	// 			Detail: "Failed to query rows.",
	// 			Cause:  listingsRows.Err().Error(),
	// 		}
	// 	}

	// 	// Ensure to delete records from the current_listings table.
	// 	currentListingsDeleteCt, err := extutils.PgxExec(
	// 		r,
	// 		consts.ConfigKeyPostgresDatapipe,
	// 		`delete from current_listings where attom_id = any($1)`,
	// 		[]any{listingsAttomIds[:recordCount]},
	// 	)
	// 	if err != nil {
	// 		return errors.Forward(err, "902d197e-aca3-49be-a12a-d9722ff6c8eb")
	// 	}

	// 	currentListingsDeleteCount += currentListingsDeleteCt.RowsAffected()

	// 	// Create the current_listings record.
	// 	currentListingsInsertCt, err := extutils.PgxExec(
	// 		r,
	// 		consts.ConfigKeyPostgresDatapipe,
	// 		`
	// 		insert into current_listings
	// 		with ranked_listings as (
	// 			select
	// 				ad_df_listing.am_id,
	// 				ad_df_listing.attom_id,
	// 				ad_df_listing.status_change_date,
	// 				ad_df_listing.mls_property_type,
	// 				ad_df_listing.mls_property_sub_type,
	// 				ad_df_listing.bedrooms_total,
	// 				ad_df_listing.bathrooms_full,
	// 				ad_df_listing.latest_listing_price,
	// 				ad_df_listing.living_area_square_feet,
	// 				ad_df_listing.year_built,
	// 				ad_df_listing.listing_status,
	// 				ad_df_listing.mls_source,
	// 				ad_df_listing.mls_number,
	// 				ad_df_listing.attom_property_sub_type,
	// 				row_number() over (
	// 					partition by ad_df_listing.attom_id
	// 					order by ad_df_listing.status_change_date desc
	// 				) as row_num
	// 			from ad_df_listing
	// 			where
	// 				ad_df_listing.am_id = any($1)
	// 				and ad_df_listing.status_change_date >= '2020-01-01'
	// 		)
	// 		select
	// 			gen_random_uuid(),
	// 			now(),
	// 			now(),
	// 			null,
	// 			ranked_listings.am_id as am_listing_id,
	// 			ranked_listings.attom_id,
	// 			properties.id,
	// 			ranked_listings.status_change_date,
	// 			mp.ouid,
	// 			case when mp.ouid = 'M00000146' and length(ranked_listings.mls_number) > 3 then
	// 				substring(ranked_listings.mls_number from 4)
	// 			else
	// 				regexp_replace(ranked_listings.mls_number, '[^0-9]', '', 'g')
	// 			end as mls_number,
	// 			ranked_listings.mls_property_type,
	// 			ranked_listings.mls_property_sub_type,
	// 			ranked_listings.bedrooms_total,
	// 			ranked_listings.bathrooms_full,
	// 			ranked_listings.latest_listing_price,
	// 			ranked_listings.living_area_square_feet,
	// 			ranked_listings.year_built,
	// 			ranked_listings.listing_status,
	// 			addresses.zip5,
	// 			null,
	// 			null,
	// 			ad_geom.location_3857,
	// 			ad_geom.location as location_4326,
	// 			ranked_listings.attom_property_sub_type
	// 		from ranked_listings
	// 		join ad_geom on ranked_listings.attom_id = ad_geom.attom_id
	// 		join properties on ranked_listings.attom_id = properties.ad_attom_id
	// 		join addresses on properties.address_id = addresses.id
	// 		left join mls_providers mp on ranked_listings.mls_source = mp.mls_source
	// 		where ranked_listings.row_num = 1
	// 		`,
	// 		[]any{listingsAmIds[:recordCount]},
	// 	)
	// 	if err != nil {
	// 		return errors.Forward(err, "23506f05-090d-4307-904e-c22f5dea7c9e")
	// 	}

	// 	currentListingsInsertCount += currentListingsInsertCt.RowsAffected()

	// 	updateObjectOut, err := in.UpdateDataFileObjectFunc(r, &entities.UpdateDataFileObjectInput{
	// 		Id:          dfObject.Id,
	// 		UpdatedAt:   time.Now(),
	// 		RecordCount: dfObject.RecordCount + recordCount,
	// 		Status:      entities.DataFileObjectStatusInProgress,
	// 	})
	// 	if err != nil {
	// 		extutils.RollbackPgxTx(r.Context(), tx, "e09b4218-fda5-4edc-8915-465547cdeeba")
	// 		return errors.Forward(err, "b7c4ea4b-4bf9-44e5-ada9-2e54f75428fb")
	// 	}

	// 	if err := tx.Commit(r.Context()); err != nil {
	// 		return &errors.Object{
	// 			Id:     "e9bebe32-d1e0-4b91-9ffa-96641f2ac217",
	// 			Code:   errors.Code_UNKNOWN,
	// 			Detail: "Failed to commit transaction.",
	// 			Cause:  err.Error(),
	// 		}
	// 	}

	// 	dfObject = updateObjectOut.Entity
	// 	listingsBuilder = newListingsBuilder()
	// 	recordCount = 0

	// 	return nil
	// }

	// for scanner.Scan() {
	// 	if recordCount == in.BatchSize {
	// 		if err := loadRecords(); err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	fields := strings.Split(scanner.Text(), in.FieldSeparator)

	// 	record, err := in.DataRecord.New(in.Headers, fields)
	// 	if err != nil {
	// 		return nil, errors.Forward(err, "9ab1fc8c-80ed-4da7-b080-6da9c53e595c")
	// 	}

	// 	recordValues, err := record.SQLValues()
	// 	if err != nil {
	// 		return nil, errors.Forward(err, "3d50348c-0d8b-4541-956c-036d2e0626df")
	// 	}

	// 	listingsBuilder = listingsBuilder.Values(recordValues...)

	// 	recordCount++
	// 	processedRecords++
	// }

	// if err := loadRecords(); err != nil {
	// 	return nil, err
	// }

	// log.Info().
	// 	Int64("currentListingsDeleteCount", currentListingsDeleteCount).
	// 	Int64("currentListingsInsertCount", currentListingsInsertCount).
	// 	Int64("processedRecords", processedRecords).
	// 	Send()

	out := &entities.LoadDataRecordOutput{
		ProcessedRecords: processedRecords,
	}

	return out, nil
}
