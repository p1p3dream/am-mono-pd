-- +migrate Up

alter table ad_df_assessor
	alter column deed_last_sale_transaction_id type bigint,
	alter column last_ownership_transfer_transaction_id type bigint
;

-- +migrate Down

alter table ad_df_assessor
	alter column deed_last_sale_transaction_id type integer,
	alter column last_ownership_transfer_transaction_id type integer
;
