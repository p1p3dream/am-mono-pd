-- +migrate Up

-- Attom Data.

alter table ad_df_assessor
	add column am_deleted_at timestamp with time zone;

create index ad_df_assessor_deleted
	on ad_df_assessor (attomid, am_deleted_at)
	where am_deleted_at is not null
;

create table ad_assessor_history (
	like ad_df_assessor including all
);

alter table ad_assessor_history
	add column am_archived_at timestamp with time zone not null;

-- First American.

alter table fa_df_assessor
	add column am_deleted_at timestamp with time zone;

create index fa_df_assessor_deleted
	on fa_df_assessor (property_id, am_deleted_at)
	where am_deleted_at is not null
;

create table fa_assessor_history (
	like fa_df_assessor including all
);

alter table fa_assessor_history
	add column am_archived_at timestamp with time zone not null;

-- +migrate Down

drop table fa_assessor_history;

alter table fa_df_assessor
	drop column am_deleted_at;

drop table ad_assessor_history;

alter table ad_df_assessor
	drop column am_deleted_at;
