-- +migrate Up

alter table ad_df_listing
	add column if not exists current_status bool;

-- +migrate Down

alter table ad_df_listing
	drop column if exists current_status;
