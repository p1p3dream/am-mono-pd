-- +migrate Up

alter table fa_df_address add column location geometry(point, 4326);

update fa_df_address
set location = st_setsrid(st_makepoint(longitude, latitude), 4326);

create index idx_fa_df_address_location on fa_df_address using gist(location);

-- +migrate Down

alter table fa_df_address drop column location;
