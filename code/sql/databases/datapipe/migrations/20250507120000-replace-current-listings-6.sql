-- +migrate Up

drop table current_listings;

alter table current_listings_6
	rename to current_listings;

-- +migrate Down
