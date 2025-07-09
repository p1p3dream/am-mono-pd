do $$
declare
	current_fips text;
	total_fips integer := (select count(*) from fips);

	current_county text;
	addresses_row_count integer;
	new_address_id uuid;

	current_index integer := 0;
	current_address record;
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

		current_county := (
			select county
			from fips
			where fips = current_fips
		);

		addresses_row_count := 0;

		for current_address in
			select
				properties.id as property_id,
				addresses.id as address_id,
				fa_df_assessor.am_id as fa_assessor_id
			from properties
			join addresses on properties.address_id = addresses.id
			join fa_df_assessor on properties.fa_property_id = fa_df_assessor.property_id
			where
				upper(addresses.fips) = current_fips
				and addresses.data_source = 'fa_df_address'
				and upper(addresses.full_street_address) <> upper(fa_df_assessor.situs_full_street_address)
		loop
			-- Ensure a clean slate for this property's address.
			delete from addresses
			where id = current_address.address_id;

			insert into addresses (
				id,
				created_at,
				updated_at,
				meta,
				data_source,
				city,
				county,
				fips,
				full_street_address,
				house_number,
				state,
				street_name,
				street_pos_direction,
				street_pre_direction,
				street_suffix,
				unit_nbr,
				unit_type,
				zip5
			)
			select
				gen_random_uuid() as id,
				now() as created_at,
				now() as updated_at,
				jsonb_build_object('property_id', properties.id) as meta,
				'fa_df_assessor' as data_source,
				normalize_address(fa_df_assessor.situs_city),
				current_county as county,
				fa_df_assessor.fips,
				normalize_address(fa_df_assessor.situs_full_street_address) as full_street_address,
				fa_df_assessor.situs_house_nbr as house_number,
				fa_df_assessor.situs_state,
				normalize_address(fa_df_assessor.situs_street) as street_name,
				fa_df_assessor.situs_direction_right as street_pos_direction,
				fa_df_assessor.situs_direction_left as street_pre_direction,
				initcap(fa_df_assessor.situs_mode) as street_suffix,
				fa_df_assessor.situs_unit_nbr,
				initcap(fa_df_assessor.situs_unit_type),
				fa_df_assessor.situs_zip5
			from properties
			join fa_df_assessor on properties.fa_property_id = fa_df_assessor.property_id
			where fa_df_assessor.am_id = current_address.fa_assessor_id
			returning
				id
			into new_address_id
			;

			update properties
			set address_id = new_address_id
			where id = current_address.property_id;

			addresses_row_count := addresses_row_count+1;
		end loop;

		raise notice '% Address count: %.', clock_timestamp(), addresses_row_count;

		commit;
	end loop;
end $$;
