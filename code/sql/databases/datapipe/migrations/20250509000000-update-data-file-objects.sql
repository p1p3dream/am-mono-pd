-- +migrate Up

create index idx_data_file_objects_queue
	on data_file_objects (
		status,
		coalesce(cardinality(priorities), 0),
		parent_file_id
	)
;

alter table data_file_objects
	add column file_size bigint;

-- +migrate Down

alter table data_file_objects
	drop column file_size;

drop index idx_data_file_objects_queue;
