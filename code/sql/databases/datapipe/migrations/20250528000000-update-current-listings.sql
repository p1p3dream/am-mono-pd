-- +migrate Up

update current_listings
set asr_property_sub_type = (
	select ad_df_listing.attom_property_sub_type
	from ad_df_listing
	where current_listings.am_listing_id = ad_df_listing.am_id
);

-- +migrate Down

update current_listings
set asr_property_sub_type = null;
