-- +migrate Up

alter table data_file_objects
	add column parent_file_id uuid
;

-- +migrate Down

alter table data_file_objects
	drop column parent_file_id
;
