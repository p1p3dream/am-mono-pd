-- +migrate Up

drop index idx_properties_fa_property_id;

alter table properties
	alter column fa_property_id drop not null
;

create index idx_properties_fa_property_id
	on properties (fa_property_id)
;

alter table properties
	add column ad_attom_id bigint
;

create index idx_properties_ad_attom_id
	on properties (ad_attom_id)
;

-- +migrate Down

drop index idx_properties_ad_attom_id;

alter table properties
	drop column ad_attom_id
;

drop index idx_properties_fa_property_id;

alter table properties
	alter column fa_property_id set not null
;

create unique index idx_properties_fa_property_id
	on properties (fa_property_id)
;
