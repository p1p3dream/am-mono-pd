do $$
declare
	current_fips text;
	total_fips integer := (select count(*) from fips);

	current_index integer := 0;
begin
	raise notice '% Start.', clock_timestamp();

	for current_fips in
		select fips
		from fips
		order by fips
		limit 500
		offset 0
	loop
		current_index := current_index + 1;

		raise notice
			'% Processing fips % - %/% (% %%).',
			clock_timestamp(),
			current_fips,
			current_index,
			total_fips,
			round((current_index::numeric/total_fips::numeric)*100, 2)
		;

		with property_city as (
			select
				addresses.id as address_id,
				(
					select normalize_address(situs_city)
					from fa_df_assessor
					where property_id = properties.fa_property_id
				) as city
			from properties
			join addresses on properties.address_id = addresses.id
			where
				upper(addresses.fips) = current_fips
				and upper(addresses.city) is null
				and addresses.data_source = 'fa_df_address'
		)
		update addresses
		set city = property_city.city
		from property_city
		where addresses.id = property_city.address_id
		;

		commit;
	end loop;
end $$;
