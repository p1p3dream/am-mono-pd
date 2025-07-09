-- +migrate Up

--------------------------------------------------------------------------------
-- Roles.
--------------------------------------------------------------------------------

create table roles (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	organization_id uuid references organizations (id),
	name            text not null,

	unique(name)
);

create index roles_organization_id on roles (organization_id);

insert into roles
	(
		id,
		created_at,
		updated_at,
		name
	)
	values
	(
		'01958128-c52d-722d-9840-24ea153f0c10',
		now(),
		now(),
		'root'
	),
	(
		'01958128-c52d-724f-b759-07ed08e86c19',
		now(),
		now(),
		'api_user'
	),
	(
		'01958128-c52d-7255-9aac-1a536d22e486',
		now(),
		now(),
		'saas_user'
	),
	(
		'01958128-c52d-7259-ba49-6ace1bf0b9c5',
		now(),
		now(),
		'api_whitelabel_user'
	),
	(
		'01958128-c52d-726b-8ea7-5cf58ba89b00',
		now(),
		now(),
		'saas_whitelabel_user'
	),
	(
		'01958167-1fe9-72ff-a78a-2f3c194aaa0a',
		now(),
		now(),
		'system_auth_check_user'
	)
;

alter table users
	add column role_id uuid references roles (id),
	add column external_id text
;

create index users_organization_id_external_id on users (organization_id, external_id);

insert into organizations
	(
		id,
		created_at,
		updated_at,
		name
	)
	values
	(
		'01956620-1a03-7b8a-888d-63ce64763a17',
		now(),
		now(),
		'Sell2Rent'
	),
	(
		'01956623-5d4f-7dfd-9d45-c86baa4da2b7',
		now(),
		now(),
		'Nhimble'
	),
	(
		'01956624-b06d-7242-b623-454807841b0e',
		now(),
		now(),
		'Roberto'
	),
	(
		'01956624-e036-758b-bf44-a720537ac5cd',
		now(),
		now(),
		'Ben'
	),
	(
		'01956624-f91a-7487-8b0e-3a05144eee46',
		now(),
		now(),
		'Palm Agent'
	)
;

delete from users where id = '019543cd-de46-7afa-9f7f-1e0c94fd3fea';

-- +migrate Down

delete from users where organization_id in (
	'01956624-f91a-7487-8b0e-3a05144eee46',
	'01956624-e036-758b-bf44-a720537ac5cd',
	'01956624-b06d-7242-b623-454807841b0e',
	'01956623-5d4f-7dfd-9d45-c86baa4da2b7',
	'01956620-1a03-7b8a-888d-63ce64763a17'
);

delete from organizations where id in (
	'01956624-f91a-7487-8b0e-3a05144eee46',
	'01956624-e036-758b-bf44-a720537ac5cd',
	'01956624-b06d-7242-b623-454807841b0e',
	'01956623-5d4f-7dfd-9d45-c86baa4da2b7',
	'01956620-1a03-7b8a-888d-63ce64763a17'
);

alter table users
	drop column role_id,
	drop column external_id
;

drop table roles;
