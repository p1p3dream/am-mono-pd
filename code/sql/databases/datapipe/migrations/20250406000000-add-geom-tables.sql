-- +migrate Up

create table ad_geom (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	attom_id bigint unique not null,
	location geometry(Point, 4326)
);

create index idx_ad_geom_attom_id
	on ad_geom (attom_id);

create index idx_ad_geom_location
	on ad_geom using gist (location);

insert into ad_geom (
	id,
	created_at,
	updated_at,
	attom_id,
	location
)
select
	gen_random_uuid(),
	now(),
	now(),
	attomid,
	st_setsrid(st_makepoint(property_longitude, property_latitude), 4326)
from ad_df_assessor
where
	property_latitude is not null
	and property_longitude is not null
;

select property_latitude, property_longitude from public.ad_df_assessor;

create table fa_geom (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	property_id bigint unique not null,
	location    geometry(Point, 4326)
);

create index idx_fa_geom_property_id
	on fa_geom (property_id);

create index idx_fa_geom_location
	on fa_geom using gist (location);

insert into fa_geom (
	id,
	created_at,
	updated_at,
	property_id,
	location
)
select
	gen_random_uuid(),
	now(),
	now(),
	property_id,
	st_setsrid(st_makepoint(situs_longitude, situs_latitude), 4326)
from fa_df_assessor
where
	situs_latitude is not null
	and situs_longitude is not null
;

-- +migrate Down

drop table fa_geom;
drop table ad_geom;
