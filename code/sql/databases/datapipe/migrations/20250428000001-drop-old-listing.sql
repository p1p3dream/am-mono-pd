-- +migrate Up

drop table ad_df_listing;

alter table ad_df_listing_2
	rename to ad_df_listing;

-- +migrate Down
