-- +migrate Up

drop index idx_data_file_objects_locator;

alter table data_file_objects
	alter column directory_id drop not null
;

create index idx_data_file_objects_locator
	on data_file_objects (directory_id, file_type, hash, parent_file_id)
;

-- +migrate Down

drop index idx_data_file_objects_locator;

alter table data_file_objects
	alter column directory_id set not null
;

create index idx_data_file_objects_locator
	on data_file_objects (directory_id, file_type, hash)
;
