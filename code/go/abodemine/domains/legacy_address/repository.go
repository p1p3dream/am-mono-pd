package legacy_address

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
)

type Repository interface {
	SearchByAddress(r *arc.Request, input *SearchRecordInput, index []string) (*SearchRecordOutput, error)
}

type repository struct{}

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source SearchRecordOutput `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (repo *repository) SearchByAddress(r *arc.Request, input *SearchRecordInput, index []string) (*SearchRecordOutput, error) {
	c, err := r.Dom().SelectOpenSearch(consts.ConfigKeyOpenSearchLegacyApi)
	if err != nil {
		return nil, errors.Forward(err, "3cfa9b1d-3fd8-4141-b275-36703fa3102c")
	}

	query := buildFuzzySearchQuery(input)

	queryJSON, err := json.Marshal(&query)
	if err != nil {
		return nil, &errors.Object{
			Id:     "5087e9c3-6431-49c8-bf52-f6562b591267",
			Code:   errors.Code_UNKNOWN,
			Detail: "Error marshalling OpenSearch query",
			Cause:  err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchRequest := opensearchapi.SearchRequest{
		Index: index,
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := searchRequest.Do(ctx, c)
	if err != nil {
		return nil, &errors.Object{
			Id:     "2bb65958-56e7-478c-b3e4-44e76c628af9",
			Code:   errors.Code_UNKNOWN,
			Detail: "Error searching OpenSearch",
			Cause:  err.Error(),
		}
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, &errors.Object{
			Id:     "ffc9b6a5-110b-418c-bfed-3c71202dc88d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Error searching OpenSearch",
			Cause:  fmt.Sprintf("status=%d, body=%s", res.StatusCode, string(body)),
		}
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, &errors.Object{
			Id:     "d4faf4a4-efd3-4fc9-83c3-836300325513",
			Code:   errors.Code_UNKNOWN,
			Detail: "Error decoding OpenSearch response",
			Cause:  err.Error(),
		}
	}

	if len(searchResponse.Hits.Hits) == 0 {
		return nil, nil
	}

	return &searchResponse.Hits.Hits[0].Source, nil
}

// buildFuzzySearchQuery builds a query for approximate search based on the filled fields
func buildFuzzySearchQuery(input *SearchRecordInput) map[string]any {
	shouldQueries := []map[string]any{}
	mustQueries := []map[string]any{}

	// Fields that require exact matches go into must queries
	if input.FullStreetAddress != "" {
		mustQueries = append(mustQueries, buildMatchQuery("SitusFullStreetAddress", input.FullStreetAddress))
	}

	if input.ZIP5 != "" {
		mustQueries = append(mustQueries, buildMatchQuery("SitusZIP5", input.ZIP5))
	}

	if input.HouseNbr != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusHouseNbr", input.HouseNbr))
	}

	if input.Street != "" {
		shouldQueries = append(shouldQueries, buildFuzzyQuery("SitusStreet", input.Street))
	}

	if input.City != "" {
		shouldQueries = append(shouldQueries, buildFuzzyQuery("SitusCity", input.City))
	}

	if input.State != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusState", input.State))
	}

	if input.ZIP4 != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusZIP4", input.ZIP4))
	}

	if input.UnitType != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusUnitType", input.UnitType))
	}

	if input.UnitNbr != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusUnitNbr", input.UnitNbr))
	}

	if input.DirectionLeft != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusDirectionLeft", input.DirectionLeft))
	}

	if input.DirectionRight != "" {
		shouldQueries = append(shouldQueries, buildMatchQuery("SitusDirectionRight", input.DirectionRight))
	}

	// Build query structure
	boolQuery := map[string]any{}

	// Add must clauses if we have any
	if len(mustQueries) > 0 {
		boolQuery["must"] = mustQueries
	}

	// Add should clauses if we have any
	if len(shouldQueries) > 0 {
		boolQuery["should"] = shouldQueries
		// Only use minimum_should_match if we have should clauses
		if len(mustQueries) == 0 {
			// If we don't have any must clauses, at least one should clause must match
			boolQuery["minimum_should_match"] = 1
		}
	}

	query := map[string]any{
		"query": map[string]any{
			"bool": boolQuery,
		},
		"size": input.Limit,
	}

	if len(input.DesiredFields) > 0 {
		query["_source"] = input.DesiredFields
	}

	return query
}

// buildMatchQuery creates an exact match query
func buildMatchQuery(field, value string) map[string]any {
	return map[string]any{
		"match": map[string]any{
			field: value,
		},
	}
}

// buildFuzzyQuery creates a fuzzy query to allow approximate matches
func buildFuzzyQuery(field, value string) map[string]any {
	return map[string]any{
		"fuzzy": map[string]any{
			field: map[string]any{
				"value":          value,
				"fuzziness":      "AUTO",
				"prefix_length":  2,
				"max_expansions": 50,
			},
		},
	}
}
