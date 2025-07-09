package avm

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
)

type Repository interface {
	SelectRentalAvmRecord(r *arc.Request, in *SelectRentalAvmRecordInput) (*SelectRentalAvmRecordOutput, error)
	SelectSaleAvmRecord(r *arc.Request, in *SelectSaleAvmRecordInput) (*SelectSaleAvmRecordOutput, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

type SelectSaleAvmRecordInput struct {
	Aupid *uuid.UUID
}

type SelectSaleAvmRecordOutput struct {
	Records []*entities.SaleAvm
}

func (repo *repository) SelectSaleAvmRecord(r *arc.Request, in *SelectSaleAvmRecordInput) (*SelectSaleAvmRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"fa_df_avm_power.fips",
			"fa_df_avm_power.apn",
			"fa_df_avm_power.situs_full_street_address",
			"fa_df_avm_power.situs_house_nbr",
			"fa_df_avm_power.situs_house_nbr_suffix",
			"fa_df_avm_power.situs_direction_left",
			"fa_df_avm_power.situs_street",
			"fa_df_avm_power.situs_mode",
			"fa_df_avm_power.situs_direction_right",
			"fa_df_avm_power.situs_unit_type",
			"fa_df_avm_power.situs_unit_nbr",
			"fa_df_avm_power.situs_city",
			"fa_df_avm_power.situs_state",
			"fa_df_avm_power.zip5",
			"fa_df_avm_power.zip4",
			"fa_df_avm_power.situs_carrier_code",
			"fa_df_avm_power.final_value",
			"fa_df_avm_power.high_value",
			"fa_df_avm_power.low_value",
			"fa_df_avm_power.confidence_score",
			"fa_df_avm_power.standard_deviation",
			"fa_df_avm_power.valuation_date",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp1_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp2_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp3_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp4_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp5_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp6_property_id)",
			"(select id from properties where fa_property_id = fa_df_avm_power.comp7_property_id)",
		).
		From("properties").
		Join("fa_df_avm_power on properties.fa_property_id = fa_df_avm_power.property_id").
		Where("properties.id = ?", in.Aupid).
		OrderBy("fa_df_avm_power.valuation_date desc").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "b70cd6bf-c0e3-41c9-8ea2-eeb18ca4861e",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "e5d063e5-4c35-45fd-a4a7-d91ca3a98211")
	}
	defer rows.Close()

	out := &SelectSaleAvmRecordOutput{}

	for rows.Next() {
		record := &entities.SaleAvm{}

		var (
			comp1Aupid *uuid.UUID
			comp2Aupid *uuid.UUID
			comp3Aupid *uuid.UUID
			comp4Aupid *uuid.UUID
			comp5Aupid *uuid.UUID
			comp6Aupid *uuid.UUID
			comp7Aupid *uuid.UUID
		)

		if err := rows.Scan(
			&record.Fips,
			&record.Apn,
			&record.FullStreetAddress,
			&record.HouseNumber,
			&record.HouseNumberSuffix,
			&record.StreetPreDirection,
			&record.StreetName,
			&record.StreetSuffix,
			&record.StreetPostDirection,
			&record.UnitType,
			&record.UnitNumber,
			&record.City,
			&record.State,
			&record.Zip5,
			&record.Zip4,
			&record.CarrierCode,
			&record.FinalValue,
			&record.HighValue,
			&record.LowValue,
			&record.ConfidenceScore,
			&record.StandardDeviation,
			&record.ValuationDate,
			&comp1Aupid,
			&comp2Aupid,
			&comp3Aupid,
			&comp4Aupid,
			&comp5Aupid,
			&comp6Aupid,
			&comp7Aupid,
		); err != nil {
			return nil, &errors.Object{
				Id:     "982e2ffe-f1f9-4c65-9c64-bb76dd6b0e32",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		if comp1Aupid != nil {
			record.Comps = append(record.Comps, comp1Aupid)
		}

		if comp2Aupid != nil {
			record.Comps = append(record.Comps, comp2Aupid)
		}

		if comp3Aupid != nil {
			record.Comps = append(record.Comps, comp3Aupid)
		}

		if comp4Aupid != nil {
			record.Comps = append(record.Comps, comp4Aupid)
		}

		if comp5Aupid != nil {
			record.Comps = append(record.Comps, comp5Aupid)
		}

		if comp6Aupid != nil {
			record.Comps = append(record.Comps, comp6Aupid)
		}

		if comp7Aupid != nil {
			record.Comps = append(record.Comps, comp7Aupid)
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}
