package models

type ApiSearchAddress struct {
	FullStreetAddress   string `json:"fullStreetAddress"`
	State               string `json:"state"`
	City                string `json:"city"`
	Zip5                string `json:"zip5"`
	StreetPreDirection  string `json:"streetPreDirection"`
	HouseNumber         string `json:"houseNumber"`
	StreetName          string `json:"streetName"`
	StreetPostDirection string `json:"streetPostDirection"`
	StreetSuffix        string `json:"streetSuffix"`
	UnitType            string `json:"unitType"`
	UnitNumber          string `json:"unitNumber"`
}
