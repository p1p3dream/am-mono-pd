-- +migrate Up

create index idx_properties_address_id
	on properties (address_id);

-- +migrate Down

drop index if exists idx_properties_address_id;
