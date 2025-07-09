-- +migrate Up

create table current_listings_9 (
	am_id         uuid not null,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	am_listing_id uuid not null,
	attom_id      bigint not null,
	aupid         uuid not null,

	status_change_date      date not null,
	ouid                    text,
	mls_number              text,
	mls_property_type       text,
	mls_property_sub_type   text,
	bedrooms_total          integer,
	bathrooms_full          numeric(5, 1),
	latest_listing_price    integer,
	living_area_square_feet integer,
	year_built              integer,
	listing_status          text,
	zip5                    text,

	am_property_type      text,
	am_property_sub_type  text,

	location_3857 geometry,
	location_4326 geometry,

	primary key (am_id)
);

insert into current_listings_9
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
		l.mls_source,
		l.mls_number,
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
	mp.ouid,
	regexp_replace(rl.mls_number, '[^0-9]', '', 'g'),
	rl.mls_property_type,
	rl.mls_property_sub_type,
	rl.bedrooms_total,
	rl.bathrooms_full,
	rl.latest_listing_price,
	rl.living_area_square_feet,
	rl.year_built,
	rl.listing_status,
	addresses.zip5,
	null,
	null,
	ad_geom.location_3857,
	ad_geom.location as location_4326
from ranked_listings rl
join ad_geom on rl.attom_id = ad_geom.attom_id
join properties on rl.attom_id = properties.ad_attom_id
join addresses on properties.address_id = addresses.id
left join mls_providers mp on rl.mls_source = mp.mls_source
where rl.row_num = 1
;

create index idx_current_listings_9_attom_id
	on current_listings_9 (attom_id);

create index idx_current_listings_9_aupid
	on current_listings_9 (aupid);

create index idx_current_listings_9_ouid
	on current_listings_9 (ouid);

create index idx_current_listings_9_mls_number
	on current_listings_9 (mls_number);

create index idx_current_listings_9_status_change_date
	on current_listings_9 (status_change_date);

create index idx_current_listings_9_mls_property_type
	on current_listings_9 (mls_property_type);

create index idx_current_listings_9_mls_property_sub_type
	on current_listings_9 (mls_property_sub_type);

create index idx_current_listings_9_bedrooms_total
	on current_listings_9 (bedrooms_total);

create index idx_current_listings_9_bathrooms_full
	on current_listings_9 (bathrooms_full);

create index idx_current_listings_9_latest_listing_price
	on current_listings_9 (latest_listing_price);

create index idx_current_listings_9_living_area_square_feet
	on current_listings_9 (living_area_square_feet);

create index idx_current_listings_9_year_built
	on current_listings_9 (year_built);

create index idx_current_listings_9_listing_status
	on current_listings_9 (listing_status);

create index idx_current_listings_9_zip5
	on current_listings_9 (zip5);

create index idx_current_listings_9_am_property_type
	on current_listings_9 (am_property_type);

create index idx_current_listings_9_am_property_sub_type
	on current_listings_9 (am_property_sub_type);

create index idx_current_listings_9_location_3857
	on current_listings_9 using gist (location_3857);

create index idx_current_listings_9_location_4326
	on current_listings_9 using gist (location_4326);

-- +migrate Down

drop table current_listings_9;
