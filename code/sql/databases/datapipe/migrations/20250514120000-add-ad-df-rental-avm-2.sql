-- +migrate Up

create table ad_df_rental_avm_2 (
	am_id         uuid not null,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	attomid                                bigint not null,
	property_address_full                  text,
	property_address_house_number          text,
	property_address_street_direction      text,
	property_address_street_name           text,
	property_address_street_suffix         text,
	property_address_street_post_direction text,
	property_address_unit_prefix           text,
	property_address_unit_value            text,
	property_address_city                  text,
	property_address_state                 text,
	property_address_zip                   text,
	property_address_zip4                  text,
	property_address_crrt                  text,
	property_use_group                     text,
	property_use_standardized              integer,
	estimated_rental_value                 integer,
	estimated_min_rental_value             integer,
	estimated_max_rental_value             integer,
	valuation_date                         date,
	publication_date                       date,

	primary key (am_id)
);

create index idx_ad_df_rental_avm_2_attomid
	on ad_df_rental_avm_2 (attomid);

create index idx_ad_df_rental_avm_2_valuation_date
	on ad_df_rental_avm_2 (valuation_date);

-- +migrate Down

drop table ad_df_rental_avm_2;
