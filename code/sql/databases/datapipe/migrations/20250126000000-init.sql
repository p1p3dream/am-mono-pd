-- +migrate Up

--------------------------------------------------------------------------------
-- FirstAmerican - Data File - Address.
--------------------------------------------------------------------------------

create table fa_df_address (
	am_id         uuid primary key,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	fips                       text not null,
	state                      text,
	county                     text,
	zip5                       text not null,
	zip4                       text,
	full_street_address        text,
	pre_directional            text,
	street_number              text,
	street                     text,
	post_directional           text,
	street_type                text,
	unit_type                  text,
	unit_nbr                   text,
	vacant_indicator           text,
	non_usps_address_indicator text,
	not_currently_deliverable  text,
	community_name             text,
	municipality               text,
	postal_community           text,
	place_name                 text,
	subdivision_name           text,
	latitude                   double precision not null,
	longitude                  double precision not null,
	property_class_id          text,
	address_type               text,
	property_id                bigint,
	address_master_id          bigint not null,
	last_update                date,
	effective_date             date,
	expiration_date            date,
	dpv_footnotes              text,
	delivery_point_check_digit integer,
	delivery_point_code        text,
	dpv_count                  integer
);

create index idx_fa_df_address_address_master_id
	on fa_df_address (address_master_id);
create index idx_fa_df_address_property_id
	on fa_df_address (property_id);

--------------------------------------------------------------------------------
-- FirstAmerican - Data File - AVMPower.
--------------------------------------------------------------------------------

create table fa_df_avm_power (
	am_id         uuid,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	fips                      text not null,
	property_id               bigint not null,
	apn                       text not null,
	situs_full_street_address text,
	situs_house_nbr           text,
	situs_house_nbr_suffix    text,
	situs_direction_left      text,
	situs_street              text,
	situs_mode                text,
	situs_direction_right     text,
	situs_unit_type           text,
	situs_unit_nbr            text,
	situs_city                text,
	situs_state               text,
	zip5                      text,
	zip4                      text,
	situs_carrier_code        text,
	final_value               numeric(20, 4),
	high_value                numeric(20, 4),
	low_value                 numeric(20, 4),
	confidence_score          double precision,
	standard_deviation        double precision,
	valuation_date            date,
	comp1_property_id         bigint,
	comp2_property_id         bigint,
	comp3_property_id         bigint,
	comp4_property_id         bigint,
	comp5_property_id         bigint,
	comp6_property_id         bigint,
	comp7_property_id         bigint,

	primary key (am_id, valuation_date)
)
partition by range (valuation_date)
;

create table fa_df_avm_power_202406 partition of fa_df_avm_power
	for values from ('2024-06-01') to ('2024-07-01');
create index idx_fa_df_avm_power_202406
	on fa_df_avm_power_202406 (property_id);

create table fa_df_avm_power_202407 partition of fa_df_avm_power
	for values from ('2024-07-01') to ('2024-08-01');
create index idx_fa_df_avm_power_202407
	on fa_df_avm_power_202407 (property_id);

create table fa_df_avm_power_202408 partition of fa_df_avm_power
	for values from ('2024-08-01') to ('2024-09-01');
create index idx_fa_df_avm_power_202408
	on fa_df_avm_power_202408 (property_id);

create table fa_df_avm_power_202409 partition of fa_df_avm_power
	for values from ('2024-09-01') to ('2024-10-01');
create index idx_fa_df_avm_power_202409
	on fa_df_avm_power_202409 (property_id);

create table fa_df_avm_power_202410 partition of fa_df_avm_power
	for values from ('2024-10-01') to ('2024-11-01');
create index idx_fa_df_avm_power_202410
	on fa_df_avm_power_202410 (property_id);

create table fa_df_avm_power_202411 partition of fa_df_avm_power
	for values from ('2024-11-01') to ('2024-12-01');
create index idx_fa_df_avm_power_202411
	on fa_df_avm_power_202411 (property_id);

create table fa_df_avm_power_202412 partition of fa_df_avm_power
	for values from ('2024-12-01') to ('2025-01-01');
create index idx_fa_df_avm_power_202412
	on fa_df_avm_power_202412 (property_id);

create table fa_df_avm_power_202501 partition of fa_df_avm_power
	for values from ('2025-01-01') to ('2025-02-01');
create index idx_fa_df_avm_power_202501
	on fa_df_avm_power_202501 (property_id);

create table fa_df_avm_power_202502 partition of fa_df_avm_power
	for values from ('2025-02-01') to ('2025-03-01');
create index idx_fa_df_avm_power_202502
	on fa_df_avm_power_202502 (property_id);

create table fa_df_avm_power_202503 partition of fa_df_avm_power
	for values from ('2025-03-01') to ('2025-04-01');
create index idx_fa_df_avm_power_202503
	on fa_df_avm_power_202503 (property_id);

create table fa_df_avm_power_202504 partition of fa_df_avm_power
	for values from ('2025-04-01') to ('2025-05-01');
create index idx_fa_df_avm_power_202504
	on fa_df_avm_power_202504 (property_id);

create table fa_df_avm_power_202505 partition of fa_df_avm_power
	for values from ('2025-05-01') to ('2025-06-01');
create index idx_fa_df_avm_power_202505
	on fa_df_avm_power_202505 (property_id);

create table fa_df_avm_power_202506 partition of fa_df_avm_power
	for values from ('2025-06-01') to ('2025-07-01');
create index idx_fa_df_avm_power_202506
	on fa_df_avm_power_202506 (property_id);

-- +migrate Down

drop table fa_df_avm_power;
drop table fa_df_address;
