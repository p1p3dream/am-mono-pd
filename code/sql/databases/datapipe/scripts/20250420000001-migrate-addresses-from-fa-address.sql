do $$
declare
	current_fips text;
	total_fips integer := (select count(*) from fips);

	addresses_row_count integer;

	current_index integer := 0;
begin
	raise notice '% Starting addresses population from fa_df_address.', clock_timestamp();

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
			'fa_df_address' as data_source,
			fa_df_address.place_name as city,
			fa_df_address.county,
			fa_df_address.fips,
			regexp_replace(fa_df_address.full_street_address, '\s+', ' ', 'g') as full_street_address,
			fa_df_address.street_number as house_number,
			fa_df_address.state,
			fa_df_address.street as street_name,
			fa_df_address.post_directional as street_pos_direction,
			fa_df_address.pre_directional as street_pre_direction,
			fa_df_address.street_type as street_suffix,
			fa_df_address.unit_nbr,
			fa_df_address.unit_type,
			fa_df_address.zip5
		from properties
		join fa_df_address on properties.fa_address_master_id = fa_df_address.address_master_id
		where
			properties.fa_address_master_id is not null
			and properties.fips = current_fips
		;

		get diagnostics addresses_row_count = row_count;
		raise notice '% Address count: %.', clock_timestamp(), addresses_row_count;

		commit;

		raise notice '% Updating properties with address_id.', clock_timestamp();

		with created_addresses as (
			select
				id,
				(meta->>'property_id')::uuid as property_id
			from addresses
			where
				fips = current_fips
				and data_source = 'fa_df_address'
				and meta->>'property_id' is not null
		)
		update properties
		set address_id = created_addresses.id
		from created_addresses
		where
			properties.fips = current_fips
			and properties.id = created_addresses.property_id
		;

		commit;
	end loop;
end $$;
