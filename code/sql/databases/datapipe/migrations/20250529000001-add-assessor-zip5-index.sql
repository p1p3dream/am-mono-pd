-- +migrate Up

create index idx_ad_df_assessor_zip
	on ad_df_assessor(property_address_zip)
;

create index idx_fa_df_assessor_zip
	on fa_df_assessor(situs_zip5)
;

create table zip5 (
	zip5 text not null primary key
);

insert into zip5 (
	select distinct(situs_zip5)
	from fa_df_assessor
	where situs_zip5 is not null
);

insert into zip5 (
	select distinct(property_address_zip)
	from ad_df_assessor
	where
		property_address_zip is not null
		and not exists (
			select 1
			from zip5
			where zip5.zip5 = ad_df_assessor.property_address_zip
		)
);

-- +migrate Down

drop table zip5;
drop index idx_fa_df_assessor_zip;
drop index idx_ad_df_assessor_zip;
