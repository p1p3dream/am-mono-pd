-- +migrate Up

create index idx_ad_df_assessor_upper_property_address_full_2
	on ad_df_assessor(upper(property_address_full));

drop index if exists idx_ad_df_assessor_upper_property_address_full;

alter index idx_ad_df_assessor_upper_property_address_full_2
	rename to idx_ad_df_assessor_upper_property_address_full;

-- +migrate Down
