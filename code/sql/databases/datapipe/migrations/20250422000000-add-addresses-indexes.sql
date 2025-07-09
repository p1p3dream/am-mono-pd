-- +migrate Up

create index idx_addresses_fips_id
	on addresses (fips, id);

-- +migrate Down

drop index idx_addresses_fips_id;
