-- +migrate Up

alter table data_file_objects
	add column worker_id uuid;

-- +migrate Down

alter table data_file_objects
	drop column worker_id;
