-- +migrate Up

alter table addresses
	drop column street_number,
	drop column street_type;

drop index idx_fa_df_assessor_upper_situs_full_street_address;

create index idx_fa_df_assessor_upper_situs_full_street_address
	on fa_df_assessor (upper(situs_full_street_address));

create index idx_fa_df_assessor_upper_situs_city
	on fa_df_assessor (upper(situs_city));

create index idx_fa_df_assessor_upper_situs_state
	on fa_df_assessor (upper(situs_state));

create index idx_addresses_meta_property_id
	on addresses ((meta->>'property_id'));

create index idx_addresses_meta_batch_id
	on addresses ((meta->>'batch_id'));

create index idx_addresses_data_source
	on addresses (data_source);

-- +migrate Down

drop index idx_addresses_data_source;

drop index idx_addresses_meta_batch_id;

drop index idx_addresses_meta_property_id;

drop index idx_fa_df_assessor_upper_situs_state;

drop index idx_fa_df_assessor_upper_situs_city;

drop index idx_fa_df_assessor_upper_situs_full_street_address;

create index idx_fa_df_assessor_upper_situs_full_street_address
	on fa_df_assessor (upper(trim(situs_full_street_address)))

alter table addresses
	add column street_number text,
	add column street_type text;
