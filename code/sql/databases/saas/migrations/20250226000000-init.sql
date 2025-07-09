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
		'019543c8-8fc8-7ab2-9d6b-982e4ccb11f5',
		now(),
		now(),
		'AbodeMine'
	),
	(
		'019543c9-f9d9-7c48-b1ac-35bd1f09b87d',
		now(),
		now(),
		'Easy Street Offers'
	)
;

insert into users
	(
		id,
		created_at,
		updated_at,
		organization_id,
		username
	)
	values
	(
		'019543cc-1893-746f-a0de-b7410f281e5e',
		now(),
		now(),
		'019543c8-8fc8-7ab2-9d6b-982e4ccb11f5',
		'abodeminebot'
	),
	(
		'019543cd-de46-7afa-9f7f-1e0c94fd3fea',
		now(),
		now(),
		'019543c9-f9d9-7c48-b1ac-35bd1f09b87d',
		'kevinloo'
	)
;

-- +migrate Down

delete from users where id in (
	'019543cc-1893-746f-a0de-b7410f281e5e',
	'019543c9-f9d9-7c48-b1ac-35bd1f09b87d'
);

delete from organizations where id in (
	'019543c8-8fc8-7ab2-9d6b-982e4ccb11f5',
	'019543c9-f9d9-7c48-b1ac-35bd1f09b87d'
);
