-- +migrate Up

create table current_listings_3 (
	am_id         uuid not null,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	am_listing_id uuid not null,
	attom_id      bigint not null,
	aupid         uuid not null,

	status_change_date      date not null,
	mls_property_type       text,
	mls_property_sub_type   text,
	bedrooms_total          integer,
	bathrooms_full          numeric(5, 1),
	latest_listing_price    integer,
	living_area_square_feet integer,
	year_built              integer,
	listing_status          text,

	am_property_type      text,
	am_property_sub_type  text,

	location_3857 geometry,
	location_4326 geometry,

	primary key (am_id)
);

insert into current_listings_3
with ranked_listings as (
	select
		l.am_id,
		l.attom_id,
		l.status_change_date,
		l.mls_property_type,
		l.mls_property_sub_type,
		l.bedrooms_total,
		l.bathrooms_full,
		l.latest_listing_price,
		l.living_area_square_feet,
		l.year_built,
		l.listing_status,
		row_number() over (
			partition by l.attom_id
			order by
				l.status_change_date desc
		) as row_num
	from ad_df_listing l
	where l.status_change_date >= '2020-01-01'
)
select
	gen_random_uuid(),
	now(),
	now(),
	null,
	rl.am_id as am_listing_id,
	rl.attom_id,
	properties.id,
	rl.status_change_date,
	rl.mls_property_type,
	rl.mls_property_sub_type,
	rl.bedrooms_total,
	rl.bathrooms_full,
	rl.latest_listing_price,
	rl.living_area_square_feet,
	rl.year_built,
	rl.listing_status,
	null,
	null,
	ad_geom.location_3857,
	ad_geom.location as location_4326
from ranked_listings rl
join ad_geom on rl.attom_id = ad_geom.attom_id
join properties on rl.attom_id = properties.ad_attom_id
where rl.row_num = 1
;

create index idx_current_listings_3_attom_id
	on current_listings_3 (attom_id);

create index idx_current_listings_3_aupid
	on current_listings_3 (aupid);

create index idx_current_listings_3_status_change_date
	on current_listings_3 (status_change_date);

create index idx_current_listings_3_mls_property_type
	on current_listings_3 (mls_property_type);

create index idx_current_listings_3_mls_property_sub_type
	on current_listings_3 (mls_property_sub_type);

create index idx_current_listings_3_bedrooms_total
	on current_listings_3 (bedrooms_total);

create index idx_current_listings_3_bathrooms_full
	on current_listings_3 (bathrooms_full);

create index idx_current_listings_3_latest_listing_price
	on current_listings_3 (latest_listing_price);

create index idx_current_listings_3_living_area_square_feet
	on current_listings_3 (living_area_square_feet);

create index idx_current_listings_3_year_built
	on current_listings_3 (year_built);

create index idx_current_listings_3_listing_status
	on current_listings_3 (listing_status);

create index idx_current_listings_3_am_property_type
	on current_listings_3 (am_property_type);

create index idx_current_listings_3_am_property_sub_type
	on current_listings_3 (am_property_sub_type);

create index idx_current_listings_3_location_3857
	on current_listings_3 using gist (location_3857);

create index idx_current_listings_3_location_4326
	on current_listings_3 using gist (location_4326);

-- +migrate Down

drop table current_listings_3;
