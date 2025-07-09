-- +migrate Up

--------------------------------------------------------------------------------
-- Api Quotas.
--------------------------------------------------------------------------------

create table api_quotas (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid not null,

	daily_quota   bigint not null,
	monthly_quota bigint not null,

	address_lo_enabled       boolean not null default false,
	assessor_lo_enabled      boolean not null default false,
	comps_lo_enabled         boolean not null default false,
	listing_lo_enabled       boolean not null default false,
	recorder_lo_enabled      boolean not null default false,
	rent_estimate_lo_enabled boolean not null default false,
	sale_estimate_lo_enabled boolean not null default false
);

create index api_quotas_organization_id
	on api_quotas (organization_id);

-- Insert base quotas for the AbodeMine org.
insert into api_quotas (
	id,
	created_at,
	updated_at,
	meta,
	organization_id,
	daily_quota,
	monthly_quota,
	address_lo_enabled,
	assessor_lo_enabled,
	comps_lo_enabled,
	listing_lo_enabled,
	recorder_lo_enabled,
	rent_estimate_lo_enabled,
	sale_estimate_lo_enabled
) values (
	gen_random_uuid(),
	now(),
	now(),
	null,
	'019543c8-8fc8-7ab2-9d6b-982e4ccb11f5',
	100000,
	3100000,
	true,
	true,
	true,
	true,
	true,
	true,
	true
);

--------------------------------------------------------------------------------
-- Api Quota Transactions.
--------------------------------------------------------------------------------

create table api_quota_transactions (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid not null,
	api_key_id      uuid,
	trx_timestamp   timestamp with time zone not null,
	description     text,

	base_req_amount           integer not null default 0,
	address_lo_amount         integer not null default 0,
	assessor_lo_amount        integer not null default 0,
	comps_lo_amount           integer not null default 0,
	listing_lo_amount         integer not null default 0,
	recorder_lo_amount        integer not null default 0,
	rent_estimate_lo_amount   integer not null default 0,
	sale_estimate_lo_amount   integer not null default 0
);

create index api_quota_transactions_org_ts
	on api_quota_transactions (organization_id, trx_timestamp);

create index api_quota_transactions_description
	on api_quota_transactions (description);

-- +migrate Down

drop table api_quota_transactions;
drop table api_quotas;
