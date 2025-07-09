-- +migrate Up

create table fa_df_avm_power_2025h2 partition of fa_df_avm_power
	for values from ('2025-07-01') to ('2026-01-01');
create index idx_fa_df_avm_power_2025h2
	on fa_df_avm_power_2025h2 (property_id);

create table fa_df_avm_power_2026h1 partition of fa_df_avm_power
	for values from ('2026-01-01') to ('2026-07-01');
create index idx_fa_df_avm_power_2026h1
	on fa_df_avm_power_2026h1 (property_id);

create table fa_df_avm_power_2026h2 partition of fa_df_avm_power
	for values from ('2026-07-01') to ('2027-01-01');
create index idx_fa_df_avm_power_2026h2
	on fa_df_avm_power_2026h2 (property_id);

-- +migrate Down
