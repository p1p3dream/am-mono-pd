-- +migrate Up

create table if not exists mls_providers (
	id         uuid primary key default gen_random_uuid(),
	ouid       text not null,
	mls_source text not null
);

create index if not exists idx_mls_providers_ouid on mls_providers (ouid);
create index if not exists idx_mls_providers_mls_source on mls_providers (mls_source);

-- +migrate Down

drop table mls_providers;
