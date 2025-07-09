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

type SelectRentalAvmRecordInput struct {
	Aupid *uuid.UUID
}

type SelectRentalAvmRecordOutput struct {
	Records []*entities.RentalAvm
}

func (repo *repository) SelectRentalAvmRecord(r *arc.Request, in *SelectRentalAvmRecordInput) (*SelectRentalAvmRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"ad_df_rental_avm.estimated_rental_value",
			"ad_df_rental_avm.estimated_min_rental_value",
			"ad_df_rental_avm.estimated_max_rental_value",
			"ad_df_rental_avm.valuation_date",
		).
		From("ad_df_rental_avm").
		Join("properties on properties.ad_attom_id = ad_df_rental_avm.attomid").
		Where("properties.id = ?", in.Aupid).
		OrderBy("ad_df_rental_avm.valuation_date desc").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "30a445da-77f5-4d83-9e6d-b72219f53009",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "e409fd79-0391-445d-a985-870417d66bc2")
	}
	defer rows.Close()

	out := &SelectRentalAvmRecordOutput{}

	for rows.Next() {
		record := &entities.RentalAvm{}

		if err := rows.Scan(
			&record.EstimatedRentalValue,
			&record.EstimatedMinRentalValue,
			&record.EstimatedMaxRentalValue,
			&record.ValuationDate,
		); err != nil {
			return nil, &errors.Object{
				Id:     "3665cbec-88e5-47af-b40f-c0e6f004dbae",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}
