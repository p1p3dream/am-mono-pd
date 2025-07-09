-- +migrate Up

create index idx_fa_df_address_place_name on fa_df_address (upper(place_name));

-- +migrate Down

drop index idx_fa_df_address_place_name;
