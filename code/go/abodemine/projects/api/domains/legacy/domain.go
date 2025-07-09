package legacy

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/ptr"
	"abodemine/lib/val"
	"abodemine/projects/api/conf"
	"abodemine/projects/api/domains/auth"
	"abodemine/projects/api/domains/legacy/models"
)

type Domain interface {
	Search(r *arc.Request, request []models.PropertySearchRequests) ([]any, error)
}

type domain struct {
	config *conf.Config

	authDomain auth.Domain
}

type NewDomainInput struct {
	Config *conf.Config

	AuthDomain auth.Domain
}

func NewDomain(in *NewDomainInput) *domain {
	return &domain{
		config:     in.Config,
		authDomain: in.AuthDomain,
	}
}

func (s *domain) Search(r *arc.Request, request []models.PropertySearchRequests) ([]any, error) {
	interimProps := request
	var doctypesearch []map[string]any

	client := wireup.BuildUSStreetAPIClient(
		wireup.SecretKeyCredential(ss_authid, ss_authtoken),
	)

	batch := street.NewBatch()
	parseCount := s.prepareBatchForAddressParsing(request, batch)

	if parseCount > 0 {
		var err error
		interimProps, err = s.ParseCount(client, batch, interimProps, request)
		if err != nil {
			return nil, err
		}
	}

	if len(interimProps) == 0 {
		return nil, errors.New("bad request")
	}

	return s.processProperties(r, interimProps, request, doctypesearch)
}

func (s *domain) prepareBatchForAddressParsing(request []models.PropertySearchRequests, batch *street.Batch) int {
	parseCount := 0

	for _, prop := range request {
		if !CheckSearchHasRequiredField(prop) {
			continue
		}

		if !NeedToCallAddressParser(prop) {
			continue
		}

		inputID := prop.InputID
		if inputID == "" {
			inputID = strings.TrimSpace(prop.FullAddress)
		}

		lookup := &street.Lookup{
			InputID:       inputID,
			Street:        strings.ToUpper(strings.TrimSpace(prop.FullAddress)),
			MaxCandidates: 1,
		}

		batch.Append(lookup)
		parseCount++
	}

	return parseCount
}

func (s *domain) processProperties(r *arc.Request, interimProps []models.PropertySearchRequests,
	request []models.PropertySearchRequests, doctypesearch []map[string]any) ([]any, error) {

	var responseProps []any

	if err := s.enrichPropertiesData(r.Context(), interimProps); err != nil {
		return nil, err
	}

	for i, obj := range interimProps {
		obj = s.GenerateOldAupidForProps(obj)

		doctypesearch = s.prepareDocTypeSearch(obj, doctypesearch)

		includedLayouts := filterLayouts(models.Address{
			IncludeLayouts: obj.IncludeLayouts,
			OmitLayouts:    obj.OmitLayouts,
		})

		needToFetch := obj.Aupid == "" && obj.OldAupid != ""
		if needToFetch {
			query := buildElasticsearchQuery(obj)
			result, err := s.MakeElasticsearchRequest(r.Context(), query, os_recorderindex5, MakeHeaders())
			if err != nil {
				log.Error().Err(&errors.Object{
					Id:     "13ed66fd-e228-430f-bca1-b5abf3f3f3c0",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to query elasticsearch for property data.",
					Cause:  err.Error(),
					Meta:   map[string]any{"aupid": obj.OldAupid},
				})
				continue
			}

			source := getScore(result)
			obj.Aupid = getStringOrEmpty(source, "aupid")
			obj.Zip5 = getStringOrEmpty(source, "SitusZIP5")
			obj.StreetDirection = getStringOrEmpty(source, "SitusDirectionLeft")
			obj.StreetName = getStringOrEmpty(source, "SitusStreet")
			obj.StreetSuffix = getStringOrEmpty(source, "SitusMode")
			obj.StreetPostDirection = getStringOrEmpty(source, "SitusDirectionRight")
			obj.HouseNumber = getStringOrEmpty(source, "SitusHouseNbr")
			obj.City = getStringOrEmpty(source, "SitusCity")
			obj.State = getStringOrEmpty(source, "SitusState")

		}

		es_requests, err := s.ElasticSearch(r.Context(), includedLayouts, obj, doctypesearch)
		if err != nil {
			return nil, &errors.Object{
				Id:     "a4d95e32-eaf3-4fe2-a8da-b8e5b064bdd9",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query elasticsearch for property data.",
				Cause:  err.Error(),
			}
		}

		returnProp := s.buildPropertyResponse(obj)

		includedLayouts = s.updateIncludedLayouts(obj, request, i, includedLayouts)

		if err := s.processEnabledLayouts(r, includedLayouts, es_requests, obj, request, &returnProp); err != nil {
			return nil, err
		}

		responseProps = append(responseProps, returnProp)
	}

	return responseProps, nil
}

func (s *domain) enrichPropertiesData(ctx context.Context, props []models.PropertySearchRequests) error {
	for i, obj := range props {
		if obj.Aupid == "" {
			continue
		}

		queryAupid := buildElasticsearchQueryAupid(obj.Aupid)
		result, err := s.MakeElasticsearchRequest(ctx, queryAupid, os_recorderindex5, MakeHeaders())
		if err != nil {
			return &errors.Object{
				Id:     "14769607-e3be-4f46-8235-3331351c4394",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query elasticsearch for property data",
				Cause:  err.Error(),
				Meta:   map[string]any{"aupid": obj.Aupid},
			}
		}

		source := getScore(result)
		props[i].Zip5 = getStringOrEmpty(source, "SitusZIP5")
		props[i].StreetDirection = getStringOrEmpty(source, "SitusDirectionLeft")
		props[i].StreetName = getStringOrEmpty(source, "SitusStreet")
		props[i].StreetSuffix = getStringOrEmpty(source, "SitusMode")
		props[i].StreetPostDirection = getStringOrEmpty(source, "SitusDirectionRight")
		props[i].HouseNumber = getStringOrEmpty(source, "SitusHouseNbr")
		props[i].City = getStringOrEmpty(source, "SitusCity")
		props[i].State = getStringOrEmpty(source, "SitusState")
	}

	return nil
}

func (s *domain) prepareDocTypeSearch(obj models.PropertySearchRequests, doctypesearch []map[string]any) []map[string]any {
	doctypesearch = append(doctypesearch, map[string]any{
		"match_phrase": map[string]any{
			"aupid": map[string]any{
				"query": obj.OldAupid,
			},
		},
	})

	if len(obj.RecorderDocTypes) > 0 {
		var dtypes []map[string]any
		for _, element := range obj.RecorderDocTypes {
			matchPhrase := map[string]any{
				"match_phrase": map[string]any{
					"documenttypecode": map[string]any{
						"query": element,
					},
				},
			}
			dtypes = append(dtypes, matchPhrase)
		}

		doctypesearch = append(doctypesearch, map[string]any{
			"bool": map[string]any{
				"should": dtypes,
			},
		})
	}

	return doctypesearch
}

func (s *domain) buildPropertyResponse(obj models.PropertySearchRequests) models.PropertySearchResponse {
	return models.PropertySearchResponse{
		HouseNumber:         obj.HouseNumber,
		StreetDirection:     obj.StreetDirection,
		StreetName:          obj.StreetName,
		StreetSuffix:        obj.StreetSuffix,
		StreetPostDirection: obj.StreetPostDirection,
		Zip5:                obj.Zip5,
		City:                obj.City,
		State:               obj.State,
		Aupid:               obj.Aupid,
	}
}

func (s *domain) updateIncludedLayouts(obj models.PropertySearchRequests, request []models.PropertySearchRequests,
	index int, includedLayouts []string) []string {

	if len(obj.IncludeLayouts) > 0 {
		includedLayouts = obj.IncludeLayouts
	} else if index < len(request) && request[index].IncludeLayouts != nil && len(request[index].IncludeLayouts) > 0 {
		includedLayouts = request[index].IncludeLayouts
	} else {
		for _, r := range request {
			if customLayouts, ok := r.Custom["_includelayouts"].([]string); ok && len(customLayouts) > 0 {
				includedLayouts = customLayouts
				break
			}
		}
	}

	var excludedLayouts []string
	if len(obj.OmitLayouts) > 0 {
		excludedLayouts = obj.OmitLayouts
	} else if index < len(request) && request[index].OmitLayouts != nil && len(request[index].OmitLayouts) > 0 {
		excludedLayouts = request[index].OmitLayouts
	} else {
		for _, r := range request {
			if customOmitLayouts, ok := r.Custom["_omitlayouts"].([]string); ok && len(customOmitLayouts) > 0 {
				excludedLayouts = customOmitLayouts
				break
			}
		}
	}

	if len(includedLayouts) == 0 {
		includedLayouts = defaultLayouts
	}

	if len(excludedLayouts) > 0 {
		filteredLayouts := make([]string, 0, len(includedLayouts))
		for _, item := range includedLayouts {
			if !slices.Contains(excludedLayouts, item) {
				filteredLayouts = append(filteredLayouts, item)
			}
		}
		includedLayouts = filteredLayouts
	}

	return includedLayouts
}

func (s *domain) processEnabledLayouts(r *arc.Request, includedLayouts []string, es_requests []map[string]any,
	obj models.PropertySearchRequests, request []models.PropertySearchRequests,
	returnProp *models.PropertySearchResponse) error {

	if slices.Contains(includedLayouts, "assessor") {
		if !r.HasFlag(flags.ApiAssessorLayoutEnabled) {
			return &errors.Object{
				Id:     "5e117393-03a2-49a1-9143-365423d48eea",
				Code:   errors.Code_PERMISSION_DENIED,
				Detail: "The assessor layout is not enabled for this organization.",
			}
		}

		if err := s.ProcessAssessorLayout(es_requests, obj, returnProp, request); err != nil {
			return &errors.Object{
				Id:     "28434294-cfe2-4de8-a2d9-0ba04a95ad45",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to process assessor layout.",
				Cause:  err.Error(),
			}
		}
	}

	var totalRecorderHits int64 = 0

	if slices.Contains(includedLayouts, "recorder") {
		if !r.HasFlag(flags.ApiRecorderLayoutEnabled) {
			return &errors.Object{
				Id:     "462b2721-2d99-4474-a140-0b76e862f084",
				Code:   errors.Code_PERMISSION_DENIED,
				Detail: "The recorder layout is not enabled for this organization.",
			}
		}

		var recorderHits int64 = 0
		var recorderHitsList []map[string]any
		if err := s.ProcessRecorderLayout(es_requests, obj, recorderHits, request, recorderHitsList, returnProp); err != nil {
			return &errors.Object{
				Id:     "1bc949a9-1a39-4595-b857-5d8484870fb2",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to process recorder layout.",
				Cause:  err.Error(),
			}
		}

		totalRecorderHits = returnProp.Hits
	}

	if slices.Contains(includedLayouts, "rentEstimate") {
		if !r.HasFlag(flags.ApiRentEstimateLayoutEnabled) {
			return &errors.Object{
				Id:     "5ccab1f5-b23a-41d5-889d-b419806c31ff",
				Code:   errors.Code_PERMISSION_DENIED,
				Detail: "The rentEstimate layout is not enabled for this organization.",
			}
		}

		if err := s.ProcessRentEstimateLayout(es_requests, obj, returnProp, request); err != nil {
			return &errors.Object{
				Id:     "29d1ed69-8fa8-46ab-8518-e68ca70bc6cf",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to process rentEstimate layout.",
				Cause:  err.Error(),
			}
		}
	}

	saleEstimateRequested := slices.Contains(includedLayouts, "saleEstimate")
	saleEstimateWithCompsRequested := slices.Contains(includedLayouts, "saleEstimateWithComps")

	if saleEstimateRequested || saleEstimateWithCompsRequested {
		if !r.HasFlag(flags.ApiSaleEstimateLayoutEnabled) {
			return &errors.Object{
				Id:     "97d0fadf-a8bb-44ba-8e00-f71aa9adb719",
				Code:   errors.Code_PERMISSION_DENIED,
				Detail: "The saleEstimate layout is not enabled for this organization.",
			}
		}

		if saleEstimateWithCompsRequested && !r.HasFlag(flags.ApiCompsLayoutEnabled) {
			return &errors.Object{
				Id:     "afb17ed4-9ce6-4f3e-a984-5bd0c455af97",
				Code:   errors.Code_PERMISSION_DENIED,
				Detail: "The comps layout is not enabled for this organization.",
			}
		}

		if err := s.ProcessSaleEstimateLayout(r.Context(), es_requests, includedLayouts, obj, request, returnProp); err != nil {
			return &errors.Object{
				Id:     "a3e400da-f405-42e1-b338-5ebd84e6cf8b",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to process saleEstimateWithComps layout.",
				Cause:  err.Error(),
			}
		}
	} else {
		returnProp.SaleEstimate = map[string]any{}
	}

	_, err := s.authDomain.InsertApiQuotaTransaction(r, &auth.InsertApiQuotaTransactionInput{
		Entity: &arc.ApiQuotaTransaction{
			Description:         ptr.String("/api/v2/search"),
			AddressLayoutAmount: 1,
			AssessorLayoutAmount: int32(val.Ternary(
				len(returnProp.Assessor) > 0,
				1,
				0,
			)),
			CompsLayoutAmount: int32(val.Ternary(
				// Simplified check because this endpoint will be deprecated soon.
				saleEstimateWithCompsRequested && len(returnProp.SaleEstimate) > 0,
				1,
				0,
			)),
			RecorderLayoutAmount: int32(val.Ternary(
				len(returnProp.Recorder) > 0,
				totalRecorderHits,
				0,
			)),
			RentEstimateLayoutAmount: int32(val.Ternary(
				len(returnProp.RentEstimate) > 0,
				1,
				0,
			)),
			SaleEstimateLayoutAmount: int32(val.Ternary(
				len(returnProp.SaleEstimate) > 0,
				1,
				0,
			)),
		},
	})
	if err != nil {
		return errors.Forward(err, "c805bef3-1d0b-4c6e-98b4-b3b6379ce4e7")
	}

	return nil
}

func (s *domain) GenerateOldAupidForProps(props models.PropertySearchRequests) models.PropertySearchRequests {
	if props.OldAupid == "" {
		streetParts := []string{props.StreetDirection, props.StreetName, props.StreetSuffix, props.StreetPostDirection}
		streetPart := strings.Join(strings.Fields(strings.Join(streetParts, " ")), " ")
		streetPart = strings.ToUpper(streetPart)

		props.OldAupid = strings.ToUpper(fmt.Sprintf("840-%s-%s-%s",
			props.Zip5, streetPart, props.HouseNumber))
	}

	return props
}

func (s *domain) ProcessAssessorLayout(es_requests []map[string]any, obj models.PropertySearchRequests,
	returnProp *models.PropertySearchResponse, request []models.PropertySearchRequests) error {

	assessorHit := s.findHitByPathAndAupid(es_requests,
		"/abodemine_assessor_latest_20240728_007/_search", obj.OldAupid)

	if assessorHit == nil {
		return nil
	}

	for key, value := range assessorHit {
		if value == nil {
			continue
		}

		if typeToConvert, exists := assessorRetype[key]; exists {
			strVal, isString := value.(string)
			if !isString {
				continue
			}

			trimmedValue := strings.TrimSpace(strVal)
			if trimmedValue == "" {
				continue
			}

			s.convertFieldType(assessorHit, key, trimmedValue, typeToConvert)
		}
	}

	s.removeExcludedFields(assessorHit, alwaysExclude)
	ProcessFields(request, []models.PropertySearchRequests{obj}, assessorHit)

	returnProp.Assessor = assessorHit
	return nil
}

func (s *domain) findHitByPathAndAupid(es_requests []map[string]any, path string, aupid string) map[string]any {
	for _, es_request := range es_requests {
		if es_request["url.path"] != path {
			continue
		}

		hits, ok := es_request["hits"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "c0435cf9-b55d-4dcc-8cbf-1afb2e5fd814",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format for property data",
				Meta:   map[string]any{"aupid": aupid},
			})
			continue
		}

		hitsArray, ok := hits["hits"].([]any)
		if !ok || len(hitsArray) == 0 {
			log.Error().Err(&errors.Object{
				Id:     "1dc5a35c-2a23-4f8b-90ab-de01535a61e7",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format for property data",
				Meta:   map[string]any{"aupid": aupid},
			})
			continue
		}

		// If the AUPID was not specified, return the first hit
		if aupid == "" {
			firstHit, ok := hitsArray[0].(map[string]any)
			if !ok {
				log.Error().Err(&errors.Object{
					Id:     "a65584f6-61bf-44b6-ac74-75a9af4dac5c",
					Code:   errors.Code_UNKNOWN,
					Detail: "Invalid hits array format for property data",
				})
				continue
			}

			source, ok := firstHit["_source"].(map[string]any)
			if !ok {
				log.Error().Err(&errors.Object{
					Id:     "ffcacfbc-da2b-427d-a76b-2f8eedf59a6e",
					Code:   errors.Code_UNKNOWN,
					Detail: "Invalid _source format for property data",
					Meta:   map[string]any{"aupid": aupid},
				})
				continue
			}

			return source
		}

		// Otherwise, search for the specific AUPID
		for _, hit := range hitsArray {
			hitMap, ok := hit.(map[string]any)
			if !ok {
				log.Error().Err(&errors.Object{
					Id:     "9c017ea2-ed23-47b5-8599-3bfb6876bcf3",
					Code:   errors.Code_UNKNOWN,
					Detail: "Invalid hits format for property data",
					Meta:   map[string]any{"aupid": aupid},
				})
				continue
			}

			source, ok := hitMap["_source"].(map[string]any)
			if !ok {
				log.Error().Err(&errors.Object{
					Id:     "263513a2-9ef0-4132-bb9f-d7ef37db2334",
					Code:   errors.Code_UNKNOWN,
					Detail: "Invalid _source format for property data",
					Meta:   map[string]any{"aupid": aupid},
				})
				continue
			}

			if source["aupid"] == aupid {
				return source
			}
		}
	}

	return nil
}

func (s *domain) convertFieldType(data map[string]any, key, value, typeName string) {
	switch typeName {
	case "float":
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			data[key] = floatVal
		}
	case "integer":
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			data[key] = intVal
		}
	case "string":
		data[key] = value
	}
}

func (s *domain) removeExcludedFields(data map[string]any, excludeList []string) {
	for key := range data {
		if slices.Contains(excludeList, key) || data[key] == nil {
			delete(data, key)
		}
	}
}

func (s *domain) ProcessRecorderLayout(es_requests []map[string]any, obj models.PropertySearchRequests,
	recorderHits int64, request []models.PropertySearchRequests, recorderHitsList []map[string]any,
	returnProp *models.PropertySearchResponse) error {

	var recorderHitz []map[string]any

	for _, es_request := range es_requests {
		if es_request["url.path"] != "/recorder_vendor_b_1_0/_search" {
			continue
		}

		hits, ok := es_request["hits"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "bc9d4790-5fd9-41e6-9c29-83b2e1994fdc",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format for recorder data",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		hitsArray, ok := hits["hits"].([]any)
		if !ok || len(hitsArray) == 0 {
			log.Error().Err(&errors.Object{
				Id:     "7b696611-77e7-4190-81e8-20df0a498972",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format for recorder data",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		if total, ok := hits["total"].(map[string]any); ok {
			if value, ok := total["value"].(float64); ok {
				recorderHits = int64(value)
			}
		}

		for _, hit := range hitsArray {
			if hitMap, ok := hit.(map[string]any); ok {
				if hitSource, ok := hitMap["_source"].(map[string]any); ok {
					recorderHitz = append(recorderHitz, hitSource)
				}
			}
		}
	}

	for _, recorderHit := range recorderHitz {
		s.removeExcludedFields(recorderHit, alwaysExclude)

		for key, value := range recorderHit {
			if typeToConvert, exists := recorderRetype[key]; exists {
				strVal, isString := value.(string)
				if !isString || value == nil {
					continue
				}

				trimmedValue := strings.TrimSpace(strVal)
				if trimmedValue == "" {
					continue
				}

				s.convertFieldType(recorderHit, key, trimmedValue, typeToConvert)
			}
		}

		s.filterFieldsByInclusionExclusion(recorderHit, request)

		recorderHitsList = append(recorderHitsList, recorderHit)
	}

	returnProp.Hits = recorderHits
	returnProp.Recorder = [][]map[string]any{recorderHitsList}

	return nil
}

func (s *domain) filterFieldsByInclusionExclusion(data map[string]any, request []models.PropertySearchRequests) {

	var includedFields, excludedFields []string

	for _, req := range request {
		if len(req.IncludeFields) > 0 && len(includedFields) == 0 {
			includedFields = req.IncludeFields
		}
		if len(req.OmitFields) > 0 && len(excludedFields) == 0 {
			excludedFields = req.OmitFields
		}
	}

	if aupid, exists := data["aupid"].(string); exists {
		for _, prop := range request {
			if prop.OldAupid == aupid {
				if len(includedFields) == 0 && len(prop.IncludeFields) > 0 {
					includedFields = prop.IncludeFields
				}
				if len(excludedFields) == 0 && len(prop.OmitFields) > 0 {
					excludedFields = prop.OmitFields
				}
				break
			}
		}
	}

	if len(includedFields) > 0 {
		for key := range data {
			if !slices.Contains(includedFields, key) {
				delete(data, key)
			}
		}
	} else if len(excludedFields) > 0 {
		for key := range data {
			if slices.Contains(excludedFields, key) {
				delete(data, key)
			}
		}
	}
}

func (s *domain) ProcessRentEstimateLayout(es_requests []map[string]any, obj models.PropertySearchRequests,
	returnProp *models.PropertySearchResponse, request []models.PropertySearchRequests) error {

	rentEstimateHit := s.findHitByPathAndAupid(es_requests, "/at_rentavm_2_0/_search", obj.OldAupid)

	if rentEstimateHit == nil {
		return nil
	}

	s.removeExcludedFields(rentEstimateHit, alwaysExclude)

	for key, value := range rentEstimateHit {
		if typeToConvert, exists := rentEstimateRetype[key]; exists {
			strVal, isString := value.(string)
			if !isString || value == nil {
				continue
			}

			trimmedValue := strings.TrimSpace(strVal)
			if trimmedValue == "" {
				continue
			}

			s.convertFieldType(rentEstimateHit, key, trimmedValue, typeToConvert)

			if strValue, ok := rentEstimateHit[key].(string); ok {
				if key == "propertyaddresszip" {
					rentEstimateHit[key] = fmt.Sprintf("%05s", strValue)
				} else if key == "propertyaddresszip4" {
					rentEstimateHit[key] = fmt.Sprintf("%04s", strValue)
				}
			}
		}
	}

	s.filterFieldsByInclusionExclusion(rentEstimateHit, request)

	returnProp.RentEstimate = rentEstimateHit
	return nil
}

func (s *domain) ProcessSaleEstimateLayout(ctx context.Context, es_requests []map[string]any, includedLayouts []string,
	obj models.PropertySearchRequests, request []models.PropertySearchRequests,
	returnProp *models.PropertySearchResponse) error {

	var saleEstimateHit map[string]any

	for _, es_request := range es_requests {
		urlPath, hasPath := es_request["url.path"]
		if !hasPath || urlPath != "/fa_avm_20250217/_search" {
			continue
		}

		hits, ok := es_request["hits"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "a776cc98-24b1-4cae-bec1-c14952ba0d4a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format in saleEstimate request.",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		hitsArray, ok := hits["hits"].([]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "299fcdbc-a179-40a0-89a0-d6b1a4eb3742",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits array format in saleEstimate request.",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		if len(hitsArray) == 0 {
			log.Error().Err(&errors.Object{
				Id:     "18194b1a-2450-491f-a800-8cd3ba1b4564",
				Code:   errors.Code_UNKNOWN,
				Detail: "No hits found for saleEstimate.",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		firstHit, ok := hitsArray[0].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "0d34b6d6-58d3-4292-a309-21b45575370c",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid first hit format for saleEstimate.",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		source, ok := firstHit["_source"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "3336b67b-7fcc-4ac5-9953-bd0861c06640",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid _source format for saleEstimate.",
				Meta:   map[string]any{"aupid": obj.OldAupid},
			})
			continue
		}

		saleEstimateHit = source
		break
	}

	if saleEstimateHit == nil {
		log.Error().Err(&errors.Object{
			Id:     "2d533a6c-7466-4694-9574-5e135412bef1",
			Code:   errors.Code_NOT_FOUND,
			Detail: "No saleEstimate found for property.",
			Meta:   map[string]any{"aupid": obj.OldAupid},
		})
		return nil
	}

	saleEstimateHitLc := make(map[string]any, len(saleEstimateHit))
	for key, value := range saleEstimateHit {
		if value != nil && strings.TrimSpace(fmt.Sprintf("%v", value)) != "" {
			lowerKey := strings.ToLower(key)
			saleEstimateHitLc[lowerKey] = value
		}
	}

	saleEstimateHit = saleEstimateHitLc

	if slices.Contains(includedLayouts, "saleEstimateWithComps") {
		if err := s.processSaleEstimateComps(ctx, saleEstimateHit); err != nil {
			return &errors.Object{
				Id:     "577b5ce5-d09b-4b42-b1f0-b1254249aed4",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to process saleEstimate comps.",
				Cause:  err.Error(),
			}
		}
	}

	for key := range saleEstimateHit {
		if slices.Contains(alwaysExclude, key) {
			delete(saleEstimateHit, key)
		}
	}

	for key, value := range saleEstimateHit {
		if typeStr, exists := saleEstimateRetype[key]; exists {
			if strVal, ok := value.(string); ok && value != nil && strings.TrimSpace(strVal) != "" {
				s.convertFieldType(saleEstimateHit, key, strVal, typeStr)
			}

			if strVal, ok := saleEstimateHit[key].(string); ok {
				switch key {
				case "situszip5":
					saleEstimateHit[key] = fmt.Sprintf("%05s", strVal)
				case "situszip4":
					saleEstimateHit[key] = fmt.Sprintf("%04s", strVal)
				}
			}
		}
	}

	returnProp.SaleEstimate = saleEstimateHit
	return nil
}

func (s *domain) processSaleEstimateComps(ctx context.Context, saleEstimateHit map[string]any) error {
	compPropertyIds := []any{
		saleEstimateHit["comp1propertyid"],
		saleEstimateHit["comp2propertyid"],
		saleEstimateHit["comp3propertyid"],
		saleEstimateHit["comp4propertyid"],
		saleEstimateHit["comp5propertyid"],
		saleEstimateHit["comp6propertyid"],
		saleEstimateHit["comp7propertyid"],
	}

	compNum := 1
	headers := MakeHeaders()

	for idx, compId := range compPropertyIds {
		if compId == nil {
			continue
		}

		strCompId := fmt.Sprintf("%v", compId)
		originalCompNum := idx + 1

		queryComp := faPropertyIdLookup(strCompId)
		result, err := s.MakeElasticsearchRequest(ctx, queryComp, os_recorderindex6, headers)
		if err != nil {
			log.Error().Err(&errors.Object{
				Id:     "6e398189-1d37-4a7f-8d79-dc997e316dca",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query elasticsearch for comp data.",
				Cause:  err.Error(),
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		hits, ok := result["hits"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "e61d5e4e-dc2d-419e-bdc3-c488059ca8e3a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid hits format for comp data.",
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		hitsArray, ok := hits["hits"].([]any)
		if !ok || len(hitsArray) == 0 {
			log.Error().Err(&errors.Object{
				Id:     "5f226933-94db-420b-bfeb-0febc50c751d",
				Code:   errors.Code_UNKNOWN,
				Detail: "No hits found for comp data.",
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		firstHit, ok := hitsArray[0].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "dcb24299-a257-4e6e-84d2-17ebd550ef5c",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid first hit format for comp data.",
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		source, ok := firstHit["_source"].(map[string]any)
		if !ok {
			log.Error().Err(&errors.Object{
				Id:     "e492a98b-925f-4489-990a-5ccb76d1a6f1",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid _source format for comp data.",
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		compAupid, ok := source["id"].(string)
		if !ok || strings.TrimSpace(compAupid) == "" {
			log.Error().Err(&errors.Object{
				Id:     "2f8cca0d-2a42-4b84-aea6-38fa65121c0da",
				Code:   errors.Code_UNKNOWN,
				Detail: "Invalid comp AUPID format.",
				Meta:   map[string]any{"comp_id": strCompId},
			})
			continue
		}

		valueavmcomp := map[string]any{"aupid": strCompId}

		var wg sync.WaitGroup
		var compAssr, compRecorder map[string]any
		var mu sync.Mutex

		wg.Add(2)

		go func() {
			defer wg.Done()
			assrData := s.makeAssr(ctx, compAupid, headers)
			if assrData != nil {
				mu.Lock()
				compAssr = assrData
				mu.Unlock()
			}
		}()

		go func() {
			defer wg.Done()
			recorderData := s.makeCompRecorder(ctx, compAupid, headers)
			if recorderData != nil {
				mu.Lock()
				compRecorder = recorderData
				mu.Unlock()
			}
		}()

		wg.Wait()

		updateCompValues(source, compAssr, compRecorder, valueavmcomp)

		key := fmt.Sprintf("valueavmcomp%d", originalCompNum)
		saleEstimateHit[key] = valueavmcomp

		compNum++
	}

	return nil
}

func (s *domain) makeCompRecorder(ctx context.Context, compAupid string, headers map[string]any) map[string]any {
	es_queryRecorder := map[string]any{
		"from": 0,
		"size": 1,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"match_phrase": map[string]any{
							"aupid": map[string]any{
								"query": compAupid,
							},
						},
					},
				},
				"filter":   []map[string]any{{"match_all": map[string]any{}}},
				"should":   []any{},
				"must_not": []any{},
			},
		},
		"sort": []map[string]any{
			{
				"recordingdate.keyword": map[string]any{"order": "desc"},
			},
			{
				"_score": map[string]any{"order": "desc"},
			},
		},
	}

	result, err := s.MakeElasticsearchRequest(ctx, es_queryRecorder, os_recorder_index, headers)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", "b4c5d6e7-f8a9-4b4c-5d6e-7f8a9b0c1d2e").
			Str("comp_aupid", compAupid).
			Str("index", os_recorder_index).
			Msg("Failed to execute elasticsearch request for comp recorder.")
		return nil
	}

	hits, ok := result["hits"].(map[string]any)
	if !ok {
		return nil
	}

	hitsArray, ok := hits["hits"].([]any)
	if !ok || len(hitsArray) == 0 {
		return nil
	}

	firstHit, ok := hitsArray[0].(map[string]any)
	if !ok {
		return nil
	}

	source, ok := firstHit["_source"].(map[string]any)
	if !ok {
		return nil
	}

	return source
}

func (s *domain) makeAssr(ctx context.Context, compAupid string, headers map[string]any) map[string]any {
	es_queryAssessor := map[string]any{
		"from": 0,
		"size": 1,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"match_phrase": map[string]any{
							"aupid": map[string]any{
								"query": compAupid,
							},
						},
					},
				},
				"filter":   []map[string]any{{"match_all": map[string]any{}}},
				"should":   []any{},
				"must_not": []any{},
			},
		},
		"sort": []map[string]any{
			{
				"assrlastupdated.keyword": map[string]any{"order": "desc"},
			},
			{
				"publicationdate.keyword": map[string]any{"order": "desc"},
			},
			{
				"taxyearassessed.keyword": map[string]any{"order": "desc"},
			},
			{
				"taxfiscalyear.keyword": map[string]any{"order": "desc"},
			},
			{
				"_score": map[string]any{"order": "desc"},
			},
		},
	}

	result, err := s.MakeElasticsearchRequest(ctx, es_queryAssessor, os_assessorindex, headers)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", "c5d6e7f8-a9b0-4c5d-6e7f-8a9b0c1d2e3f").
			Str("comp_aupid", compAupid).
			Str("index", os_assessorindex).
			Msg("Failed to execute elasticsearch request for comp assessor.")
		return nil
	}

	hits, ok := result["hits"].(map[string]any)
	if !ok {
		return nil
	}

	hitsArray, ok := hits["hits"].([]any)
	if !ok || len(hitsArray) == 0 {
		return nil
	}

	firstHit, ok := hitsArray[0].(map[string]any)
	if !ok {
		return nil
	}

	source, ok := firstHit["_source"].(map[string]any)
	if !ok {
		return nil
	}

	return source
}

func (s *domain) ElasticSearch(ctx context.Context, includedLayouts []string, address models.PropertySearchRequests,
	doctypesearch []map[string]any) ([]map[string]any, error) {

	var es_requests []map[string]any
	headers := MakeHeaders()

	expectedCapacity := 0
	if slices.Contains(includedLayouts, "recorder") {
		expectedCapacity++
	}
	if slices.Contains(includedLayouts, "rentEstimate") {
		expectedCapacity++
	}
	if slices.Contains(includedLayouts, "saleEstimate") || slices.Contains(includedLayouts, "saleEstimateWithComps") {
		expectedCapacity++
	}
	if slices.Contains(includedLayouts, "assessor") {
		expectedCapacity++
	}

	es_requests = make([]map[string]any, 0, expectedCapacity)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastError error

	if slices.Contains(includedLayouts, "recorder") {
		wg.Add(1)
		go func() {
			defer wg.Done()

			esQueryRecorder := map[string]any{
				"from": 0,
				"size": 10,
				"query": map[string]any{
					"bool": map[string]any{
						"must": doctypesearch,
					},
				},
				"sort": []map[string]any{
					{
						"recordingdate.keyword": map[string]any{"order": "desc"},
					},
				},
			}

			resp, err := s.MakeElasticsearchRequest(ctx, esQueryRecorder, os_recorder_index, headers)
			if err != nil {
				mu.Lock()
				lastError = err
				mu.Unlock()
				return
			}

			resp["url.path"] = "/recorder_vendor_b_1_0/_search"

			mu.Lock()
			es_requests = append(es_requests, resp)
			mu.Unlock()
		}()
	}

	if slices.Contains(includedLayouts, "rentEstimate") {
		wg.Add(1)
		go func() {
			defer wg.Done()

			es_queryRentEstimate := map[string]any{
				"from": 0,
				"size": 1,
				"query": map[string]any{
					"bool": map[string]any{
						"must": []map[string]any{
							{
								"match_phrase": map[string]any{
									"aupid": map[string]any{
										"query": address.OldAupid,
									},
								},
							},
						},
					},
				},
				"sort": []map[string]any{
					{
						"_score": map[string]any{"order": "desc"},
					},
					{
						"valuationdate": map[string]any{"order": "desc"},
					},
					{
						"publicationdate": map[string]any{"order": "desc"},
					},
				},
			}

			resp, err := s.MakeElasticsearchRequest(ctx, es_queryRentEstimate, os_rent_estimate_index, headers)
			if err != nil {
				mu.Lock()
				lastError = err
				mu.Unlock()
				return
			}

			resp["url.path"] = "/at_rentavm_2_0/_search"

			mu.Lock()
			es_requests = append(es_requests, resp)
			mu.Unlock()
		}()
	}

	if slices.Contains(includedLayouts, "saleEstimate") || slices.Contains(includedLayouts, "saleEstimateWithComps") {
		wg.Add(1)
		go func() {
			defer wg.Done()

			es_querySaleEstimate := map[string]any{
				"from": 0,
				"size": 10,
				"query": map[string]any{
					"bool": map[string]any{
						"must": []map[string]any{
							{
								"match_phrase": map[string]any{
									"aupid": map[string]any{
										"query": address.Aupid,
									},
								},
							},
						},
					},
				},
				"sort": []map[string]any{
					{
						"_score": map[string]any{"order": "desc"},
					},
					{
						"ValuationDate": map[string]any{"order": "desc"},
					},
					{
						"SitusUnitNbr.keyword": map[string]any{"order": "desc"},
					},
				},
			}

			resp, err := s.MakeElasticsearchRequest(ctx, es_querySaleEstimate, os_recorderindex5, headers)
			if err != nil {
				mu.Lock()
				lastError = err
				mu.Unlock()
				return
			}

			resp["url.path"] = "/fa_avm_20250217/_search"

			mu.Lock()
			es_requests = append(es_requests, resp)
			mu.Unlock()
		}()
	}

	if slices.Contains(includedLayouts, "assessor") {
		wg.Add(1)
		go func() {
			defer wg.Done()

			es_queryAssessor := map[string]any{
				"from": 0,
				"size": 1,
				"query": map[string]any{
					"bool": map[string]any{
						"must": []map[string]any{
							{
								"match_phrase": map[string]any{
									"aupid": map[string]any{
										"query": address.OldAupid,
									},
								},
							},
						},
					},
				},
				"sort": []map[string]any{
					{
						"assrlastupdated.keyword": map[string]any{"order": "desc"},
					},
					{
						"publicationdate.keyword": map[string]any{"order": "desc"},
					},
					{
						"taxyearassessed.keyword": map[string]any{"order": "desc"},
					},
					{
						"_score": map[string]any{"order": "desc"},
					},
				},
			}

			resp, err := s.MakeElasticsearchRequest(ctx, es_queryAssessor, os_assessorindex, headers)
			if err != nil {
				mu.Lock()
				lastError = err
				mu.Unlock()
				return
			}

			resp["url.path"] = "/abodemine_assessor_latest_20240728_007/_search"

			mu.Lock()
			es_requests = append(es_requests, resp)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if lastError != nil {
		return nil, lastError
	}

	return es_requests, nil
}

func (s *domain) ParseCount(client *street.Client, batch *street.Batch, interimProps []models.PropertySearchRequests,
	request []models.PropertySearchRequests) ([]models.PropertySearchRequests, error) {

	err := client.SendBatch(batch)
	if err != nil {
		return nil, &errors.Object{
			Id:     "617baf1d-f2cb-4888-8f60-2bd8fd55faf4",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to send batch to Street.",
			Cause:  err.Error(),
		}
	}

	for _, res := range batch.Records() {
		var indexToReplace int = -1

		for i, prop := range interimProps {
			if prop.FullAddress == res.InputID {
				indexToReplace = i
				break
			}
		}

		if len(res.Results) == 0 {
			return nil, &errors.Object{
				Id:     "79493e50-9067-4c33-9b66-09ca80f77c01",
				Code:   errors.Code_NOT_FOUND,
				Detail: fmt.Sprintf("No results found for address: %s", res.InputID),
			}
		}

		newObj := s.createUpdatedProperty(res)

		if indexToReplace != -1 {
			newObj.IncludeFields = interimProps[indexToReplace].IncludeFields
			newObj.OmitFields = interimProps[indexToReplace].OmitFields
			newObj.IncludeLayouts = interimProps[indexToReplace].IncludeLayouts

			interimProps[indexToReplace] = newObj
		}
	}

	return interimProps, nil
}

func (s *domain) createUpdatedProperty(res *street.Lookup) models.PropertySearchRequests {
	newObj := models.PropertySearchRequests{}
	newObj.HouseNumber = res.Results[0].Components.PrimaryNumber
	newObj.Unit = res.Results[0].Components.SecondaryNumber
	newObj.StreetDirection = res.Results[0].Components.StreetPredirection
	newObj.StreetName = res.Results[0].Components.StreetName
	newObj.StreetSuffix = res.Results[0].Components.StreetSuffix
	newObj.StreetPostDirection = res.Results[0].Components.StreetPostdirection
	newObj.Zip5 = res.Results[0].Components.ZIPCode
	newObj.City = res.Results[0].Components.CityName
	newObj.State = res.Results[0].Components.StateAbbreviation

	streetPart := strings.TrimSpace(fmt.Sprintf("%s %s %s %s",
		newObj.StreetDirection, newObj.StreetName, newObj.StreetSuffix, newObj.StreetPostDirection))
	streetPart = strings.Join(strings.Fields(streetPart), " ")
	streetPart = strings.ToUpper(streetPart)

	newObj.OldAupid = strings.ToUpper(fmt.Sprintf("840-%s-%s-%s", newObj.Zip5, streetPart, newObj.HouseNumber))

	if newObj.Unit != "" && newObj.Unit != "undefined" && strings.TrimSpace(newObj.Unit) != "" {
		newObj.OldAupid = newObj.OldAupid + "-" + strings.TrimSpace(newObj.Unit)
	}

	return newObj
}
