-- +migrate Up

alter table data_file_directories
	add column priorities integer[];
create index idx_data_file_directories_status_priorities on data_file_directories (partner_id, status, priorities);

alter table data_file_objects
	add column priorities integer[];
create index idx_data_file_objects_status_priorities on data_file_objects (directory_id, status, priorities);

-- +migrate Down

drop index idx_data_file_objects_status_priorities;
alter table data_file_objects drop column priorities;

drop index idx_data_file_directories_status_priorities;
alter table data_file_directories drop column priorities;
