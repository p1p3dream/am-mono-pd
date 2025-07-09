package first_american

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type AVMPower struct {
	AMId        uuid.UUID      `json:"am_id"`
	AMCreatedAt time.Time      `json:"am_created_at"`
	AMUpdatedAt time.Time      `json:"am_updated_at,omitempty"`
	AMMeta      map[string]any `json:"am_meta,omitempty"`

	Fips                   string
	PropertyID             int64
	APN                    string
	SitusFullStreetAddress *string
	SitusHouseNbr          *string
	SitusHouseNbrSuffix    *string
	SitusDirectionLeft     *string
	SitusStreet            *string
	SitusMode              *string
	SitusDirectionRight    *string
	SitusUnitType          *string
	SitusUnitNbr           *string
	SitusCity              *string
	SitusState             *string
	SitusZIP5              *string
	SitusZIP4              *string
	SitusCarrierCode       *string
	FinalValue             *decimal.Decimal
	HighValue              *decimal.Decimal
	LowValue               *decimal.Decimal
	ConfidenceScore        *float64
	StandardDeviation      *float64
	ValuationDate          *time.Time
	Comp1PropertyID        *int64
	Comp2PropertyID        *int64
	Comp3PropertyID        *int64
	Comp4PropertyID        *int64
	Comp5PropertyID        *int64
	Comp6PropertyID        *int64
	Comp7PropertyID        *int64
}

func (dr *AVMPower) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(AVMPower)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "Fips":
			if field == "" {
				return nil, &errors.Object{
					Id:     "49137f44-a998-4899-ac01-1ddee4e37144",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Fips is required.",
				}
			}
			record.Fips = field
		case "PropertyID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "0f366e8e-e5cf-4e87-ac51-d2cfd645c998",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": v,
					},
				}
			}
			record.PropertyID = v
		case "APN":
			if field == "" {
				return nil, &errors.Object{
					Id:     "82bf1474-57c0-491e-bf9b-927ff4a259c6",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "APN is required.",
				}
			}
			record.APN = field
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
		case "FinalValue":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "b4527ee3-8f5d-4b2b-99bc-60dafc718750")
			}
			record.FinalValue = v
		case "HighValue":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bfa558c1-4671-4156-91d0-79e4c6901bd8")
			}
			record.HighValue = v
		case "LowValue":
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "1ff1a7db-a40a-460e-86f2-7251b584a6b0")
			}
			record.LowValue = v
		case "ConfidenceScore":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bac9f6cf-bd6f-4d82-9448-5818b9d467b4")
			}
			record.ConfidenceScore = v
		case "StandardDeviation":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ba042650-5ad8-4ca8-8e5b-7257d28aa733")
			}
			record.StandardDeviation = v
		case "ValuationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "c1e860fa-b4f2-4513-9366-4a6b7106c4a1")
			}
			record.ValuationDate = v
		case "Comp1PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "16756c76-07fe-4851-8fbc-8d7aeef5e187")
			}
			record.Comp1PropertyID = v
		case "Comp2PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "399043e7-d4b9-4f13-a6fb-891c0fa68100")
			}
			record.Comp2PropertyID = v
		case "Comp3PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "550fe721-c40c-4482-8beb-df6f7f39802d")
			}
			record.Comp3PropertyID = v
		case "Comp4PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "e1ce80be-b0c2-47ee-b9fe-1c259df0f809")
			}
			record.Comp4PropertyID = v
		case "Comp5PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "ce09b9d4-0877-446e-a627-3e3558b272b6")
			}
			record.Comp5PropertyID = v
		case "Comp6PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "6e4e9843-cb57-4207-b5c1-b61f9c5c66be")
			}
			record.Comp6PropertyID = v
		case "Comp7PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "21d2f182-2fe4-476e-9bbd-68feada576d2")
			}
			record.Comp7PropertyID = v
		default:
			return nil, &errors.Object{
				Id:     "3c7b1c73-507b-4694-b5e8-62251a392ec9",
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

func (dr *AVMPower) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"fips",
		"property_id",
		"apn",
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
		"zip5",
		"zip4",
		"situs_carrier_code",
		"final_value",
		"high_value",
		"low_value",
		"confidence_score",
		"standard_deviation",
		"valuation_date",
		"comp1_property_id",
		"comp2_property_id",
		"comp3_property_id",
		"comp4_property_id",
		"comp5_property_id",
		"comp6_property_id",
		"comp7_property_id",
	}
}

func (dr *AVMPower) SQLTable() string {
	return "fa_df_avm_power"
}

func (dr *AVMPower) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "98572dcd-4b2e-44bc-ba38-e1024d5fbd65",
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
		dr.Fips,
		dr.PropertyID,
		dr.APN,
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
		dr.FinalValue,
		dr.HighValue,
		dr.LowValue,
		dr.ConfidenceScore,
		dr.StandardDeviation,
		dr.ValuationDate,
		dr.Comp1PropertyID,
		dr.Comp2PropertyID,
		dr.Comp3PropertyID,
		dr.Comp4PropertyID,
		dr.Comp5PropertyID,
		dr.Comp6PropertyID,
		dr.Comp7PropertyID,
	}

	return values, nil
}

func (dr *AVMPower) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		Mode: entities.DataRecordModeBatchInsert,
	}
}
