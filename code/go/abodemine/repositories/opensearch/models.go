package opensearch

type Document interface {
	OpenSearchId() string
}

type opensearchActionMetadata struct {
	Index *opensearchActionMetadataIndex `json:"index"`
}

type opensearchActionMetadataIndex struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}

type IndexCreateBody struct {
	Settings map[string]*IndexCreateBodySetting `json:"settings,omitempty"`
	Mappings map[string]any                     `json:"mappings,omitempty"`
}

type IndexCreateBodySetting struct {
	NumberOfShards   int `json:"number_of_shards,omitempty"`
	NumberOfReplicas int `json:"number_of_replicas,omitempty"`
}

type IndexCreateBodyMappingProperty struct {
	Type  string `json:"type,omitempty"`
	Index *bool  `json:"index,omitempty"`
}

// type BulkResponse struct {
// 	Items []*BulkResponseItem `json:"items,omitempty"`
// }

// type BulkResponseItem struct {
// }
