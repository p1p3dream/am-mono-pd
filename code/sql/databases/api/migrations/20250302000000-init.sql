-- +migrate Up

--------------------------------------------------------------------------------
-- Api keys.
--------------------------------------------------------------------------------

create table api_keys (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid not null,
	user_id         uuid,
	role_id         uuid not null,
	-- We cache role_name to prevent a database lookup.
	-- This shouldn't be updated frequently.
	role_name       text not null,
	key_type        smallint not null,
	key_status      smallint not null,
	expires_at      timestamp with time zone,
	last_used_at    timestamp with time zone,
	revoked_at      timestamp with time zone,
	revoked_by      uuid,
	key_hash        text not null,
	name            text,
	description     text,

	unique (key_type, key_hash)
);

create index api_keys_organization_id_user_id on api_keys (organization_id, user_id);

--------------------------------------------------------------------------------
-- Clients.
--------------------------------------------------------------------------------

create table clients (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid not null,
	name            text not null,
	description     text,
	redirect_uri    text
);

create index clients_organization_id on clients (organization_id);

-- +migrate Down

drop table clients;
drop table api_keys;
