package auth

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/val"
)

type Repository interface {
	SelectApiKeyRecord(r *arc.Request, in *SelectApiKeyRecordInput) (*SelectApiKeyRecordOutput, error)

	InsertApiQuotaTransactionRecord(r *arc.Request, in *InsertApiQuotaTransactionRecordInput) (*InsertApiQuotaTransactionRecordOutput, error)
	SelectApiQuotaAvailability(r *arc.Request, in *SelectApiQuotaAvailabilityInput) (*SelectApiQuotaAvailabilityOutput, error)

	SelectClientRecord(r *arc.Request, in *SelectClientRecordInput) (*SelectClientRecordOutput, error)
}

type repository struct{}

type SelectApiKeyRecordInput struct {
	KeyType   ApiKeyType
	KeyHash   string
	KeyStatus ApiKeyStatus

	SelectExpired bool
}

type SelectApiKeyRecordOutput struct {
	Record *ApiKey
}

func (rep *repository) SelectApiKeyRecord(r *arc.Request, in *SelectApiKeyRecordInput) (*SelectApiKeyRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("api_keys").
		Set("updated_at", "now()").
		Set("last_used_at", "now()").
		Set("key_status", squirrel.Expr(
			`
				case
				when expires_at is not null and expires_at < now() then ?
				else key_status
				end
			`,
			ApiKeyStatusExpired,
		)).
		Where("key_type = ?", in.KeyType).
		Where("key_hash = ?", in.KeyHash).
		Suffix(`
			returning
				id,
				organization_id,
				user_id,
				role_id,
				role_name,
				key_type,
				key_status,
				expires_at
		`)

	if in.KeyStatus != 0 {
		builder = builder.Where("key_status = ?", in.KeyStatus)
	}

	if !in.SelectExpired {
		builder = builder.Where("(expires_at is null or expires_at < now())")
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "0b78b5b5-bfc2-44d1-b663-581fbe02f51b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresApi, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "72f13837-f79b-445e-bfea-fe00fcf3f08c")
	}

	var (
		userId    pgtype.UUID
		expiresAt pgtype.Timestamptz
	)

	out := &SelectApiKeyRecordOutput{}
	record := new(ApiKey)

	if err := row.Scan(
		&record.Id,
		&record.OrganizationId,
		&userId,
		&record.RoleId,
		&record.RoleName,
		&record.KeyType,
		&record.KeyStatus,
		&expiresAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}

		return nil, &errors.Object{
			Id:     "c59b8fac-7d93-4a26-8f0a-4630adee1ebf",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	if userId.Valid {
		u, err := val.UUIDFromBytes(userId.Bytes[:])
		if err != nil {
			return nil, errors.Forward(err, "b520a378-25dc-42d3-a872-287fb94c9c9e")
		}

		record.UserId = u
	}

	if expiresAt.Valid {
		record.ExpiresAt = expiresAt.Time
	}

	out.Record = record

	return out, nil
}

type InsertApiQuotaTransactionRecordInput struct {
	Record *arc.ApiQuotaTransaction
}

type InsertApiQuotaTransactionRecordOutput struct {
	DailyQuota    int64
	DailyUsage    int64
	HasDailyQuota bool

	MonthlyQuota    int64
	MonthlyUsage    int64
	HasMonthlyQuota bool

	TrxLayoutSum int64
}

func (repo *repository) InsertApiQuotaTransactionRecord(r *arc.Request, in *InsertApiQuotaTransactionRecordInput) (*InsertApiQuotaTransactionRecordOutput, error) {
	record := in.Record

	apiQuotaAvailabilityPrefixOut, err := repo.buildApiQuotaAvailabilityPrefix(r, &buildApiQuotaAvailabilityPrefixInput{
		OrganizationId:            record.OrganizationId,
		ApiQuotaTransactionRecord: record,
	})
	if err != nil {
		return nil, errors.Forward(err, "a5da13f3-6127-4986-ac96-79b37cec5592")
	}

	now := time.Now()

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("api_quota_transactions").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"organization_id",
			"api_key_id",
			"trx_timestamp",
			"description",
			"base_req_amount",
			"address_lo_amount",
			"assessor_lo_amount",
			"comps_lo_amount",
			"listing_lo_amount",
			"recorder_lo_amount",
			"rent_estimate_lo_amount",
			"sale_estimate_lo_amount",
		).
		PrefixExpr(apiQuotaAvailabilityPrefixOut.SquirrelExpr).
		Values(
			record.Id,
			now,
			now,
			record.Meta,
			record.OrganizationId,
			record.ApiKeyId,
			now,
			record.Description,
			squirrel.Expr("case when (select ok from has_quota) and ? then 0 else 1 end", apiQuotaAvailabilityPrefixOut.TrxLayoutSum > 0),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.AddressLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.AssessorLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.CompsLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.ListingLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.RecorderLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.RentEstimateLayoutAmount),
			squirrel.Expr("case when (select ok from has_quota) then ? else 0 end", record.SaleEstimateLayoutAmount),
		).
		Suffix(
			`
			returning
				coalesce((select daily_quota from api_quotas where organization_id = ?), 0),
				(select ending_daily_usage.total_sum from ending_daily_usage),
				(select ok from has_daily_quota),
				coalesce((select monthly_quota from api_quotas where organization_id = ?), 0),
				(select ending_monthly_usage.total_sum from ending_monthly_usage),
				(select ok from has_monthly_quota)
			`,
			record.OrganizationId,
			record.OrganizationId,
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "b253a828-082a-4c47-b52d-7cbf636afbaf",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresApi, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "15b465ef-4b51-473a-94f0-f238a7f04836")
	}

	out := &InsertApiQuotaTransactionRecordOutput{
		TrxLayoutSum: apiQuotaAvailabilityPrefixOut.TrxLayoutSum,
	}

	if err := row.Scan(
		&out.DailyQuota,
		&out.DailyUsage,
		&out.HasDailyQuota,
		&out.MonthlyQuota,
		&out.MonthlyUsage,
		&out.HasMonthlyQuota,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}

		return nil, &errors.Object{
			Id:     "4df3b53e-bace-435a-bba9-905ecd0c5d7f",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	return out, nil
}

type SelectApiQuotaAvailabilityInput struct {
	OrganizationId uuid.UUID
}

type SelectApiQuotaAvailabilityOutput struct {
	EnabledLayouts  []string
	HasDailyQuota   bool
	HasMonthlyQuota bool
}

func (repo *repository) SelectApiQuotaAvailability(r *arc.Request, in *SelectApiQuotaAvailabilityInput) (*SelectApiQuotaAvailabilityOutput, error) {
	// Build against an empty transaction to check the quota availability.
	apiQuotaAvailabilityPrefixOut, err := repo.buildApiQuotaAvailabilityPrefix(r, &buildApiQuotaAvailabilityPrefixInput{
		OrganizationId:            in.OrganizationId,
		ApiQuotaTransactionRecord: &arc.ApiQuotaTransaction{},
	})
	if err != nil {
		return nil, errors.Forward(err, "50717fcd-e1a7-4674-83e0-a3ffbc120573")
	}

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select().
		Column(squirrel.Expr("(select ok from has_daily_quota)")).
		Column(squirrel.Expr("(select ok from has_monthly_quota)")).
		Column(squirrel.Expr(
			`(
			select
				case when api_quotas.address_lo_enabled then ARRAY['API_ADDRESS_LAYOUT_ENABLED'] end ||
				case when api_quotas.assessor_lo_enabled then ARRAY['API_ASSESSOR_LAYOUT_ENABLED'] end ||
				case when api_quotas.comps_lo_enabled then ARRAY['API_COMPS_LAYOUT_ENABLED'] end ||
				case when api_quotas.listing_lo_enabled then ARRAY['API_LISTING_LAYOUT_ENABLED'] end ||
				case when api_quotas.recorder_lo_enabled then ARRAY['API_RECORDER_LAYOUT_ENABLED'] end ||
				case when api_quotas.rent_estimate_lo_enabled then ARRAY['API_RENT_ESTIMATE_LAYOUT_ENABLED'] end ||
				case when api_quotas.sale_estimate_lo_enabled then ARRAY['API_SALE_ESTIMATE_LAYOUT_ENABLED'] end
			from api_quotas
			where organization_id = ?
			)`,
			in.OrganizationId,
		)).
		PrefixExpr(apiQuotaAvailabilityPrefixOut.SquirrelExpr)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "783f5704-1c80-4d82-84fb-fa995a28eaca",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresApi, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "e64f9082-eda0-4e85-a17a-6242dbe69548")
	}

	out := &SelectApiQuotaAvailabilityOutput{}

	if err := row.Scan(
		&out.HasDailyQuota,
		&out.HasMonthlyQuota,
		&out.EnabledLayouts,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}

		return nil, &errors.Object{
			Id:     "2c5f4208-f291-4639-bd69-757f0b483eea",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	return out, nil
}

type buildApiQuotaAvailabilityPrefixInput struct {
	OrganizationId            uuid.UUID
	ApiQuotaTransactionRecord *arc.ApiQuotaTransaction
}

type buildApiQuotaAvailabilityPrefixOutput struct {
	SquirrelExpr squirrel.Sqlizer
	TrxLayoutSum int64
}

func (repo *repository) buildApiQuotaAvailabilityPrefix(_ *arc.Request, in *buildApiQuotaAvailabilityPrefixInput) (*buildApiQuotaAvailabilityPrefixOutput, error) {
	orgId := in.OrganizationId
	trxRecord := in.ApiQuotaTransactionRecord
	trxLayoutSum := int64(trxRecord.AddressLayoutAmount) +
		int64(trxRecord.AssessorLayoutAmount) +
		int64(trxRecord.CompsLayoutAmount) +
		int64(trxRecord.ListingLayoutAmount) +
		int64(trxRecord.RecorderLayoutAmount) +
		int64(trxRecord.RentEstimateLayoutAmount) +
		int64(trxRecord.SaleEstimateLayoutAmount)

	originalTrxLayoutSum := trxLayoutSum

	if originalTrxLayoutSum == 0 {
		// Add the baseRequestAmount to the layout sum.
		trxLayoutSum++
	}

	expr := squirrel.Expr(
		`
		with current_daily_usage as (
			select
				extract(day from trx_timestamp) as day,
				sum(base_req_amount) as base_req_sum,
				sum(address_lo_amount) as address_lo_sum,
				sum(assessor_lo_amount) as assessor_lo_sum,
				sum(comps_lo_amount) as comps_lo_sum,
				sum(listing_lo_amount) as listing_lo_sum,
				sum(recorder_lo_amount) as recorder_lo_sum,
				sum(rent_estimate_lo_amount) as rent_estimate_lo_sum,
				sum(sale_estimate_lo_amount) as sale_estimate_lo_sum
			from api_quota_transactions
			where
				organization_id = ?
				and trx_timestamp >= date_trunc('month', current_date)
			group by day
		), ending_daily_usage as (
			select sum(
				? +
				coalesce(
					(
						select sum(
							coalesce(base_req_sum, 0) +
							coalesce(address_lo_sum, 0) +
							coalesce(assessor_lo_sum, 0) +
							coalesce(comps_lo_sum, 0) +
							coalesce(listing_lo_sum, 0) +
							coalesce(recorder_lo_sum, 0) +
							coalesce(rent_estimate_lo_sum, 0) +
							coalesce(sale_estimate_lo_sum, 0)
						)
						from current_daily_usage
						where day = extract(day from current_date)
					),
					0
				)
			) as total_sum
		), has_daily_quota as (
			select coalesce(
				total_sum <= (
					select daily_quota
					from api_quotas
					where organization_id = ?
				),
				false
			) as ok
			from ending_daily_usage
		), current_monthly_usage as (
			select
				sum(base_req_sum) as base_req_sum,
				sum(address_lo_sum) as address_lo_sum,
				sum(assessor_lo_sum) as assessor_lo_sum,
				sum(comps_lo_sum) as comps_lo_sum,
				sum(listing_lo_sum) as listing_lo_sum,
				sum(recorder_lo_sum) as recorder_lo_sum,
				sum(rent_estimate_lo_sum) as rent_estimate_lo_sum,
				sum(sale_estimate_lo_sum) as sale_estimate_lo_sum
			from current_daily_usage
		), ending_monthly_usage as (
			select
				sum(
					coalesce(base_req_sum, 0) +
					coalesce(address_lo_sum, 0) +
					coalesce(assessor_lo_sum, 0) +
					coalesce(comps_lo_sum, 0) +
					coalesce(listing_lo_sum, 0) +
					coalesce(recorder_lo_sum, 0) +
					coalesce(rent_estimate_lo_sum, 0) +
					coalesce(sale_estimate_lo_sum, 0) +
					?
				) as total_sum
			from current_monthly_usage
		), has_monthly_quota as (
			select coalesce(
				total_sum <= (
					select monthly_quota
					from api_quotas
					where organization_id = ?
				),
				false
			) as ok
			from ending_monthly_usage
		), has_quota as (
			select (
				(select ok from has_daily_quota)
				and (select ok from has_monthly_quota)
			) as ok
		)
		`,
		orgId,
		trxLayoutSum,
		orgId,
		trxLayoutSum,
		orgId,
	)

	out := &buildApiQuotaAvailabilityPrefixOutput{
		SquirrelExpr: expr,
		TrxLayoutSum: originalTrxLayoutSum,
	}

	return out, nil
}

type SelectClientRecordInput struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
}

type SelectClientRecordOutput struct {
	Record *Client
}

func (rep *repository) SelectClientRecord(r *arc.Request, in *SelectClientRecordInput) (*SelectClientRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(`
			id,
			organization_id,
			redirect_uri
		`).
		From("clients").
		Where("id = ?", in.Id).
		Where("organization_id = ?", in.OrganizationId)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "f7c40c28-b649-4b0d-95fa-24e04b931054",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresApi, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "0d6d16ed-0eb0-4061-b78c-f7cfa7cf18a4")
	}

	var redirectUri pgtype.Text

	out := &SelectClientRecordOutput{}
	record := new(Client)

	if err := row.Scan(
		&record.Id,
		&record.OrganizationId,
		&redirectUri,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}

		return nil, &errors.Object{
			Id:     "b8a2c454-597f-42ce-8686-404ed81b1772",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	if redirectUri.Valid {
		record.RedirectUri = redirectUri.String
	}

	out.Record = record

	return out, nil
}
