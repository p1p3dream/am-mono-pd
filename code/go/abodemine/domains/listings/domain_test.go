package listings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"abodemine/lib/errors"
	"abodemine/lib/ptr"
)

func TestValidateListingsSearchFilters(t *testing.T) {
	tests := []struct {
		name        string
		filter      *SearchListingsInput
		expectError bool
		errorLabel  string
	}{
		{
			name:        "Nil Filter",
			filter:      nil,
			expectError: true,
			errorLabel:  "invalid_search_filters",
		},
		{
			name: "Invalid Bedroom Range",
			filter: &SearchListingsInput{
				MinBeds: ptr.Int(5),
				MaxBeds: ptr.Int(3),
			},
			expectError: true,
			errorLabel:  "invalid_bedroom_range",
		},
		{
			name: "Invalid Bathroom Range",
			filter: &SearchListingsInput{
				MinBaths: ptr.Float64(3.5),
				MaxBaths: ptr.Float64(2.0),
			},
			expectError: true,
			errorLabel:  "invalid_bathroom_range",
		},
		{
			name: "Invalid Price Range",
			filter: &SearchListingsInput{
				MinPrice: ptr.Int(500000),
				MaxPrice: ptr.Int(300000),
			},
			expectError: true,
			errorLabel:  "invalid_price_range",
		},
		{
			name: "Invalid Square Footage Range",
			filter: &SearchListingsInput{
				MinSqFt: ptr.Int(3000),
				MaxSqFt: ptr.Int(2000),
			},
			expectError: true,
			errorLabel:  "invalid_square_footage_range",
		},
		{
			name: "Invalid Year Built Range",
			filter: &SearchListingsInput{
				MinYearBuilt: ptr.Int(2020),
				MaxYearBuilt: ptr.Int(2010),
			},
			expectError: true,
			errorLabel:  "invalid_year_built_range",
		},
		{
			name: "Invalid Days on Market Range",
			filter: &SearchListingsInput{
				MinDOM: ptr.Int(30),
				MaxDOM: ptr.Int(15),
			},
			expectError: true,
			errorLabel:  "invalid_days_on_market_range",
		},
		{
			name: "Invalid Sale Date Format",
			filter: &SearchListingsInput{
				MinSaleDate: ptr.String("2023-13-01"),
				MaxSaleDate: ptr.String("2023-12-31"),
			},
			expectError: true,
			errorLabel:  "invalid_sale_date_format",
		},
		{
			name: "Invalid Sale Date Range",
			filter: &SearchListingsInput{
				MinSaleDate: ptr.String("2023-12-31"),
				MaxSaleDate: ptr.String("2023-01-01"),
			},
			expectError: true,
			errorLabel:  "invalid_sale_date_range",
		},
		{
			name: "Invalid Sale Price Range",
			filter: &SearchListingsInput{
				MinSalePrice: ptr.Int(600000),
				MaxSalePrice: ptr.Int(400000),
			},
			expectError: true,
			errorLabel:  "invalid_sale_price_range",
		},
		{
			name: "Invalid Lot Size Range",
			filter: &SearchListingsInput{
				MinLotSqFt: ptr.Float64(20000),
				MaxLotSqFt: ptr.Float64(10000),
			},
			expectError: true,
			errorLabel:  "invalid_lot_size_range",
		},
		{
			name: "Invalid Lot Acres Range",
			filter: &SearchListingsInput{
				MinLotAcres: ptr.Float64(1.0),
				MaxLotAcres: ptr.Float64(0.5),
			},
			expectError: true,
			errorLabel:  "invalid_lot_acres_range",
		},
		{
			name: "Invalid Garage Spaces Range",
			filter: &SearchListingsInput{
				MinGarageSpaces: ptr.Int(3),
				MaxGarageSpaces: ptr.Int(1),
			},
			expectError: true,
			errorLabel:  "invalid_garage_spaces_range",
		},
		{
			name: "Invalid Stories Range",
			filter: &SearchListingsInput{
				MinStories: ptr.Int(3),
				MaxStories: ptr.Int(1),
			},
			expectError: true,
			errorLabel:  "invalid_stories_range",
		},
		{
			name: "Valid Filter",
			filter: &SearchListingsInput{
				MinBeds:      ptr.Int(3),
				MaxBeds:      ptr.Int(5),
				MinBaths:     ptr.Float64(2.0),
				MaxBaths:     ptr.Float64(3.5),
				MinPrice:     ptr.Int(300000),
				MaxPrice:     ptr.Int(500000),
				MinSqFt:      ptr.Int(1500),
				MaxSqFt:      ptr.Int(3000),
				MinYearBuilt: ptr.Int(2000),
				MaxYearBuilt: ptr.Int(2020),
				MinDOM:       ptr.Int(0),
				MaxDOM:       ptr.Int(30),
				MinSaleDate:  ptr.String("2023-01-01"),
				MaxSaleDate:  ptr.String("2023-12-31"),
			},
			expectError: false,
		},
		{
			name: "Single Min Beds Filter",
			filter: &SearchListingsInput{
				MinBeds: ptr.Int(3),
			},
			expectError: false,
		},
		{
			name: "Single Max Beds Filter",
			filter: &SearchListingsInput{
				MaxBeds: ptr.Int(5),
			},
			expectError: false,
		},
		{
			name: "Zero Min Beds Filter",
			filter: &SearchListingsInput{
				MinBeds: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Baths Filter",
			filter: &SearchListingsInput{
				MinBaths: ptr.Float64(2.0),
			},
			expectError: false,
		},
		{
			name: "Single Max Baths Filter",
			filter: &SearchListingsInput{
				MaxBaths: ptr.Float64(3.5),
			},
			expectError: false,
		},
		{
			name: "Zero Min Baths Filter",
			filter: &SearchListingsInput{
				MinBaths: ptr.Float64(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Price Filter",
			filter: &SearchListingsInput{
				MinPrice: ptr.Int(300000),
			},
			expectError: false,
		},
		{
			name: "Single Max Price Filter",
			filter: &SearchListingsInput{
				MaxPrice: ptr.Int(500000),
			},
			expectError: false,
		},
		{
			name: "Zero Min Price Filter",
			filter: &SearchListingsInput{
				MinPrice: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min SqFt Filter",
			filter: &SearchListingsInput{
				MinSqFt: ptr.Int(1500),
			},
			expectError: false,
		},
		{
			name: "Single Max SqFt Filter",
			filter: &SearchListingsInput{
				MaxSqFt: ptr.Int(3000),
			},
			expectError: false,
		},
		{
			name: "Zero Min SqFt Filter",
			filter: &SearchListingsInput{
				MinSqFt: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Year Built Filter",
			filter: &SearchListingsInput{
				MinYearBuilt: ptr.Int(2000),
			},
			expectError: false,
		},
		{
			name: "Single Max Year Built Filter",
			filter: &SearchListingsInput{
				MaxYearBuilt: ptr.Int(2020),
			},
			expectError: false,
		},
		{
			name: "Zero Min Year Built Filter",
			filter: &SearchListingsInput{
				MinYearBuilt: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min DOM Filter",
			filter: &SearchListingsInput{
				MinDOM: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Max DOM Filter",
			filter: &SearchListingsInput{
				MaxDOM: ptr.Int(30),
			},
			expectError: false,
		},
		{
			name: "Single Min Sale Date Filter",
			filter: &SearchListingsInput{
				MinSaleDate: ptr.String("2023-01-01"),
			},
			expectError: false,
		},
		{
			name: "Single Max Sale Date Filter",
			filter: &SearchListingsInput{
				MaxSaleDate: ptr.String("2023-12-31"),
			},
			expectError: false,
		},
		{
			name: "Single Min Sale Price Filter",
			filter: &SearchListingsInput{
				MinSalePrice: ptr.Int(300000),
			},
			expectError: false,
		},
		{
			name: "Single Max Sale Price Filter",
			filter: &SearchListingsInput{
				MaxSalePrice: ptr.Int(500000),
			},
			expectError: false,
		},
		{
			name: "Zero Min Sale Price Filter",
			filter: &SearchListingsInput{
				MinSalePrice: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Lot SqFt Filter",
			filter: &SearchListingsInput{
				MinLotSqFt: ptr.Float64(5000),
			},
			expectError: false,
		},
		{
			name: "Single Max Lot SqFt Filter",
			filter: &SearchListingsInput{
				MaxLotSqFt: ptr.Float64(10000),
			},
			expectError: false,
		},
		{
			name: "Zero Min Lot SqFt Filter",
			filter: &SearchListingsInput{
				MinLotSqFt: ptr.Float64(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Lot Acres Filter",
			filter: &SearchListingsInput{
				MinLotAcres: ptr.Float64(0.25),
			},
			expectError: false,
		},
		{
			name: "Single Max Lot Acres Filter",
			filter: &SearchListingsInput{
				MaxLotAcres: ptr.Float64(1.0),
			},
			expectError: false,
		},
		{
			name: "Zero Min Lot Acres Filter",
			filter: &SearchListingsInput{
				MinLotAcres: ptr.Float64(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Garage Spaces Filter",
			filter: &SearchListingsInput{
				MinGarageSpaces: ptr.Int(1),
			},
			expectError: false,
		},
		{
			name: "Single Max Garage Spaces Filter",
			filter: &SearchListingsInput{
				MaxGarageSpaces: ptr.Int(3),
			},
			expectError: false,
		},
		{
			name: "Zero Min Garage Spaces Filter",
			filter: &SearchListingsInput{
				MinGarageSpaces: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Single Min Stories Filter",
			filter: &SearchListingsInput{
				MinStories: ptr.Int(1),
			},
			expectError: false,
		},
		{
			name: "Single Max Stories Filter",
			filter: &SearchListingsInput{
				MaxStories: ptr.Int(3),
			},
			expectError: false,
		},
		{
			name: "Zero Min Stories Filter",
			filter: &SearchListingsInput{
				MinStories: ptr.Int(0),
			},
			expectError: false,
		},
		{
			name: "Min Status Change Date Filter",
			filter: &SearchListingsInput{
				MinStatusChangeDate: ptr.String("2023-12-31"),
			},
			expectError: false,
		},
		{
			name: "Invalid Min Status Change Date",
			filter: &SearchListingsInput{
				MinStatusChangeDate: ptr.String("202-12-31"),
			},
			expectError: true,
			errorLabel:  "invalid_status_change_date_format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateListingsSearchFilters(tt.filter)
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errorLabel, err.(*errors.Object).Label)
					assert.Equal(t, errors.Code_INVALID_ARGUMENT, err.(*errors.Object).Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
