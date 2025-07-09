-- +migrate Up

drop index if exists idx_ad_df_listing_upper_property_address_full;
drop index if exists idx_ad_df_listing_mls_listing_county_fips;

-- +migrate Down

create index if not exists idx_ad_df_listing_mls_listing_county_fips on ad_df_listing (mls_listing_county_fips);
create index if not exists idx_ad_df_listing_upper_property_address_full on ad_df_listing (upper(trim(property_address_full)));
