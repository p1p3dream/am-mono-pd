package first_american

import (
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type Address struct {
	AMId        uuid.UUID      `json:"am_id"`
	AMCreatedAt time.Time      `json:"am_created_at"`
	AMUpdatedAt time.Time      `json:"am_updated_at,omitempty"`
	AMMeta      map[string]any `json:"am_meta,omitempty"`

	FIPS                    string     `json:"fips"`
	State                   *string    `json:"state"`
	County                  *string    `json:"county"`
	ZIP5                    string     `json:"zip5"`
	ZIP4                    *string    `json:"zip4"`
	FullStreetAddress       *string    `json:"full_street_address"`
	PreDirectional          *string    `json:"pre_directional"`
	StreetNumber            *string    `json:"street_number"`
	Street                  *string    `json:"street"`
	PostDirectional         *string    `json:"post_directional"`
	StreetType              *string    `json:"street_type"`
	UnitType                *string    `json:"unit_type"`
	UnitNbr                 *string    `json:"unit_nbr"`
	VacantIndicator         *string    `json:"vacant_indicator"`
	NonUSPSAddressIndicator *string    `json:"non_usps_address_indicator"`
	NotCurrentlyDeliverable *string    `json:"not_currently_deliverable"`
	CommunityName           *string    `json:"community_name"`
	Municipality            *string    `json:"municipality"`
	PostalCommunity         *string    `json:"postal_community"`
	PlaceName               *string    `json:"place_name"`
	SubdivisionName         *string    `json:"subdivision_name"`
	Latitude                *float64   `json:"latitude"`
	Longitude               *float64   `json:"longitude"`
	PropertyClassID         *string    `json:"property_class_id"`
	AddressType             *string    `json:"address_type"`
	PropertyID              *int64     `json:"property_id"`
	AddressMasterID         int64      `json:"address_master_id"`
	LastUpdate              *time.Time `json:"last_update"`
	EffectiveDate           *time.Time `json:"effective_date"`
	ExpirationDate          *time.Time `json:"expiration_date"`
	DPVFootnotes            *string    `json:"dpv_footnotes"`
	DeliveryPointCheckDigit *int       `json:"delivery_point_check_digit"`
	DeliveryPointCode       *string    `json:"delivery_point_code"`
	DPVCount                *int       `json:"dpv_count"`
}

func (dr *Address) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(Address)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "FIPS":
			if field == "" {
				return nil, &errors.Object{
					Id:     "51bb7a42-745e-4689-b552-a58c75636e8a",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "FIPS is required.",
				}
			}
			record.FIPS = field
		case "State":
			record.State = val.StringPtrIfNonZero(field)
		case "County":
			record.County = val.StringPtrIfNonZero(field)
		case "ZIP5":
			if field == "" {
				return nil, &errors.Object{
					Id:     "cf195972-41fa-477e-b3c7-728a3849783c",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "ZIP5 is required.",
				}
			}
			record.ZIP5 = field
		case "ZIP4":
			record.ZIP4 = val.StringPtrIfNonZero(field)
		case "FullStreetAddress":
			record.FullStreetAddress = val.StringPtrIfNonZero(field)
		case "PreDirectional":
			record.PreDirectional = val.StringPtrIfNonZero(field)
		case "StreetNumber":
			record.StreetNumber = val.StringPtrIfNonZero(field)
		case "Street":
			record.Street = val.StringPtrIfNonZero(field)
		case "PostDirectional":
			record.PostDirectional = val.StringPtrIfNonZero(field)
		case "StreetType":
			record.StreetType = val.StringPtrIfNonZero(field)
		case "UnitType":
			record.UnitType = val.StringPtrIfNonZero(field)
		case "UnitNbr":
			record.UnitNbr = val.StringPtrIfNonZero(field)
		case "VacantIndicator":
			record.VacantIndicator = val.StringPtrIfNonZero(field)
		case "NonUSPSAddressIndicator":
			record.NonUSPSAddressIndicator = val.StringPtrIfNonZero(field)
		case "NotCurrentlyDeliverable":
			record.NotCurrentlyDeliverable = val.StringPtrIfNonZero(field)
		case "CommunityName":
			record.CommunityName = val.StringPtrIfNonZero(field)
		case "Municipality":
			record.Municipality = val.StringPtrIfNonZero(field)
		case "PostalCommunity":
			record.PostalCommunity = val.StringPtrIfNonZero(field)
		case "PlaceName":
			record.PlaceName = val.StringPtrIfNonZero(field)
		case "SubdivisionName":
			record.SubdivisionName = val.StringPtrIfNonZero(field)
		case "Latitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "4145b028-363a-45d5-9a15-2a5739d4dec4")
			}
			record.Latitude = v
		case "Longitude":
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "9631a32d-c007-4a31-9d42-99a8fce81a84")
			}
			record.Longitude = v
		case "PropertyClassID":
			record.PropertyClassID = val.StringPtrIfNonZero(field)
		case "AddressType":
			record.AddressType = val.StringPtrIfNonZero(field)
		case "PropertyID":
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "da5ca757-8ebc-432f-870c-2219ab860fe8")
			}
			record.PropertyID = v
		case "AddressMasterID":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:     "78741220-72f9-4f98-85a0-83424563245d",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to parse AddressMasterID.",
					Cause:  err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.AddressMasterID = v
		case "LastUpdate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "dbf6e246-428f-4a4a-821f-28ff6cad3ce2")
			}
			record.LastUpdate = v
		case "EffectiveDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "8e79ece1-549b-4549-bf47-249e56e395d6")
			}
			record.EffectiveDate = v
		case "ExpirationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.IntegerDate, field)
			if err != nil {
				return nil, errors.Forward(err, "afb7a879-aae8-46af-8a70-e6d281e4e880")
			}
			record.ExpirationDate = v
		case "DPVFootnotes":
			record.DPVFootnotes = val.StringPtrIfNonZero(field)
		case "DeliveryPointCheckDigit":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "2739bc35-f0fa-4c08-bd8e-65464c0f3434")
			}
			record.DeliveryPointCheckDigit = v
		case "DeliveryPointCode":
			record.DeliveryPointCode = val.StringPtrIfNonZero(field)
		case "DPVCount":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "169919a7-314d-417e-835a-b628744054db")
			}
			record.DPVCount = v
		default:
			return nil, &errors.Object{
				Id:     "973a9506-222e-4a10-b6aa-45fe1780354d",
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

func (dr *Address) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"fips",
		"state",
		"county",
		"zip5",
		"zip4",
		"full_street_address",
		"pre_directional",
		"street_number",
		"street",
		"post_directional",
		"street_type",
		"unit_type",
		"unit_nbr",
		"vacant_indicator",
		"non_usps_address_indicator",
		"not_currently_deliverable",
		"community_name",
		"municipality",
		"postal_community",
		"place_name",
		"subdivision_name",
		"latitude",
		"longitude",
		"property_class_id",
		"address_type",
		"property_id",
		"address_master_id",
		"last_update",
		"effective_date",
		"expiration_date",
		"dpv_footnotes",
		"delivery_point_check_digit",
		"delivery_point_code",
		"dpv_count",
		"location",
	}
}

func (dr *Address) SQLTable() string {
	return "fa_df_address"
}

func (dr *Address) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "bb27bb7a-498f-45e1-a33b-be16dcf8362b",
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

	var location squirrel.Sqlizer

	if dr.Latitude != nil && dr.Longitude != nil {
		location = squirrel.Expr("st_setsrid(st_makepoint(?, ?), 4326)", dr.Longitude, dr.Latitude)
	} else {
		location = squirrel.Expr("null")
	}

	values := []any{
		dr.AMId,
		dr.AMCreatedAt,
		now,
		dr.AMMeta,
		dr.FIPS,
		dr.State,
		dr.County,
		dr.ZIP5,
		dr.ZIP4,
		dr.FullStreetAddress,
		dr.PreDirectional,
		dr.StreetNumber,
		dr.Street,
		dr.PostDirectional,
		dr.StreetType,
		dr.UnitType,
		dr.UnitNbr,
		dr.VacantIndicator,
		dr.NonUSPSAddressIndicator,
		dr.NotCurrentlyDeliverable,
		dr.CommunityName,
		dr.Municipality,
		dr.PostalCommunity,
		dr.PlaceName,
		dr.SubdivisionName,
		dr.Latitude,
		dr.Longitude,
		dr.PropertyClassID,
		dr.AddressType,
		dr.PropertyID,
		dr.AddressMasterID,
		dr.LastUpdate,
		dr.EffectiveDate,
		dr.ExpirationDate,
		dr.DPVFootnotes,
		dr.DeliveryPointCheckDigit,
		dr.DeliveryPointCode,
		dr.DPVCount,
		location,
	}

	return values, nil
}

func (dr *Address) LoadParams() *entities.DataRecordLoadParams {
	return nil
}
