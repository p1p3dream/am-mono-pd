// Package flags contains global flags that are sent upstream
// with requests. The purpose of the package is facilitate
// flag validation and selection.
package flags

import (
	"abodemine/lib/errors"
)

func Select(name string) (uint, bool) {
	v, ok := flagByName[name]
	return v, ok
}

func ValidateMany(flags []string) error {
	for _, flag := range flags {
		if len(flag) == 0 {
			return &errors.Object{
				Id:     "3986b394-da2b-4ad1-9278-2a67325b4d27",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Empty flag.",
			}
		}

		if _, ok := flagByName[flag]; !ok {
			return &errors.Object{
				Id:     "1186a8fe-ba73-4ac1-8bbc-e704770b659f",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Unknown flag.",
				Meta: map[string]any{
					"flag": flag,
				},
			}
		}
	}

	return nil
}

const (
	SearchUsingLegacyAddressTable = 1 + iota
	SearchUsingAddressesTable
	SearchUsingOpenSearch

	ApiAddressLayoutEnabled
	ApiAssessorLayoutEnabled
	ApiCompsLayoutEnabled
	ApiListingLayoutEnabled
	ApiRecorderLayoutEnabled
	ApiRentEstimateLayoutEnabled
	ApiSaleEstimateLayoutEnabled
)

var flagByName = map[string]uint{
	"SEARCH_USING_LEGACY_ADDRESS_TABLE": SearchUsingLegacyAddressTable,
	"SEARCH_USING_ADDRESSES_TABLE":      SearchUsingAddressesTable,
	"SEARCH_USING_OPENSEARCH":           SearchUsingOpenSearch,

	"API_ADDRESS_LAYOUT_ENABLED":       ApiAddressLayoutEnabled,
	"API_ASSESSOR_LAYOUT_ENABLED":      ApiAssessorLayoutEnabled,
	"API_COMPS_LAYOUT_ENABLED":         ApiCompsLayoutEnabled,
	"API_LISTING_LAYOUT_ENABLED":       ApiListingLayoutEnabled,
	"API_RECORDER_LAYOUT_ENABLED":      ApiRecorderLayoutEnabled,
	"API_RENT_ESTIMATE_LAYOUT_ENABLED": ApiRentEstimateLayoutEnabled,
	"API_SALE_ESTIMATE_LAYOUT_ENABLED": ApiSaleEstimateLayoutEnabled,
}
