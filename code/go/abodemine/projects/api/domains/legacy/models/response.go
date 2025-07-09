package models

type PropertySearchResponse struct {
	HouseNumber         string `json:"housenumber"`
	StreetDirection     string `json:"streetdirection"`
	StreetName          string `json:"streetname"`
	StreetSuffix        string `json:"streetsuffix"`
	StreetPostDirection string `json:"streetpostdirection"`
	Zip5                string `json:"zip5"`
	City                string `json:"city"`
	State               string `json:"state"`
	Aupid               string `json:"aupid"`

	Assessor     map[string]any     `json:"assessor,omitempty"`
	Hits         int64              `json:"hits,omitempty"`
	Recorder     [][]map[string]any `json:"recorder,omitempty"`
	RentEstimate map[string]any     `json:"rentestimate"`
	SaleEstimate map[string]any     `json:"saleestimate"`
}
