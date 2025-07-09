-- +migrate Up

insert into organizations
	(
		id,
		created_at,
		updated_at,
		name
	)
	values
	(
		'01959623-9404-7661-b173-a014029666d0',
		now(),
		now(),
		'Blanket Homes'
	),
	(
		'01959623-9404-769d-8be8-316e485ebb1e',
		now(),
		now(),
		'Inspectify'
	)
;

-- +migrate Down

delete from organizations where id in (
	'01959623-9404-7661-b173-a014029666d0',
	'01959623-9404-769d-8be8-316e485ebb1e'
);
