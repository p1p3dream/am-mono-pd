-- +migrate Up

--------------------------------------------------------------------------------
-- Organizations.
--------------------------------------------------------------------------------

create table organizations (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	name text not null
);

--------------------------------------------------------------------------------
-- Users.
--------------------------------------------------------------------------------

create table users (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid references organizations (id),
	username        text not null,
	email           text,
	pass_algo       smallint,
	pass_hash       bytea,

	unique (username)
);

create index users_organization_id on users (organization_id);

-- +migrate Down

drop table users;
drop table organizations;
