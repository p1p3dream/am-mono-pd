-- +migrate Up

alter table properties
	add column fa_master_address_id bigint,
	add column fips text,
	add column location geometry(point, 4326)
;
-- +migrate Down

alter table properties
	drop column fa_master_address_id,
	drop column fips,
	drop column location
;
