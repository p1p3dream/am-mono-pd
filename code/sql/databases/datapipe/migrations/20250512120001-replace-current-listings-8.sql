-- +migrate Up

drop table current_listings;

alter table current_listings_8
	rename to current_listings;

-- +migrate Down
