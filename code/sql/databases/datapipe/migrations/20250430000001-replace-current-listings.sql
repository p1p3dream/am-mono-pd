-- +migrate Up

drop table current_listings;

alter table current_listings_2
	rename to current_listings;

-- +migrate Down
