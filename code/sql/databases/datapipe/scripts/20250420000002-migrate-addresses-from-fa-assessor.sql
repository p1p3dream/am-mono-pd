do $$
declare
	current_fips text;
	total_fips integer := (select count(*) from fips);

	current_county text;
	addresses_row_count integer;

	batch_id text;
	current_index integer := 0;
begin
	raise notice '% Starting addresses population from fa_df_assessor.', clock_timestamp();

	for current_fips in
		select fips
		from fips
		order by fips
		limit 500
		offset 0
	loop
		current_index := current_index + 1;
		batch_id := gen_random_uuid()::text;

		raise notice
			'% Processing fips % - %/% (% %%). BatchId: %.',
			clock_timestamp(),
			current_fips,
			current_index,
			total_fips,
			round((current_index::numeric/total_fips::numeric)*100, 2),
			batch_id
		;

		current_county := (
			select county
			from fips
			where fips = current_fips
		);

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
			jsonb_build_object(
				'batch_id',
				batch_id,
				'property_id',
				properties.id
			) as meta,
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
		where
			properties.fips = current_fips
			and properties.address_id is null
			and properties.fa_property_id is not null
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
			where meta->>'batch_id' = batch_id
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
