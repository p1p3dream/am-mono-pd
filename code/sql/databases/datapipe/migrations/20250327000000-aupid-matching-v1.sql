-- +migrate Up

create index if not exists idx_fa_df_address_fips on fa_df_address (fips);
create index if not exists idx_fa_df_address_upper_full_street_address on fa_df_address (upper(trim(full_street_address)));

create index if not exists idx_ad_df_listing_mls_listing_county_fips on ad_df_listing (mls_listing_county_fips);
create index if not exists idx_ad_df_listing_upper_property_address_full on ad_df_listing (upper(trim(property_address_full)));

-- +migrate Down

drop index if exists idx_ad_df_listing_upper_property_address_full;
drop index if exists idx_ad_df_listing_mls_listing_county_fips;

drop index if exists idx_fa_df_address_upper_full_street_address;
drop index if exists idx_fa_df_address_fips;
