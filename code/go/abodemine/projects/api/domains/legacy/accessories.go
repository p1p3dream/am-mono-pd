package legacy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/projects/api/domains/legacy/models"
)

var (
	os_recorderindex5      = os.Getenv("ELASTICSEARCH_INDEX_AVM")
	os_recorderindex6      = os.Getenv("ELASTICSEARCH_INDEX_ADDRESS")
	os_assessorindex       = os.Getenv("ELASTICSEARCH_INDEX_ASSESSOR")
	os_rent_estimate_index = os.Getenv("ELASTICSEARCH_INDEX_RENTESTIMATE")
	os_recorder_index      = os.Getenv("ELASTICSEARCH_INDEX_RECORDER")

	//elasticSearch
	ss_authid    = os.Getenv("SS_AUTHID")
	ss_authtoken = os.Getenv("SS_AUTHTOKEN")
)

var (
	ErrElasticsearchRequest = errors.New("elasticsearch request failed")
	ErrInvalidResponse      = errors.New("invalid elasticsearch response")
)

var (
	defaultLayouts = []string{"recorder", "assessor"}

	alwaysExclude = []string{
		"log",
		"type",
		"host",
		"@version",
		"event",
		"@timestamp",
		"Index",
		"bucket",
		"path",
		"filename",
		"_1",
		"id",
		"latitude",
		"longitude",
		"location",
		"address_details",
		"comp1propertyid",
		"comp2propertyid",
		"comp3propertyid",
		"comp4propertyid",
		"comp5propertyid",
		"comp6propertyid",
		"comp7propertyid",
		"propertyid",
		"index",
	}

	fieldsToReturnInRoot = []string{"aupid", "recorder", "assessor", "rentEstimate", "saleEstimate", "hits", "a"}

	searchFields = []string{
		"aupid",
		"housenumber",
		"streetdirection",
		"streetname",
		"streetsuffix",
		"streetpostdirection",
		"zip5",
		"fulladdress",
	}

	addressFields = map[string]string{
		"FullStreetAddress": "propertyaddressfull",
		"StreetNumber":      "propertyaddresshousenumber",
		"PreDirectional":    "propertyaddressstreetdirection",
		"Street":            "propertyaddressstreetname",
		"StreetType":        "propertyaddressstreetsuffix",
		"PostDirectional":   "propertyaddressstreetpostdirection",
		"UnitType":          "propertyaddressunitprefix",
		"UnitNbr":           "propertyaddressunitvalue",
		"PostalCommunity":   "propertyaddresscity",
		"State":             "propertyaddressstate",
		"ZIP5":              "propertyaddresszip",
		"ZIP4":              "propertyaddresszip4",
	}

	assrFields = map[string]string{
		"bathcount":          "bathcount",
		"bathpartialcount":   "bathpartialcount",
		"bedroomscount":      "bedroomscount",
		"areagross":          "areagross",
		"yearbuilt":          "yearbuilt",
		"yearbuilteffective": "yearbuilteffective",
	}

	recorderFields = map[string]string{
		"recordingdate":  "recordingdate",
		"transferamount": "transferamount",
	}

	saleEstimateRetype = map[string]string{
		"aupid":                  "string",
		"situsfullstreetaddress": "string",
		"situshousenbr":          "string",
		"situshousenbrsuffix":    "string",
		"situsdirectionleft":     "string",
		"situsstreet":            "string",
		"situsmode":              "string",
		"situsdirectionright":    "string",
		"situsunittype":          "string",
		"situsunitnbr":           "string",
		"situscity":              "string",
		"situsstate":             "string",
		"situszip5":              "string",
		"situszip4":              "string",
		"situscarriercode":       "string",
		"finalvalue":             "integer",
		"highvalue":              "integer",
		"lowvalue":               "integer",
		"confidencescore":        "integer",
		"standarddeviation":      "integer",
		"valuationdate":          "string",
	}

	rentEstimateRetype = map[string]string{
		"aupid":                              "string",
		"propertyaddressfull":                "string",
		"propertyaddresshousenumber":         "string",
		"propertyaddressstreetdirection":     "string",
		"propertyaddressstreetname":          "string",
		"propertyaddressstreetsuffix":        "string",
		"propertyaddressstreetpostdirection": "string",
		"propertyaddressunitprefix":          "string",
		"propertyaddressunitvalue":           "string",
		"propertyaddresscity":                "string",
		"propertyaddressstate":               "string",
		"propertyaddresszip":                 "string",
		"propertyaddresszip4":                "string",
		"propertyaddresscrrt":                "string",
		"propertyusegroup":                   "string",
		"propertyusestandardized":            "string",
		"estimatedrentalvalue":               "integer",
		"estimatedminrentalvalue":            "integer",
		"estimatedmaxrentalvalue":            "integer",
		"valuationdate":                      "string",
		"publicationdate":                    "string",
	}

	recorderRetype = map[string]string{
		"transactionid":                                     "integer",
		"instrumentdate":                                    "float",
		"propertyaddressinfoprivacy":                        "integer",
		"recordingdate":                                     "string",
		"transferinfopurchasetypecode":                      "integer",
		"foreclosureauctionsale":                            "integer",
		"transferinfodistresscircumstancecode":              "integer",
		"quitclaimflag":                                     "integer",
		"transferinfomultiparcelflag":                       "integer",
		"armslengthflag":                                    "integer",
		"transferamount":                                    "float",
		"transfertaxtotal":                                  "float",
		"transfertaxcity":                                   "float",
		"transfertaxcounty":                                 "float",
		"transferinfopurchasedownpayment":                   "integer",
		"transferinfopurchaseloantovalue":                   "float",
		"lastupdated":                                       "float",
		"publicationdate":                                   "float",
		"grantoraddressinfoprivacy":                         "integer",
		"granteeinfoentitycount":                            "integer",
		"granteeinvestorflag":                               "integer",
		"granteemailaddressinfoprivacy":                     "integer",
		"mortgage1recordingdate":                            "float",
		"mortgage1amount":                                   "integer",
		"mortgage1lendercode":                               "integer",
		"mortgage1lenderinfosellercarrybackflag":            "integer",
		"mortgage1term":                                     "integer",
		"mortgage1termdate":                                 "float",
		"mortgage1interestrate":                             "float",
		"mortgage1infointeresttypechangeyear":               "integer",
		"mortgage1infointeresttypechangemonth":              "integer",
		"mortgage1infointeresttypechangeday":                "integer",
		"mortgage1interestrateminfirstchangerateconversion": "float",
		"mortgage1interestratemaxfirstchangerateconversion": "float",
		"mortgage1interestmargin":                           "integer",
		"mortgage1interestindex":                            "float",
		"mortgage1interestratemax":                          "float",
		"mortgage1interestonlyflag":                         "integer",
		"mortgage2recordingdate":                            "float",
		"mortgage2amount":                                   "integer",
		"mortgage2lendercode":                               "integer",
		"mortgage2lenderinfosellercarrybackflag":            "integer",
		"mortgage2term":                                     "integer",
		"mortgage2termdate":                                 "float",
		"mortgage2interestrate":                             "float",
		"mortgage2infointeresttypechangeyear":               "integer",
		"mortgage2infointeresttypechangemonth":              "integer",
		"mortgage2infointeresttypechangeday":                "integer",
		"mortgage2interestrateminfirstchangerateconversion": "float",
		"mortgage2interestratemaxfirstchangerateconversion": "float",
		"mortgage2interestmargin":                           "integer",
		"mortgage2interestindex":                            "float",
		"mortgage2interestratemax":                          "float",
		"mortgage2interestonlyflag":                         "integer",
	}

	assessorRetype = map[string]string{
		"accessabilityelevatorflag":          "integer",
		"accessabilityhandicapflag":          "integer",
		"apnyearadded":                       "integer",
		"arborpergolaflag":                   "integer",
		"area1stfloor":                       "integer",
		"area2ndfloor":                       "integer",
		"areabuilding":                       "integer",
		"areagross":                          "integer",
		"arealotacres":                       "float",
		"arealotdepth":                       "float",
		"arealotsf":                          "float",
		"arealotwidth":                       "float",
		"areaupperfloors":                    "integer",
		"assessorlastsaleamount":             "integer",
		"assessorpriorsaleamount":            "integer",
		"balconyarea":                        "integer",
		"bathcount":                          "float",
		"bathhousearea":                      "integer",
		"bathhouseflag":                      "integer",
		"bathpartialcount":                   "integer",
		"bedroomscount":                      "integer",
		"boataccessflag":                     "integer",
		"boathousearea":                      "integer",
		"boathouseflag":                      "integer",
		"boatliftflag":                       "integer",
		"breezewayflag":                      "integer",
		"buildingscount":                     "integer",
		"cabinarea":                          "integer",
		"cabinflag":                          "integer",
		"canopyarea":                         "integer",
		"canopyflag":                         "integer",
		"censusblock":                        "integer",
		"censusblockgroup":                   "integer",
		"censustract":                        "integer",
		"centralvacuumflag":                  "integer",
		"communityrecroomflag":               "integer",
		"congressionaldistricthouse":         "integer",
		"contentintercomflag":                "integer",
		"contentoverheaddoorflag":            "integer",
		"contentsaunaflag":                   "integer",
		"contentsoundsystemflag":             "integer",
		"contentstormshutterflag":            "integer",
		"courtyardarea":                      "integer",
		"courtyardflag":                      "integer",
		"deckarea":                           "integer",
		"deckflag":                           "integer",
		"deedlastsaleprice":                  "integer",
		"deedlastsaletransactionid":          "integer",
		"drivewayarea":                       "integer",
		"escalatorflag":                      "integer",
		"featurebalconyflag":                 "integer",
		"fencearea":                          "integer",
		"fireplacecount":                     "integer",
		"flooringmaterialprimary":            "integer",
		"gazeboarea":                         "integer",
		"gazeboflag":                         "integer",
		"golfcoursegreenflag":                "integer",
		"graineryarea":                       "integer",
		"graineryflag":                       "integer",
		"greenhousearea":                     "integer",
		"greenhouseflag":                     "integer",
		"guesthousearea":                     "integer",
		"guesthouseflag":                     "integer",
		"kennelarea":                         "integer",
		"kennelflag":                         "integer",
		"lastownershiptransfertransactionid": "integer",
		"leantoarea":                         "integer",
		"leantoflag":                         "integer",
		"loadingplatformarea":                "integer",
		"loadingplatformflag":                "integer",
		"milkhousearea":                      "integer",
		"milkhouseflag":                      "integer",
		"outdoorkitchenfireplaceflag":        "integer",
		"parcelnumberyearchange":             "integer",
		"parcelshellrecord":                  "integer",
		"parkingcarportarea":                 "integer",
		"parkinggaragearea":                  "integer",
		"parkingrvparkingflag":               "integer",
		"parkingspacecount":                  "integer",
		"patioarea":                          "integer",
		"plumbingfixturescount":              "integer",
		"polestructurearea":                  "integer",
		"polestructureflag":                  "integer",
		"pondflag":                           "integer",
		"pool":                               "integer",
		"poolarea":                           "integer",
		"poolhousearea":                      "integer",
		"poolhouseflag":                      "integer",
		"porcharea":                          "integer",
		"poultryhousearea":                   "integer",
		"poultryhouseflag":                   "integer",
		"previousassessedvalue":              "integer",
		"propertyaddressinfoprivacy":         "integer",
		"propertylatitude":                   "float",
		"propertylongitude":                  "float",
		"quonsetarea":                        "integer",
		"quonsetflag":                        "integer",
		"roomsatticarea":                     "integer",
		"roomsatticflag":                     "integer",
		"roomsbasementarea":                  "integer",
		"roomsbasementareafinished":          "integer",
		"roomsbasementareaunfinished":        "integer",
		"roomsbonusroomflag":                 "integer",
		"roomsbreakfastnookflag":             "integer",
		"roomscellarflag":                    "integer",
		"roomscellarwineflag":                "integer",
		"roomscount":                         "integer",
		"roomsexerciseflag":                  "integer",
		"roomsgameflag":                      "integer",
		"roomsgreatflag":                     "integer",
		"roomshobbyflag":                     "integer",
		"roomslaundryflag":                   "integer",
		"roomsmediaflag":                     "integer",
		"roomsmudflag":                       "integer",
		"roomsofficearea":                    "integer",
		"roomsofficeflag":                    "integer",
		"roomssaferoomflag":                  "integer",
		"roomssittingflag":                   "integer",
		"roomsstormshelter":                  "integer",
		"roomsstudyflag":                     "integer",
		"roomssunroomflag":                   "integer",
		"roomsutilityarea":                   "integer",
		"safetyfiresprinklersflag":           "integer",
		"securityalarmflag":                  "integer",
		"shedarea":                           "integer",
		"siloarea":                           "integer",
		"siloflag":                           "integer",
		"sportscourtflag":                    "integer",
		"sprinklersflag":                     "integer",
		"stablearea":                         "integer",
		"stableflag":                         "integer",
		"storagebuildingarea":                "integer",
		"storagebuildingflag":                "integer",
		"storiescount":                       "integer",
		"structurestyle":                     "integer",
		"taxassessedimprovementsperc":        "float",
		"taxassessedvalueimprovements":       "integer",
		"taxassessedvalueland":               "integer",
		"taxassessedvaluetotal":              "integer",
		"taxbilledamount":                    "float",
		"taxdelinquentyear":                  "integer",
		"taxfiscalyear":                      "integer",
		"taxmarketimprovementsperc":          "float",
		"taxmarketvalueimprovements":         "integer",
		"taxmarketvalueland":                 "integer",
		"taxmarketvaluetotal":                "integer",
		"taxmarketvalueyear":                 "integer",
		"taxyearassessed":                    "integer",
		"tenniscourtflag":                    "integer",
		"topographycode":                     "integer",
		"unitscount":                         "integer",
		"utilitiesmobilehomehookupflag":      "integer",
		"utilitybuildingarea":                "integer",
		"utilitybuildingflag":                "integer",
		"waterfeatureflag":                   "integer",
		"wetbarflag":                         "integer",
		"yearbuilt":                          "integer",
		"yearbuilteffective":                 "integer",
	}
)

// NeedToCallAddressParser determines whether the address parser needs to be called
// based on the fields present in the PropertySearchRequests object.
//
// It checks if the "fulladdress" field is present and non-empty, and if the "aupid"
// field is either absent or empty.
func NeedToCallAddressParser(obj models.PropertySearchRequests) bool {
	v := reflect.ValueOf(obj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldsToReturnInRoot = append(fieldsToReturnInRoot, t.Field(i).Name)
	}

	hasFullAddressInSearch := false
	for _, field := range searchFields {
		if strings.EqualFold(field, "fulladdress") {
			hasFullAddressInSearch = true
			break
		}
	}

	hasAupidInSearch := false
	for _, field := range searchFields {
		if strings.EqualFold(field, "aupid") {
			hasAupidInSearch = true
			break
		}
	}

	if hasFullAddressInSearch && strings.TrimSpace(obj.FullAddress) != "" && (!hasAupidInSearch || strings.TrimSpace(obj.OldAupid) == "") {
		return true
	}

	return false
}

func CheckSearchHasRequiredField(obj models.PropertySearchRequests) bool {
	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		for _, searchField := range searchFields {
			if strings.EqualFold(field.Name, searchField) {
				value := v.Field(i)
				if !value.IsZero() {
					return true
				}
			}
		}
	}

	return false
}

func faPropertyIdLookup(propertyId string) map[string]any {
	query := map[string]any{
		"from": 0,
		"size": 1,
		"query": map[string]any{
			"terms": map[string]any{
				"PropertyID": map[string]any{
					"index": "am_properties",
					"id":    propertyId,
					"path":  "fa_property_id",
				},
			},
		},
		"sort": []map[string]any{
			{
				"_score": map[string]any{"order": "desc"},
			},
		},
	}

	return query
}

func ToString(v any) any {
	if str, ok := v.(string); ok && str != "" {
		return strings.ToUpper(str)
	}
	return nil
}

func PadWithZeros(v any, length int) any {
	if v == nil {
		return nil
	}
	str := fmt.Sprintf("%v", v)
	if len(str) >= length {
		return strings.ToUpper(str)
	}

	return strings.ToUpper(strings.Repeat("0", length-len(str)) + str)
}

// processFields processes the fields to be included or excluded based on the provided property search requests
// and assessor hit data. It returns the included and excluded fields.
func ProcessFields(_req []models.PropertySearchRequests, interimProps []models.PropertySearchRequests, assessorHit map[string]any) (includedFields []string, excludedFields []string) {
	for _, req := range _req {
		if len(req.IncludeFields) > 0 {
			includedFields = req.IncludeFields
			break
		}
	}

	for _, req := range _req {
		if len(req.OmitFields) > 0 {
			excludedFields = req.OmitFields
			break
		}
	}

	if len(includedFields) == 0 {
		if aupid, exists := assessorHit["aupid"]; exists {
			for _, prop := range interimProps {
				if prop.OldAupid == aupid.(string) {
					if len(prop.IncludeFields) > 0 {
						includedFields = prop.IncludeFields
						break
					}
				}
			}
		}
	}

	if len(excludedFields) == 0 {
		if aupid, exists := assessorHit["aupid"]; exists {
			for _, prop := range interimProps {
				if prop.OldAupid == aupid.(string) {
					if len(prop.OmitFields) > 0 {
						excludedFields = prop.OmitFields
						break
					}
				}
			}
		}
	}

	if len(includedFields) > 0 {
		for key := range assessorHit {
			if !slices.Contains(includedFields, key) {
				delete(assessorHit, key)
			}
		}
	} else if len(excludedFields) > 0 {
		for key := range assessorHit {
			if slices.Contains(excludedFields, key) {
				delete(assessorHit, key)
			}
		}
	}

	return includedFields, excludedFields
}

// MakeElasticsearchRequest performs a search query against Elasticsearch
// Parameters:
//   - query: map containing the Elasticsearch query
//   - index: target Elasticsearch index (defaults to os_recorderindex5 if empty)
//   - headers: additional HTTP headers to include in the request
//
// Returns: search results as map[string]any and error if any
// MakeElasticsearchRequest performs a search query against Elasticsearch
// Parameters:
//   - query: map containing the Elasticsearch query
//   - index: target Elasticsearch index (defaults to os_recorderindex5 if empty)
//   - headers: additional HTTP headers to include in the request
//
// Returns: search results as map[string]any and error if any
func (dom *domain) MakeElasticsearchRequest(ctx context.Context, query map[string]any, index string, headers map[string]any) (map[string]any, error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}

	jsonBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	if index == "" {
		index = os_recorderindex5
	}

	osConfig, ok := dom.config.File.OpenSearch[consts.ConfigKeyOpenSearchLegacyApi]
	if !ok {
		return nil, errors.New("internal error")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/%s/_search", osConfig.Addresses[0], index), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}
	req.SetBasicAuth(osConfig.Username, osConfig.Password)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var esError models.ElasticsearchError
		if err := json.Unmarshal(body, &esError); err != nil {
			return nil, fmt.Errorf("%w: status code %d", ErrElasticsearchRequest, resp.StatusCode)
		}
		return nil, fmt.Errorf("%w: %v", ErrElasticsearchRequest, esError)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Request != nil && resp.Request.URL != nil {
		result["url.path"] = resp.Request.URL.Path
	}

	return result, nil
}

// CountByIndex sends a count request to an Elasticsearch server for the provided index.
func (dom *domain) CountByIndex(index string, headers map[string]any) (float64, error) {
	query := map[string]any{
		"query": map[string]any{
			"match_all": map[string]any{},
		},
	}

	jsonBody, err := json.Marshal(query)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling query")
		return 0, err
	}

	if index == "" {
		index = os_recorderindex5
	}

	osConfig, ok := dom.config.File.OpenSearch[consts.ConfigKeyOpenSearchLegacyApi]
	if !ok {
		return 0, errors.New("internal error")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/_count", osConfig.Addresses[0], index), bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error().Err(err).Msg("Error creating request")
		return 0, err
	}

	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}
	req.SetBasicAuth(osConfig.Username, osConfig.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error().Err(err).Msg("Error decoding response")
		return 0, err
	}

	count := result["count"].(float64)
	return count, nil
}

// MakeHeaders generates a map of HTTP headers with predefined values.
// The headers include:
// - "Content-Type": set to "application/json"
// - "X-Amz-Date": set to the current UTC time plus 5 minutes, formatted as "20060102T150405Z"
// This is necessary for making requests to the Elasticsearch server.
func MakeHeaders() map[string]any {
	headers := map[string]any{
		"Content-Type": "application/json",
		"X-Amz-Date":   time.Now().Add(5 * time.Minute).UTC().Format("20060102T150405Z"),
	}

	return headers
}

// filterLayouts filters the layouts to include based on the provided address.
func filterLayouts(address models.Address) []string {
	if address.IncludeLayouts == nil {
		return []string{}
	}

	includedLayouts := make([]string, len(address.IncludeLayouts))
	copy(includedLayouts, address.IncludeLayouts)

	if len(address.OmitLayouts) == 0 {
		return includedLayouts
	}

	excludedMap := make(map[string]bool, len(address.OmitLayouts))
	for _, layout := range address.OmitLayouts {
		excludedMap[layout] = true
	}

	filteredLayouts := make([]string, 0, len(includedLayouts))
	for _, item := range includedLayouts {
		if !excludedMap[item] {
			filteredLayouts = append(filteredLayouts, item)
		}
	}

	return filteredLayouts
}

// buildElasticsearchQuery constructs an Elasticsearch query based on the provided property search criteria.
// It creates a boolean query with "must" clauses for each non-empty field in the PropertySearchRequests.
// The fields matched include street name, house number, ZIP code, street direction, street suffix,
// street post direction, city, and state.
// The function returns a map representing the complete Elasticsearch query with:
//   - Pagination parameters (from: 0, size: 1)
//   - The boolean query with all applicable field matches
//   - Sort order by ValuationDate (descending), SitusUnitNbr.keyword (ascending), and relevance score (descending)
func buildElasticsearchQuery(prop models.PropertySearchRequests) map[string]any {
	mustMatch := []map[string]any{}

	if prop.StreetName != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusStreet": map[string]any{
					"query": prop.StreetName,
				},
			},
		})
	}

	if prop.HouseNumber != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusHouseNbr": map[string]any{
					"query": prop.HouseNumber,
				},
			},
		})
	}

	if prop.Zip5 != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusZIP5": map[string]any{
					"query": prop.Zip5,
				},
			},
		})
	}

	if prop.StreetDirection != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusDirectionLeft": map[string]any{
					"query": prop.StreetDirection,
				},
			},
		})
	}

	if prop.StreetSuffix != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusMode": map[string]any{
					"query": prop.StreetSuffix,
				},
			},
		})
	}

	if prop.StreetPostDirection != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusDirectionRight": map[string]any{
					"query": prop.StreetPostDirection,
				},
			},
		})
	}

	if prop.City != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusCity": map[string]any{
					"query": prop.City,
				},
			},
		})
	}

	if prop.State != "" {
		mustMatch = append(mustMatch, map[string]any{
			"match": map[string]any{
				"SitusState": map[string]any{
					"query": prop.State,
				},
			},
		})
	}

	query := map[string]any{
		"from": 0,
		"size": 1,
		"query": map[string]any{
			"bool": map[string]any{
				"must": mustMatch,
			},
		},
		"sort": []map[string]any{
			{
				"ValuationDate": map[string]any{"order": "desc"},
			},
			{
				"SitusUnitNbr.keyword": map[string]any{"order": "asc"},
			},
			{
				"_score": map[string]any{"order": "desc"},
			},
		},
	}

	return query
}

func buildElasticsearchQueryAupid(aupid string) map[string]any {
	query := map[string]any{
		"from": 0,
		"size": 1,
		"query": map[string]any{
			"match": map[string]any{
				"aupid": map[string]any{
					"query": aupid,
				},
			},
		},
	}

	return query
}

func updateCompValues(fa_address_response map[string]any, compAssr map[string]any, compRecorder map[string]any, valueavmcomp map[string]any) {
	for sourceKey, destKey := range addressFields {
		if value := fa_address_response[sourceKey]; value != nil && value != "" {
			if sourceKey == "ZIP5" || sourceKey == "ZIP4" {
				padding := 5
				if sourceKey == "ZIP4" {
					padding = 4
				}
				valueavmcomp[destKey] = PadWithZeros(value, padding)
			} else {
				valueavmcomp[destKey] = ToString(value)
			}
		}
	}

	if compAssr != nil {
		for sourceKey, destKey := range assrFields {
			if value := compAssr[sourceKey]; value != nil && value != "" {
				valueavmcomp[destKey] = value
			}
		}
	}

	if compRecorder != nil {
		for sourceKey, destKey := range recorderFields {
			if value := compRecorder[sourceKey]; value != nil && value != "" {
				valueavmcomp[destKey] = value
			}
		}
	}
}

func getStringOrEmpty(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getScore(result map[string]any) map[string]any {
	if result == nil {
		return nil
	}

	hits, ok := result["hits"].(map[string]any)
	if !ok {
		log.Error().Err(&errors.Object{
			Id:     "29e5342b-4051-4374-8bc1-fc79e65fac7c",
			Code:   errors.Code_UNKNOWN,
			Detail: "Invalid hits format for property data",
		})
		return nil
	}

	hitsArray, ok := hits["hits"].([]any)
	if !ok || len(hitsArray) == 0 {
		log.Error().Err(&errors.Object{
			Id:     "7a908ed0-8785-4535-b364-43ea273fbe3a",
			Code:   errors.Code_UNKNOWN,
			Detail: "No hits found for property data.",
		})
		return nil
	}

	firstHit, ok := hitsArray[0].(map[string]any)
	if !ok {
		log.Error().Err(&errors.Object{
			Id:     "9f4db5e0-dffc-4339-a16d-13c577bdc60b",
			Code:   errors.Code_UNKNOWN,
			Detail: "First hit is not a map for property data.",
		})
		return nil
	}

	source, ok := firstHit["_source"].(map[string]any)
	if !ok {
		log.Error().Err(&errors.Object{
			Id:     "17fb0929-de13-48ba-825e-b652f8567cbe",
			Code:   errors.Code_UNKNOWN,
			Detail: "Could not extract _source from first hit for property data.",
		})
		return nil
	}

	return source

}
