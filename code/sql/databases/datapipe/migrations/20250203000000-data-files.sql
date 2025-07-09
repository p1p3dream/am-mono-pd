-- +migrate Up

--------------------------------------------------------------------------------
-- Data File Directories.
--------------------------------------------------------------------------------

create table data_file_directories (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	partner_id          uuid not null,
	parent_directory_id uuid,
	status              integer not null,
	path                text not null,
	name                text not null
);

create index idx_data_file_directories_path
	on data_file_directories (partner_id, path);

--------------------------------------------------------------------------------
-- Data File Objects.
--------------------------------------------------------------------------------

create table data_file_objects (
	id         uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	meta       jsonb,

	directory_id uuid not null,
	file_type    integer not null,
	hash         bytea not null,
	status       integer not null,
	record_count integer not null,
	file_dir     text,
	file_name    text
);

create index idx_data_file_objects_locator
	on data_file_objects (directory_id, file_type, hash);

-- +migrate Down

drop table data_file_objects;
drop table data_file_directories;
