-- +migrate Up

--------------------------------------------------------------------------------
-- Properties.
--------------------------------------------------------------------------------

create table properties (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	fa_property_id bigint
);

create unique index idx_properties_fa_property_id
	on properties (fa_property_id);

-- +migrate Down

drop table properties;
