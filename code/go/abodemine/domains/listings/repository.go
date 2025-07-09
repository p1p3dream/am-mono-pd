package listings

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
)

const (
	countThreshold    = 20
	maxResultsPerPage = 100
	sold              = "Sold"
	underContract     = "Under Contract"
)

type Repository interface {
	SearchMlsListings(r *arc.Request, in *SearchListingsInput, radius float64) (*SearchListingsOutput, error)

	SelectListingRecord(r *arc.Request, in *SelectListingRecordInput) (*SelectListingRecordOutput, error)
	GetLatLonFromAupid(r *arc.Request, aupid string) (float64, float64, error)
	GetLatLonFromAddressId(r *arc.Request, addressId uuid.UUID) (float64, float64, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (rep *repository) SearchMlsListings(r *arc.Request, in *SearchListingsInput, radius float64) (*SearchListingsOutput, error) {
	sql, args, err := searchByRadius(in, radius)
	if err != nil {
		return nil, &errors.Object{
			Id:     "7c4e89c3-04a1-4115-b7f8-6a7eec072276",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "9badd085-65a2-4043-a7f5-ee3ec0c1b813")
	}
	defer rows.Close()

	out := &SearchListingsOutput{}
	for rows.Next() {
		var property ListingComparison

		err := rows.Scan(
			&property.StatusChangeDate,
			&property.FullStreetAddress,
			&property.HouseNumber,
			&property.StreetPreDirection,
			&property.StreetName,
			&property.StreetSuffix,
			&property.StreetPostDirection,
			&property.UnitType,
			&property.UnitNumber,
			&property.City,
			&property.State,
			&property.Zip5,
			&property.Zip4,
			&property.County,
			&property.Township,
			&property.MlsListingAddress,
			&property.MlsListingCity,
			&property.MlsListingState,
			&property.MlsListingZip,
			&property.MlsListingCountyFips,
			&property.MlsNumber,
			&property.MlsSource,
			&property.Ouid,
			&property.ListingStatus,
			&property.MlsSoldDate,
			&property.MlsSoldPrice,
			&property.AssessorLastSaleDate,
			&property.AssessorLastSaleAmount,
			&property.MarketValue,
			&property.MarketValueDate,
			&property.AvgMarketPricePerSqFt,
			&property.ListingDate,
			&property.LatestListingPrice,
			&property.PreviousListingPrice,
			&property.LatestPriceChangeDate,
			&property.PendingDate,
			&property.SpecialListingConditions,
			&property.OriginalListingDate,
			&property.OriginalListingPrice,
			&property.LeaseOption,
			&property.LeaseTerm,
			&property.LeaseIncludes,
			&property.Concessions,
			&property.ConcessionsAmount,
			&property.ConcessionsComments,
			&property.ContingencyDate,
			&property.ContingencyDescription,
			&property.AsrPropertySubType,
			&property.OwnershipDescription,
			&property.Latitude,
			&property.Longitude,
			&property.Apn,
			&property.LegalDescription,
			&property.LegalSubdivision,
			&property.DaysOnMarket,
			&property.CumulativeDaysOnMarket,
			&property.ListingAgentFullName,
			&property.ListingAgentMlsId,
			&property.ListingAgentStateLicense,
			&property.ListingAgentAor,
			&property.ListingAgentPreferredPhone,
			&property.ListingAgentEmail,
			&property.ListingOfficeName,
			&property.ListingOfficeMlsId,
			&property.ListingOfficeAor,
			&property.ListingOfficePhone,
			&property.ListingOfficeEmail,
			&property.ListingCoAgentFullName,
			&property.ListingCoAgentMlsId,
			&property.ListingCoAgentStateLicense,
			&property.ListingCoAgentAor,
			&property.ListingCoAgentPreferredPhone,
			&property.ListingCoAgentEmail,
			&property.ListingCoAgentOfficeName,
			&property.ListingCoAgentOfficeMlsId,
			&property.ListingCoAgentOfficeAor,
			&property.ListingCoAgentOfficePhone,
			&property.ListingCoAgentOfficeEmail,
			&property.BuyerAgentFullName,
			&property.BuyerAgentMlsId,
			&property.BuyerAgentStateLicense,
			&property.BuyerAgentAor,
			&property.BuyerAgentPreferredPhone,
			&property.BuyerAgentEmail,
			&property.BuyerOfficeName,
			&property.BuyerOfficeMlsId,
			&property.BuyerOfficeAor,
			&property.BuyerOfficePhone,
			&property.BuyerOfficeEmail,
			&property.BuyerCoAgentFullName,
			&property.BuyerCoAgentMlsId,
			&property.BuyerCoAgentStateLicense,
			&property.BuyerCoAgentAor,
			&property.BuyerCoAgentPreferredPhone,
			&property.BuyerCoAgentEmail,
			&property.BuyerCoAgentOfficeName,
			&property.BuyerCoAgentOfficeMlsId,
			&property.BuyerCoAgentOfficeAor,
			&property.BuyerCoAgentOfficePhone,
			&property.BuyerCoAgentOfficeEmail,
			&property.PublicListingRemarks,
			&property.HasHomeWarranty,
			&property.TaxYearAssessed,
			&property.TaxAssessedValueTotal,
			&property.TaxAmount,
			&property.TaxAnnualOther,
			&property.OwnerName,
			&property.OwnerVesting,
			&property.YearBuilt,
			&property.YearBuiltEffective,
			&property.YearBuiltSource,
			&property.IsNewConstruction,
			&property.BuilderName,
			&property.HasAdditionalParcels,
			&property.NumberOfLots,
			&property.LivingAreaSquareFeet,
			&property.LotSizeAcres,
			&property.LotSizeSource,
			&property.LotDimensions,
			&property.LotFeatureList,
			&property.FrontageLength,
			&property.FrontageType,
			&property.FrontageRoadType,
			&property.LivingAreaSquareFeet,
			&property.LivingAreaSource,
			&property.Levels,
			&property.Stories,
			&property.BuildingStoriesTotal,
			&property.BuildingKeywords,
			&property.BuildingAreaTotal,
			&property.NumberOfUnitsTotal,
			&property.NumberOfBuildings,
			&property.HasPropertyAttached,
			&property.OtherStructures,
			&property.RoomsTotal,
			&property.BedroomsTotal,
			&property.BathroomsFull,
			&property.BathroomsHalf,
			&property.BathroomsQuarter,
			&property.BathroomsThreeQuarters,
			&property.BasementFeatures,
			&property.BelowGradeSquareFeet,
			&property.BasementTotalSqFt,
			&property.BasementFinishedSqFt,
			&property.BasementUnfinishedSqFt,
			&property.PropertyCondition,
			&property.HasNeededRepairs,
			&property.RepairsDescription,
			&property.Disclosures,
			&property.ConstructionMaterials,
			&property.HasGarage,
			&property.HasAttachedGarage,
			&property.GarageSpaces,
			&property.HasCarport,
			&property.CarportSpaces,
			&property.ParkingFeatures,
			&property.ParkingOther,
			&property.OpenParkingSpaces,
			&property.ParkingTotal,
			&property.HasPrivatePool,
			&property.PoolFeatures,
			&property.Occupancy,
			&property.HasView,
			&property.View,
			&property.Topography,
			&property.HasHeating,
			&property.HeatingFeatures,
			&property.HasCooling,
			&property.Cooling,
			&property.HasFireplace,
			&property.Fireplace,
			&property.FireplaceNumber,
			&property.FoundationFeatures,
			&property.Roof,
			&property.ArchitecturalStyleFeatures,
			&property.PatioAndPorchFeatures,
			&property.Utilities,
			&property.ElectricIncluded,
			&property.ElectricDescription,
			&property.WaterIncluded,
			&property.WaterSource,
			&property.Sewer,
			&property.GasDescription,
			&property.OtherEquipmentIncluded,
			&property.LaundryFeatures,
			&property.Appliances,
			&property.InteriorFeatures,
			&property.ExteriorFeatures,
			&property.FencingFeatures,
			&property.PetsAllowed,
			&property.HorseZoning,
			&property.SeniorCommunity,
			&property.WaterbodyName,
			&property.IsWaterfront,
			&property.WaterfrontFeatures,
			&property.ZoningCode,
			&property.ZoningDescription,
			&property.CurrentUse,
			&property.PossibleUse,
			&property.HasAssociation,
			&property.Association1Name,
			&property.Association1Phone,
			&property.Association1Fee,
			&property.Association1FeeFrequency,
			&property.Association2Name,
			&property.Association2Phone,
			&property.Association2Fee,
			&property.Association2FeeFrequency,
			&property.AssociationFeeIncludes,
			&property.AssociationAmenities,
			&property.SchoolElementary,
			&property.SchoolElementaryDistrict,
			&property.SchoolMiddle,
			&property.SchoolMiddleDistrict,
			&property.SchoolHigh,
			&property.SchoolHighDistrict,
			&property.HasGreenVerification,
			&property.GreenBuildingVerificationType,
			&property.GreenEnergyEfficient,
			&property.GreenEnergyGeneration,
			&property.GreenIndoorAirQuality,
			&property.GreenLocation,
			&property.GreenSustainability,
			&property.GreenWaterConservation,
			&property.HasLandLease,
			&property.LandLeaseAmount,
			&property.LandLeaseAmountFrequency,
			&property.LandLeaseExpirationDate,
			&property.CapRate,
			&property.GrossIncome,
			&property.IncomeIncludes,
			&property.GrossScheduledIncome,
			&property.NetOperatingIncome,
			&property.TotalActualRent,
			&property.ExistingLeaseType,
			&property.FinancialDataSource,
			&property.HasRentControl,
			&property.UnitTypeDescription,
			&property.UnitTypeFurnished,
			&property.NumberOfUnitsLeased,
			&property.NumberOfUnitsMoMo,
			&property.NumberOfUnitsVacant,
			&property.VacancyAllowance,
			&property.VacancyAllowanceRate,
			&property.OperatingExpense,
			&property.CableTvExpense,
			&property.ElectricExpense,
			&property.FuelExpense,
			&property.FurnitureReplacementExpense,
			&property.GardenerExpense,
			&property.InsuranceExpense,
			&property.OperatingExpenseIncludes,
			&property.LicensesExpense,
			&property.MaintenanceExpense,
			&property.ManagerExpense,
			&property.NewTaxesExpense,
			&property.OtherExpense,
			&property.PestControlExpense,
			&property.PoolExpense,
			&property.ProfessionalManagementExpense,
			&property.SuppliesExpense,
			&property.TrashExpense,
			&property.WaterSewerExpense,
			&property.WorkmansCompensationExpense,
			&property.OwnerPays,
			&property.TenantPays,
			&property.ListingMarketingUrl,
			&property.PhotosCount,
			&property.PhotoKey,
			&property.PhotoUrlPrefix,
		)
		if err != nil {
			return nil, &errors.Object{
				Id:     "5a6b7c8d-9e0f-1g2h-3i4j-5k6l7m8n9o0p",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to scan property row.",
				Cause:  err.Error(),
			}
		}

		out.PropertyListing = append(out.PropertyListing, property)
	}

	if err := rows.Err(); err != nil {
		return nil, &errors.Object{
			Id:     "bcb618db-bed9-4aa9-9e6f-890f57d49fcb",
			Code:   errors.Code_UNKNOWN,
			Detail: "Error iterating through result rows.",
			Cause:  err.Error(),
		}
	}

	result, err := countSearchResults(r, in, radius)
	if err != nil {
		return nil, errors.Forward(err, "2a249bb5-1edc-4e1e-b231-f35c8e211a0c")
	}
	out.Pagination = &PaginationInfo{
		TotalResults: result,
		PageLimit:    *in.PageLimit,
		PageNumber:   in.PageNumber,
	}

	return out, nil
}

func searchByRadius(filter *SearchListingsInput, radius float64) (string, []any, error) {
	offset := (filter.PageNumber - 1) * *filter.PageLimit

	cte, initialWhere := searchByGeometry(filter.GeoFilter, radius)
	query := fmt.Sprintf(`
		 %s
			SELECT
				ad_df_listing.status_change_date,
				ad_df_listing.property_address_full,
				ad_df_listing.property_address_house_number,
				ad_df_listing.property_address_street_direction,
				ad_df_listing.property_address_street_name,
				ad_df_listing.property_address_street_suffix,
				ad_df_listing.property_address_street_post_direction,
				ad_df_listing.property_address_unit_prefix,
				ad_df_listing.property_address_unit_value,
				ad_df_listing.property_address_city,
				ad_df_listing.property_address_state,
				ad_df_listing.property_address_zip,
				ad_df_listing.property_address_zip4,
				ad_df_listing.situs_county,
				ad_df_listing.township,
				ad_df_listing.mls_listing_address,
				ad_df_listing.mls_listing_city,
				ad_df_listing.mls_listing_state,
				ad_df_listing.mls_listing_zip,
				ad_df_listing.mls_listing_county_fips,
				cl.mls_number,
				ad_df_listing.mls_source,
				cl.ouid,
				ad_df_listing.listing_status,
				ad_df_listing.mls_sold_date,
				ad_df_listing.mls_sold_price,
				ad_df_listing.assessor_last_sale_date,
				ad_df_listing.assessor_last_sale_amount,
				ad_df_listing.market_value,
				ad_df_listing.market_value_date,
				ad_df_listing.avg_market_price_per_sq_ft,
				ad_df_listing.listing_date,
				ad_df_listing.latest_listing_price,
				ad_df_listing.previous_listing_price,
				ad_df_listing.latest_price_change_date,
				ad_df_listing.pending_date,
				ad_df_listing.special_listing_conditions,
				ad_df_listing.original_listing_date,
				ad_df_listing.original_listing_price,
				ad_df_listing.lease_option,
				ad_df_listing.lease_term,
				ad_df_listing.lease_includes,
				ad_df_listing.concessions,
				ad_df_listing.concessions_amount,
				ad_df_listing.concessions_comments,
				ad_df_listing.contingency_date,
				ad_df_listing.contingency_description,
				ad_df_listing.attom_property_sub_type,
				ad_df_listing.ownership_description,
				ad_df_listing.latitude,
				ad_df_listing.longitude,
				ad_df_listing.apn_formatted,
				ad_df_listing.legal_description,
				ad_df_listing.legal_subdivision,
				ad_df_listing.days_on_market,
				ad_df_listing.cumulative_days_on_market,
				ad_df_listing.listing_agent_full_name,
				ad_df_listing.listing_agent_mls_id,
				ad_df_listing.listing_agent_state_license,
				ad_df_listing.listing_agent_aor,
				ad_df_listing.listing_agent_preferred_phone,
				ad_df_listing.listing_agent_email,
				ad_df_listing.listing_office_name,
				ad_df_listing.listing_office_mls_id,
				ad_df_listing.listing_office_aor,
				ad_df_listing.listing_office_phone,
				ad_df_listing.listing_office_email,
				ad_df_listing.listing_co_agent_full_name,
				ad_df_listing.listing_co_agent_mls_id,
				ad_df_listing.listing_co_agent_state_license,
				ad_df_listing.listing_co_agent_aor,
				ad_df_listing.listing_co_agent_preferred_phone,
				ad_df_listing.listing_co_agent_email,
				ad_df_listing.listing_co_agent_office_name,
				ad_df_listing.listing_co_agent_office_mls_id,
				ad_df_listing.listing_co_agent_office_aor,
				ad_df_listing.listing_co_agent_office_phone,
				ad_df_listing.listing_co_agent_office_email,
				ad_df_listing.buyer_agent_full_name,
				ad_df_listing.buyer_agent_mls_id,
				ad_df_listing.buyer_agent_state_license,
				ad_df_listing.buyer_agent_aor,
				ad_df_listing.buyer_agent_preferred_phone,
				ad_df_listing.buyer_agent_email,
				ad_df_listing.buyer_office_name,
				ad_df_listing.buyer_office_mls_id,
				ad_df_listing.buyer_office_aor,
				ad_df_listing.buyer_office_phone,
				ad_df_listing.buyer_office_email,
				ad_df_listing.buyer_co_agent_full_name,
				ad_df_listing.buyer_co_agent_mls_id,
				ad_df_listing.buyer_co_agent_state_license,
				ad_df_listing.buyer_co_agent_aor,
				ad_df_listing.buyer_co_agent_preferred_phone,
				ad_df_listing.buyer_co_agent_email,
				ad_df_listing.buyer_co_agent_office_name,
				ad_df_listing.buyer_co_agent_office_mls_id,
				ad_df_listing.buyer_co_agent_office_aor,
				ad_df_listing.buyer_co_agent_office_phone,
				ad_df_listing.buyer_co_agent_office_email,
				ad_df_listing.public_listing_remarks,
				ad_df_listing.home_warranty_yn,
				ad_df_listing.tax_year_assessed,
				ad_df_listing.tax_assessed_value_total,
				ad_df_listing.tax_amount,
				ad_df_listing.tax_annual_other,
				ad_df_listing.owner_name,
				ad_df_listing.owner_vesting,
				ad_df_listing.year_built,
				ad_df_listing.year_built_effective,
				ad_df_listing.year_built_source,
				ad_df_listing.new_construction_yn,
				ad_df_listing.builder_name,
				ad_df_listing.additional_parcels_yn,
				ad_df_listing.number_of_lots,
				ad_df_listing.lot_size_square_feet,
				ad_df_listing.lot_size_acres,
				ad_df_listing.lot_size_source,
				ad_df_listing.lot_dimensions,
				ad_df_listing.lot_feature_list,
				ad_df_listing.frontage_length,
				ad_df_listing.frontage_type,
				ad_df_listing.frontage_road_type,
				ad_df_listing.living_area_square_feet,
				ad_df_listing.living_area_source,
				ad_df_listing.levels,
				ad_df_listing.stories,
				ad_df_listing.building_stories_total,
				ad_df_listing.building_keywords,
				ad_df_listing.building_area_total,
				ad_df_listing.number_of_units_total,
				ad_df_listing.number_of_buildings,
				ad_df_listing.property_attached_yn,
				ad_df_listing.other_structures,
				ad_df_listing.rooms_total,
				ad_df_listing.bedrooms_total,
				ad_df_listing.bathrooms_full,
				ad_df_listing.bathrooms_half,
				ad_df_listing.bathrooms_quarter,
				ad_df_listing.bathrooms_three_quarters,
				ad_df_listing.basement_features,
				ad_df_listing.below_grade_square_feet,
				ad_df_listing.basement_total_sq_ft,
				ad_df_listing.basement_finished_sq_ft,
				ad_df_listing.basement_unfinished_sq_ft,
				ad_df_listing.property_condition,
				ad_df_listing.repairs_yn,
				ad_df_listing.repairs_description,
				ad_df_listing.disclosures,
				ad_df_listing.construction_materials,
				ad_df_listing.garage_yn,
				ad_df_listing.attached_garage_yn,
				ad_df_listing.garage_spaces,
				ad_df_listing.carport_yn,
				ad_df_listing.carport_spaces,
				ad_df_listing.parking_features,
				ad_df_listing.parking_other,
				ad_df_listing.open_parking_spaces,
				ad_df_listing.parking_total,
				ad_df_listing.pool_private_yn,
				ad_df_listing.pool_features,
				ad_df_listing.occupancy,
				ad_df_listing.view_yn,
				ad_df_listing.view_col,
				ad_df_listing.topography,
				ad_df_listing.heating_yn,
				ad_df_listing.heating_features,
				ad_df_listing.cooling_yn,
				ad_df_listing.cooling,
				ad_df_listing.fireplace_yn,
				ad_df_listing.fireplace,
				ad_df_listing.fireplace_number,
				ad_df_listing.foundation_features,
				ad_df_listing.roof,
				ad_df_listing.architectural_style_features,
				ad_df_listing.patio_and_porch_features,
				ad_df_listing.utilities,
				ad_df_listing.electric_included,
				ad_df_listing.electric_description,
				ad_df_listing.water_included,
				ad_df_listing.water_source,
				ad_df_listing.sewer,
				ad_df_listing.gas_description,
				ad_df_listing.other_equipment_included,
				ad_df_listing.laundry_features,
				ad_df_listing.appliances,
				ad_df_listing.interior_features,
				ad_df_listing.exterior_features,
				ad_df_listing.fencing_features,
				ad_df_listing.pets_allowed,
				ad_df_listing.horse_zoning_yn,
				ad_df_listing.senior_community_yn,
				ad_df_listing.waterbody_name,
				ad_df_listing.waterfront_yn,
				ad_df_listing.waterfront_features,
				ad_df_listing.zoning_code,
				ad_df_listing.zoning_description,
				ad_df_listing.current_use,
				ad_df_listing.possible_use,
				ad_df_listing.association_yn,
				ad_df_listing.association1_name,
				ad_df_listing.association1_phone,
				ad_df_listing.association1_fee,
				ad_df_listing.association1_fee_frequency,
				ad_df_listing.association2_name,
				ad_df_listing.association2_phone,
				ad_df_listing.association2_fee,
				ad_df_listing.association2_fee_frequency,
				ad_df_listing.association_fee_includes,
				ad_df_listing.association_amenities,
				ad_df_listing.school_elementary,
				ad_df_listing.school_elementary_district,
				ad_df_listing.school_middle,
				ad_df_listing.school_middle_district,
				ad_df_listing.school_high,
				ad_df_listing.school_high_district,
				ad_df_listing.green_verification_yn,
				ad_df_listing.green_building_verification_type,
				ad_df_listing.green_energy_efficient,
				ad_df_listing.green_energy_generation,
				ad_df_listing.green_indoor_air_quality,
				ad_df_listing.green_location,
				ad_df_listing.green_sustainability,
				ad_df_listing.green_water_conservation,
				ad_df_listing.land_lease_yn,
				ad_df_listing.land_lease_amount,
				ad_df_listing.land_lease_amount_frequency,
				ad_df_listing.land_lease_expiration_date,
				ad_df_listing.cap_rate,
				ad_df_listing.gross_income,
				ad_df_listing.income_includes,
				ad_df_listing.gross_scheduled_income,
				ad_df_listing.net_operating_income,
				ad_df_listing.total_actual_rent,
				ad_df_listing.existing_lease_type,
				ad_df_listing.financial_data_source,
				ad_df_listing.rent_control_yn,
				ad_df_listing.unit_type_description,
				ad_df_listing.unit_type_furnished,
				ad_df_listing.number_of_units_leased,
				ad_df_listing.number_of_units_mo_mo,
				ad_df_listing.number_of_units_vacant,
				ad_df_listing.vacancy_allowance,
				ad_df_listing.vacancy_allowance_rate,
				ad_df_listing.operating_expense,
				ad_df_listing.cable_tv_expense,
				ad_df_listing.electric_expense,
				ad_df_listing.fuel_expense,
				ad_df_listing.furniture_replacement_expense,
				ad_df_listing.gardener_expense,
				ad_df_listing.insurance_expense,
				ad_df_listing.operating_expense_includes,
				ad_df_listing.licenses_expense,
				ad_df_listing.maintenance_expense,
				ad_df_listing.manager_expense,
				ad_df_listing.new_taxes_expense,
				ad_df_listing.other_expense,
				ad_df_listing.pest_control_expense,
				ad_df_listing.pool_expense,
				ad_df_listing.professional_management_expense,
				ad_df_listing.supplies_expense,
				ad_df_listing.trash_expense,
				ad_df_listing.water_sewer_expense,
				ad_df_listing.workmans_compensation_expense,
				ad_df_listing.owner_pays,
				ad_df_listing.tenant_pays,
				ad_df_listing.listing_marketing_url,
				ad_df_listing.photos_count,
				ad_df_listing.photo_key,
				ad_df_listing.photo_url_prefix
			FROM ad_df_listing
			INNER JOIN (SELECT am_listing_id, ouid, mls_number
			FROM current_listings
			%s
			LIMIT $1 OFFSET $2) cl ON ad_df_listing.am_id = cl.am_listing_id;
	`, cte, buildWhereStr(filter, initialWhere))

	args := []any{
		filter.PageLimit,
		offset,
	}

	return query, args, nil
}

func buildWhereStr(filter *SearchListingsInput, where string) string {
	if filter == nil {
		return where
	}

	var clauses []string
	wasEmpty := (where == "")

	if filter.MinBeds != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.bedrooms_total >= %d", *filter.MinBeds))
	}
	if filter.MaxBeds != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.bedrooms_total <= %d", *filter.MaxBeds))
	}

	if filter.MinBaths != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.bathrooms_full >= %f", *filter.MinBaths))
	}
	if filter.MaxBaths != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.bathrooms_full <= %f", *filter.MaxBaths))
	}

	if filter.MinPrice != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.latest_listing_price >= %d", *filter.MinPrice))
	}
	if filter.MaxPrice != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.latest_listing_price <= %d", *filter.MaxPrice))
	}

	if filter.MinSqFt != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.living_area_square_feet >= %d", *filter.MinSqFt))
	}
	if filter.MaxSqFt != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.living_area_square_feet <= %d", *filter.MaxSqFt))
	}

	if filter.MinYearBuilt != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.year_built >= %d", *filter.MinYearBuilt))
	}
	if filter.MaxYearBuilt != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.year_built <= %d", *filter.MaxYearBuilt))
	}

	if filter.AsrPropertySubType != nil {
		cleaned := strings.ToUpper(strings.ReplaceAll(*filter.AsrPropertySubType, "'", ""))
		clauses = append(clauses, fmt.Sprintf(
			"UPPER(current_listings.asr_property_sub_type) = '%s'",
			cleaned,
		))
	}

	if len(filter.Statuses) > 0 {
		statusValues := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			statusValues[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(status, "'", ""))
		}
		clauses = append(clauses, fmt.Sprintf("current_listings.listing_status IN (%s)", strings.Join(statusValues, ", ")))
	}

	if len(filter.Zip5Codes) > 0 {
		zipValues := make([]string, len(filter.Zip5Codes))
		for i, zip := range filter.Zip5Codes {
			zipValues[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(zip, "'", ""))
		}
		clauses = append(clauses, fmt.Sprintf("current_listings.zip5 IN (%s)", strings.Join(zipValues, ", ")))
	}

	if filter.MinStatusChangeDate != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.status_change_date >= '%s'", strings.ReplaceAll(*filter.MinStatusChangeDate, "'", "")))
	}

	if filter.Ouid != nil {
		clauses = append(clauses, fmt.Sprintf("current_listings.ouid = '%s'", strings.ReplaceAll(*filter.Ouid, "'", "")))
	}

	if len(filter.MlsNumbers) > 0 {
		mlsValues := make([]string, len(filter.MlsNumbers))
		for i, mls := range filter.MlsNumbers {
			mlsValues[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(mls, "'", ""))
		}

		clauses = append(clauses, fmt.Sprintf("current_listings.mls_number IN (%s)", strings.Join(mlsValues, ", ")))
	}

	if len(clauses) > 0 {
		if wasEmpty {
			where = "WHERE " + strings.Join(clauses, " AND ")
		} else {
			where += " AND " + strings.Join(clauses, " AND ")
		}
	}

	if !wasEmpty {
		where += " ORDER BY current_listings.location_3857 <-> reference_geom.geom"
	}

	return where
}

func searchByGeometry(geoFilter *GeoFilter, radius float64) (string, string) {
	if geoFilter == nil {
		return "", ""
	}

	if geoFilter.GeoDistance != nil {
		cteClause := fmt.Sprintf("WITH reference_geom AS (SELECT ST_Transform(ST_SetSRID(ST_MakePoint(%f, %f), 4326), 3857) AS geom)",
			geoFilter.GeoDistance.Location.Lon,
			geoFilter.GeoDistance.Location.Lat)

		radiusInMeters := radius * 1609.34 // Convert miles to meters
		whereClause := fmt.Sprintf(", reference_geom WHERE location_3857 && ST_Expand(reference_geom.geom, %f)", radiusInMeters)

		return cteClause, whereClause
	}

	if geoFilter.GeoPolygon != nil {
		var wktParts []string
		for _, p := range geoFilter.GeoPolygon.Points {
			wktParts = append(wktParts, fmt.Sprintf("%f %f", p.Lon, p.Lat))
		}

		wktPolygon := fmt.Sprintf("POLYGON((%s))", strings.Join(wktParts, ", "))
		cteClause := fmt.Sprintf(
			"WITH reference_geom AS (SELECT ST_Transform(ST_GeomFromText('%s', 4326), 3857) AS geom)",
			wktPolygon,
		)
		whereClause := ", reference_geom WHERE ST_Within(location_3857, reference_geom.geom)"
		return cteClause, whereClause
	}

	return "", ""
}

func countSearchResults(r *arc.Request, filter *SearchListingsInput, radius float64) (int, error) {
	var countResultsQuery string

	if filter.GeoFilter != nil && filter.GeoFilter.GeoDistance != nil {
		radiusInMeters := radius * 1609.34
		whereClause := fmt.Sprintf("current_listings.location_3857 && ST_Expand(reference_geom.geom, %f)", radiusInMeters)
		whereClause = buildWhereStr(filter, whereClause)

		countResultsQuery = fmt.Sprintf(`
			SELECT COUNT(*)
				FROM (
				SELECT 1
				FROM current_listings,
					(
						SELECT ST_Transform(ST_SetSRID(ST_MakePoint(%f,%f), 4326), 3857) AS geom
					) AS reference_geom
				WHERE %s
				) AS limited_results`,
			filter.GeoFilter.GeoDistance.Location.Lon,
			filter.GeoFilter.GeoDistance.Location.Lat,
			whereClause,
		)
	} else if filter.GeoFilter != nil && filter.GeoFilter.GeoPolygon != nil {
		var wktParts []string
		for _, p := range filter.GeoFilter.GeoPolygon.Points {
			wktParts = append(wktParts, fmt.Sprintf("%f %f", p.Lon, p.Lat))
		}

		wktPolygon := fmt.Sprintf("POLYGON((%s))", strings.Join(wktParts, ", "))
		whereClause := "ST_Within(current_listings.location_3857, reference_geom.geom)"
		whereClause = buildWhereStr(filter, whereClause)

		countResultsQuery = fmt.Sprintf(`
			SELECT COUNT(*)
				FROM (
				SELECT 1
				FROM current_listings,
					(
						SELECT ST_Transform(ST_GeomFromText('%s', 4326), 3857) AS geom
					) AS reference_geom
				WHERE %s
				) AS limited_results`, wktPolygon, whereClause)
	} else {
		countResultsQuery = fmt.Sprintf(`
			SELECT COUNT(*)
				FROM (
					SELECT am_listing_id
					FROM current_listings
					%s
				)`, buildWhereStr(filter, ""))
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, countResultsQuery, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, errors.Forward(err, "cfe6a8d5-4d98-48ae-b288-af5ea0e1c2c8")
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.Forward(err, "73715f5b-2844-4a8e-8a0c-32b55df8b81f")
	}

	return count, nil

}

type SelectListingRecordInput struct {
	Aupid *uuid.UUID
}

type SelectListingRecordOutput struct {
	Records []*entities.Listing
}

func (repo *repository) SelectListingRecord(r *arc.Request, in *SelectListingRecordInput) (*SelectListingRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"ad_df_listing.mls_record_id",
			"ad_df_listing.mls_listing_id",
			"ad_df_listing.status_change_date",
			"ad_df_listing.property_address_full",
			"ad_df_listing.property_address_house_number",
			"ad_df_listing.property_address_street_direction",
			"ad_df_listing.property_address_street_name",
			"ad_df_listing.property_address_street_suffix",
			"ad_df_listing.property_address_street_post_direction",
			"ad_df_listing.property_address_unit_prefix",
			"ad_df_listing.property_address_unit_value",
			"ad_df_listing.property_address_city",
			"ad_df_listing.property_address_state",
			"ad_df_listing.property_address_zip",
			"ad_df_listing.property_address_zip4",
			"ad_df_listing.situs_county",
			"ad_df_listing.township",
			"ad_df_listing.mls_listing_address",
			"ad_df_listing.mls_listing_city",
			"ad_df_listing.mls_listing_state",
			"ad_df_listing.mls_listing_zip",
			"ad_df_listing.mls_listing_county_fips",
			"current_listings.mls_number",
			"ad_df_listing.mls_source",
			"current_listings.ouid",
			"ad_df_listing.listing_status",
			"ad_df_listing.mls_sold_date",
			"ad_df_listing.mls_sold_price",
			"ad_df_listing.assessor_last_sale_date",
			"ad_df_listing.assessor_last_sale_amount",
			"ad_df_listing.market_value",
			"ad_df_listing.market_value_date",
			"ad_df_listing.avg_market_price_per_sq_ft",
			"ad_df_listing.listing_date",
			"ad_df_listing.latest_listing_price",
			"ad_df_listing.previous_listing_price",
			"ad_df_listing.latest_price_change_date",
			"ad_df_listing.pending_date",
			"ad_df_listing.special_listing_conditions",
			"ad_df_listing.original_listing_date",
			"ad_df_listing.original_listing_price",
			"ad_df_listing.lease_option",
			"ad_df_listing.lease_term",
			"ad_df_listing.lease_includes",
			"ad_df_listing.concessions",
			"ad_df_listing.concessions_amount",
			"ad_df_listing.concessions_comments",
			"ad_df_listing.contingency_date",
			"ad_df_listing.contingency_description",
			"ad_df_listing.mls_property_type",
			"ad_df_listing.mls_property_sub_type",
			"ad_df_listing.attom_property_type",
			"ad_df_listing.attom_property_sub_type",
			"ad_df_listing.ownership_description",
			"ad_df_listing.latitude",
			"ad_df_listing.longitude",
			"ad_df_listing.apn_formatted",
			"ad_df_listing.legal_description",
			"ad_df_listing.legal_subdivision",
			"ad_df_listing.days_on_market",
			"ad_df_listing.cumulative_days_on_market",
			"ad_df_listing.listing_agent_full_name",
			"ad_df_listing.listing_agent_mls_id",
			"ad_df_listing.listing_agent_state_license",
			"ad_df_listing.listing_agent_aor",
			"ad_df_listing.listing_agent_preferred_phone",
			"ad_df_listing.listing_agent_email",
			"ad_df_listing.listing_office_name",
			"ad_df_listing.listing_office_mls_id",
			"ad_df_listing.listing_office_aor",
			"ad_df_listing.listing_office_phone",
			"ad_df_listing.listing_office_email",
			"ad_df_listing.listing_co_agent_full_name",
			"ad_df_listing.listing_co_agent_mls_id",
			"ad_df_listing.listing_co_agent_state_license",
			"ad_df_listing.listing_co_agent_aor",
			"ad_df_listing.listing_co_agent_preferred_phone",
			"ad_df_listing.listing_co_agent_email",
			"ad_df_listing.listing_co_agent_office_name",
			"ad_df_listing.listing_co_agent_office_mls_id",
			"ad_df_listing.listing_co_agent_office_aor",
			"ad_df_listing.listing_co_agent_office_phone",
			"ad_df_listing.listing_co_agent_office_email",
			"ad_df_listing.buyer_agent_full_name",
			"ad_df_listing.buyer_agent_mls_id",
			"ad_df_listing.buyer_agent_state_license",
			"ad_df_listing.buyer_agent_aor",
			"ad_df_listing.buyer_agent_preferred_phone",
			"ad_df_listing.buyer_agent_email",
			"ad_df_listing.buyer_office_name",
			"ad_df_listing.buyer_office_mls_id",
			"ad_df_listing.buyer_office_aor",
			"ad_df_listing.buyer_office_phone",
			"ad_df_listing.buyer_office_email",
			"ad_df_listing.buyer_co_agent_full_name",
			"ad_df_listing.buyer_co_agent_mls_id",
			"ad_df_listing.buyer_co_agent_state_license",
			"ad_df_listing.buyer_co_agent_aor",
			"ad_df_listing.buyer_co_agent_preferred_phone",
			"ad_df_listing.buyer_co_agent_email",
			"ad_df_listing.buyer_co_agent_office_name",
			"ad_df_listing.buyer_co_agent_office_mls_id",
			"ad_df_listing.buyer_co_agent_office_aor",
			"ad_df_listing.buyer_co_agent_office_phone",
			"ad_df_listing.buyer_co_agent_office_email",
			"ad_df_listing.public_listing_remarks",
			"ad_df_listing.home_warranty_yn",
			"ad_df_listing.tax_year_assessed",
			"ad_df_listing.tax_assessed_value_total",
			"ad_df_listing.tax_amount",
			"ad_df_listing.tax_annual_other",
			"ad_df_listing.owner_name",
			"ad_df_listing.owner_vesting",
			"ad_df_listing.year_built",
			"ad_df_listing.year_built_effective",
			"ad_df_listing.year_built_source",
			"ad_df_listing.new_construction_yn",
			"ad_df_listing.builder_name",
			"ad_df_listing.additional_parcels_yn",
			"ad_df_listing.number_of_lots",
			"ad_df_listing.lot_size_square_feet",
			"ad_df_listing.lot_size_acres",
			"ad_df_listing.lot_size_source",
			"ad_df_listing.lot_dimensions",
			"ad_df_listing.lot_feature_list",
			"ad_df_listing.frontage_length",
			"ad_df_listing.frontage_type",
			"ad_df_listing.frontage_road_type",
			"ad_df_listing.living_area_square_feet",
			"ad_df_listing.living_area_source",
			"ad_df_listing.levels",
			"ad_df_listing.stories",
			"ad_df_listing.building_stories_total",
			"ad_df_listing.building_keywords",
			"ad_df_listing.building_area_total",
			"ad_df_listing.number_of_units_total",
			"ad_df_listing.number_of_buildings",
			"ad_df_listing.property_attached_yn",
			"ad_df_listing.other_structures",
			"ad_df_listing.rooms_total",
			"ad_df_listing.bedrooms_total",
			"ad_df_listing.bathrooms_full",
			"ad_df_listing.bathrooms_half",
			"ad_df_listing.bathrooms_quarter",
			"ad_df_listing.bathrooms_three_quarters",
			"ad_df_listing.basement_features",
			"ad_df_listing.below_grade_square_feet",
			"ad_df_listing.basement_total_sq_ft",
			"ad_df_listing.basement_finished_sq_ft",
			"ad_df_listing.basement_unfinished_sq_ft",
			"ad_df_listing.property_condition",
			"ad_df_listing.repairs_yn",
			"ad_df_listing.repairs_description",
			"ad_df_listing.disclosures",
			"ad_df_listing.construction_materials",
			"ad_df_listing.garage_yn",
			"ad_df_listing.attached_garage_yn",
			"ad_df_listing.garage_spaces",
			"ad_df_listing.carport_yn",
			"ad_df_listing.carport_spaces",
			"ad_df_listing.parking_features",
			"ad_df_listing.parking_other",
			"ad_df_listing.open_parking_spaces",
			"ad_df_listing.parking_total",
			"ad_df_listing.pool_private_yn",
			"ad_df_listing.pool_features",
			"ad_df_listing.occupancy",
			"ad_df_listing.view_yn",
			"ad_df_listing.view_col",
			"ad_df_listing.topography",
			"ad_df_listing.heating_yn",
			"ad_df_listing.heating_features",
			"ad_df_listing.cooling_yn",
			"ad_df_listing.cooling",
			"ad_df_listing.fireplace_yn",
			"ad_df_listing.fireplace",
			"ad_df_listing.fireplace_number",
			"ad_df_listing.foundation_features",
			"ad_df_listing.roof",
			"ad_df_listing.architectural_style_features",
			"ad_df_listing.patio_and_porch_features",
			"ad_df_listing.utilities",
			"ad_df_listing.electric_included",
			"ad_df_listing.electric_description",
			"ad_df_listing.water_included",
			"ad_df_listing.water_source",
			"ad_df_listing.sewer",
			"ad_df_listing.gas_description",
			"ad_df_listing.other_equipment_included",
			"ad_df_listing.laundry_features",
			"ad_df_listing.appliances",
			"ad_df_listing.interior_features",
			"ad_df_listing.exterior_features",
			"ad_df_listing.fencing_features",
			"ad_df_listing.pets_allowed",
			"ad_df_listing.horse_zoning_yn",
			"ad_df_listing.senior_community_yn",
			"ad_df_listing.waterbody_name",
			"ad_df_listing.waterfront_yn",
			"ad_df_listing.waterfront_features",
			"ad_df_listing.zoning_code",
			"ad_df_listing.zoning_description",
			"ad_df_listing.current_use",
			"ad_df_listing.possible_use",
			"ad_df_listing.association_yn",
			"ad_df_listing.association1_name",
			"ad_df_listing.association1_phone",
			"ad_df_listing.association1_fee",
			"ad_df_listing.association1_fee_frequency",
			"ad_df_listing.association2_name",
			"ad_df_listing.association2_phone",
			"ad_df_listing.association2_fee",
			"ad_df_listing.association2_fee_frequency",
			"ad_df_listing.association_fee_includes",
			"ad_df_listing.association_amenities",
			"ad_df_listing.school_elementary",
			"ad_df_listing.school_elementary_district",
			"ad_df_listing.school_middle",
			"ad_df_listing.school_middle_district",
			"ad_df_listing.school_high",
			"ad_df_listing.school_high_district",
			"ad_df_listing.green_verification_yn",
			"ad_df_listing.green_building_verification_type",
			"ad_df_listing.green_energy_efficient",
			"ad_df_listing.green_energy_generation",
			"ad_df_listing.green_indoor_air_quality",
			"ad_df_listing.green_location",
			"ad_df_listing.green_sustainability",
			"ad_df_listing.green_water_conservation",
			"ad_df_listing.land_lease_yn",
			"ad_df_listing.land_lease_amount",
			"ad_df_listing.land_lease_amount_frequency",
			"ad_df_listing.land_lease_expiration_date",
			"ad_df_listing.cap_rate",
			"ad_df_listing.gross_income",
			"ad_df_listing.income_includes",
			"ad_df_listing.gross_scheduled_income",
			"ad_df_listing.net_operating_income",
			"ad_df_listing.total_actual_rent",
			"ad_df_listing.existing_lease_type",
			"ad_df_listing.financial_data_source",
			"ad_df_listing.rent_control_yn",
			"ad_df_listing.unit_type_description",
			"ad_df_listing.unit_type_furnished",
			"ad_df_listing.number_of_units_leased",
			"ad_df_listing.number_of_units_mo_mo",
			"ad_df_listing.number_of_units_vacant",
			"ad_df_listing.vacancy_allowance",
			"ad_df_listing.vacancy_allowance_rate",
			"ad_df_listing.operating_expense",
			"ad_df_listing.cable_tv_expense",
			"ad_df_listing.electric_expense",
			"ad_df_listing.fuel_expense",
			"ad_df_listing.furniture_replacement_expense",
			"ad_df_listing.gardener_expense",
			"ad_df_listing.insurance_expense",
			"ad_df_listing.operating_expense_includes",
			"ad_df_listing.licenses_expense",
			"ad_df_listing.maintenance_expense",
			"ad_df_listing.manager_expense",
			"ad_df_listing.new_taxes_expense",
			"ad_df_listing.other_expense",
			"ad_df_listing.pest_control_expense",
			"ad_df_listing.pool_expense",
			"ad_df_listing.professional_management_expense",
			"ad_df_listing.supplies_expense",
			"ad_df_listing.trash_expense",
			"ad_df_listing.water_sewer_expense",
			"ad_df_listing.workmans_compensation_expense",
			"ad_df_listing.owner_pays",
			"ad_df_listing.tenant_pays",
			"ad_df_listing.listing_marketing_url",
			"ad_df_listing.photos_count",
			"ad_df_listing.photo_key",
			"ad_df_listing.photo_url_prefix",
		).
		From("properties").
		Join("ad_df_listing on properties.ad_attom_id = ad_df_listing.attom_id").
		LeftJoin("current_listings on ad_df_listing.attom_id = current_listings.attom_id").
		Where("properties.id = ?", in.Aupid).
		OrderBy(
			"ad_df_listing.status_change_date desc",
			"ad_df_listing.mls_record_id desc",
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "d4c4354c-dc64-490e-8202-eb49ea06f204",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "e52b6749-4981-4a4f-b60d-0dcab5571ab1")
	}
	defer rows.Close()

	out := &SelectListingRecordOutput{}

	for rows.Next() {
		record := &entities.Listing{}

		if err := rows.Scan(
			&record.MlsRecordId,
			&record.MlsListingId,
			&record.StatusChangeDate,
			&record.FullStreetAddress,
			&record.HouseNumber,
			&record.StreetPreDirection,
			&record.StreetName,
			&record.StreetSuffix,
			&record.StreetPostDirection,
			&record.UnitType,
			&record.UnitNumber,
			&record.City,
			&record.State,
			&record.Zip5,
			&record.Zip4,
			&record.County,
			&record.Township,
			&record.MlsListingAddress,
			&record.MlsListingCity,
			&record.MlsListingState,
			&record.MlsListingZip,
			&record.MlsListingCountyFips,
			&record.MlsNumber,
			&record.MlsSource,
			&record.Ouid,
			&record.ListingStatus,
			&record.MlsSoldDate,
			&record.MlsSoldPrice,
			&record.AssessorLastSaleDate,
			&record.AssessorLastSaleAmount,
			&record.MarketValue,
			&record.MarketValueDate,
			&record.AvgMarketPricePerSqFt,
			&record.ListingDate,
			&record.LatestListingPrice,
			&record.PreviousListingPrice,
			&record.LatestPriceChangeDate,
			&record.PendingDate,
			&record.SpecialListingConditions,
			&record.OriginalListingDate,
			&record.OriginalListingPrice,
			&record.LeaseOption,
			&record.LeaseTerm,
			&record.LeaseIncludes,
			&record.Concessions,
			&record.ConcessionsAmount,
			&record.ConcessionsComments,
			&record.ContingencyDate,
			&record.ContingencyDescription,
			&record.MlsPropertyType,
			&record.MlsPropertySubType,
			&record.AsrPropertyType,
			&record.AsrPropertySubType,
			&record.OwnershipDescription,
			&record.Latitude,
			&record.Longitude,
			&record.Apn,
			&record.LegalDescription,
			&record.LegalSubdivision,
			&record.DaysOnMarket,
			&record.CumulativeDaysOnMarket,
			&record.ListingAgentFullName,
			&record.ListingAgentMlsId,
			&record.ListingAgentStateLicense,
			&record.ListingAgentAor,
			&record.ListingAgentPreferredPhone,
			&record.ListingAgentEmail,
			&record.ListingOfficeName,
			&record.ListingOfficeMlsId,
			&record.ListingOfficeAor,
			&record.ListingOfficePhone,
			&record.ListingOfficeEmail,
			&record.ListingCoAgentFullName,
			&record.ListingCoAgentMlsId,
			&record.ListingCoAgentStateLicense,
			&record.ListingCoAgentAor,
			&record.ListingCoAgentPreferredPhone,
			&record.ListingCoAgentEmail,
			&record.ListingCoAgentOfficeName,
			&record.ListingCoAgentOfficeMlsId,
			&record.ListingCoAgentOfficeAor,
			&record.ListingCoAgentOfficePhone,
			&record.ListingCoAgentOfficeEmail,
			&record.BuyerAgentFullName,
			&record.BuyerAgentMlsId,
			&record.BuyerAgentStateLicense,
			&record.BuyerAgentAor,
			&record.BuyerAgentPreferredPhone,
			&record.BuyerAgentEmail,
			&record.BuyerOfficeName,
			&record.BuyerOfficeMlsId,
			&record.BuyerOfficeAor,
			&record.BuyerOfficePhone,
			&record.BuyerOfficeEmail,
			&record.BuyerCoAgentFullName,
			&record.BuyerCoAgentMlsId,
			&record.BuyerCoAgentStateLicense,
			&record.BuyerCoAgentAor,
			&record.BuyerCoAgentPreferredPhone,
			&record.BuyerCoAgentEmail,
			&record.BuyerCoAgentOfficeName,
			&record.BuyerCoAgentOfficeMlsId,
			&record.BuyerCoAgentOfficeAor,
			&record.BuyerCoAgentOfficePhone,
			&record.BuyerCoAgentOfficeEmail,
			&record.PublicListingRemarks,
			&record.HasHomeWarranty,
			&record.TaxYearAssessed,
			&record.TaxAssessedValueTotal,
			&record.TaxAmount,
			&record.TaxAnnualOther,
			&record.OwnerName,
			&record.OwnerVesting,
			&record.YearBuilt,
			&record.YearBuiltEffective,
			&record.YearBuiltSource,
			&record.IsNewConstruction,
			&record.BuilderName,
			&record.HasAdditionalParcels,
			&record.NumberOfLots,
			&record.LotSizeSquareFeet,
			&record.LotSizeAcres,
			&record.LotSizeSource,
			&record.LotDimensions,
			&record.LotFeatureList,
			&record.FrontageLength,
			&record.FrontageType,
			&record.FrontageRoadType,
			&record.LivingAreaSquareFeet,
			&record.LivingAreaSource,
			&record.Levels,
			&record.Stories,
			&record.BuildingStoriesTotal,
			&record.BuildingKeywords,
			&record.BuildingAreaTotal,
			&record.NumberOfUnitsTotal,
			&record.NumberOfBuildings,
			&record.HasPropertyAttached,
			&record.OtherStructures,
			&record.RoomsTotal,
			&record.BedroomsTotal,
			&record.BathroomsFull,
			&record.BathroomsHalf,
			&record.BathroomsQuarter,
			&record.BathroomsThreeQuarters,
			&record.BasementFeatures,
			&record.BelowGradeSquareFeet,
			&record.BasementTotalSqFt,
			&record.BasementFinishedSqFt,
			&record.BasementUnfinishedSqFt,
			&record.PropertyCondition,
			&record.HasNeededRepairs,
			&record.RepairsDescription,
			&record.Disclosures,
			&record.ConstructionMaterials,
			&record.HasGarage,
			&record.HasAttachedGarage,
			&record.GarageSpaces,
			&record.HasCarport,
			&record.CarportSpaces,
			&record.ParkingFeatures,
			&record.ParkingOther,
			&record.OpenParkingSpaces,
			&record.ParkingTotal,
			&record.HasPrivatePool,
			&record.PoolFeatures,
			&record.Occupancy,
			&record.HasView,
			&record.View,
			&record.Topography,
			&record.HasHeating,
			&record.HeatingFeatures,
			&record.HasCooling,
			&record.Cooling,
			&record.HasFireplace,
			&record.Fireplace,
			&record.FireplaceNumber,
			&record.FoundationFeatures,
			&record.Roof,
			&record.ArchitecturalStyleFeatures,
			&record.PatioAndPorchFeatures,
			&record.Utilities,
			&record.ElectricIncluded,
			&record.ElectricDescription,
			&record.WaterIncluded,
			&record.WaterSource,
			&record.Sewer,
			&record.GasDescription,
			&record.OtherEquipmentIncluded,
			&record.LaundryFeatures,
			&record.Appliances,
			&record.InteriorFeatures,
			&record.ExteriorFeatures,
			&record.FencingFeatures,
			&record.PetsAllowed,
			&record.HorseZoning,
			&record.SeniorCommunity,
			&record.WaterbodyName,
			&record.IsWaterfront,
			&record.WaterfrontFeatures,
			&record.ZoningCode,
			&record.ZoningDescription,
			&record.CurrentUse,
			&record.PossibleUse,
			&record.HasAssociation,
			&record.Association1Name,
			&record.Association1Phone,
			&record.Association1Fee,
			&record.Association1FeeFrequency,
			&record.Association2Name,
			&record.Association2Phone,
			&record.Association2Fee,
			&record.Association2FeeFrequency,
			&record.AssociationFeeIncludes,
			&record.AssociationAmenities,
			&record.SchoolElementary,
			&record.SchoolElementaryDistrict,
			&record.SchoolMiddle,
			&record.SchoolMiddleDistrict,
			&record.SchoolHigh,
			&record.SchoolHighDistrict,
			&record.HasGreenVerification,
			&record.GreenBuildingVerificationType,
			&record.GreenEnergyEfficient,
			&record.GreenEnergyGeneration,
			&record.GreenIndoorAirQuality,
			&record.GreenLocation,
			&record.GreenSustainability,
			&record.GreenWaterConservation,
			&record.HasLandLease,
			&record.LandLeaseAmount,
			&record.LandLeaseAmountFrequency,
			&record.LandLeaseExpirationDate,
			&record.CapRate,
			&record.GrossIncome,
			&record.IncomeIncludes,
			&record.GrossScheduledIncome,
			&record.NetOperatingIncome,
			&record.TotalActualRent,
			&record.ExistingLeaseType,
			&record.FinancialDataSource,
			&record.HasRentControl,
			&record.UnitTypeDescription,
			&record.UnitTypeFurnished,
			&record.NumberOfUnitsLeased,
			&record.NumberOfUnitsMoMo,
			&record.NumberOfUnitsVacant,
			&record.VacancyAllowance,
			&record.VacancyAllowanceRate,
			&record.OperatingExpense,
			&record.CableTvExpense,
			&record.ElectricExpense,
			&record.FuelExpense,
			&record.FurnitureReplacementExpense,
			&record.GardenerExpense,
			&record.InsuranceExpense,
			&record.OperatingExpenseIncludes,
			&record.LicensesExpense,
			&record.MaintenanceExpense,
			&record.ManagerExpense,
			&record.NewTaxesExpense,
			&record.OtherExpense,
			&record.PestControlExpense,
			&record.PoolExpense,
			&record.ProfessionalManagementExpense,
			&record.SuppliesExpense,
			&record.TrashExpense,
			&record.WaterSewerExpense,
			&record.WorkmansCompensationExpense,
			&record.OwnerPays,
			&record.TenantPays,
			&record.ListingMarketingUrl,
			&record.PhotosCount,
			&record.PhotoKey,
			&record.PhotoUrlPrefix,
		); err != nil {
			return nil, &errors.Object{
				Id:     "717d8c81-0e12-47b8-9777-74909d6250a8",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}

func (rep *repository) GetLatLonFromAupid(r *arc.Request, aupid string) (float64, float64, error) {
	// build the SQL query using squirrel
	// using the auPid to filter the results
	// and limit the results to 1 due to the full history

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("ad_df_assessor.property_latitude", "ad_df_assessor.property_longitude").
		From("properties").
		Join("ad_df_assessor ON ad_df_assessor.attomid = properties.ad_attom_id").
		Where("properties.id = ?", aupid).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, 0, &errors.Object{
			Id:     "e61355e3-0667-45ec-ac28-6ae8e95ac7dd",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if row == nil {
		return 0, 0, errors.Forward(err, "299b84cb-54e3-48ab-856c-f8ccb14b824b")
	}

	var lat float64
	var lon float64
	if err := row.Scan(&lat, &lon); err != nil {

		// Check if the error is a "no rows" error
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, &errors.Object{
				Id:     "6eda6efc-d666-467f-8e60-695ad1b4d8c8",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Aupid not found.",
				Cause:  err.Error(),
			}
		}

		return 0, 0, &errors.Object{
			Id:     "e8d63290-5389-4948-9d07-2362ad490c6e",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to scan row.",
			Cause:  err.Error(),
		}
	}

	return lat, lon, nil
}

func (rep *repository) GetLatLonFromAddressId(r *arc.Request, addressId uuid.UUID) (float64, float64, error) {
	// build the SQL query using squirrel
	// using the addressId to filter the results
	// and limit the results to 1 due to the full history

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("ad_df_assessor.property_latitude", "ad_df_assessor.property_longitude").
		From("properties").
		Join("ad_df_assessor ON ad_df_assessor.attomid = properties.ad_attom_id").
		Where("properties.address_id = ?", addressId).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, 0, &errors.Object{
			Id:     "bc635a0f-6a06-4487-a32b-37626b0156a0",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if row == nil {
		return 0, 0, &errors.Object{
			Id:     "e462543a-e8e6-4123-a304-745e9b0c5499",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to query row.",
			Cause:  err.Error(),
		}
	}

	var lat float64
	var lon float64
	if err := row.Scan(&lat, &lon); err != nil {

		// Check if the error is a "no rows" error
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, &errors.Object{
				Id:     "5b730ccd-6834-493f-8910-03cb7cff9e5e",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Address not found.",
				Cause:  err.Error(),
			}
		}

		return 0, 0, &errors.Object{
			Id:     "ef040e1f-95ff-42fb-b91f-0648a8e75edc",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to scan row.",
			Cause:  err.Error(),
		}
	}

	return lat, lon, nil
}
