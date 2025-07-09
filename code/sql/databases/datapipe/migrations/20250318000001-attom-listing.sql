-- +migrate Up

--------------------------------------------------------------------------------
-- AttomData - Data File - Listing.
--------------------------------------------------------------------------------

create table ad_df_listing (
	am_id         uuid,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,

	attom_id                               bigint not null,
	mls_record_id                          bigint not null,
	mls_listing_id                         bigint not null,
	status_change_date                     date not null,
	property_address_full                  text,
	property_address_house_number          text,
	property_address_street_direction      text,
	property_address_street_name           text,
	property_address_street_suffix         text,
	property_address_street_post_direction text,

	property_address_unit_prefix text,
	property_address_unit_value  text,
	property_address_city        text,
	property_address_state       text,
	property_address_zip         text, --is nullable
	property_address_zip4        text,
	situs_county                 text,
	township                     text,
	mls_listing_address          text,
	mls_listing_city             text, --is nullable

	mls_listing_state         text, --is nullable
	mls_listing_zip           text, --is nullable
	mls_listing_county_fips   text, --is nullable
	mls_number                text,
	mls_source                text,
	listing_status            text,
	mls_sold_date             date,
	mls_sold_price            integer,
	assessor_last_sale_date   date,
	assessor_last_sale_amount text,

	market_value               integer,
	market_value_date          date,
	avg_market_price_per_sq_ft integer,
	listing_date               date,
	latest_listing_price       integer,
	previous_listing_price     integer,
	latest_price_change_date   date,
	pending_date               date,
	special_listing_conditions text,
	original_listing_date      date,

	original_listing_price  integer,
	lease_option            text,
	lease_term              text,
	lease_includes          text,
	concessions             text,
	concessions_amount      integer,
	concessions_comments    text,
	contingency_date        date,
	contingency_description text,
	mls_property_type       text,

	mls_property_sub_type   text,
	attom_property_type     text,
	attom_property_sub_type text,
	ownership_description   text,
	latitude                double precision, --is nullable
	longitude               double precision, --is nullable
	apn_formatted           text, --is nullable
	legal_description       text,
	legal_subdivision       text,
	days_on_market          integer,

	cumulative_days_on_market     integer,
	listing_agent_full_name       text,
	listing_agent_mls_id          text,
	listing_agent_state_license   text,
	listing_agent_aor             text,
	listing_agent_preferred_phone text,
	listing_agent_email           text,
	listing_office_name           text,
	listing_office_mls_id         text,
	listing_office_aor            text,

	listing_office_phone             text,
	listing_office_email             text,
	listing_co_agent_full_name       text,
	listing_co_agent_mls_id          text,
	listing_co_agent_state_license   text,
	listing_co_agent_aor             text,
	listing_co_agent_preferred_phone text,
	listing_co_agent_email           text,
	listing_co_agent_office_name     text,
	listing_co_agent_office_mls_id   text,

	listing_co_agent_office_aor   text,
	listing_co_agent_office_phone text,
	listing_co_agent_office_email text,
	buyer_agent_full_name         text,
	buyer_agent_mls_id            text,
	buyer_agent_state_license     text,
	buyer_agent_aor               text,
	buyer_agent_preferred_phone   text,
	buyer_agent_email             text,
	buyer_office_name             text,

	buyer_office_mls_id            text,
	buyer_office_aor               text,
	buyer_office_phone             text,
	buyer_office_email             text,
	buyer_co_agent_full_name       text,
	buyer_co_agent_mls_id          text,
	buyer_co_agent_state_license   text,
	buyer_co_agent_aor             text,
	buyer_co_agent_preferred_phone text,
	buyer_co_agent_email           text,

	buyer_co_agent_office_name   text,
	buyer_co_agent_office_mls_id text,
	buyer_co_agent_office_aor    text,
	buyer_co_agent_office_phone  text,
	buyer_co_agent_office_email  text,
	public_listing_remarks       text,
	home_warranty_yn             boolean,
	tax_year_assessed            integer,
	tax_assessed_value_total     integer,
	tax_amount                   integer,

	tax_annual_other      integer,
	owner_name            text,
	owner_vesting         text,
	year_built            integer,
	year_built_effective  integer,
	year_built_source     text,
	new_construction_yn   boolean,
	builder_name          text,
	additional_parcels_yn boolean,
	number_of_lots        integer,

	lot_size_square_feet    double precision,
	lot_size_acres          double precision,
	lot_size_source         text,
	lot_dimensions          text,
	lot_feature_list        text,
	frontage_length         text,
	frontage_type           text,
	frontage_road_type      text,
	living_area_square_feet integer,
	living_area_source      text,

	levels                 text,
	stories                numeric(8,2),
	building_stories_total numeric(5,1),
	building_keywords      text,
	building_area_total    integer,
	number_of_units_total  integer,
	number_of_buildings    integer,
	property_attached_yn   boolean,
	other_structures       text,
	rooms_total            integer,

	bedrooms_total            integer,
	bathrooms_full            numeric(5,1),
	bathrooms_half            integer,
	bathrooms_quarter         integer,
	bathrooms_three_quarters  integer,
	basement_features         text,
	below_grade_square_feet   integer,
	basement_total_sq_ft      integer,
	basement_finished_sq_ft   integer,
	basement_unfinished_sq_ft integer,

	property_condition     text,
	repairs_yn             boolean,
	repairs_description    text,
	disclosures            text,
	construction_materials text,
	garage_yn              boolean,
	attached_garage_yn     boolean,
	garage_spaces          numeric(7,1),
	carport_yn             boolean,
	carport_spaces         double precision,

	parking_features    text,
	parking_other       text,
	open_parking_spaces double precision,
	parking_total       double precision,
	pool_private_yn     boolean,
	pool_features       text,
	occupancy           text,
	view_yn             boolean,
	view_col            text,
	topography          text,

	heating_yn                   boolean,
	heating_features             text,
	cooling_yn                   boolean,
	cooling                      text,
	fireplace_yn                 boolean,
	fireplace                    text,
	fireplace_number             double precision,
	foundation_features          text,
	roof                         text,
	architectural_style_features text,

	patio_and_porch_features text,
	utilities                text,
	electric_included        boolean,
	electric_description     text,
	water_included           boolean,
	water_source             text,
	sewer                    text,
	gas_description          text,
	other_equipment_included text,
	laundry_features         text,

	appliances          text,
	interior_features   text,
	exterior_features   text,
	fencing_features    text,
	pets_allowed        text,
	horse_zoning_yn     boolean,
	senior_community_yn boolean,
	waterbody_name      text,
	waterfront_yn       boolean,
	waterfront_features text,

	zoning_code                text,
	zoning_description         text,
	current_use                text,
	possible_use               text,
	association_yn             boolean,
	association1_name          text,
	association1_phone         text,
	association1_fee           integer,
	association1_fee_frequency text,
	association2_name          text,

	association2_phone         text,
	association2_fee           integer,
	association2_fee_frequency text,
	association_fee_includes   text,
	association_amenities      text,
	school_elementary          text,
	school_elementary_district text,
	school_middle              text,
	school_middle_district     text,
	school_high                text,

	school_high_district             text,
	green_verification_yn            boolean,
	green_building_verification_type text,
	green_energy_efficient           text,
	green_energy_generation          text,
	green_indoor_air_quality         text,
	green_location                   text,
	green_sustainability             text,
	green_water_conservation         text,
	land_lease_yn                    boolean,

	land_lease_amount           numeric(20,4),
	land_lease_amount_frequency text,
	land_lease_expiration_date  date,
	cap_rate                    numeric(20,4),
	gross_income                numeric(20,4),
	income_includes             text,
	gross_scheduled_income      numeric(20,4),
	net_operating_income        numeric(20,4),
	total_actual_rent           numeric(20,4),
	existing_lease_type         text,

	financial_data_source  text,
	rent_control_yn        boolean,
	unit_type_description  text,
	unit_type_furnished    text,
	number_of_units_leased double precision,
	number_of_units_mo_mo  double precision,
	number_of_units_vacant double precision,
	vacancy_allowance      double precision,
	vacancy_allowance_rate double precision,
	operating_expense      numeric(20,4),

	cable_tv_expense              numeric(20,4),
	electric_expense              numeric(20,4),
	fuel_expense                  numeric(20,4),
	furniture_replacement_expense numeric(20,4),
	gardener_expense              numeric(20,4),
	insurance_expense             numeric(20,4),
	operating_expense_includes    text,
	licenses_expense              numeric(20,4),
	maintenance_expense           numeric(20,4),
	manager_expense               numeric(20,4),

	new_taxes_expense               numeric(20,4),
	other_expense                   numeric(20,4),
	pest_control_expense            numeric(20,4),
	pool_expense                    numeric(20,4),
	professional_management_expense numeric(20,4),
	supplies_expense                numeric(20,4),
	trash_expense                   numeric(20,4),
	water_sewer_expense             numeric(20,4),
	workmans_compensation_expense   numeric(20,4),
	owner_pays                      text,

	tenant_pays           text,
	listing_marketing_url text,
	photos_count          integer,
	photo_key             text,
	photo_url_prefix      text,

	primary key (am_id, status_change_date)
)
partition by range (status_change_date)
;

create table ad_df_listing_lt_2000 partition of ad_df_listing
	for values from ('1800-01-01') to ('2000-01-01');
create index idx_ad_df_listing_lt_2000_attom_id
	on ad_df_listing_lt_2000 (attom_id);
create index idx_ad_df_listing_lt_2000_mls_listing_id
	on ad_df_listing_lt_2000 (mls_listing_id);

create table ad_df_listing_2000_2004 partition of ad_df_listing
	for values from ('2000-01-01') to ('2005-01-01');
create index idx_ad_df_listing_2000_2004_attom_id
	on ad_df_listing_2000_2004 (attom_id);
create index idx_ad_df_listing_2000_2004_mls_listing_id
	on ad_df_listing_2000_2004 (mls_listing_id);

create table ad_df_listing_2005_2007 partition of ad_df_listing
	for values from ('2005-01-01') to ('2008-01-01');
create index idx_ad_df_listing_2005_2007_attom_id
	on ad_df_listing_2005_2007 (attom_id);
create index idx_ad_df_listing_2005_2007_mls_listing_id
	on ad_df_listing_2005_2007 (mls_listing_id);

create table ad_df_listing_2008_2009 partition of ad_df_listing
	for values from ('2008-01-01') to ('2010-01-01');
create index idx_ad_df_listing_2008_2009_attom_id
	on ad_df_listing_2008_2009 (attom_id);
create index idx_ad_df_listing_2008_2009_mls_listing_id
	on ad_df_listing_2008_2009 (mls_listing_id);

create table ad_df_listing_2010_2011 partition of ad_df_listing
	for values from ('2010-01-01') to ('2012-01-01');
create index idx_ad_df_listing_2010_2011_attom_id
	on ad_df_listing_2010_2011 (attom_id);
create index idx_ad_df_listing_2010_2011_mls_listing_id
	on ad_df_listing_2010_2011 (mls_listing_id);

create table ad_df_listing_2012_2013 partition of ad_df_listing
	for values from ('2012-01-01') to ('2014-01-01');
create index idx_ad_df_listing_2012_2013_attom_id
	on ad_df_listing_2012_2013 (attom_id);
create index idx_ad_df_listing_2012_2013_mls_listing_id
	on ad_df_listing_2012_2013 (mls_listing_id);

create table ad_df_listing_2014 partition of ad_df_listing
	for values from ('2014-01-01') to ('2015-01-01');
create index idx_ad_df_listing_2014_attom_id
	on ad_df_listing_2014 (attom_id);
create index idx_ad_df_listing_2014_mls_listing_id
	on ad_df_listing_2014 (mls_listing_id);

create table ad_df_listing_2015 partition of ad_df_listing
	for values from ('2015-01-01') to ('2016-01-01');
create index idx_ad_df_listing_2015_attom_id
	on ad_df_listing_2015 (attom_id);
create index idx_ad_df_listing_2015_mls_listing_id
	on ad_df_listing_2015 (mls_listing_id);

create table ad_df_listing_2016 partition of ad_df_listing
	for values from ('2016-01-01') to ('2017-01-01');
create index idx_ad_df_listing_2016_attom_id
	on ad_df_listing_2016 (attom_id);
create index idx_ad_df_listing_2016_mls_listing_id
	on ad_df_listing_2016 (mls_listing_id);

create table ad_df_listing_2017 partition of ad_df_listing
	for values from ('2017-01-01') to ('2018-01-01');
create index idx_ad_df_listing_2017_attom_id
	on ad_df_listing_2017 (attom_id);
create index idx_ad_df_listing_2017_mls_listing_id
	on ad_df_listing_2017 (mls_listing_id);

create table ad_df_listing_2018 partition of ad_df_listing
	for values from ('2018-01-01') to ('2019-01-01');
create index idx_ad_df_listing_2018_attom_id
	on ad_df_listing_2018 (attom_id);
create index idx_ad_df_listing_2018_mls_listing_id
	on ad_df_listing_2018 (mls_listing_id);

create table ad_df_listing_2019 partition of ad_df_listing
	for values from ('2019-01-01') to ('2020-01-01');
create index idx_ad_df_listing_2019_attom_id
	on ad_df_listing_2019 (attom_id);
create index idx_ad_df_listing_2019_mls_listing_id
	on ad_df_listing_2019 (mls_listing_id);

create table ad_df_listing_2020 partition of ad_df_listing
	for values from ('2020-01-01') to ('2021-01-01');
create index idx_ad_df_listing_2020_attom_id
	on ad_df_listing_2020 (attom_id);
create index idx_ad_df_listing_2020_mls_listing_id
	on ad_df_listing_2020 (mls_listing_id);

create table ad_df_listing_2021 partition of ad_df_listing
	for values from ('2021-01-01') to ('2022-01-01');
create index idx_ad_df_listing_2021_attom_id
	on ad_df_listing_2021 (attom_id);
create index idx_ad_df_listing_2021_mls_listing_id
	on ad_df_listing_2021 (mls_listing_id);

create table ad_df_listing_2022 partition of ad_df_listing
	for values from ('2022-01-01') to ('2023-01-01');
create index idx_ad_df_listing_2022_attom_id
	on ad_df_listing_2022 (attom_id);
create index idx_ad_df_listing_2022_mls_listing_id
	on ad_df_listing_2022 (mls_listing_id);

create table ad_df_listing_2023 partition of ad_df_listing
	for values from ('2023-01-01') to ('2024-01-01');
create index idx_ad_df_listing_2023_attom_id
	on ad_df_listing_2023 (attom_id);
create index idx_ad_df_listing_2023_mls_listing_id
	on ad_df_listing_2023 (mls_listing_id);

create table ad_df_listing_2024 partition of ad_df_listing
	for values from ('2024-01-01') to ('2025-01-01');
create index idx_ad_df_listing_2024_attom_id
	on ad_df_listing_2024 (attom_id);
create index idx_ad_df_listing_2024_mls_listing_id
	on ad_df_listing_2024 (mls_listing_id);

create table ad_df_listing_2025 partition of ad_df_listing
	for values from ('2025-01-01') to ('2026-01-01');
create index idx_ad_df_listing_2025_attom_id
	on ad_df_listing_2025 (attom_id);
create index idx_ad_df_listing_2025_mls_listing_id
	on ad_df_listing_2025 (mls_listing_id);

create table ad_df_listing_2026 partition of ad_df_listing
	for values from ('2026-01-01') to ('2027-01-01');
create index idx_ad_df_listing_2026_attom_id
	on ad_df_listing_2026 (attom_id);
create index idx_ad_df_listing_2026_mls_listing_id
	on ad_df_listing_2026 (mls_listing_id);

create table ad_df_listing_gt_2026 partition of ad_df_listing
	for values from ('2027-01-01') to ('9999-12-31');
create index idx_ad_df_listing_gt_2026_attom_id
	on ad_df_listing_gt_2026 (attom_id);
create index idx_ad_df_listing_gt_2026_mls_listing_id
	on ad_df_listing_gt_2026 (mls_listing_id);

-- +migrate Down

drop table ad_df_listing;
