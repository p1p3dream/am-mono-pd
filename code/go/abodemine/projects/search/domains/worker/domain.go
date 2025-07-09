package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/rs/zerolog/log"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/projects/search/conf"
)

// Domain represents the worker domain
type Domain struct {
	config *conf.Config
}

// NewDomain creates a new worker domain
func NewDomain(config *conf.Config) *Domain {
	return &Domain{
		config: config,
	}
}

// SyncAddressesInput represents the input for syncing addresses
type SyncAddressesInput struct {
	BatchSize int
	IndexName string
	Zip5      string // Optional ZIP5 code to filter addresses by
	FIPS      string // Optional FIPS code to filter addresses by
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Count    int
	Duration time.Duration
	Error    error
}

// AddressQueryParams contains parameters for building address queries
type AddressQueryParams struct {
	LastID    string
	BatchSize int
	Zip5      string
	FIPS      string
}

// OpenSearchIndex contains information about an OpenSearch index
type OpenSearchIndex struct {
	Name     string
	IsExists bool
}

// FIPSProcessingStatus represents the status of FIPS code processing
type FIPSProcessingStatus struct {
	FIPS       string    `json:"fips" mapping_type:"keyword"`
	Processed  bool      `json:"processed" mapping_type:"boolean"`
	Count      int       `json:"count" mapping_type:"integer"`
	StartedAt  time.Time `json:"started_at" mapping_type:"date"`
	FinishedAt time.Time `json:"finished_at" mapping_type:"date"`
	Error      *string   `json:"error" mapping_type:"text"`
}

// SyncAddresses is the main entry point for syncing addresses
func (d *Domain) SyncAddresses(r *arc.Request, input *SyncAddressesInput) (int, error) {
	ctx := r.Context()

	// Get clients
	osClient, pgxPool, err := getClients(r)
	if err != nil {
		return 0, err
	}

	// Generate current mapping from struct
	currentMapping := GenerateAddressMapping()

	// Check if there's an existing index with the alias
	hasNewMapping, currentIndex, err := checkMappingChanges(ctx, osClient, input.IndexName, currentMapping)
	if err != nil {
		return 0, err
	}

	var indexToUse string

	if hasNewMapping {
		// Mapping has changed, create new index
		log.Info().Msg("Mapping changes detected, creating new index")
		newIndex, oldIndex, err := setupIndex(ctx, osClient, input.IndexName)
		if err != nil {
			return 0, err
		}
		indexToUse = newIndex.Name
		currentIndex = oldIndex
	} else {
		// No mapping changes, use existing index
		log.Info().
			Str("index", currentIndex.Name).
			Msg("No mapping changes detected, using existing index")
		indexToUse = currentIndex.Name
	}

	// Create FIPS processing status index if it doesn't exist
	createFIPSIndexReq := opensearchapi.IndicesCreateRequest{
		Index: "fips_processing_status",
		Body:  strings.NewReader(GenerateFIPSStatusMapping()),
	}

	createFIPSIndexRes, err := createFIPSIndexReq.Do(ctx, osClient)
	if err != nil && !strings.Contains(err.Error(), "resource_already_exists_exception") {
		return 0, err
	}
	if createFIPSIndexRes != nil {
		defer createFIPSIndexRes.Body.Close()
	}

	// Check if alias exists, if not create it
	getAliasReq := opensearchapi.IndicesGetAliasRequest{
		Name: []string{input.IndexName},
	}

	aliasRes, err := getAliasReq.Do(ctx, osClient)
	if err != nil || aliasRes.StatusCode != 200 {
		// Alias doesn't exist, create it
		log.Info().Msg("Alias doesn't exist, creating it")
		aliasActions := []map[string]interface{}{
			{
				"add": map[string]interface{}{
					"index": indexToUse,
					"alias": input.IndexName,
				},
			},
		}

		aliasActionBody := map[string]interface{}{
			"actions": aliasActions,
		}

		aliasActionJSON, err := json.Marshal(aliasActionBody)
		if err != nil {
			return 0, err
		}

		updateAliasReq := opensearchapi.IndicesUpdateAliasesRequest{
			Body: strings.NewReader(string(aliasActionJSON)),
		}

		updateAliasRes, err := updateAliasReq.Do(ctx, osClient)
		if err != nil {
			return 0, err
		}
		defer updateAliasRes.Body.Close()

		if updateAliasRes.IsError() {
			responseBody, _ := io.ReadAll(updateAliasRes.Body)
			return 0, fmt.Errorf("failed to create alias: %s", responseBody)
		}

		log.Info().
			Str("alias", input.IndexName).
			Str("index", indexToUse).
			Msg("Successfully created alias")
	}

	// Get FIPS codes to process
	var fipsCodes []string
	if input.FIPS != "" {
		// If specific FIPS is provided, use only that
		fipsCodes = []string{input.FIPS}
	} else {
		// Get unprocessed FIPS codes for today
		unprocessedFIPS, err := getUnprocessedFIPS(ctx, osClient, pgxPool)
		if err != nil {
			return 0, err
		}

		if len(unprocessedFIPS) == 0 {
			log.Info().Msg("No unprocessed FIPS codes found, all FIPS codes have been successfully processed in the last 24 hours")
			return 0, nil
		}

		fipsCodes = unprocessedFIPS
	}

	totalProcessed := 0
	const parallelFIPS = 4 // Number of FIPS codes to process in parallel
	startTime := time.Now()
	totalFIPS := len(fipsCodes)

	// Process FIPS codes in batches of parallelFIPS
	for i := 0; i < len(fipsCodes); i += parallelFIPS {
		end := i + parallelFIPS
		if end > len(fipsCodes) {
			end = len(fipsCodes)
		}

		batch := fipsCodes[i:end]
		log.Info().
			Int("batch_start", i).
			Int("batch_end", end).
			Int("batch_size", len(batch)).
			Msg("Starting batch of FIPS codes")

		// Create a wait group to wait for all FIPS in this batch
		var wg sync.WaitGroup
		results := make(chan struct {
			fips  string
			count int
			err   error
		}, len(batch))

		// Process each FIPS in the batch in parallel
		for _, fips := range batch {
			wg.Add(1)
			go func(fips string) {
				defer wg.Done()

				// Create FIPS status
				status := &FIPSProcessingStatus{
					FIPS:      fips,
					Processed: false,
					StartedAt: time.Now(),
				}

				// Update status to started
				if err := updateFIPSStatus(ctx, osClient, status); err != nil {
					log.Error().Err(err).Str("fips", fips).Msg("Failed to update FIPS status")
					results <- struct {
						fips  string
						count int
						err   error
					}{fips: fips, count: 0, err: err}
					return
				}

				log.Info().Str("fips", fips).Msg("Starting to process FIPS")

				// Create channels for collecting results
				addressChan := make(chan SyncResult)

				// Start indexing process for this FIPS
				go runAddressSync(ctx, osClient, pgxPool, indexToUse, input.BatchSize, input.Zip5, fips, addressChan)

				// Collect results
				addressResult := <-addressChan

				// Update status
				status.Processed = true
				status.FinishedAt = time.Now()
				status.Count = addressResult.Count
				if addressResult.Error != nil {
					errorMsg := addressResult.Error.Error()
					status.Error = &errorMsg
				}

				if err := updateFIPSStatus(ctx, osClient, status); err != nil {
					log.Error().Err(err).Str("fips", fips).Msg("Failed to update FIPS status")
				}

				results <- struct {
					fips  string
					count int
					err   error
				}{fips: fips, count: addressResult.Count, err: addressResult.Error}
			}(fips)
		}

		// Wait for all FIPS in this batch to complete
		wg.Wait()
		close(results)

		// Process results from this batch
		for result := range results {
			if result.err != nil {
				log.Error().Err(result.err).Str("fips", result.fips).Msg("Error processing FIPS")
				continue
			}

			totalProcessed += result.count
			log.Info().
				Str("fips", result.fips).
				Int("count", result.count).
				Int("total_processed", totalProcessed).
				Msg("Finished processing FIPS")
		}

		// Calculate progress and time estimates
		elapsedTime := time.Since(startTime)
		completedFIPS := end
		remainingFIPS := totalFIPS - completedFIPS
		progressPercent := float64(completedFIPS) / float64(totalFIPS) * 100

		// Calculate average time per FIPS and estimated remaining time
		avgTimePerFIPS := elapsedTime / time.Duration(completedFIPS)
		estimatedRemainingTime := avgTimePerFIPS * time.Duration(remainingFIPS)
		estimatedTotalTime := elapsedTime + estimatedRemainingTime

		log.Info().
			Int("batch_start", i).
			Int("batch_end", end).
			Int("total_processed", totalProcessed).
			Float64("progress_percent", progressPercent).
			Int("completed_fips", completedFIPS).
			Int("remaining_fips", remainingFIPS).
			Dur("elapsed_time", elapsedTime).
			Dur("avg_time_per_fips", avgTimePerFIPS).
			Dur("estimated_remaining_time", estimatedRemainingTime).
			Dur("estimated_total_time", estimatedTotalTime).
			Msg("Finished batch of FIPS codes")
	}

	// If we created a new index and alias didn't exist before, we don't need to switch the alias
	// since we already created it at the beginning
	if hasNewMapping && (aliasRes == nil || aliasRes.StatusCode != 200) {
		log.Info().Msg("Alias was created at the beginning, no need to switch")
	} else if hasNewMapping {
		// Only switch alias if it existed before and we created a new index
		if err := updateAlias(ctx, osClient, currentIndex, OpenSearchIndex{Name: indexToUse, IsExists: true}, input.IndexName); err != nil {
			return totalProcessed, err
		}
		log.Info().Msg("Alias updated to point to new index")
	}

	// Log final statistics
	totalElapsedTime := time.Since(startTime)
	log.Info().
		Int("total_fips_processed", totalFIPS).
		Int("total_addresses_processed", totalProcessed).
		Dur("total_elapsed_time", totalElapsedTime).
		Dur("avg_time_per_fips", totalElapsedTime/time.Duration(totalFIPS)).
		Msg("Address sync completed successfully")

	return totalProcessed, nil
}

// runAddressSync runs the address sync process
func runAddressSync(ctx context.Context, osClient *opensearch.Client, pgxPool *pgxpool.Pool,
	indexName string, batchSize int, zip5 string, fips string, resultChan chan<- SyncResult) {

	syncStartTime := time.Now()
	totalProcessed := 0

	// Reduce batch size to prevent bulk request size issues
	adjustedBatchSize := batchSize
	if adjustedBatchSize < 100 {
		adjustedBatchSize = 100 // Minimum batch size
	}

	// Get total count before starting
	totalCount, err := getTotalAddressCount(ctx, pgxPool, zip5)
	if err != nil {
		resultChan <- SyncResult{
			Count:    0,
			Duration: time.Since(syncStartTime),
			Error:    err,
		}
		return
	}

	log.Info().
		Int("total_addresses", totalCount).
		Int("batch_size", adjustedBatchSize).
		Str("zip5", zip5).
		Str("fips", fips).
		Msg("Starting Address sync process")

	// Create a semaphore to limit parallel operations
	maxParallel := 10
	semaphore := make(chan struct{}, maxParallel)
	errorChan := make(chan error, 1)

	// Start PostgreSQL data retrieval
	lastID := ""
	for {
		// Fetch addresses from the addresses table
		addresses, err := fetchAddressesFromTable(ctx, pgxPool, AddressQueryParams{
			LastID:    lastID,
			BatchSize: adjustedBatchSize,
			Zip5:      zip5,
			FIPS:      fips,
		})

		if err != nil {
			resultChan <- SyncResult{
				Count:    totalProcessed,
				Duration: time.Since(syncStartTime),
				Error:    err,
			}
			return
		}

		// If no records were found, we're done
		if len(addresses) == 0 {
			break
		}

		// Update lastID for next batch
		lastID = addresses[len(addresses)-1].AMId

		// Acquire semaphore before processing
		select {
		case semaphore <- struct{}{}:
		case <-ctx.Done():
			resultChan <- SyncResult{
				Count:    totalProcessed,
				Duration: time.Since(syncStartTime),
				Error:    ctx.Err(),
			}
			return
		}

		// Process addresses in a goroutine
		go func(addrs []*address.AddressDocument, batchLastID string) {
			defer func() { <-semaphore }() // Release semaphore when done

			// Format address text fields
			formatAddressFields(addrs)

			// Index the addresses
			count, _, err := indexAddresses(ctx, addrs, osClient, indexName)
			if err != nil {
				select {
				case errorChan <- err:
				case <-ctx.Done():
				}
				return
			}

			// Update total processed
			totalProcessed += count

			// Log batch completion
			log.Info().
				Str("process", "Address").
				Str("last_id", batchLastID).
				Str("batch_size", fmt.Sprintf("%d", len(addrs))).
				Str("total_processed", fmt.Sprintf("%d", totalProcessed)).
				Msg("Batch processing completed")
		}(addresses, lastID)
	}

	// Wait for all semaphore slots to be released
	for i := 0; i < maxParallel; i++ {
		semaphore <- struct{}{}
	}

	// Check for any errors
	select {
	case err := <-errorChan:
		resultChan <- SyncResult{
			Count:    totalProcessed,
			Duration: time.Since(syncStartTime),
			Error:    err,
		}
		return
	default:
	}

	// Send success result
	log.Info().
		Str("total_processed", fmt.Sprintf("%d", totalProcessed)).
		Str("total_count", fmt.Sprintf("%d", totalCount)).
		Str("total_duration", fmt.Sprintf("%dms", time.Since(syncStartTime).Milliseconds())).
		Msg("Address sync completed successfully")

	resultChan <- SyncResult{
		Count:    totalProcessed,
		Duration: time.Since(syncStartTime),
		Error:    nil,
	}
}

// getTotalAddressCount gets the total number of addresses in the database
func getTotalAddressCount(ctx context.Context, pgxPool *pgxpool.Pool, zip5 string) (int, error) {
	var count int
	if zip5 != "" {
		err := pgxPool.QueryRow(ctx, "SELECT COUNT(*) FROM addresses WHERE zip5 = $1", zip5).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		err := pgxPool.QueryRow(ctx, "SELECT COUNT(*) FROM addresses").Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

// getClients retrieves the necessary clients for the sync operation
func getClients(r *arc.Request) (*opensearch.Client, *pgxpool.Pool, error) {
	// Get OpenSearch client
	osClient, err := r.Dom().SelectOpenSearch(consts.ConfigKeyOpenSearchSearch)
	if err != nil {
		return nil, nil, &errors.Object{
			Id:     "1e2c5cfa-f1d6-4ecb-b968-2e219b3a5b62",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get OpenSearch client",
			Cause:  err.Error(),
		}
	}

	// Get database client
	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, nil, &errors.Object{
			Id:     "1d4b1ee4-2db9-4cba-a8db-9fbe2972e502",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get database client",
			Cause:  err.Error(),
		}
	}

	return osClient, pgxPool, nil
}

// setupIndex creates a new index and checks for an existing one
func setupIndex(ctx context.Context, osClient *opensearch.Client, baseIndexName string) (OpenSearchIndex, OpenSearchIndex, error) {
	// Generate mapping from struct
	mapping := GenerateAddressMapping()
	timestamp := time.Now().Format("20060102_150405")
	newIndexName := fmt.Sprintf("%s_%s", baseIndexName, timestamp)

	log.Info().
		Str("base_index_name", baseIndexName).
		Str("new_index_name", newIndexName).
		Msg("Starting zero-downtime sync process")

	// Create new index with mapping
	createReq := opensearchapi.IndicesCreateRequest{
		Index: newIndexName,
		Body:  strings.NewReader(mapping),
	}

	createRes, err := createReq.Do(ctx, osClient)
	if err != nil {
		return OpenSearchIndex{}, OpenSearchIndex{}, &errors.Object{
			Id:     "7ee88808-c73d-4392-9c58-e354180b870e",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create index",
			Cause:  err.Error(),
		}
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		responseBody, _ := io.ReadAll(createRes.Body)
		return OpenSearchIndex{}, OpenSearchIndex{}, &errors.Object{
			Id:     "d6f662f9-0da3-40ca-8867-3994110f355b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create index",
			Cause:  fmt.Sprintf("failed to create index: %s", responseBody),
		}
	}

	log.Info().Str("index", newIndexName).Msg("Created new index")

	// Get current alias information if it exists
	currentIndex := OpenSearchIndex{}
	getAliasReq := opensearchapi.IndicesGetAliasRequest{
		Name: []string{baseIndexName},
	}

	aliasRes, err := getAliasReq.Do(ctx, osClient)

	// Parse current alias if it exists
	if err != nil || aliasRes.StatusCode != 200 {
		log.Info().Msg("No existing alias found, will create new one")
	} else {
		defer aliasRes.Body.Close()
		var aliasData map[string]interface{}
		if err := json.NewDecoder(aliasRes.Body).Decode(&aliasData); err != nil {
			log.Warn().Err(err).Msg("Failed to parse alias response, will continue regardless")
		} else {
			// Get current index from alias
			for indexName := range aliasData {
				currentIndex.Name = indexName
				currentIndex.IsExists = true
				break
			}
			log.Info().
				Bool("has_existing_index", currentIndex.IsExists).
				Str("current_index", currentIndex.Name).
				Msg("Current alias information")
		}
	}

	return OpenSearchIndex{Name: newIndexName, IsExists: true}, currentIndex, nil
}

// formatAddressFields formats the text fields of addresses
func formatAddressFields(addresses []*address.AddressDocument) {
	for _, addr := range addresses {
		if addr.Street != nil && *addr.Street != "" {
			formattedStreet := FormatAddressText(*addr.Street)
			addr.Street = &formattedStreet
		}
		if addr.PreDirectional != nil && *addr.PreDirectional != "" {
			formattedPreDir := FormatAddressText(*addr.PreDirectional)
			addr.PreDirectional = &formattedPreDir
		}
		if addr.PostDirectional != nil && *addr.PostDirectional != "" {
			formattedPostDir := FormatAddressText(*addr.PostDirectional)
			addr.PostDirectional = &formattedPostDir
		}
		if addr.StreetType != nil && *addr.StreetType != "" {
			formattedStreetType := FormatAddressText(*addr.StreetType)
			addr.StreetType = &formattedStreetType
		}
		if addr.UnitType != nil && *addr.UnitType != "" {
			formattedUnitType := FormatAddressText(*addr.UnitType)
			addr.UnitType = &formattedUnitType
		}
		if addr.City != nil && *addr.City != "" {
			formattedCity := FormatAddressText(*addr.City)
			addr.City = &formattedCity
		}
	}
}

// updateAlias updates the alias to point to the new index
func updateAlias(ctx context.Context, osClient *opensearch.Client, currentIndex, newIndex OpenSearchIndex, aliasName string) error {
	var aliasActions []map[string]interface{}

	// Get all indices with the addresses_ prefix
	getIndicesReq := opensearchapi.CatIndicesRequest{
		Format: "json",
		Index:  []string{"addresses_*"},
	}

	indicesRes, err := getIndicesReq.Do(ctx, osClient)
	if err != nil {
		return &errors.Object{
			Id:     "bdc94a77-031c-4427-94f2-0d91562bd091",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get indices",
			Cause:  err.Error(),
		}
	}
	defer indicesRes.Body.Close()

	var indices []map[string]interface{}
	if err := json.NewDecoder(indicesRes.Body).Decode(&indices); err != nil {
		return &errors.Object{
			Id:     "3e95f0f9-f2e0-4dcb-a814-5ea3e9fe4286",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to decode indices response",
			Cause:  err.Error(),
		}
	}

	// Delete all addresses_* indices except the new one
	for _, index := range indices {
		indexName := index["index"].(string)
		if indexName != newIndex.Name {
			deleteReq := opensearchapi.IndicesDeleteRequest{
				Index: []string{indexName},
			}
			deleteRes, err := deleteReq.Do(ctx, osClient)
			if err != nil {
				log.Warn().
					Err(err).
					Str("index", indexName).
					Msg("Failed to delete old index")
				continue
			}
			defer deleteRes.Body.Close()

			if deleteRes.IsError() {
				responseBody, _ := io.ReadAll(deleteRes.Body)
				log.Warn().
					Str("index", indexName).
					Str("response", string(responseBody)).
					Msg("Failed to delete old index")
				continue
			}

			log.Info().
				Str("index", indexName).
				Msg("Successfully deleted old index")
		}
	}

	// Add the new index to the alias
	aliasActions = append(aliasActions, map[string]interface{}{
		"add": map[string]interface{}{
			"index": newIndex.Name,
			"alias": aliasName,
		},
	})

	// Build the alias update request
	aliasActionBody := map[string]interface{}{
		"actions": aliasActions,
	}

	aliasActionJSON, err := json.Marshal(aliasActionBody)
	if err != nil {
		return &errors.Object{
			Id:     "0287e078-88f0-4f43-976c-e0159f0fce41",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to marshal alias action",
			Cause:  err.Error(),
		}
	}

	// Execute the alias update
	updateAliasReq := opensearchapi.IndicesUpdateAliasesRequest{
		Body: strings.NewReader(string(aliasActionJSON)),
	}

	updateAliasRes, err := updateAliasReq.Do(ctx, osClient)
	if err != nil {
		return &errors.Object{
			Id:     "6cbf60fb-712a-4a3b-a70b-0e7f40598a3a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update alias",
			Cause:  err.Error(),
		}
	}
	defer updateAliasRes.Body.Close()

	if updateAliasRes.IsError() {
		responseBody, _ := io.ReadAll(updateAliasRes.Body)
		return &errors.Object{
			Id:     "589c2d8d-6078-44b3-abc4-f35550536f98",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update alias",
			Cause:  fmt.Sprintf("failed to update alias: %s", responseBody),
		}
	}

	log.Info().
		Str("alias", aliasName).
		Str("new_index", newIndex.Name).
		Msg("Successfully updated alias to point to new index")

	return nil
}

// GenerateAddressMapping creates an OpenSearch mapping based on AddressDocument struct
func GenerateAddressMapping() string {
	t := reflect.TypeOf(address.AddressDocument{})
	properties := make(map[string]map[string]string)

	// Process each field in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip fields that are not exported or have no json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse the json tag to get the field name
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			continue
		}

		// Get the mapping type from the tag
		mappingType := field.Tag.Get("mapping_type")
		if mappingType == "" {
			// Default to keyword for strings if not specified
			if field.Type.Kind() == reflect.String ||
				(field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String) {
				mappingType = "keyword"
			} else {
				// Skip fields without mapping type
				continue
			}
		}

		if mappingType == "store" {
			properties[jsonName] = map[string]string{"type": "keyword", "index": "false"}
		} else {
			properties[jsonName] = map[string]string{"type": mappingType}
		}
	}

	// Create the mapping JSON
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
		},
		"mappings": map[string]interface{}{
			"dynamic":    "strict",
			"properties": properties,
		},
	}

	// Marshal to JSON string
	bytes, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal address mapping")
		// Return a basic mapping as fallback
		return `{"settings":{"index":{"number_of_shards":1,"number_of_replicas":1}},"mappings":{"dynamic":"strict","properties":{"am_id":{"type":"keyword"}}}}`
	}

	return string(bytes)
}

// indexAddresses indexes a list of addresses in OpenSearch
func indexAddresses(ctx context.Context, addresses []*address.AddressDocument, osClient *opensearch.Client, indexName string) (int, time.Time, error) {
	count := 0
	latest := time.Now()

	// Create a buffer for bulk indexing
	var bulkBody strings.Builder

	// Process each address
	for _, address := range addresses {
		// Track the latest update time
		if address.AMUpdatedAt.After(latest) {
			latest = address.AMUpdatedAt
		}

		// Add geo_point field for geospatial queries
		if address.Latitude != 0 || address.Longitude != 0 {
			address.Location = fmt.Sprintf("%f,%f", address.Latitude, address.Longitude)
		}

		// Create action line for bulk indexing
		actionLine := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    address.AMId,
			},
		}
		actionLineJSON, err := json.Marshal(actionLine)
		if err != nil {
			log.Error().Err(err).Msg("Error marshaling action line")
			continue
		}

		// Create document line
		documentJSON, err := json.Marshal(address)
		if err != nil {
			log.Error().Err(err).Msg("Error marshaling address document")
			continue
		}

		// Add to bulk request
		bulkBody.Write(actionLineJSON)
		bulkBody.WriteString("\n")
		bulkBody.Write(documentJSON)
		bulkBody.WriteString("\n")

		count++
	}

	// Send all documents in one batch
	if bulkBody.Len() > 0 {
		// Retry configuration
		maxRetries := 5
		baseDelay := 1 * time.Second
		maxDelay := 30 * time.Second

		var lastErr error
		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				// Calculate exponential backoff delay
				delay := baseDelay * time.Duration(1<<uint(attempt-1))
				if delay > maxDelay {
					delay = maxDelay
				}
				log.Info().
					Int("attempt", attempt+1).
					Dur("delay", delay).
					Msg("Retrying bulk index request after delay")
				time.Sleep(delay)
			}

			req := opensearchapi.BulkRequest{
				Body: strings.NewReader(bulkBody.String()),
			}

			res, err := req.Do(ctx, osClient)
			if err != nil {
				lastErr = err
				log.Error().Err(err).Msg("Error executing bulk request")
				continue
			}
			defer res.Body.Close()

			if res.IsError() {
				responseBody, _ := io.ReadAll(res.Body)
				lastErr = fmt.Errorf("bulk request failed: %s", responseBody)

				// Check if it's a 429 error (too many requests)
				if res.StatusCode == 429 {
					log.Warn().
						Int("attempt", attempt+1).
						Int("status_code", res.StatusCode).
						Str("response", string(responseBody)).
						Msg("OpenSearch rate limit exceeded, will retry")
					continue
				}

				log.Error().
					Int("attempt", attempt+1).
					Int("status_code", res.StatusCode).
					Str("response", string(responseBody)).
					Msg("Bulk request failed")
				continue
			}

			// Parse the response
			var bulkResponse map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
				log.Warn().Err(err).Msg("Failed to parse bulk response")
			} else {
				// Extract useful metrics from response
				took, _ := bulkResponse["took"].(float64)
				errors, hasErrors := bulkResponse["errors"].(bool)

				log.Info().
					Float64("opensearch_took_ms", took).
					Bool("has_errors", hasErrors).
					Msg("Bulk response processed")

				if hasErrors && errors {
					// Log details about the errors
					if items, ok := bulkResponse["items"].([]interface{}); ok {
						errorCount := 0
						for _, item := range items {
							if itemMap, ok := item.(map[string]interface{}); ok {
								if indexInfo, ok := itemMap["index"].(map[string]interface{}); ok {
									if _, hasError := indexInfo["error"]; hasError {
										errorCount++
									}
								}
							}
						}
						log.Warn().
							Int("error_count", errorCount).
							Int("total_items", len(items)).
							Msg("Bulk operation had errors")
					}
				}
			}

			// If we get here, the request was successful
			return count, latest, nil
		}

		// If we've exhausted all retries, return the last error
		return count, latest, &errors.Object{
			Id:     "0789711f-71a6-4406-86b6-b3b895971195",
			Code:   errors.Code_UNKNOWN,
			Detail: "Bulk request failed after retries",
			Cause:  lastErr.Error(),
		}
	}

	return count, latest, nil
}

// fetchAddressesFromTable fetches addresses from the addresses table
func fetchAddressesFromTable(ctx context.Context, pgxPool *pgxpool.Pool, params AddressQueryParams) ([]*address.AddressDocument, error) {
	// Build address query
	query := buildAddressesTableQuery(params)
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// Execute query
	rows, err := pgxPool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	addresses, err := scanAddressesTable(rows)
	if err != nil {
		return nil, err
	}

	log.Info().
		Int("records_found", len(addresses)).
		Msg("Completed scanning rows into address structs")

	return addresses, nil
}

// buildAddressesTableQuery creates an SQL query for the addresses table
func buildAddressesTableQuery(params AddressQueryParams) squirrel.SelectBuilder {
	query := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"a.id",
			"a.created_at",
			"a.updated_at",
			"a.city",
			"a.county",
			"a.data_source",
			"a.fips",
			"a.full_street_address",
			"a.house_number",
			"a.state",
			"a.street_name",
			"a.street_number",
			"a.street_pos_direction",
			"a.street_pre_direction",
			"a.street_suffix",
			"a.street_type",
			"a.unit_nbr",
			"a.unit_type",
			"a.zip5",
			"p.id",
			"p.ad_attom_id",
			"p.fa_property_id",
		).
		From("addresses a").
		LeftJoin("properties p ON a.id = p.address_id")

	// Add ZIP5 filter if specified
	if params.Zip5 != "" {
		query = query.Where(squirrel.Eq{"a.zip5": params.Zip5})
	}

	// Add FIPS filter if specified
	if params.FIPS != "" {
		query = query.Where(squirrel.Eq{"a.fips": params.FIPS})
	}

	// Add cursor-based pagination
	if params.LastID != "" {
		query = query.Where(squirrel.Gt{"a.id": params.LastID})
	}

	return query.OrderBy("a.id").Limit(uint64(params.BatchSize))
}

// scanAddressesTable scans all rows from addresses table into AddressDocument objects
func scanAddressesTable(rows pgx.Rows) ([]*address.AddressDocument, error) {
	var addresses []*address.AddressDocument

	for rows.Next() {
		var address address.AddressDocument
		// Use pointers for all string fields to handle NULL values
		var id string
		var createdAt, updatedAt time.Time
		var city, county, dataSource, fips, fullStreetAddress, houseNumber, state *string
		var streetName, streetNumber, streetPosDirection, streetPreDirection, streetSuffix, streetType, unitNbr, unitType, zip5 *string
		var propertyId *string
		var adAttomId, faPropertyId *int64

		err := rows.Scan(
			&id,
			&createdAt,
			&updatedAt,
			&city,
			&county,
			&dataSource,
			&fips,
			&fullStreetAddress,
			&houseNumber,
			&state,
			&streetName,
			&streetNumber,
			&streetPosDirection,
			&streetPreDirection,
			&streetSuffix,
			&streetType,
			&unitNbr,
			&unitType,
			&zip5,
			&propertyId,
			&adAttomId,
			&faPropertyId,
		)

		if err != nil {
			log.Error().Err(err).
				Str("process", "scanAddressesTable").
				Msg("Error scanning address row")
			continue
		}

		// Assign values to the AddressDocument, ensuring we don't dereference nil pointers
		address.AMId = id
		address.AMUpdatedAt = updatedAt

		// Handle potentially NULL values
		if fips != nil {
			address.FIPS = *fips
		} else {
			address.FIPS = ""
		}

		if zip5 != nil {
			address.ZIP5 = *zip5
		} else {
			address.ZIP5 = ""
		}

		if propertyId != nil {
			address.PropertyId = propertyId
		}

		if adAttomId != nil {
			address.ADAttomId = adAttomId
		}

		if faPropertyId != nil {
			address.FAPropertyId = faPropertyId
		}

		// Assign pointer fields (already nullable)
		address.City = city
		address.County = county
		address.State = state
		address.FullAddress = fullStreetAddress
		address.StreetNumber = streetNumber
		address.Street = streetName
		address.PreDirectional = streetPreDirection
		address.PostDirectional = streetPosDirection
		address.StreetType = streetType
		address.UnitType = unitType
		address.UnitNbr = unitNbr

		// Set source field
		if dataSource != nil {
			source := *dataSource
			address.Source = &source
		}

		// Get full US state name (only if state is not nil)
		if state != nil {
			fullStateName := getFullStateName(state)
			address.StateFullName = &fullStateName
		}

		// Calculate location if needed (skipped since we don't have lat/long in addresses table)

		addresses = append(addresses, &address)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).
			Str("process", "scanAddressesTable").
			Msg("Error during row iteration")
		return nil, err
	}

	return addresses, nil
}

// checkMappingChanges checks if the current mapping is different from the existing index mapping
// Returns true if mappings are different, the current index name, and any error
func checkMappingChanges(ctx context.Context, osClient *opensearch.Client, aliasName, currentMapping string) (bool, OpenSearchIndex, error) {
	currentIndex := OpenSearchIndex{}

	// Get current alias information if it exists
	getAliasReq := opensearchapi.IndicesGetAliasRequest{
		Name: []string{aliasName},
	}

	aliasRes, err := getAliasReq.Do(ctx, osClient)

	// Check if alias doesn't exist yet
	if err != nil || aliasRes.StatusCode != 200 {
		log.Info().Msg("No existing alias found, will create new index")
		return true, currentIndex, nil
	}

	defer aliasRes.Body.Close()
	var aliasData map[string]interface{}
	if err := json.NewDecoder(aliasRes.Body).Decode(&aliasData); err != nil {
		log.Warn().Err(err).Msg("Failed to parse alias response, will create new index")
		return true, currentIndex, nil
	}

	// Get current index from alias
	for indexName := range aliasData {
		currentIndex.Name = indexName
		currentIndex.IsExists = true
		break
	}

	if !currentIndex.IsExists {
		log.Info().Msg("No existing index found for alias, will create new index")
		return true, currentIndex, nil
	}

	// Get the mapping from the existing index
	getMappingReq := opensearchapi.IndicesGetMappingRequest{
		Index: []string{currentIndex.Name},
	}

	mappingRes, err := getMappingReq.Do(ctx, osClient)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get mapping, will create new index")
		return true, currentIndex, nil
	}
	defer mappingRes.Body.Close()

	var mappingData map[string]interface{}
	if err := json.NewDecoder(mappingRes.Body).Decode(&mappingData); err != nil {
		log.Warn().Err(err).Msg("Failed to decode mapping response, will create new index")
		return true, currentIndex, nil
	}

	// Extract mapping JSON
	indexMapping, ok := mappingData[currentIndex.Name].(map[string]interface{})
	if !ok {
		log.Warn().Msg("Failed to parse index mapping, will create new index")
		return true, currentIndex, nil
	}

	// Compare mappings - first convert current mapping to comparable format
	var currentMappingObj, existingMappingObj map[string]interface{}
	if err := json.Unmarshal([]byte(currentMapping), &currentMappingObj); err != nil {
		log.Warn().Err(err).Msg("Failed to parse current mapping, will create new index")
		return true, currentIndex, nil
	}

	// Serialize and deserialize to normalize the JSON for comparison
	existingMappingJSON, err := json.Marshal(indexMapping["mappings"])
	if err != nil {
		log.Warn().Err(err).Msg("Failed to serialize existing mapping, will create new index")
		return true, currentIndex, nil
	}

	if err := json.Unmarshal(existingMappingJSON, &existingMappingObj); err != nil {
		log.Warn().Err(err).Msg("Failed to parse existing mapping, will create new index")
		return true, currentIndex, nil
	}

	// Compare only the mappings section
	currentMappingsJSON, err := json.Marshal(currentMappingObj["mappings"])
	if err != nil {
		log.Warn().Err(err).Msg("Failed to serialize current mappings, will create new index")
		return true, currentIndex, nil
	}

	existingMappingsJSON, err := json.Marshal(existingMappingObj)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to serialize existing mappings, will create new index")
		return true, currentIndex, nil
	}

	// Compare the normalized mappings
	if string(currentMappingsJSON) != string(existingMappingsJSON) {
		log.Info().Msg("Mapping has changed, will create new index")
		return true, currentIndex, nil
	}

	// Mappings are the same
	log.Info().Msg("Mapping is unchanged, will use existing index")
	return false, currentIndex, nil
}

// GenerateFIPSStatusMapping creates an OpenSearch mapping for FIPS processing status
func GenerateFIPSStatusMapping() string {
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
		},
		"mappings": map[string]interface{}{
			"dynamic": "strict",
			"properties": map[string]interface{}{
				"fips": map[string]string{
					"type": "keyword",
				},
				"processed": map[string]string{
					"type": "boolean",
				},
				"count": map[string]string{
					"type": "integer",
				},
				"started_at": map[string]string{
					"type": "date",
				},
				"finished_at": map[string]string{
					"type": "date",
				},
				"error": map[string]string{
					"type": "text",
				},
			},
		},
	}

	bytes, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal FIPS status mapping")
		return `{"settings":{"index":{"number_of_shards":1,"number_of_replicas":1}},"mappings":{"dynamic":"strict","properties":{"fips":{"type":"keyword"},"processed":{"type":"boolean"},"count":{"type":"integer"},"started_at":{"type":"date"},"finished_at":{"type":"date"},"error":{"type":"text"}}}}`
	}

	return string(bytes)
}

// getDistinctFIPS retrieves all distinct FIPS codes from the addresses table
func getDistinctFIPS(ctx context.Context, pgxPool *pgxpool.Pool) ([]string, error) {
	rows, err := pgxPool.Query(ctx, "SELECT DISTINCT fips FROM addresses WHERE fips IS NOT NULL ORDER BY fips")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fipsCodes []string
	for rows.Next() {
		var fips string
		if err := rows.Scan(&fips); err != nil {
			return nil, err
		}
		fipsCodes = append(fipsCodes, fips)
	}

	return fipsCodes, nil
}

// getUnprocessedFIPS retrieves FIPS codes that haven't been successfully processed in the last 24 hours
func getUnprocessedFIPS(ctx context.Context, osClient *opensearch.Client, pgxPool *pgxpool.Pool) ([]string, error) {
	// Get all distinct FIPS codes from the database
	allFIPS, err := getDistinctFIPS(ctx, pgxPool)
	if err != nil {
		return nil, err
	}

	// Get FIPS codes that were successfully processed in the last 24 hours
	oneDayAgo := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"processed": true,
						},
					},
					{
						"range": map[string]interface{}{
							"finished_at": map[string]interface{}{
								"gte": oneDayAgo,
							},
						},
					},
				},
			},
		},
		"size":    10000, // Adjust based on expected number of FIPS codes
		"_source": []string{"fips"},
	}

	searchReq := opensearchapi.SearchRequest{
		Index: []string{"fips_processing_status"},
		Body:  strings.NewReader(mustMarshal(query)),
	}

	resp, err := searchReq.Do(ctx, osClient)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source struct {
					FIPS string `json:"fips"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Create a set of recently processed FIPS codes
	recentlyProcessedFIPS := make(map[string]bool)
	for _, hit := range result.Hits.Hits {
		recentlyProcessedFIPS[hit.Source.FIPS] = true
	}

	// Find FIPS codes that haven't been processed in the last 24 hours
	var unprocessedFIPS []string
	for _, fips := range allFIPS {
		if !recentlyProcessedFIPS[fips] {
			unprocessedFIPS = append(unprocessedFIPS, fips)
		}
	}

	log.Info().
		Int("total_fips", len(allFIPS)).
		Int("recently_processed", len(recentlyProcessedFIPS)).
		Int("to_process", len(unprocessedFIPS)).
		Msg("FIPS processing status")

	return unprocessedFIPS, nil
}

// updateFIPSStatus updates the processing status of a FIPS code
func updateFIPSStatus(ctx context.Context, osClient *opensearch.Client, status *FIPSProcessingStatus) error {
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return err
	}

	indexReq := opensearchapi.IndexRequest{
		Index:      "fips_processing_status",
		DocumentID: status.FIPS,
		Body:       strings.NewReader(string(statusJSON)),
	}

	resp, err := indexReq.Do(ctx, osClient)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to update FIPS status: %s", resp.String())
	}

	return nil
}

// mustMarshal is a helper function that marshals data to JSON and panics on error
func mustMarshal(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return string(bytes)
}
