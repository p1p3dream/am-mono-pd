-- +migrate Up

create index idx_data_file_directories_meta
	on data_file_directories using gin (meta);

create index idx_data_file_objects_meta
	on data_file_objects using gin (meta);

-- +migrate Down

drop index if exists idx_data_file_objects_meta;
drop index if exists idx_data_file_directories_meta;
