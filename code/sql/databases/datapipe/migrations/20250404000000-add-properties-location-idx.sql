-- +migrate Up

alter table properties
	rename column fa_master_address_id to fa_address_master_id
;

create index idx_properties_fa_address_master_id on properties (fa_address_master_id);
create index idx_properties_fips on properties (fips);
create index idx_properties_location on properties using gist (location);

-- +migrate Down

drop index idx_properties_location;
drop index idx_properties_fips;
drop index idx_properties_fa_address_master_id;

alter table properties
	rename column fa_address_master_id to fa_master_address_id
;
