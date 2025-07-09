-- +migrate Up

create index idx_fa_df_address_state on fa_df_address (upper(state));
create index idx_fa_df_address_county on fa_df_address (upper(county));
create index idx_fa_df_address_zip5 on fa_df_address (upper(zip5));
create index idx_fa_df_address_street_number on fa_df_address (upper(street_number));
create index idx_fa_df_address_street on fa_df_address (upper(street));
create index idx_fa_df_address_street_type on fa_df_address (upper(street_type));
create index idx_fa_df_address_unit_type on fa_df_address (upper(unit_type));
create index idx_fa_df_address_unit_nbr on fa_df_address (upper(unit_nbr));
create index idx_fa_df_address_pre_directional on fa_df_address (upper(pre_directional));
create index idx_fa_df_address_post_directional on fa_df_address (upper(post_directional));

-- +migrate Down

drop index idx_fa_df_address_post_directional;
drop index idx_fa_df_address_pre_directional;
drop index idx_fa_df_address_unit_number;
drop index idx_fa_df_address_unit_type;
drop index idx_fa_df_address_street_type;
drop index idx_fa_df_address_street;
drop index idx_fa_df_address_street_number;
drop index idx_fa_df_address_zip5;
drop index idx_fa_df_address_county;
drop index idx_fa_df_address_state;
