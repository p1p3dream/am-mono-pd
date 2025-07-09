package listings

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"abodemine/lib/ptr"
)

func TestCheckHasGeoFilter(t *testing.T) {
	tests := []struct {
		name      string
		geoFilter *GeoFilter
		expected  bool
	}{
		{
			name:      "Nil GeoFilter",
			geoFilter: nil,
			expected:  false,
		},
		{
			name:      "Empty GeoFilter",
			geoFilter: &GeoFilter{},
			expected:  false,
		},
		{
			name: "Has GeoDistance",
			geoFilter: &GeoFilter{
				GeoDistance: &GeoDistance{
					Radius:   10,
					Location: GeoPoint{Lat: 37.7749, Lon: -122.4194},
				},
			},
			expected: true,
		},
		{
			name: "Has GeoBoundingBox",
			geoFilter: &GeoFilter{
				GeoBoundingBox: &GeoBoundingBox{
					Location: &BoundingBoxLocation{
						TopLeft:     GeoPoint{Lat: 37.7914, Lon: -122.4587},
						BottomRight: GeoPoint{Lat: 37.7749, Lon: -122.4194},
					},
				},
			},
			expected: true,
		},
		{
			name: "Has GeoPolygon",
			geoFilter: &GeoFilter{
				GeoPolygon: &GeoPolygon{
					Points: []GeoPoint{
						{Lat: 37.7749, Lon: -122.4194},
						{Lat: 37.7847, Lon: -122.4587},
						{Lat: 37.7914, Lon: -122.4089},
					},
				},
			},
			expected: true,
		},
		{
			name: "Has GeoBoundingBox but no Location",
			geoFilter: &GeoFilter{
				GeoBoundingBox: &GeoBoundingBox{},
			},
			expected: false,
		},
		{
			name: "Has GeoPolygon but no Points",
			geoFilter: &GeoFilter{
				GeoPolygon: &GeoPolygon{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasGeoFilter := tt.geoFilter != nil &&
				((tt.geoFilter.GeoDistance != nil) ||
					(tt.geoFilter.GeoBoundingBox != nil && tt.geoFilter.GeoBoundingBox.Location != nil) ||
					(tt.geoFilter.GeoPolygon != nil && len(tt.geoFilter.GeoPolygon.Points) > 0))

			assert.Equal(t, tt.expected, hasGeoFilter)
		})
	}
}

func Test_buildWhereStr(t *testing.T) {
	whereWithGeo := ", reference_geom WHERE location_3857 && ST_Expand(reference_geom.geom, 16093.400000)"

	type args struct {
		filter *SearchListingsInput
		where  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No filter",
			args: args{
				filter: &SearchListingsInput{},
				where:  "",
			},
			want: "",
		},
		{
			name: "No geo filter",
			args: args{
				filter: &SearchListingsInput{
					GeoFilter: nil,
					MinBeds:   ptr.Int(1),
				},
				where: "",
			},
			want: "WHERE current_listings.bedrooms_total >= 1",
		},
		{
			name: "No geo filter with min and max beds",
			args: args{
				filter: &SearchListingsInput{
					GeoFilter: nil,
					MinBeds:   ptr.Int(1),
					MaxBeds:   ptr.Int(3),
				},
				where: "",
			},
			want: "WHERE current_listings.bedrooms_total >= 1 AND current_listings.bedrooms_total <= 3",
		},
		{
			name: "No geo filter with zip5codes",
			args: args{
				filter: &SearchListingsInput{
					GeoFilter: nil,
					Zip5Codes: []string{"12345", "67890"},
				},
				where: "",
			},
			want: "WHERE current_listings.zip5 IN ('12345', '67890')",
		},
		{
			name: "Has where and has zip5codes",
			args: args{
				filter: &SearchListingsInput{
					Zip5Codes: []string{"12345", "67890"},
				},
				where: whereWithGeo,
			},
			want: fmt.Sprintf("%s AND current_listings.zip5 IN ('12345', '67890') ORDER BY current_listings.location_3857 <-> reference_geom.geom", whereWithGeo),
		},
		{
			name: "Has where and has 2 filters",
			args: args{
				filter: &SearchListingsInput{
					MinBeds:   ptr.Int(1),
					MaxBeds:   ptr.Int(3),
					Zip5Codes: []string{"12345", "67890"},
				},
				where: whereWithGeo,
			},
			want: fmt.Sprintf("%s AND current_listings.bedrooms_total >= 1 AND current_listings.bedrooms_total <= 3 AND current_listings.zip5 IN ('12345', '67890') ORDER BY current_listings.location_3857 <-> reference_geom.geom", whereWithGeo),
		},
		{
			name: "MLS property type with single quote",
			args: args{
				filter: &SearchListingsInput{
					MlsPropertyType: ptr.String("Residential'"),
				},
				where: "",
			},
			want: "WHERE current_listings.mls_property_type = 'Residential'",
		},
		{
			name: "MinStatusChangeDate only",
			args: args{
				filter: &SearchListingsInput{
					MinStatusChangeDate: ptr.String("2024-01-01"),
				},
				where: "",
			},
			want: "WHERE current_listings.status_change_date >= '2024-01-01'",
		},
		{
			name: "All filters with existing WHERE",
			args: args{
				filter: &SearchListingsInput{
					MinBeds:             ptr.Int(1),
					MaxBeds:             ptr.Int(3),
					MinBaths:            ptr.Float64(1.5),
					MaxBaths:            ptr.Float64(2.5),
					MinPrice:            ptr.Int(100000),
					MaxPrice:            ptr.Int(500000),
					MinSqFt:             ptr.Int(800),
					MaxSqFt:             ptr.Int(1200),
					MinYearBuilt:        ptr.Int(1990),
					MaxYearBuilt:        ptr.Int(2020),
					MlsPropertyType:     ptr.String("Residential"),
					MlsPropertySubType:  ptr.String("Condominium"),
					Statuses:            []string{"Active", "Pending"},
					Zip5Codes:           []string{"12345"},
					MinStatusChangeDate: ptr.String("2024-01-01"),
				},
				where: whereWithGeo,
			},
			want: fmt.Sprintf(`%s AND current_listings.bedrooms_total >= 1 AND current_listings.bedrooms_total <= 3 AND current_listings.bathrooms_full >= 1.500000 AND current_listings.bathrooms_full <= 2.500000 AND current_listings.latest_listing_price >= 100000 AND current_listings.latest_listing_price <= 500000 AND current_listings.living_area_square_feet >= 800 AND current_listings.living_area_square_feet <= 1200 AND current_listings.year_built >= 1990 AND current_listings.year_built <= 2020 AND current_listings.mls_property_type = 'Residential' AND current_listings.mls_property_sub_type = 'Condominium' AND current_listings.listing_status IN ('Active', 'Pending') AND current_listings.zip5 IN ('12345') AND current_listings.status_change_date >= '2024-01-01' ORDER BY current_listings.location_3857 <-> reference_geom.geom`, whereWithGeo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildWhereStr(tt.args.filter, tt.args.where)
			got = strings.TrimSpace(got)
			if got != tt.want {
				t.Errorf("%s \n got  = %v \n want = %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_searchByGeometry(t *testing.T) {
	type args struct {
		geoFilter *GeoFilter
		radius    float64
	}
	tests := []struct {
		name      string
		args      args
		wantCte   string
		wantWhere string
	}{
		{
			name: "No geo filter",
			args: args{
				geoFilter: nil,
				radius:    0,
			},
			wantCte:   "",
			wantWhere: "",
		},
		{
			name: "GeoDistance",
			args: args{
				geoFilter: &GeoFilter{
					GeoDistance: &GeoDistance{
						Radius:   10,
						Location: GeoPoint{Lat: 37.7749, Lon: -122.4194},
					},
				},
				radius: 10,
			},
			wantCte:   "WITH reference_geom AS (SELECT ST_Transform(ST_SetSRID(ST_MakePoint(-122.419400, 37.774900), 4326), 3857) AS geom)",
			wantWhere: ", reference_geom WHERE location_3857 && ST_Expand(reference_geom.geom, 16093.400000)",
		},
		{
			name: "GeoPolygon",
			args: args{
				geoFilter: &GeoFilter{
					GeoPolygon: &GeoPolygon{
						Points: []GeoPoint{
							{Lat: 37.7749, Lon: -122.4194},
							{Lat: 37.7847, Lon: -122.4587},
							{Lat: 37.7914, Lon: -122.4089},
						},
					},
				},
				radius: 1,
			},
			wantCte:   "WITH reference_geom AS (SELECT ST_Transform(ST_GeomFromText('POLYGON((-122.419400 37.774900, -122.458700 37.784700, -122.408900 37.791400))', 4326), 3857) AS geom)",
			wantWhere: ", reference_geom WHERE ST_Within(location_3857, reference_geom.geom)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCet, gotWhere := searchByGeometry(tt.args.geoFilter, tt.args.radius)
			if gotCet != tt.wantCte {
				t.Errorf("%s \ngotCte =  %v\nwantCte = %v", tt.name, gotCet, tt.wantCte)
			}
			if gotWhere != tt.wantWhere {
				t.Errorf("%s \ngotWhere =  %v\nwantWhere = %v", tt.name, gotWhere, tt.wantWhere)
			}
		})
	}
}
