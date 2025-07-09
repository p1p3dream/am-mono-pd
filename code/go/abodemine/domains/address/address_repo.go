package address

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/geog"
	"abodemine/lib/opensearchutils"
	"abodemine/lib/ptr"
	"abodemine/lib/sqlutil"
	"abodemine/models"
	"abodemine/repositories/opensearch"
)

type Repository interface {
	SearchPropertyAddressRecord(r *arc.Request, in *SearchPropertyAddressRecordInput) (*SearchPropertyAddressRecordOutput, error)
	SelectPropertyAddressRecord(r *arc.Request, in *SelectPropertyAddressRecordInput) (*SelectPropertyAddressRecordOutput, error)

	SelectFipsRecord(r *arc.Request, in *SelectFipsRecordInput) (*SelectFipsRecordOutput, error)

	SelectZip5Record(r *arc.Request, in *SelectZip5RecordInput) (*SelectZip5RecordOutput, error)
	UpdateZip5Table(r *arc.Request) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

type SearchPropertyAddressRecordInput struct {
	IncludePropertyRefs bool

	ApiSearchAddress *models.ApiSearchAddress

	Limit int32
}

type SearchPropertyAddressRecordOutput struct {
	Records      []*entities.PropertyAddress
	PropertyRefs []*entities.PropertyRef
}

func (repo *repository) SearchPropertyAddressRecord(r *arc.Request, in *SearchPropertyAddressRecordInput) (*SearchPropertyAddressRecordOutput, error) {
	out := &SearchPropertyAddressRecordOutput{}

	if in.ApiSearchAddress == nil {
		return out, nil
	}

	osClient, err := r.Dom().SelectOpenSearch(consts.ConfigKeyOpenSearchSearch)
	if err != nil {
		return nil, &errors.Object{
			Id:     "edbfdfbb-5717-44df-9e19-c10ae6dacfc1",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get OpenSearch client",
			Cause:  err.Error(),
		}
	}

	osIndex, err := r.Dom().Values().String(consts.ConfigKeyOpenSearchIndexAddresses)
	if err != nil {
		return nil, errors.Forward(err, "4eef4a7a-7cd0-42ba-b966-1298a8557d42")
	}

	query := PropertyAddressSearchQuery{}

	if in.Limit > 0 {
		query.Size = in.Limit
	}

	addr := in.ApiSearchAddress
	must := query.Query.Bool.Must

	if addr.FullStreetAddress != "" {
		must = append(must, map[string]any{
			"match_phrase_prefix": map[string]any{
				"fullStreetAddress": map[string]any{
					"query": strings.ToUpper(addr.FullStreetAddress),
				},
			},
		})
	}

	if addr.State != "" {
		must = append(must, opensearchutils.NewMatchQuery("state", addr.State).ToMap())
	}

	if addr.City != "" {
		must = append(must, opensearchutils.NewMatchQuery("city", addr.City).ToMap())
	}

	if addr.Zip5 != "" {
		must = append(must, opensearchutils.NewTermQuery("zip5", addr.Zip5).ToMap())
	}

	if addr.HouseNumber != "" {
		must = append(must, opensearchutils.NewMatchQuery("houseNumber", addr.HouseNumber).ToMap())
	}

	if addr.StreetName != "" {
		must = append(must, opensearchutils.NewMatchQuery("streetName", addr.StreetName).ToMap())
	}

	if addr.StreetSuffix != "" {
		must = append(must, opensearchutils.NewMatchQuery("streetSuffix", addr.StreetSuffix).ToMap())
	}

	if addr.UnitType != "" {
		must = append(must, opensearchutils.NewMatchQuery("unitType", addr.UnitType).ToMap())
	}

	if addr.UnitNumber != "" {
		must = append(must, opensearchutils.NewMatchQuery("unitNumber", addr.UnitNumber).ToMap())
	}

	if addr.StreetPreDirection != "" {
		must = append(must, opensearchutils.NewMatchQuery("streetPreDirection", addr.StreetPreDirection).ToMap())
	}

	if addr.StreetPostDirection != "" {
		must = append(must, opensearchutils.NewMatchQuery("streetPostDirection", addr.StreetPostDirection).ToMap())
	}

	query.Query.Bool.Must = must

	queryBody, err := json.Marshal(query)
	if err != nil {
		return nil, &errors.Object{
			Id:     "9401408e-d145-4db4-aa86-3b1bff858cf3",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to marshal query.",
			Cause:  err.Error(),
		}
	}

	search := opensearchapi.SearchRequest{
		Index: []string{osIndex},
		Body:  bytes.NewReader(queryBody),
	}

	resp, err := search.Do(r.Context(), osClient)
	if err != nil {
		return nil, &errors.Object{
			Id:     "faf16ca0-c6a2-4f26-a153-061688eb9976",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to search addresses",
			Cause:  err.Error(),
		}
	}
	defer resp.Body.Close()

	var result PropertyAddressSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &errors.Object{
			Id:     "16e71479-aa7f-4db5-9411-c7865b5a6483",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to decode search response",
			Cause:  err.Error(),
		}
	}

	if len(result.Hits.Hits) == 0 {
		return out, nil
	}

	for _, hit := range result.Hits.Hits {
		out.Records = append(out.Records, &hit.Source)

		if in.IncludePropertyRefs {
			out.PropertyRefs = append(out.PropertyRefs, &entities.PropertyRef{
				Aupid: hit.Source.Aupid,
			})
		}
	}

	return out, nil
}

type ScanPropertyAddressRecordInput struct {
	Columns []string
	Rows    pgx.Rows

	IncludePropertyRefs      bool
	IncludeStateFullName     bool
	ReturnOpenSearchDocument bool
}

type ScanPropertyAddressRecordOutput struct {
	Documents    []opensearch.Document
	PropertyRefs []*entities.PropertyRef
	Records      []*entities.PropertyAddress
}

func (repo *repository) ScanPropertyAddressRecord(r *arc.Request, in *ScanPropertyAddressRecordInput) (*ScanPropertyAddressRecordOutput, error) {
	cols := make(map[string]struct{}, len(in.Columns))

	for _, column := range in.Columns {
		cols[column] = struct{}{}
	}

	out := &ScanPropertyAddressRecordOutput{}

	// Columns MUST be scanned in alphabetical order.
	for in.Rows.Next() {
		propertyRef := &entities.PropertyRef{}
		record := &entities.PropertyAddress{}
		scans := make([]any, len(cols))
		col := 0

		if _, ok := cols["addresses.city"]; ok {
			scans[col] = &record.City
			col++
		}

		if _, ok := cols["addresses.county"]; ok {
			scans[col] = &record.County
			col++
		}

		if _, ok := cols["addresses.created_at"]; ok {
			scans[col] = &record.CreatedAt
			col++
		}

		if _, ok := cols["addresses.fips"]; ok {
			scans[col] = &record.Fips
			col++
		}

		if _, ok := cols["addresses.full_street_address"]; ok {
			scans[col] = &record.FullStreetAddress
			col++
		}

		if _, ok := cols["addresses.house_number"]; ok {
			scans[col] = &record.HouseNumber
			col++
		}

		if _, ok := cols["addresses.id"]; ok {
			scans[col] = &record.Id
			col++
		}

		if _, ok := cols["addresses.meta->>'property_id'"]; ok {
			scans[col] = &record.Aupid
			col++
		}

		if _, ok := cols["addresses.state"]; ok {
			scans[col] = &record.State
			col++
		}

		if _, ok := cols["addresses.street_name"]; ok {
			scans[col] = &record.StreetName
			col++
		}

		if _, ok := cols["addresses.street_pos_direction"]; ok {
			scans[col] = &record.StreetPostDirection
			col++
		}

		if _, ok := cols["addresses.street_pre_direction"]; ok {
			scans[col] = &record.StreetPreDirection
			col++
		}

		if _, ok := cols["addresses.street_suffix"]; ok {
			scans[col] = &record.StreetSuffix
			col++
		}

		if _, ok := cols["addresses.unit_nbr"]; ok {
			scans[col] = &record.UnitNumber
			col++
		}

		if _, ok := cols["addresses.unit_type"]; ok {
			scans[col] = &record.UnitType
			col++
		}

		if _, ok := cols["addresses.updated_at"]; ok {
			scans[col] = &record.UpdatedAt
			col++
		}

		if _, ok := cols["addresses.zip5"]; ok {
			scans[col] = &record.Zip5
			col++
		}

		if _, ok := cols["properties.ad_attom_id"]; ok {
			scans[col] = &propertyRef.ADAttomId
			col++
		}

		if _, ok := cols["properties.fa_property_id"]; ok {
			scans[col] = &propertyRef.FAPropertyId
			col++
		}

		if _, ok := cols["properties.id"]; ok {
			scans[col] = &propertyRef.Aupid
			col++
		}

		if err := in.Rows.Scan(scans...); err != nil {
			return nil, &errors.Object{
				Id:     "d3f9e8c5-5b6f-4a3e-8b3e-3c6b5f8b9e3d",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to scan row.",
				Cause:  err.Error(),
			}
		}

		if in.IncludeStateFullName && record.State != nil {
			record.StateFullName = ptr.String(geog.UsStateFullName(*record.State))
		}

		if in.ReturnOpenSearchDocument {
			out.Documents = append(out.Documents, record)
		} else {
			out.Records = append(out.Records, record)
		}

		if in.IncludePropertyRefs {
			out.PropertyRefs = append(out.PropertyRefs, propertyRef)
		}
	}

	return out, nil
}

type SelectPropertyAddressRecordInput struct {
	ApiSearchAddress *models.ApiSearchAddress
	Aupid            *uuid.UUID
	Fips             string
	IdGt             *uuid.UUID

	Columns []string
	OrderBy []string
	Limit   uint64

	IncludePropertyRefs      bool
	IncludeStateFullName     bool
	ReturnOpenSearchDocument bool
}

type SelectPropertyAddressRecordOutput struct {
	Documents    []opensearch.Document
	Records      []*entities.PropertyAddress
	PropertyRefs []*entities.PropertyRef
}

func (repo *repository) SelectPropertyAddressRecord(r *arc.Request, in *SelectPropertyAddressRecordInput) (*SelectPropertyAddressRecordOutput, error) {
	var columns []string

	if len(in.Columns) > 0 {
		columns = in.Columns
	} else {
		columns = []string{
			"addresses.full_street_address",
			"addresses.state",
			"addresses.county", // administrative region (can contain multiple cities)
			"addresses.zip5",
			"addresses.street_pre_direction",
			"addresses.house_number",
			"addresses.street_name",
			"addresses.street_pos_direction",
			"addresses.street_suffix",
			"addresses.unit_type",
			"addresses.unit_nbr",
			"addresses.city", // name of the specific city or locality
		}
	}

	var joins []*sqlutil.SquirrelJoin

	if in.IncludePropertyRefs {
		columns = append(
			columns,
			"properties.id",
			"properties.ad_attom_id",
			"properties.fa_property_id",
		)

		joins = append(joins, &sqlutil.SquirrelJoin{
			Pred: "join properties on addresses.id = properties.address_id",
		})
	}

	// Ensure columns are sorted for consistent scanning.
	// This MUST happen just before the SQL is built.
	sort.Strings(columns)

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(columns...).
		From("addresses")

	switch {
	case in.Aupid != nil:
		joins = append(joins, &sqlutil.SquirrelJoin{
			Pred: "join properties on addresses.id = properties.address_id",
		})

		builder = builder.Where("properties.id = ?", in.Aupid)
	case in.ApiSearchAddress != nil:
		// Use upper(...) to hit indexes.

		if in.ApiSearchAddress.FullStreetAddress != "" {
			builder = builder.Where(
				"upper(addresses.full_street_address) = ?",
				in.ApiSearchAddress.FullStreetAddress,
			)
		}

		if in.ApiSearchAddress.State != "" {
			builder = builder.Where(
				"upper(addresses.state) = ?",
				in.ApiSearchAddress.State,
			)
		}

		if in.ApiSearchAddress.City != "" {
			builder = builder.Where(
				"upper(addresses.city) = ?",
				in.ApiSearchAddress.City,
			)
		}

		if in.ApiSearchAddress.Zip5 != "" {
			builder = builder.Where(
				"upper(addresses.zip5) = ?",
				in.ApiSearchAddress.Zip5,
			)
		}

		if in.ApiSearchAddress.HouseNumber != "" {
			builder = builder.Where(
				"upper(addresses.house_number) = ?",
				in.ApiSearchAddress.HouseNumber,
			)
		}

		if in.ApiSearchAddress.StreetName != "" {
			builder = builder.Where(
				"upper(addresses.street_name) = ?",
				in.ApiSearchAddress.StreetName,
			)
		}

		if in.ApiSearchAddress.StreetSuffix != "" {
			builder = builder.Where(
				"upper(addresses.street_suffix) = ?",
				in.ApiSearchAddress.StreetSuffix,
			)
		}

		if in.ApiSearchAddress.UnitType != "" {
			builder = builder.Where(
				"upper(addresses.unit_type) = ?",
				in.ApiSearchAddress.UnitType,
			)
		}

		if in.ApiSearchAddress.UnitNumber != "" {
			builder = builder.Where(
				"upper(addresses.unit_nbr) = ?",
				in.ApiSearchAddress.UnitNumber,
			)
		}

		if in.ApiSearchAddress.StreetPreDirection != "" {
			builder = builder.Where(
				"upper(addresses.street_pre_direction) = ?",
				in.ApiSearchAddress.StreetPreDirection,
			)
		}

		if in.ApiSearchAddress.StreetPostDirection != "" {
			builder = builder.Where(
				"upper(addresses.street_pos_direction) = ?",
				in.ApiSearchAddress.StreetPostDirection,
			)
		}
	}

	if in.Fips != "" {
		builder = builder.Where("addresses.fips = ?", in.Fips)
	}

	if in.IdGt != nil {
		builder = builder.Where("addresses.id > ?", in.IdGt)
	}

	if len(in.OrderBy) > 0 {
		builder = builder.OrderBy(in.OrderBy...)
	}

	if in.Limit > 0 {
		builder = builder.Limit(in.Limit)
	}

	for _, join := range joins {
		builder = builder.JoinClause(join.Pred, join.Args...)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "1612fe20-b290-4c6f-9a24-b3c11332f356",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "e12d217d-8613-42fa-93a4-32e10071169c")
	}
	defer rows.Close()

	scanOut, err := repo.ScanPropertyAddressRecord(r, &ScanPropertyAddressRecordInput{
		Columns:                  columns,
		Rows:                     rows,
		IncludePropertyRefs:      in.IncludePropertyRefs,
		IncludeStateFullName:     in.IncludeStateFullName,
		ReturnOpenSearchDocument: in.ReturnOpenSearchDocument,
	})
	if err != nil {
		return nil, errors.Forward(err, "6504a671-b701-4964-8b21-e4f54ebab310")
	}

	out := &SelectPropertyAddressRecordOutput{
		Documents:    scanOut.Documents,
		PropertyRefs: scanOut.PropertyRefs,
		Records:      scanOut.Records,
	}

	return out, nil
}

type SelectFipsRecordInput struct {
	OrderBy string
}

type SelectFipsRecordOutput struct {
	Records []*Fips
}

func (repo *repository) SelectFipsRecord(r *arc.Request, in *SelectFipsRecordInput) (*SelectFipsRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("fips").
		From("fips")

	if in.OrderBy != "" {
		builder = builder.OrderBy(in.OrderBy)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "b5873e52-b1ef-4fac-9e58-2cb1ba7bed93",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "f1dc20a7-21c1-413d-9a02-a414b64caac0")
	}
	defer rows.Close()

	out := &SelectFipsRecordOutput{}

	for rows.Next() {
		record := &Fips{}

		if err := rows.Scan(
			&record.Fips,
		); err != nil {
			return nil, &errors.Object{
				Id:     "155e0e33-59a2-4f01-8fd6-eb67a14de1ae",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}

type SelectZip5RecordInput struct {
	OrderBy string
}

type SelectZip5RecordOutput struct {
	Records []*Zip5
}

func (repo *repository) SelectZip5Record(r *arc.Request, in *SelectZip5RecordInput) (*SelectZip5RecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("zip5").
		From("zip5")

	if in.OrderBy != "" {
		builder = builder.OrderBy(in.OrderBy)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "71095010-ded3-42e5-8949-c62e5cc6f368",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "7b4b5af3-9d72-4246-955c-b9e97d05c341")
	}
	defer rows.Close()

	out := &SelectZip5RecordOutput{}

	for rows.Next() {
		record := &Zip5{}

		if err := rows.Scan(
			&record.Zip5,
		); err != nil {
			return nil, &errors.Object{
				Id:     "6f23ebdd-b0b1-40d2-ad73-3aa5d24cd1df",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to select row.",
				Cause:  err.Error(),
			}
		}

		out.Records = append(out.Records, record)
	}

	return out, nil
}

func (repo *repository) UpdateZip5Table(r *arc.Request) error {
	sql := `
		insert into zip5 (
			select distinct(property_address_zip) as zip5
			from ad_df_assessor
			where
				property_address_zip is not null
				and not exists (
					select 1
					from zip5
					where zip5.zip5 = ad_df_assessor.property_address_zip
				)
			union
			select distinct(situs_zip5) as zip5
			from fa_df_assessor
			where
				situs_zip5 is not null
				and not exists (
					select 1
					from zip5
					where zip5.zip5 = fa_df_assessor.situs_zip5
				)
		)
	`

	_, err := extutils.PgxExec(r, consts.ConfigKeyPostgresDatapipe, sql, nil)
	if err != nil {
		return errors.Forward(err, "1428c154-71e5-4e11-a7f0-a36931ce192e")
	}

	return nil
}
