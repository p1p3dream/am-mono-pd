-- +migrate Up

alter table fips
	add column county text
;

create index idx_fips_upper_county on fips (upper(county));

-- +migrate Down

alter table fips
	drop column county;
