-- +migrate Up

--------------------------------------------------------------------------------
-- FirstAmerican - Fips.
--------------------------------------------------------------------------------

create table fips (
	am_id         uuid primary key,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	fips text not null unique
);

with distinct_fips as (
	select distinct fips
	from fa_df_address
)
insert into fips (
	am_id,
	am_created_at,
	am_updated_at,
	fips
)
select
	gen_random_uuid(),
	now(),
	now(),
	fips
from distinct_fips
order by fips
;

-- +migrate Down

drop table fips;
