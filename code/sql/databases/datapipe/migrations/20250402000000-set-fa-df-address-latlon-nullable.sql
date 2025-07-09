-- +migrate Up

alter table fa_df_address
	alter column latitude drop not null,
	alter column longitude drop not null
;
-- +migrate Down

alter table fa_df_address
	alter column latitude set not null,
	alter column longitude set not null
;
