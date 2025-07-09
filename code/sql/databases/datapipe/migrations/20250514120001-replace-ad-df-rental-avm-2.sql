-- +migrate Up

drop table ad_df_rental_avm;

alter table ad_df_rental_avm_2
	rename to ad_df_rental_avm;

-- +migrate Down
