package attom_data

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type Recorder struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any

	TransactionID                                     int64
	ATTOMID                                           int64
	DocumentRecordingStateCode                        *string
	DocumentRecordingCountyName                       *string
	DocumentRecordingJurisdictionName                 *string
	DocumentRecordingCountyFIPs                       *string
	DocumentTypeCode                                  *string
	DocumentNumberFormatted                           *string
	DocumentNumberLegacy                              *string
	InstrumentNumber                                  *string
	Book                                              *string
	Page                                              *string
	InstrumentDate                                    *time.Time
	RecordingDate                                     *time.Time
	TransactionType                                   *string
	TransferInfoPurchaseTypeCode                      *int
	ForeclosureAuctionSale                            *bool
	TransferInfoDistressCircumstanceCode              *int
	QuitclaimFlag                                     *bool
	TransferInfoMultiParcelFlag                       *int
	ArmsLengthFlag                                    *int
	PartialInterest                                   *string
	TransferAmount                                    *decimal.Decimal
	TransferAmountInfoAccuracy                        *string
	TransferTaxTotal                                  *decimal.Decimal
	TransferTaxCity                                   *decimal.Decimal
	TransferTaxCounty                                 *decimal.Decimal
	Grantor1NameFull                                  *string
	Grantor1NameFirst                                 *string
	Grantor1NameMiddle                                *string
	Grantor1NameLast                                  *string
	Grantor1NameSuffix                                *string
	Grantor1InfoEntityClassification                  *string
	Grantor1InfoOwnerType                             *string
	Grantor2NameFull                                  *string
	Grantor2NameFirst                                 *string
	Grantor2NameMiddle                                *string
	Grantor2NameLast                                  *string
	Grantor2NameSuffix                                *string
	Grantor2InfoEntityClassification                  *string
	Grantor2InfoOwnerType                             *string
	Grantor3NameFull                                  *string
	Grantor3NameFirst                                 *string
	Grantor3NameMiddle                                *string
	Grantor3NameLast                                  *string
	Grantor3NameSuffix                                *string
	Grantor3InfoEntityClassification                  *string
	Grantor4NameFull                                  *string
	Grantor4NameFirst                                 *string
	Grantor4NameMiddle                                *string
	Grantor4NameLast                                  *string
	Grantor4NameSuffix                                *string
	Grantor4InfoEntityClassification                  *string
	GrantorAddressFull                                *string
	GrantorAddressHouseNumber                         *string
	GrantorAddressStreetDirection                     *string
	GrantorAddressStreetName                          *string
	GrantorAddressStreetSuffix                        *string
	GrantorAddressStreetPostDirection                 *string
	GrantorAddressUnitPrefix                          *string
	GrantorAddressUnitValue                           *string
	GrantorAddressCity                                *string
	GrantorAddressState                               *string
	GrantorAddressZIP                                 *string
	GrantorAddressZIP4                                *string
	GrantorAddressCRRT                                *string
	GrantorAddressInfoFormat                          *string
	GrantorAddressInfoPrivacy                         *bool
	Grantee1NameFull                                  *string
	Grantee1NameFirst                                 *string
	Grantee1NameMiddle                                *string
	Grantee1NameLast                                  *string
	Grantee1NameSuffix                                *string
	Grantee1InfoEntityClassification                  *string
	Grantee1InfoOwnerType                             *string
	Grantee2NameFull                                  *string
	Grantee2NameFirst                                 *string
	Grantee2NameMiddle                                *string
	Grantee2NameLast                                  *string
	Grantee2NameSuffix                                *string
	Grantee2InfoEntityClassification                  *string
	GranteeInfoVesting1                               *string
	Grantee3NameFull                                  *string
	Grantee3NameFirst                                 *string
	Grantee3NameMiddle                                *string
	Grantee3NameLast                                  *string
	Grantee3NameSuffix                                *string
	Grantee3InfoEntityClassification                  *string
	Grantee4NameFull                                  *string
	Grantee4NameFirst                                 *string
	Grantee4NameMiddle                                *string
	Grantee4NameLast                                  *string
	Grantee4NameSuffix                                *string
	Grantee4InfoEntityClassification                  *string
	GranteeMailCareOfName                             *string
	GranteeInfoEntityCount                            *int
	GranteeInfoVesting2                               *string
	GranteeInvestorFlag                               *bool
	GranteeMailAddressFull                            *string
	GranteeMailAddressHouseNumber                     *string
	GranteeMailAddressStreetDirection                 *string
	GranteeMailAddressStreetName                      *string
	GranteeMailAddressStreetSuffix                    *string
	GranteeMailAddressStreetPostDirection             *string
	GranteeMailAddressUnitPrefix                      *string
	GranteeMailAddressUnitValue                       *string
	GranteeMailAddressCity                            *string
	GranteeMailAddressState                           *string
	GranteeMailAddressZIP                             *string
	GranteeMailAddressZIP4                            *string
	GranteeMailAddressCRRT                            *string
	GranteeMailAddressInfoFormat                      *string
	GranteeMailAddressInfoPrivacy                     *bool
	GranteeGrantorOwnerRelationshipCode               *string
	TitleCompanyStandardizedCode                      *string
	TitleCompanyStandardizedName                      *string
	TitleCompanyRaw                                   *string
	LegalDescriptionPart1                             *string
	LegalDescriptionPart2                             *string
	LegalDescriptionPart3                             *string
	LegalDescriptionPart4                             *string
	LegalRange                                        *string
	LegalTownship                                     *string
	LegalSection                                      *string
	LegalDistrict                                     *string
	LegalSubDivision                                  *string
	LegalTract                                        *string
	LegalBlock                                        *string
	LegalLot                                          *string
	LegalUnit                                         *string
	LegalPlatMapBook                                  *string
	LegalPlatMapPage                                  *string
	APNFormatted                                      *string
	APNOriginal                                       *string
	PropertyAddressFull                               *string
	PropertyAddressHouseNumber                        *string
	PropertyAddressStreetDirection                    *string
	PropertyAddressStreetName                         *string
	PropertyAddressStreetSuffix                       *string
	PropertyAddressStreetPostDirection                *string
	PropertyAddressUnitPrefix                         *string
	PropertyAddressUnitValue                          *string
	PropertyAddressCity                               *string
	PropertyAddressState                              *string
	PropertyAddressZIP                                *string
	PropertyAddressZIP4                               *string
	PropertyAddressCRRT                               *string
	PropertyAddressInfoFormat                         *string
	PropertyAddressInfoPrivacy                        *bool
	RecorderMapReference                              *string
	PropertyUseGroup                                  *string
	PropertyUseStandardized                           *string
	Mortgage1DocumentNumberFormatted                  *string
	Mortgage1DocumentNumberLegacy                     *string
	Mortgage1InstrumentNumber                         *string
	Mortgage1Book                                     *string
	Mortgage1Page                                     *string
	Mortgage1RecordingDate                            *time.Time
	Mortgage1Type                                     *string
	Mortgage1Amount                                   *int
	Mortgage1LenderCode                               *int
	Mortgage1LenderNameFullStandardized               *string
	Mortgage1LenderNameFirst                          *string
	Mortgage1LenderNameLast                           *string
	Mortgage1LenderAddress                            *string
	Mortgage1LenderAddressCity                        *string
	Mortgage1LenderAddressState                       *string
	Mortgage1LenderAddressZIP                         *string
	Mortgage1LenderAddressZIP4                        *string
	Mortgage1LenderInfoEntityClassification           *string
	Mortgage1LenderInfoSellerCarryBackFlag            *bool
	Mortgage1Term                                     *int
	Mortgage1TermType                                 *string
	Mortgage1TermDate                                 *time.Time
	Mortgage1InfoPrepaymentPenaltyFlag                *bool
	Mortgage1InfoPrepaymentTerm                       *string
	Mortgage1InterestRateType                         *string
	Mortgage1InterestRate                             *decimal.Decimal
	Mortgage1InterestTypeInitial                      *string
	Mortgage1FixedStepConversionRate                  *string
	Mortgage1DocumentInfoRiderAdjustableRateFlag      *bool
	Mortgage1InfoInterestTypeChangeYear               *int
	Mortgage1InfoInterestTypeChangeMonth              *int
	Mortgage1InfoInterestTypeChangeDay                *int
	Mortgage1InterestRateMinFirstChangeRateConversion *decimal.Decimal
	Mortgage1InterestRateMaxFirstChangeRateConversion *decimal.Decimal
	Mortgage1InterestChangeFrequency                  *string
	Mortgage1InterestMargin                           *int
	Mortgage1InterestIndex                            *decimal.Decimal
	Mortgage1InterestRateMax                          *decimal.Decimal
	Mortgage1AdjustableRateIndex                      *string
	Mortgage1InterestOnlyFlag                         *bool
	Mortgage1InterestOnlyPeriod                       *string
	Mortgage2DocumentNumberFormatted                  *string
	Mortgage2DocumentNumberLegacy                     *string
	Mortgage2InstrumentNumber                         *string
	Mortgage2Book                                     *string
	Mortgage2Page                                     *string
	Mortgage2RecordingDate                            *time.Time
	Mortgage2Type                                     *string
	Mortgage2Amount                                   *int
	Mortgage2LenderCode                               *int
	Mortgage2LenderNameFullStandardized               *string
	Mortgage2LenderNameFirst                          *string
	Mortgage2LenderNameLast                           *string
	Mortgage2LenderAddress                            *string
	Mortgage2LenderAddressCity                        *string
	Mortgage2LenderAddressState                       *string
	Mortgage2LenderAddressZIP                         *string
	Mortgage2LenderAddressZIP4                        *string
	Mortgage2LenderInfoEntityClassification           *string
	Mortgage2LenderInfoSellerCarryBackFlag            *bool
	Mortgage2Term                                     *int
	Mortgage2TermType                                 *string
	Mortgage2TermDate                                 *time.Time
	Mortgage2InfoPrepaymentPenaltyFlag                *bool
	Mortgage2InfoPrepaymentTerm                       *string
	Mortgage2InterestRateType                         *string
	Mortgage2InterestRate                             *decimal.Decimal
	Mortgage2InterestTypeInitial                      *string
	Mortgage2FixedStepConversionRate                  *string
	Mortgage2DocumentInfoRiderAdjustableRateFlag      *bool
	Mortgage2InfoInterestTypeChangeYear               *int
	Mortgage2InfoInterestTypeChangeMonth              *int
	Mortgage2InfoInterestTypeChangeDay                *int
	Mortgage2InterestRateMinFirstChangeRateConversion *decimal.Decimal
	Mortgage2InterestRateMaxFirstChangeRateConversion *decimal.Decimal
	Mortgage2InterestChangeFrequency                  *string
	Mortgage2InterestMargin                           *int
	Mortgage2InterestIndex                            *decimal.Decimal
	Mortgage2InterestRateMax                          *decimal.Decimal
	Mortgage2AdjustableRateIndex                      *string
	Mortgage2InterestOnlyFlag                         *bool
	Mortgage2InterestOnlyPeriod                       *string
	TransferInfoPurchaseDownPayment                   *int
	TransferInfoPurchaseLoanToValue                   *decimal.Decimal
	LastUpdated                                       *time.Time
	PublicationDate                                   *time.Time
}

func (dr *Recorder) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(Recorder)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "TransactionID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "97e255d7-90df-4e9d-a27a-35627f5609f0",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.TransactionID = v
		case "[ATTOM ID]":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "efc554f7-9a49-48e4-99a3-8781d0663a4f",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.ATTOMID = v
		case "DocumentRecordingStateCode":
			record.DocumentRecordingStateCode = val.StringPtrIfNonZero(field)
		case "DocumentRecordingCountyName":
			record.DocumentRecordingCountyName = val.StringPtrIfNonZero(field)
		case "DocumentRecordingJurisdictionName":
			record.DocumentRecordingJurisdictionName = val.StringPtrIfNonZero(field)
		case "DocumentRecordingCountyFIPs":
			record.DocumentRecordingCountyFIPs = val.StringPtrIfNonZero(field)
		case "DocumentTypeCode":
			record.DocumentTypeCode = val.StringPtrIfNonZero(field)
		case "DocumentNumberFormatted":
			record.DocumentNumberFormatted = val.StringPtrIfNonZero(field)
		case "DocumentNumberLegacy":
			record.DocumentNumberLegacy = val.StringPtrIfNonZero(field)
		case "InstrumentNumber":
			record.InstrumentNumber = val.StringPtrIfNonZero(field)
		case "Book":
			record.Book = val.StringPtrIfNonZero(field)
		case "Page":
			record.Page = val.StringPtrIfNonZero(field)
		case "InstrumentDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "8452f2fd-a260-48d1-a035-059a5b48ad11")
			}
			record.InstrumentDate = v
		case "RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "237f3be0-53f3-4b36-9acf-4c81822d749f")
			}
			record.RecordingDate = v
		case "TransactionType":
			record.TransactionType = val.StringPtrIfNonZero(field)
		case "TransferInfoPurchaseTypeCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e72d61da-b558-4a33-9a7f-e188c6cbd226")
			}
			record.TransferInfoPurchaseTypeCode = v
		case "ForeclosureAuctionSale":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "af350d6b-7ff0-4dd6-a13d-acdd6a42f53b")
			}
			record.ForeclosureAuctionSale = v
		case "TransferInfoDistressCircumstanceCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "833ea877-2693-4f69-b4b6-2bd942fcd5db")
			}
			record.TransferInfoDistressCircumstanceCode = v
		case "QuitclaimFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bb24dcb0-095a-459d-9ddc-05321ba194e7")
			}
			record.QuitclaimFlag = v
		case "TransferInfoMultiParcelFlag":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "43364ffc-7ed9-42e1-b701-5745c21b8971")
			}
			record.TransferInfoMultiParcelFlag = v
		case "ArmsLengthFlag":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4d6393c4-4475-40ff-9d37-c42bebe810fe")
			}
			record.ArmsLengthFlag = v
		case "PartialInterest":
			record.PartialInterest = val.StringPtrIfNonZero(field)
		case "TransferAmount":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "56d14e12-510f-45bb-a472-599ca1c8205b")
			}
			record.TransferAmount = v
		case "TransferAmountInfoAccuracy":
			record.TransferAmountInfoAccuracy = val.StringPtrIfNonZero(field)
		case "TransferTaxTotal":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bc9e05d2-ef02-4f22-b410-f66113498a0e")
			}
			record.TransferTaxTotal = v
		case "TransferTaxCity":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c11f84b0-d3de-4439-87f3-b0ea97d03445")
			}
			record.TransferTaxCity = v
		case "TransferTaxCounty":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "97a0c745-4870-4310-8323-67e13153316b")
			}
			record.TransferTaxCounty = v
		case "Grantor1NameFull":
			record.Grantor1NameFull = val.StringPtrIfNonZero(field)
		case "Grantor1NameFirst":
			record.Grantor1NameFirst = val.StringPtrIfNonZero(field)
		case "Grantor1NameMiddle":
			record.Grantor1NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantor1NameLast":
			record.Grantor1NameLast = val.StringPtrIfNonZero(field)
		case "Grantor1NameSuffix":
			record.Grantor1NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantor1InfoEntityClassification":
			record.Grantor1InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Grantor1InfoOwnerType":
			record.Grantor1InfoOwnerType = val.StringPtrIfNonZero(field)
		case "Grantor2NameFull":
			record.Grantor2NameFull = val.StringPtrIfNonZero(field)
		case "Grantor2NameFirst":
			record.Grantor2NameFirst = val.StringPtrIfNonZero(field)
		case "Grantor2NameMiddle":
			record.Grantor2NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantor2NameLast":
			record.Grantor2NameLast = val.StringPtrIfNonZero(field)
		case "Grantor2NameSuffix":
			record.Grantor2NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantor2InfoEntityClassification":
			record.Grantor2InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Grantor2InfoOwnerType":
			record.Grantor2InfoOwnerType = val.StringPtrIfNonZero(field)
		case "Grantor3NameFull":
			record.Grantor3NameFull = val.StringPtrIfNonZero(field)
		case "Grantor3NameFirst":
			record.Grantor3NameFirst = val.StringPtrIfNonZero(field)
		case "Grantor3NameMiddle":
			record.Grantor3NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantor3NameLast":
			record.Grantor3NameLast = val.StringPtrIfNonZero(field)
		case "Grantor3NameSuffix":
			record.Grantor3NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantor3InfoEntityClassification":
			record.Grantor3InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Grantor4NameFull":
			record.Grantor4NameFull = val.StringPtrIfNonZero(field)
		case "Grantor4NameFirst":
			record.Grantor4NameFirst = val.StringPtrIfNonZero(field)
		case "Grantor4NameMiddle":
			record.Grantor4NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantor4NameLast":
			record.Grantor4NameLast = val.StringPtrIfNonZero(field)
		case "Grantor4NameSuffix":
			record.Grantor4NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantor4InfoEntityClassification":
			record.Grantor4InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "GrantorAddressFull":
			record.GrantorAddressFull = val.StringPtrIfNonZero(field)
		case "GrantorAddressHouseNumber":
			record.GrantorAddressHouseNumber = val.StringPtrIfNonZero(field)
		case "GrantorAddressStreetDirection":
			record.GrantorAddressStreetDirection = val.StringPtrIfNonZero(field)
		case "GrantorAddressStreetName":
			record.GrantorAddressStreetName = val.StringPtrIfNonZero(field)
		case "GrantorAddressStreetSuffix":
			record.GrantorAddressStreetSuffix = val.StringPtrIfNonZero(field)
		case "GrantorAddressStreetPostDirection":
			record.GrantorAddressStreetPostDirection = val.StringPtrIfNonZero(field)
		case "GrantorAddressUnitPrefix":
			record.GrantorAddressUnitPrefix = val.StringPtrIfNonZero(field)
		case "GrantorAddressUnitValue":
			record.GrantorAddressUnitValue = val.StringPtrIfNonZero(field)
		case "GrantorAddressCity":
			record.GrantorAddressCity = val.StringPtrIfNonZero(field)
		case "GrantorAddressState":
			record.GrantorAddressState = val.StringPtrIfNonZero(field)
		case "GrantorAddressZIP":
			record.GrantorAddressZIP = val.StringPtrIfNonZero(field)
		case "GrantorAddressZIP4":
			record.GrantorAddressZIP4 = val.StringPtrIfNonZero(field)
		case "GrantorAddressCRRT":
			record.GrantorAddressCRRT = val.StringPtrIfNonZero(field)
		case "GrantorAddressInfoFormat":
			record.GrantorAddressInfoFormat = val.StringPtrIfNonZero(field)
		case "GrantorAddressInfoPrivacy":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "62fd0ac3-db18-498d-91e8-9d5bcd9b5340")
			}
			record.GrantorAddressInfoPrivacy = v
		case "Grantee1NameFull":
			record.Grantee1NameFull = val.StringPtrIfNonZero(field)
		case "Grantee1NameFirst":
			record.Grantee1NameFirst = val.StringPtrIfNonZero(field)
		case "Grantee1NameMiddle":
			record.Grantee1NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantee1NameLast":
			record.Grantee1NameLast = val.StringPtrIfNonZero(field)
		case "Grantee1NameSuffix":
			record.Grantee1NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantee1InfoEntityClassification":
			record.Grantee1InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Grantee1InfoOwnerType":
			record.Grantee1InfoOwnerType = val.StringPtrIfNonZero(field)
		case "Grantee2NameFull":
			record.Grantee2NameFull = val.StringPtrIfNonZero(field)
		case "Grantee2NameFirst":
			record.Grantee2NameFirst = val.StringPtrIfNonZero(field)
		case "Grantee2NameMiddle":
			record.Grantee2NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantee2NameLast":
			record.Grantee2NameLast = val.StringPtrIfNonZero(field)
		case "Grantee2NameSuffix":
			record.Grantee2NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantee2InfoEntityClassification":
			record.Grantee2InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "GranteeInfoVesting1":
			record.GranteeInfoVesting1 = val.StringPtrIfNonZero(field)
		case "Grantee3NameFull":
			record.Grantee3NameFull = val.StringPtrIfNonZero(field)
		case "Grantee3NameFirst":
			record.Grantee3NameFirst = val.StringPtrIfNonZero(field)
		case "Grantee3NameMiddle":
			record.Grantee3NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantee3NameLast":
			record.Grantee3NameLast = val.StringPtrIfNonZero(field)
		case "Grantee3NameSuffix":
			record.Grantee3NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantee3InfoEntityClassification":
			record.Grantee3InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Grantee4NameFull":
			record.Grantee4NameFull = val.StringPtrIfNonZero(field)
		case "Grantee4NameFirst":
			record.Grantee4NameFirst = val.StringPtrIfNonZero(field)
		case "Grantee4NameMiddle":
			record.Grantee4NameMiddle = val.StringPtrIfNonZero(field)
		case "Grantee4NameLast":
			record.Grantee4NameLast = val.StringPtrIfNonZero(field)
		case "Grantee4NameSuffix":
			record.Grantee4NameSuffix = val.StringPtrIfNonZero(field)
		case "Grantee4InfoEntityClassification":
			record.Grantee4InfoEntityClassification = val.StringPtrIfNonZero(field)
		case "GranteeMailCareOfName":
			record.GranteeMailCareOfName = val.StringPtrIfNonZero(field)
		case "GranteeInfoEntityCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "65671688-a876-4f29-a025-464db38bbe49")
			}
			record.GranteeInfoEntityCount = v
		case "GranteeInfoVesting2":
			record.GranteeInfoVesting2 = val.StringPtrIfNonZero(field)
		case "GranteeInvestorFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "afac76fa-2fc3-4db8-b486-4427c3da36b7")
			}
			record.GranteeInvestorFlag = v
		case "GranteeMailAddressFull":
			record.GranteeMailAddressFull = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressHouseNumber":
			record.GranteeMailAddressHouseNumber = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressStreetDirection":
			record.GranteeMailAddressStreetDirection = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressStreetName":
			record.GranteeMailAddressStreetName = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressStreetSuffix":
			record.GranteeMailAddressStreetSuffix = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressStreetPostDirection":
			record.GranteeMailAddressStreetPostDirection = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressUnitPrefix":
			record.GranteeMailAddressUnitPrefix = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressUnitValue":
			record.GranteeMailAddressUnitValue = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressCity":
			record.GranteeMailAddressCity = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressState":
			record.GranteeMailAddressState = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressZIP":
			record.GranteeMailAddressZIP = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressZIP4":
			record.GranteeMailAddressZIP4 = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressCRRT":
			record.GranteeMailAddressCRRT = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressInfoFormat":
			record.GranteeMailAddressInfoFormat = val.StringPtrIfNonZero(field)
		case "GranteeMailAddressInfoPrivacy":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "487f521f-8339-42e3-95cd-54f51e6b2249")
			}
			record.GranteeMailAddressInfoPrivacy = v
		case "GranteeGrantorOwnerRelationshipCode":
			record.GranteeGrantorOwnerRelationshipCode = val.StringPtrIfNonZero(field)
		case "TitleCompanyStandardizedCode":
			record.TitleCompanyStandardizedCode = val.StringPtrIfNonZero(field)
		case "TitleCompanyStandardizedName":
			record.TitleCompanyStandardizedName = val.StringPtrIfNonZero(field)
		case "TitleCompanyRaw":
			record.TitleCompanyRaw = val.StringPtrIfNonZero(field)
		case "LegalDescriptionPart1":
			record.LegalDescriptionPart1 = val.StringPtrIfNonZero(field)
		case "LegalDescriptionPart2":
			record.LegalDescriptionPart2 = val.StringPtrIfNonZero(field)
		case "LegalDescriptionPart3":
			record.LegalDescriptionPart3 = val.StringPtrIfNonZero(field)
		case "LegalDescriptionPart4":
			record.LegalDescriptionPart4 = val.StringPtrIfNonZero(field)
		case "LegalRange":
			record.LegalRange = val.StringPtrIfNonZero(field)
		case "LegalTownship":
			record.LegalTownship = val.StringPtrIfNonZero(field)
		case "LegalSection":
			record.LegalSection = val.StringPtrIfNonZero(field)
		case "LegalDistrict":
			record.LegalDistrict = val.StringPtrIfNonZero(field)
		case "LegalSubDivision":
			record.LegalSubDivision = val.StringPtrIfNonZero(field)
		case "LegalTract":
			record.LegalTract = val.StringPtrIfNonZero(field)
		case "LegalBlock":
			record.LegalBlock = val.StringPtrIfNonZero(field)
		case "LegalLot":
			record.LegalLot = val.StringPtrIfNonZero(field)
		case "LegalUnit":
			record.LegalUnit = val.StringPtrIfNonZero(field)
		case "LegalPlatMapBook":
			record.LegalPlatMapBook = val.StringPtrIfNonZero(field)
		case "LegalPlatMapPage":
			record.LegalPlatMapPage = val.StringPtrIfNonZero(field)
		case "APNFormatted":
			record.APNFormatted = val.StringPtrIfNonZero(field)
		case "APNOriginal":
			record.APNOriginal = val.StringPtrIfNonZero(field)
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
		case "PropertyAddressInfoFormat":
			record.PropertyAddressInfoFormat = val.StringPtrIfNonZero(field)
		case "PropertyAddressInfoPrivacy":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "44d69977-7ee7-4406-a761-03e6e0c28790")
			}
			record.PropertyAddressInfoPrivacy = v
		case "RecorderMapReference":
			record.RecorderMapReference = val.StringPtrIfNonZero(field)
		case "PropertyUseGroup":
			record.PropertyUseGroup = val.StringPtrIfNonZero(field)
		case "PropertyUseStandardized":
			record.PropertyUseStandardized = val.StringPtrIfNonZero(field)
		case "Mortgage1DocumentNumberFormatted":
			record.Mortgage1DocumentNumberFormatted = val.StringPtrIfNonZero(field)
		case "Mortgage1DocumentNumberLegacy":
			record.Mortgage1DocumentNumberLegacy = val.StringPtrIfNonZero(field)
		case "Mortgage1InstrumentNumber":
			record.Mortgage1InstrumentNumber = val.StringPtrIfNonZero(field)
		case "Mortgage1Book":
			record.Mortgage1Book = val.StringPtrIfNonZero(field)
		case "Mortgage1Page":
			record.Mortgage1Page = val.StringPtrIfNonZero(field)
		case "Mortgage1RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "9391c9ed-3dcd-4d39-95b6-40c7194ab8f7")
			}
			record.Mortgage1RecordingDate = v
		case "Mortgage1Type":
			record.Mortgage1Type = val.StringPtrIfNonZero(field)
		case "Mortgage1Amount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "0cf6ca63-97b7-49ed-90db-6c84d53e5449")
			}
			record.Mortgage1Amount = v
		case "Mortgage1LenderCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d8d1856d-8424-4f85-aba9-b08b5d975144")
			}
			record.Mortgage1LenderCode = v
		case "Mortgage1LenderNameFullStandardized":
			record.Mortgage1LenderNameFullStandardized = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderNameFirst":
			record.Mortgage1LenderNameFirst = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderNameLast":
			record.Mortgage1LenderNameLast = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderAddress":
			record.Mortgage1LenderAddress = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderAddressCity":
			record.Mortgage1LenderAddressCity = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderAddressState":
			record.Mortgage1LenderAddressState = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderAddressZIP":
			record.Mortgage1LenderAddressZIP = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderAddressZIP4":
			record.Mortgage1LenderAddressZIP4 = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderInfoEntityClassification":
			record.Mortgage1LenderInfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Mortgage1LenderInfoSellerCarryBackFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "568c6d67-50f7-48ba-ac73-f796195b5cad")
			}
			record.Mortgage1LenderInfoSellerCarryBackFlag = v
		case "Mortgage1Term":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f2c4724e-27e6-4616-aa46-5bfe29d1e065")
			}
			record.Mortgage1Term = v
		case "Mortgage1TermType":
			record.Mortgage1TermType = val.StringPtrIfNonZero(field)
		case "Mortgage1TermDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "379afbcd-c2ce-4eec-8ecc-af46efb5c2d6")
			}
			record.Mortgage1TermDate = v
		case "Mortgage1InfoPrepaymentPenaltyFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "7c04b9ed-6afa-426d-a8f9-738c125e0a8c")
			}
			record.Mortgage1InfoPrepaymentPenaltyFlag = v
		case "Mortgage1InfoPrepaymentTerm":
			record.Mortgage1InfoPrepaymentTerm = val.StringPtrIfNonZero(field)
		case "Mortgage1InterestRateType":
			record.Mortgage1InterestRateType = val.StringPtrIfNonZero(field)
		case "Mortgage1InterestRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "81e98e56-27e0-483a-ab23-a6bd7768dac4")
			}
			record.Mortgage1InterestRate = v
		case "Mortgage1InterestTypeInitial":
			record.Mortgage1InterestTypeInitial = val.StringPtrIfNonZero(field)
		case "Mortgage1FixedStepConversionRate":
			record.Mortgage1FixedStepConversionRate = val.StringPtrIfNonZero(field)
		case "Mortgage1DocumentInfoRiderAdjustableRateFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "d978f534-2980-48fd-a60f-d536783ed807")
			}
			record.Mortgage1DocumentInfoRiderAdjustableRateFlag = v
		case "Mortgage1InfoInterestTypeChangeYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "a535bcb4-4a39-4171-9c5d-d8f3fa9ceca9")
			}
			record.Mortgage1InfoInterestTypeChangeYear = v
		case "Mortgage1InfoInterestTypeChangeMonth":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "f15da057-67a0-438f-a109-23b1d976ef89")
			}
			record.Mortgage1InfoInterestTypeChangeMonth = v
		case "Mortgage1InfoInterestTypeChangeDay":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2c3f38a9-d5c9-477a-aef5-da316d3eb39f")
			}
			record.Mortgage1InfoInterestTypeChangeDay = v
		case "Mortgage1InterestRateMinFirstChangeRateConversion":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1a3f337d-ccd2-4b5a-b4f6-ccff0517c321")
			}
			record.Mortgage1InterestRateMinFirstChangeRateConversion = v
		case "Mortgage1InterestRateMaxFirstChangeRateConversion":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "84788b68-081a-4de9-ad3a-acca6e83fb5a")
			}
			record.Mortgage1InterestRateMaxFirstChangeRateConversion = v
		case "Mortgage1InterestChangeFrequency":
			record.Mortgage1InterestChangeFrequency = val.StringPtrIfNonZero(field)
		case "Mortgage1InterestMargin":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "77b108f9-8234-4e56-b4f3-44940f646065")
			}
			record.Mortgage1InterestMargin = v
		case "Mortgage1InterestIndex":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "542d1b47-02fe-47fe-ad2d-8f4de5f2b57a")
			}
			record.Mortgage1InterestIndex = v
		case "Mortgage1InterestRateMax":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4884581b-e885-4cdd-a88a-6f35568b18c7")
			}
			record.Mortgage1InterestRateMax = v
		case "Mortgage1AdjustableRateIndex":
			record.Mortgage1AdjustableRateIndex = val.StringPtrIfNonZero(field)
		case "Mortgage1InterestOnlyFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "11516eb8-89e5-4795-a09d-fd47214932bf")
			}
			record.Mortgage1InterestOnlyFlag = v
		case "Mortgage1InterestOnlyPeriod":
			record.Mortgage1InterestOnlyPeriod = val.StringPtrIfNonZero(field)
		case "Mortgage2DocumentNumberFormatted":
			record.Mortgage2DocumentNumberFormatted = val.StringPtrIfNonZero(field)
		case "Mortgage2DocumentNumberLegacy":
			record.Mortgage2DocumentNumberLegacy = val.StringPtrIfNonZero(field)
		case "Mortgage2InstrumentNumber":
			record.Mortgage2InstrumentNumber = val.StringPtrIfNonZero(field)
		case "Mortgage2Book":
			record.Mortgage2Book = val.StringPtrIfNonZero(field)
		case "Mortgage2Page":
			record.Mortgage2Page = val.StringPtrIfNonZero(field)
		case "Mortgage2RecordingDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "d507ff5a-b36d-4cc3-ab83-cb8d25b8b451")
			}
			record.Mortgage2RecordingDate = v
		case "Mortgage2Type":
			record.Mortgage2Type = val.StringPtrIfNonZero(field)
		case "Mortgage2Amount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "06ba17b4-1de0-4d2f-8ed9-eb50ad5f76b2")
			}
			record.Mortgage2Amount = v
		case "Mortgage2LenderCode":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5ff18982-5099-4a05-bc28-bf123fe4763a")
			}
			record.Mortgage2LenderCode = v
		case "Mortgage2LenderNameFullStandardized":
			record.Mortgage2LenderNameFullStandardized = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderNameFirst":
			record.Mortgage2LenderNameFirst = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderNameLast":
			record.Mortgage2LenderNameLast = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderAddress":
			record.Mortgage2LenderAddress = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderAddressCity":
			record.Mortgage2LenderAddressCity = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderAddressState":
			record.Mortgage2LenderAddressState = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderAddressZIP":
			record.Mortgage2LenderAddressZIP = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderAddressZIP4":
			record.Mortgage2LenderAddressZIP4 = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderInfoEntityClassification":
			record.Mortgage2LenderInfoEntityClassification = val.StringPtrIfNonZero(field)
		case "Mortgage2LenderInfoSellerCarryBackFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b37635b4-1e34-4e7a-aa7c-850b2bff08d1")
			}
			record.Mortgage2LenderInfoSellerCarryBackFlag = v
		case "Mortgage2Term":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "686b38ba-d865-4797-9b74-a89284e38bbd")
			}
			record.Mortgage2Term = v
		case "Mortgage2TermType":
			record.Mortgage2TermType = val.StringPtrIfNonZero(field)
		case "Mortgage2TermDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "a29ea865-f10e-47fe-9437-cf3001695925")
			}
			record.Mortgage2TermDate = v
		case "Mortgage2InfoPrepaymentPenaltyFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "533119f6-e6ac-461b-95fd-db3cc5e27d2c")
			}
			record.Mortgage2InfoPrepaymentPenaltyFlag = v
		case "Mortgage2InfoPrepaymentTerm":
			record.Mortgage2InfoPrepaymentTerm = val.StringPtrIfNonZero(field)
		case "Mortgage2InterestRateType":
			record.Mortgage2InterestRateType = val.StringPtrIfNonZero(field)
		case "Mortgage2InterestRate":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "61db52be-f20d-4c5d-a5c3-cdb003a45164")
			}
			record.Mortgage2InterestRate = v
		case "Mortgage2InterestTypeInitial":
			record.Mortgage2InterestTypeInitial = val.StringPtrIfNonZero(field)
		case "Mortgage2FixedStepConversionRate":
			record.Mortgage2FixedStepConversionRate = val.StringPtrIfNonZero(field)
		case "Mortgage2DocumentInfoRiderAdjustableRateFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "5a9d39ea-ef6e-4f50-8d26-44459ec00b59")
			}
			record.Mortgage2DocumentInfoRiderAdjustableRateFlag = v
		case "Mortgage2InfoInterestTypeChangeYear":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9b8a28ee-b1b6-43c9-9413-fae4285e6f94")
			}
			record.Mortgage2InfoInterestTypeChangeYear = v
		case "Mortgage2InfoInterestTypeChangeMonth":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "61c1b971-c51e-49c3-aca6-3cc4a5227100")
			}
			record.Mortgage2InfoInterestTypeChangeMonth = v
		case "Mortgage2InfoInterestTypeChangeDay":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "52b7f951-772d-4ea7-887d-eaf33321e23c")
			}
			record.Mortgage2InfoInterestTypeChangeDay = v
		case "Mortgage2InterestRateMinFirstChangeRateConversion":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b14fee14-a6e8-497f-acca-64a37893372e")
			}
			record.Mortgage2InterestRateMinFirstChangeRateConversion = v
		case "Mortgage2InterestRateMaxFirstChangeRateConversion":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "87c77a7a-2dd0-4c9c-a0ed-ff58b2fae5fc")
			}
			record.Mortgage2InterestRateMaxFirstChangeRateConversion = v
		case "Mortgage2InterestChangeFrequency":
			record.Mortgage2InterestChangeFrequency = val.StringPtrIfNonZero(field)
		case "Mortgage2InterestMargin":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "df303287-4dc0-4550-bfb8-33d8a7f551c8")
			}
			record.Mortgage2InterestMargin = v
		case "Mortgage2InterestIndex":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "614c1513-22da-42ee-948e-520b727f542e")
			}
			record.Mortgage2InterestIndex = v
		case "Mortgage2InterestRateMax":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "120b22ae-95af-4fc0-829c-576c63b726e9")
			}
			record.Mortgage2InterestRateMax = v
		case "Mortgage2AdjustableRateIndex":
			record.Mortgage2AdjustableRateIndex = val.StringPtrIfNonZero(field)
		case "Mortgage2InterestOnlyFlag":
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3d96e619-7063-41fe-873a-23f180a0e399")
			}
			record.Mortgage2InterestOnlyFlag = v
		case "Mortgage2InterestOnlyPeriod":
			record.Mortgage2InterestOnlyPeriod = val.StringPtrIfNonZero(field)
		case "TransferInfoPurchaseDownPayment":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b16fb44e-e7fa-4a2a-a816-007c2d77c342")
			}
			record.TransferInfoPurchaseDownPayment = v
		case "TransferInfoPurchaseLoanToValue":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "023bca31-4c01-457e-9008-bf8a2ce62873")
			}
			record.TransferInfoPurchaseLoanToValue = v
		case "LastUpdated":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "1654911c-51d3-4642-a5ed-9430b5062398")
			}
			record.LastUpdated = v
		case "PublicationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "b0ad4f35-aa6a-45b3-9e38-bdc6a0b5b754")
			}
			record.PublicationDate = v
		default:
			return nil, &errors.Object{
				Id:     "6f3e487a-c4ea-4af9-8e72-1ca13d2304d8",
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

func (dr *Recorder) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"transaction_id",
		"attomid",
		"document_recording_state_code",
		"document_recording_county_name",
		"document_recording_jurisdiction_name",
		"document_recording_county_fi_ps",
		"document_type_code",
		"document_number_formatted",
		"document_number_legacy",
		"instrument_number",
		"book",
		"page",
		"instrument_date",
		"recording_date",
		"transaction_type",
		"transfer_info_purchase_type_code",
		"foreclosure_auction_sale",
		"transfer_info_distress_circumstance_code",
		"quitclaim_flag",
		"transfer_info_multi_parcel_flag",
		"arms_length_flag",
		"partial_interest",
		"transfer_amount",
		"transfer_amount_info_accuracy",
		"transfer_tax_total",
		"transfer_tax_city",
		"transfer_tax_county",
		"grantor1name_full",
		"grantor1name_first",
		"grantor1name_middle",
		"grantor1name_last",
		"grantor1name_suffix",
		"grantor1info_entity_classification",
		"grantor1info_owner_type",
		"grantor2name_full",
		"grantor2name_first",
		"grantor2name_middle",
		"grantor2name_last",
		"grantor2name_suffix",
		"grantor2info_entity_classification",
		"grantor2info_owner_type",
		"grantor3name_full",
		"grantor3name_first",
		"grantor3name_middle",
		"grantor3name_last",
		"grantor3name_suffix",
		"grantor3info_entity_classification",
		"grantor4name_full",
		"grantor4name_first",
		"grantor4name_middle",
		"grantor4name_last",
		"grantor4name_suffix",
		"grantor4info_entity_classification",
		"grantor_address_full",
		"grantor_address_house_number",
		"grantor_address_street_direction",
		"grantor_address_street_name",
		"grantor_address_street_suffix",
		"grantor_address_street_post_direction",
		"grantor_address_unit_prefix",
		"grantor_address_unit_value",
		"grantor_address_city",
		"grantor_address_state",
		"grantor_address_zip",
		"grantor_address_zip4",
		"grantor_address_crrt",
		"grantor_address_info_format",
		"grantor_address_info_privacy",
		"grantee1name_full",
		"grantee1name_first",
		"grantee1name_middle",
		"grantee1name_last",
		"grantee1name_suffix",
		"grantee1info_entity_classification",
		"grantee1info_owner_type",
		"grantee2name_full",
		"grantee2name_first",
		"grantee2name_middle",
		"grantee2name_last",
		"grantee2name_suffix",
		"grantee2info_entity_classification",
		"grantee_info_vesting1",
		"grantee3name_full",
		"grantee3name_first",
		"grantee3name_middle",
		"grantee3name_last",
		"grantee3name_suffix",
		"grantee3info_entity_classification",
		"grantee4name_full",
		"grantee4name_first",
		"grantee4name_middle",
		"grantee4name_last",
		"grantee4name_suffix",
		"grantee4info_entity_classification",
		"grantee_mail_care_of_name",
		"grantee_info_entity_count",
		"grantee_info_vesting2",
		"grantee_investor_flag",
		"grantee_mail_address_full",
		"grantee_mail_address_house_number",
		"grantee_mail_address_street_direction",
		"grantee_mail_address_street_name",
		"grantee_mail_address_street_suffix",
		"grantee_mail_address_street_post_direction",
		"grantee_mail_address_unit_prefix",
		"grantee_mail_address_unit_value",
		"grantee_mail_address_city",
		"grantee_mail_address_state",
		"grantee_mail_address_zip",
		"grantee_mail_address_zip4",
		"grantee_mail_address_crrt",
		"grantee_mail_address_info_format",
		"grantee_mail_address_info_privacy",
		"grantee_grantor_owner_relationship_code",
		"title_company_standardized_code",
		"title_company_standardized_name",
		"title_company_raw",
		"legal_description_part1",
		"legal_description_part2",
		"legal_description_part3",
		"legal_description_part4",
		"legal_range",
		"legal_township",
		"legal_section",
		"legal_district",
		"legal_sub_division",
		"legal_tract",
		"legal_block",
		"legal_lot",
		"legal_unit",
		"legal_plat_map_book",
		"legal_plat_map_page",
		"apn_formatted",
		"apn_original",
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
		"property_address_info_format",
		"property_address_info_privacy",
		"recorder_map_reference",
		"property_use_group",
		"property_use_standardized",
		"mortgage1document_number_formatted",
		"mortgage1document_number_legacy",
		"mortgage1instrument_number",
		"mortgage1book",
		"mortgage1page",
		"mortgage1recording_date",
		"mortgage1type",
		"mortgage1amount",
		"mortgage1lender_code",
		"mortgage1lender_name_full_standardized",
		"mortgage1lender_name_first",
		"mortgage1lender_name_last",
		"mortgage1lender_address",
		"mortgage1lender_address_city",
		"mortgage1lender_address_state",
		"mortgage1lender_address_zip",
		"mortgage1lender_address_zip4",
		"mortgage1lender_info_entity_classification",
		"mortgage1lender_info_seller_carry_back_flag",
		"mortgage1term",
		"mortgage1term_type",
		"mortgage1term_date",
		"mortgage1info_prepayment_penalty_flag",
		"mortgage1info_prepayment_term",
		"mortgage1interest_rate_type",
		"mortgage1interest_rate",
		"mortgage1interest_type_initial",
		"mortgage1fixed_step_conversion_rate",
		"mortgage1document_info_rider_adjustable_rate_flag",
		"mortgage1info_interest_type_change_year",
		"mortgage1info_interest_type_change_month",
		"mortgage1info_interest_type_change_day",
		"mortgage1interest_rate_min_first_change_rate_conversion",
		"mortgage1interest_rate_max_first_change_rate_conversion",
		"mortgage1interest_change_frequency",
		"mortgage1interest_margin",
		"mortgage1interest_index",
		"mortgage1interest_rate_max",
		"mortgage1adjustable_rate_index",
		"mortgage1interest_only_flag",
		"mortgage1interest_only_period",
		"mortgage2document_number_formatted",
		"mortgage2document_number_legacy",
		"mortgage2instrument_number",
		"mortgage2book",
		"mortgage2page",
		"mortgage2recording_date",
		"mortgage2type",
		"mortgage2amount",
		"mortgage2lender_code",
		"mortgage2lender_name_full_standardized",
		"mortgage2lender_name_first",
		"mortgage2lender_name_last",
		"mortgage2lender_address",
		"mortgage2lender_address_city",
		"mortgage2lender_address_state",
		"mortgage2lender_address_zip",
		"mortgage2lender_address_zip4",
		"mortgage2lender_info_entity_classification",
		"mortgage2lender_info_seller_carry_back_flag",
		"mortgage2term",
		"mortgage2term_type",
		"mortgage2term_date",
		"mortgage2info_prepayment_penalty_flag",
		"mortgage2info_prepayment_term",
		"mortgage2interest_rate_type",
		"mortgage2interest_rate",
		"mortgage2interest_type_initial",
		"mortgage2fixed_step_conversion_rate",
		"mortgage2document_info_rider_adjustable_rate_flag",
		"mortgage2info_interest_type_change_year",
		"mortgage2info_interest_type_change_month",
		"mortgage2info_interest_type_change_day",
		"mortgage2interest_rate_min_first_change_rate_conversion",
		"mortgage2interest_rate_max_first_change_rate_conversion",
		"mortgage2interest_change_frequency",
		"mortgage2interest_margin",
		"mortgage2interest_index",
		"mortgage2interest_rate_max",
		"mortgage2adjustable_rate_index",
		"mortgage2interest_only_flag",
		"mortgage2interest_only_period",
		"transfer_info_purchase_down_payment",
		"transfer_info_purchase_loan_to_value",
		"last_updated",
		"publication_date",
	}
}

func (dr *Recorder) SQLTable() string {
	return "ad_df_recorder"
}

func (dr *Recorder) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "aac2b658-c74a-48b8-9192-96a689e04cfe",
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
		dr.TransactionID,
		dr.ATTOMID,
		dr.DocumentRecordingStateCode,
		dr.DocumentRecordingCountyName,
		dr.DocumentRecordingJurisdictionName,
		dr.DocumentRecordingCountyFIPs,
		dr.DocumentTypeCode,
		dr.DocumentNumberFormatted,
		dr.DocumentNumberLegacy,
		dr.InstrumentNumber,
		dr.Book,
		dr.Page,
		dr.InstrumentDate,
		dr.RecordingDate,
		dr.TransactionType,
		dr.TransferInfoPurchaseTypeCode,
		dr.ForeclosureAuctionSale,
		dr.TransferInfoDistressCircumstanceCode,
		dr.QuitclaimFlag,
		dr.TransferInfoMultiParcelFlag,
		dr.ArmsLengthFlag,
		dr.PartialInterest,
		dr.TransferAmount,
		dr.TransferAmountInfoAccuracy,
		dr.TransferTaxTotal,
		dr.TransferTaxCity,
		dr.TransferTaxCounty,
		dr.Grantor1NameFull,
		dr.Grantor1NameFirst,
		dr.Grantor1NameMiddle,
		dr.Grantor1NameLast,
		dr.Grantor1NameSuffix,
		dr.Grantor1InfoEntityClassification,
		dr.Grantor1InfoOwnerType,
		dr.Grantor2NameFull,
		dr.Grantor2NameFirst,
		dr.Grantor2NameMiddle,
		dr.Grantor2NameLast,
		dr.Grantor2NameSuffix,
		dr.Grantor2InfoEntityClassification,
		dr.Grantor2InfoOwnerType,
		dr.Grantor3NameFull,
		dr.Grantor3NameFirst,
		dr.Grantor3NameMiddle,
		dr.Grantor3NameLast,
		dr.Grantor3NameSuffix,
		dr.Grantor3InfoEntityClassification,
		dr.Grantor4NameFull,
		dr.Grantor4NameFirst,
		dr.Grantor4NameMiddle,
		dr.Grantor4NameLast,
		dr.Grantor4NameSuffix,
		dr.Grantor4InfoEntityClassification,
		dr.GrantorAddressFull,
		dr.GrantorAddressHouseNumber,
		dr.GrantorAddressStreetDirection,
		dr.GrantorAddressStreetName,
		dr.GrantorAddressStreetSuffix,
		dr.GrantorAddressStreetPostDirection,
		dr.GrantorAddressUnitPrefix,
		dr.GrantorAddressUnitValue,
		dr.GrantorAddressCity,
		dr.GrantorAddressState,
		dr.GrantorAddressZIP,
		dr.GrantorAddressZIP4,
		dr.GrantorAddressCRRT,
		dr.GrantorAddressInfoFormat,
		dr.GrantorAddressInfoPrivacy,
		dr.Grantee1NameFull,
		dr.Grantee1NameFirst,
		dr.Grantee1NameMiddle,
		dr.Grantee1NameLast,
		dr.Grantee1NameSuffix,
		dr.Grantee1InfoEntityClassification,
		dr.Grantee1InfoOwnerType,
		dr.Grantee2NameFull,
		dr.Grantee2NameFirst,
		dr.Grantee2NameMiddle,
		dr.Grantee2NameLast,
		dr.Grantee2NameSuffix,
		dr.Grantee2InfoEntityClassification,
		dr.GranteeInfoVesting1,
		dr.Grantee3NameFull,
		dr.Grantee3NameFirst,
		dr.Grantee3NameMiddle,
		dr.Grantee3NameLast,
		dr.Grantee3NameSuffix,
		dr.Grantee3InfoEntityClassification,
		dr.Grantee4NameFull,
		dr.Grantee4NameFirst,
		dr.Grantee4NameMiddle,
		dr.Grantee4NameLast,
		dr.Grantee4NameSuffix,
		dr.Grantee4InfoEntityClassification,
		dr.GranteeMailCareOfName,
		dr.GranteeInfoEntityCount,
		dr.GranteeInfoVesting2,
		dr.GranteeInvestorFlag,
		dr.GranteeMailAddressFull,
		dr.GranteeMailAddressHouseNumber,
		dr.GranteeMailAddressStreetDirection,
		dr.GranteeMailAddressStreetName,
		dr.GranteeMailAddressStreetSuffix,
		dr.GranteeMailAddressStreetPostDirection,
		dr.GranteeMailAddressUnitPrefix,
		dr.GranteeMailAddressUnitValue,
		dr.GranteeMailAddressCity,
		dr.GranteeMailAddressState,
		dr.GranteeMailAddressZIP,
		dr.GranteeMailAddressZIP4,
		dr.GranteeMailAddressCRRT,
		dr.GranteeMailAddressInfoFormat,
		dr.GranteeMailAddressInfoPrivacy,
		dr.GranteeGrantorOwnerRelationshipCode,
		dr.TitleCompanyStandardizedCode,
		dr.TitleCompanyStandardizedName,
		dr.TitleCompanyRaw,
		dr.LegalDescriptionPart1,
		dr.LegalDescriptionPart2,
		dr.LegalDescriptionPart3,
		dr.LegalDescriptionPart4,
		dr.LegalRange,
		dr.LegalTownship,
		dr.LegalSection,
		dr.LegalDistrict,
		dr.LegalSubDivision,
		dr.LegalTract,
		dr.LegalBlock,
		dr.LegalLot,
		dr.LegalUnit,
		dr.LegalPlatMapBook,
		dr.LegalPlatMapPage,
		dr.APNFormatted,
		dr.APNOriginal,
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
		dr.PropertyAddressInfoFormat,
		dr.PropertyAddressInfoPrivacy,
		dr.RecorderMapReference,
		dr.PropertyUseGroup,
		dr.PropertyUseStandardized,
		dr.Mortgage1DocumentNumberFormatted,
		dr.Mortgage1DocumentNumberLegacy,
		dr.Mortgage1InstrumentNumber,
		dr.Mortgage1Book,
		dr.Mortgage1Page,
		dr.Mortgage1RecordingDate,
		dr.Mortgage1Type,
		dr.Mortgage1Amount,
		dr.Mortgage1LenderCode,
		dr.Mortgage1LenderNameFullStandardized,
		dr.Mortgage1LenderNameFirst,
		dr.Mortgage1LenderNameLast,
		dr.Mortgage1LenderAddress,
		dr.Mortgage1LenderAddressCity,
		dr.Mortgage1LenderAddressState,
		dr.Mortgage1LenderAddressZIP,
		dr.Mortgage1LenderAddressZIP4,
		dr.Mortgage1LenderInfoEntityClassification,
		dr.Mortgage1LenderInfoSellerCarryBackFlag,
		dr.Mortgage1Term,
		dr.Mortgage1TermType,
		dr.Mortgage1TermDate,
		dr.Mortgage1InfoPrepaymentPenaltyFlag,
		dr.Mortgage1InfoPrepaymentTerm,
		dr.Mortgage1InterestRateType,
		dr.Mortgage1InterestRate,
		dr.Mortgage1InterestTypeInitial,
		dr.Mortgage1FixedStepConversionRate,
		dr.Mortgage1DocumentInfoRiderAdjustableRateFlag,
		dr.Mortgage1InfoInterestTypeChangeYear,
		dr.Mortgage1InfoInterestTypeChangeMonth,
		dr.Mortgage1InfoInterestTypeChangeDay,
		dr.Mortgage1InterestRateMinFirstChangeRateConversion,
		dr.Mortgage1InterestRateMaxFirstChangeRateConversion,
		dr.Mortgage1InterestChangeFrequency,
		dr.Mortgage1InterestMargin,
		dr.Mortgage1InterestIndex,
		dr.Mortgage1InterestRateMax,
		dr.Mortgage1AdjustableRateIndex,
		dr.Mortgage1InterestOnlyFlag,
		dr.Mortgage1InterestOnlyPeriod,
		dr.Mortgage2DocumentNumberFormatted,
		dr.Mortgage2DocumentNumberLegacy,
		dr.Mortgage2InstrumentNumber,
		dr.Mortgage2Book,
		dr.Mortgage2Page,
		dr.Mortgage2RecordingDate,
		dr.Mortgage2Type,
		dr.Mortgage2Amount,
		dr.Mortgage2LenderCode,
		dr.Mortgage2LenderNameFullStandardized,
		dr.Mortgage2LenderNameFirst,
		dr.Mortgage2LenderNameLast,
		dr.Mortgage2LenderAddress,
		dr.Mortgage2LenderAddressCity,
		dr.Mortgage2LenderAddressState,
		dr.Mortgage2LenderAddressZIP,
		dr.Mortgage2LenderAddressZIP4,
		dr.Mortgage2LenderInfoEntityClassification,
		dr.Mortgage2LenderInfoSellerCarryBackFlag,
		dr.Mortgage2Term,
		dr.Mortgage2TermType,
		dr.Mortgage2TermDate,
		dr.Mortgage2InfoPrepaymentPenaltyFlag,
		dr.Mortgage2InfoPrepaymentTerm,
		dr.Mortgage2InterestRateType,
		dr.Mortgage2InterestRate,
		dr.Mortgage2InterestTypeInitial,
		dr.Mortgage2FixedStepConversionRate,
		dr.Mortgage2DocumentInfoRiderAdjustableRateFlag,
		dr.Mortgage2InfoInterestTypeChangeYear,
		dr.Mortgage2InfoInterestTypeChangeMonth,
		dr.Mortgage2InfoInterestTypeChangeDay,
		dr.Mortgage2InterestRateMinFirstChangeRateConversion,
		dr.Mortgage2InterestRateMaxFirstChangeRateConversion,
		dr.Mortgage2InterestChangeFrequency,
		dr.Mortgage2InterestMargin,
		dr.Mortgage2InterestIndex,
		dr.Mortgage2InterestRateMax,
		dr.Mortgage2AdjustableRateIndex,
		dr.Mortgage2InterestOnlyFlag,
		dr.Mortgage2InterestOnlyPeriod,
		dr.TransferInfoPurchaseDownPayment,
		dr.TransferInfoPurchaseLoanToValue,
		dr.LastUpdated,
		dr.PublicationDate,
	}

	return values, nil
}

func (dr *Recorder) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		Mode: entities.DataRecordModeBatchInsert,
	}
}
