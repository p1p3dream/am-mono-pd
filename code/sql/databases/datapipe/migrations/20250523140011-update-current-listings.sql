-- +migrate Up

alter table current_listings
	add column asr_property_sub_type text;

create index idx_current_listings_asr_property_sub_type_upper on
	current_listings (upper(asr_property_sub_type));

-- +migrate Down

alter table current_listings
	drop column asr_property_sub_type;
