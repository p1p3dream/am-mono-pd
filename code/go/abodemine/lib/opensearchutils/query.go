package opensearchutils

// MatchQuery represents a match query in OpenSearch
type MatchQuery struct {
	Field string
	Value string
}

// NewMatchQuery creates a new match query
func NewMatchQuery(field, value string) *MatchQuery {
	return &MatchQuery{
		Field: field,
		Value: value,
	}
}

// ToMap converts the match query to a map for OpenSearch
func (m *MatchQuery) ToMap() map[string]any {
	return map[string]any{
		"match": map[string]any{
			m.Field: m.Value,
		},
	}
}

// TermQuery represents a term query in OpenSearch
type TermQuery struct {
	Field string
	Value string
}

// NewTermQuery creates a new term query
func NewTermQuery(field, value string) *TermQuery {
	return &TermQuery{
		Field: field,
		Value: value,
	}
}

// ToMap converts the term query to a map for OpenSearch
func (t *TermQuery) ToMap() map[string]any {
	return map[string]any{
		"term": map[string]any{
			t.Field: t.Value,
		},
	}
}

// ExistsQuery represents an exists query in OpenSearch
type ExistsQuery struct {
	Field string
}

// NewExistsQuery creates a new exists query
func NewExistsQuery(field string) *ExistsQuery {
	return &ExistsQuery{
		Field: field,
	}
}

// ToMap converts the exists query to a map for OpenSearch
func (e *ExistsQuery) ToMap() map[string]any {
	return map[string]any{
		"exists": map[string]any{
			"field": e.Field,
		},
	}
}

// Query represents a generic OpenSearch query that can be converted to a map
type Query interface {
	ToMap() map[string]any
}
