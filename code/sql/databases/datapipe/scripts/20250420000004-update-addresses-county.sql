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

		update addresses
		set county = (
			select county
			from fips
			where fips = current_fips
		)
		where
			upper(fips) = current_fips
			and upper(county) is null
		;

		commit;
	end loop;
end $$;
