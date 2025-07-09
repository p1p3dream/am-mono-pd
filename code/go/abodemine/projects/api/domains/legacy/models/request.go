package models

type PropertySearchRequests struct {
	Apn                 string         `json:"apn"`
	Aupid               string         `json:"aupid"`
	City                string         `json:"city"`
	Custom              map[string]any `json:"custom,omitempty"`
	FullAddress         string         `json:"fulladdress"`
	HouseNumber         string         `json:"housenumber"`
	IncludeFields       []string       `json:"_includefields"`
	IncludeLayouts      []string       `json:"_includelayouts"`
	InputID             string         `json:"inputid"`
	Limit               int            `json:"limit"`
	OldAupid            string         `json:"old_aupid"`
	OmitFields          []string       `json:"_omitfields"`
	OmitLayouts         []string       `json:"_omitlayouts"`
	Page                int            `json:"page"`
	RecorderDocTypes    []string       `json:"recorderdoctypes"`
	State               string         `json:"state"`
	StreetDirection     string         `json:"streetdirection"`
	StreetName          string         `json:"streetname"`
	StreetPostDirection string         `json:"streetpostdirection"`
	StreetSuffix        string         `json:"streetsuffix"`
	Unit                string         `json:"unit"`
	Zip5                string         `json:"zip5"`
}
