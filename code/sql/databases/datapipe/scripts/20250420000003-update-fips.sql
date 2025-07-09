do $$
declare
	current_fips text;
	total_fips integer := (select count(*) from fips);

	current_index integer := 0;
	batch_commit_size integer := 100;
begin
	raise notice '% Start.', clock_timestamp();

	for current_fips in
		select fips
		from fips
		where county is null
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

		update fips
		set county = fa_county.name
		-- Ensure we use the most common (recent) county name for the fips.
		from (
			select
				county as name,
				count(*)
			from fa_df_address
			where fips = current_fips
			group by 1
			order by 2 desc
			limit 1
		) as fa_county
		where
			fips.fips = current_fips
			and fips.county is null
		;

		if current_index % batch_commit_size = 0 then
			raise notice '% Commit.', clock_timestamp();
			commit;
		end if;
	end loop;
end $$;
