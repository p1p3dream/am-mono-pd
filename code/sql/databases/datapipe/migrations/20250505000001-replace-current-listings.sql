-- +migrate Up

drop table current_listings;

alter table current_listings_3
	rename to current_listings;

-- +migrate Down
