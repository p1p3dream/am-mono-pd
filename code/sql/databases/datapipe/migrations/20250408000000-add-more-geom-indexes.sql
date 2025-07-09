-- +migrate Up

alter table ad_geom add column location_3857 geometry;
update ad_geom set location_3857 = st_transform(location, 3857);
create index idx_ad_geom_location_3857 on ad_geom using gist(location_3857);

alter table fa_geom add column location_3857 geometry;
update fa_geom set location_3857 = st_transform(location, 3857);
create index idx_fa_geom_location_3857 on fa_geom using gist(location_3857);

-- +migrate Down

drop index idx_fa_geom_location_3857;
alter table fa_geom drop column location_3857;

drop index idx_ad_geom_location_3857;
alter table ad_geom drop column location_3857;

