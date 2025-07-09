-- +migrate Up

-- +migrate StatementBegin
create or replace function normalize_address(input_text text) returns text as $$
declare
	result text;
	word   text;
	words  text[];
	i      int;
begin
	-- Return null if input is null.
	if input_text is null then
		return null;
	end if;

	-- Convert the input to uppercase for consistent processing.
	input_text := upper(input_text);

	-- Split input into words
	words := regexp_split_to_array(input_text, '\s+');

	-- Normalize each word
	for i in 1..array_length(words, 1) loop
		if words[i] ~ '^MC[A-Z]' then
			-- Handle "MC" prefix: McDonald, McArthur, etc.
			words[i] := 'Mc' || initcap(substring(words[i] from 3));
		else
			-- Regular capitalization
			words[i] := initcap(words[i]);
		end if;
	end loop;

	-- Recombine the words.
	result := array_to_string(words, ' ');

	-- Normalize possessives: change 'S to 's.
	result := regexp_replace(result, '''S(\s|$)', '''s\1', 'g');

	return result;
end;
$$ language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
