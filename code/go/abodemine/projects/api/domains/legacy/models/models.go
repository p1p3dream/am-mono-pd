package models

type User struct {
	ApiKey     string
	Name       string
	CustomerId int
}

type ElasticsearchError struct {
	Error struct {
		RootCause []struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"root_cause"`
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}

type Address struct {
	IncludeLayouts []string
	OmitLayouts    []string
}
