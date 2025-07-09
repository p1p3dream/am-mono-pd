-- +migrate Up

--------------------------------------------------------------------------------
-- Addresses table.
--------------------------------------------------------------------------------

create table addresses (
	id                uuid primary key,
	created_at        timestamp with time zone not null,
	updated_at        timestamp with time zone not null,
	meta              jsonb,
	city                 text,
	county               text,
	data_source          text,
	fips                 text,
	full_street_address  text,
	house_number         text,
	state                text,
	street_name          text,
	street_number        text,
	street_pos_direction text,
	street_pre_direction text,
	street_suffix        text,
	street_type          text,
	unit_nbr             text,
	unit_type            text,
	zip5                 text
);

create index idx_addresses_upper_fips
	on addresses (upper(fips));

create index idx_addresses_upper_full_street_address
	on addresses (upper(full_street_address));

create index idx_addresses_upper_state
	on addresses (upper(state));

create index idx_addresses_upper_county
	on addresses (upper(county));

create index idx_addresses_upper_zip5
	on addresses (upper(zip5));

create index idx_addresses_upper_street_number
	on addresses (upper(street_number));

create index idx_addresses_upper_street_name
	on addresses (upper(street_name));

create index idx_addresses_upper_street_type
	on addresses (upper(street_type));

create index idx_addresses_upper_unit_type
	on addresses (upper(unit_type));

create index idx_addresses_upper_unit_nbr
	on addresses (upper(unit_nbr));

create index idx_addresses_upper_city
	on addresses (upper(city));

create index idx_addresses_upper_house_number
	on addresses (upper(house_number));

create index idx_addresses_upper_street_pre_direction
	on addresses (upper(street_pre_direction));

create index idx_addresses_upper_street_pos_direction
	on addresses (upper(street_pos_direction));

create index idx_addresses_upper_street_suffix
	on addresses (upper(street_suffix));

--------------------------------------------------------------------------------
-- Alter properties table to include addresses_id.
--------------------------------------------------------------------------------

alter table properties
	add column address_id uuid;

-- +migrate Down

alter table properties
	drop column if exists address_id;

drop table if exists addresses;
