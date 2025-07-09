package attom_data

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type RentalAvm struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any

	ATTOMID                            int64
	PropertyAddressFull                *string
	PropertyAddressHouseNumber         *string
	PropertyAddressStreetDirection     *string
	PropertyAddressStreetName          *string
	PropertyAddressStreetSuffix        *string
	PropertyAddressStreetPostDirection *string
	PropertyAddressUnitPrefix          *string
	PropertyAddressUnitValue           *string
	PropertyAddressCity                *string
	PropertyAddressState               *string
	PropertyAddressZIP                 *string
	PropertyAddressZIP4                *string
	PropertyAddressCRRT                *string
	PropertyUseGroup                   *string
	PropertyUseStandardized            *int
	EstimatedRentalValue               *int
	EstimatedMinRentalValue            *int
	EstimatedMaxRentalValue            *int
	ValuationDate                      *time.Time
	PublicationDate                    *time.Time
}

func (dr *RentalAvm) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new(RentalAvm)

	for k, header := range headers {
		field := fields[k]

		switch header {
		case "[ATTOM ID]":
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "6004e19f-6c48-4ff5-a486-04f91e633b9c",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.ATTOMID = v
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
		case "PropertyUseGroup":
			record.PropertyUseGroup = val.StringPtrIfNonZero(field)
		case "PropertyUseStandardized":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "c29cab6c-fc0f-4204-b33c-dde3a63d9fb7")
			}
			record.PropertyUseStandardized = v
		case "EstimatedRentalValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "3ed0567b-0b2e-48a3-9da8-54761fbafd62")
			}
			record.EstimatedRentalValue = v
		case "EstimatedMinRentalValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "530b16c0-2c57-4938-8069-654cb20c8ef0")
			}
			record.EstimatedMinRentalValue = v
		case "EstimatedMaxRentalValue":
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "bb5f5702-5efb-4125-abf8-4ae6ae210d02")
			}
			record.EstimatedMaxRentalValue = v
		case "ValuationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "cfa9dce2-457f-4fd5-b760-de54588a8920")
			}
			record.ValuationDate = v
		case "PublicationDate":
			v, err := val.TimePtrFromStringIfNonZero(consts.RFC3339Date, field)
			if err != nil {
				return nil, errors.Forward(err, "4c10246a-4b26-4a99-bf4f-4b45b4f5b29f")
			}
			record.PublicationDate = v
		default:
			return nil, &errors.Object{
				Id:     "d48e3097-f939-41db-a241-c92889890203",
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

func (dr *RentalAvm) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
		"attomid",
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
		"property_use_group",
		"property_use_standardized",
		"estimated_rental_value",
		"estimated_min_rental_value",
		"estimated_max_rental_value",
		"valuation_date",
		"publication_date",
	}
}

func (dr *RentalAvm) SQLTable() string {
	return "ad_df_rental_avm"
}

func (dr *RentalAvm) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "509af6f7-7b0a-44c3-b6cc-d3b47783fd03",
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
		dr.PropertyUseGroup,
		dr.PropertyUseStandardized,
		dr.EstimatedRentalValue,
		dr.EstimatedMinRentalValue,
		dr.EstimatedMaxRentalValue,
		dr.ValuationDate,
		dr.PublicationDate,
	}

	return values, nil
}

func (dr *RentalAvm) LoadParams() *entities.DataRecordLoadParams {
	return &entities.DataRecordLoadParams{
		Mode: entities.DataRecordModeBatchInsert,
	}
}
