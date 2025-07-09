package main

const dataFileTemplate = `
{{- $PackageName := .PackageName -}}
package {{ $PackageName }}

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

type {{ .StructName }} struct {
	AMId        uuid.UUID
	AMCreatedAt time.Time
	AMUpdatedAt time.Time
	AMMeta      map[string]any
{{ range .Fields }}
	{{ .GoName }} {{ if .GoPtr }}*{{ end }}{{ .GoType }}
{{- end }}
}

func (dr *{{ .StructName }}) New(headers map[int]string, fields []string) (entities.DataRecord, error) {
	record := new({{ .StructName }})

	for k, header := range headers {
		field := fields[k]

		switch header {
{{- range .Fields }}
		case "{{ .ColumnName }}":
{{- if eq .GoValidationType "BOOL_PTR_FROM_STRING_IF_NON_ZERO" }}
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := val.BoolPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "DECIMAL_PTR_FROM_STRING_IF_NON_ZERO" }}
			v, err := val.DecimalPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "FLOAT64_PTR_FROM_STRING_IF_NON_ZERO" }}
			v, err := val.Float64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "INT_PTR_FROM_STRING_IF_NON_ZERO" }}
			v, err := val.IntPtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "INT64_PTR_FROM_STRING_IF_NON_ZERO" }}
			v, err := val.Int64PtrFromStringIfNonZero(field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_BOOL" }}
			switch strings.ToUpper(field) {
			case "Y":
				field = "t"
			case "N":
				field = "f"
			}

			v, err := strconv.ParseBool(field)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_DECIMAL" }}
			v, err := decimal.NewFromString(field)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_FLOAT64" }}
			v, err := strconv.ParseFloat(field, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_INT" }}
			v, err := strconv.Atoi(field)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_INT64" }}
			v, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "MUST_STRING" }}
			if field == "" {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Detail: "{{ .ColumnName }} is required.",
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = field
{{- else if eq .GoValidationType "MUST_TIME" }}
			v, err := time.Parse({{ .GoTimeFormat }}, field)
			if err != nil {
				return nil, &errors.Object{
					Id:    "{{ uuidv4 }}",
					Code:  errors.Code_INVALID_ARGUMENT,
					Cause: err.Error(),
					Meta: map[string]any{
						"value": field,
					},
				}
			}
			record.{{ .GoName }} = v
{{- else if eq .GoValidationType "STRING_PTR_IF_NON_ZERO" }}
			record.{{ .GoName }} = val.StringPtrIfNonZero(field)
{{- else if eq .GoValidationType "TIME_PTR_FROM_STRING_IF_NON_ZERO" }}
			v, err := val.TimePtrFromStringIfNonZero({{ .GoTimeFormat }}, field)
			if err != nil {
				return nil, errors.Forward(err, "{{ uuidv4 }}")
			}
			record.{{ .GoName }} = v
{{- else }}
			UNKNOWN_VALIDATION_TYPE: {{ .GoValidationType }}
{{- end }}
{{- end }}
		default:
			return nil, &errors.Object{
				Id:     "{{ uuidv4}}",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Unknown header.",
				Meta: map[string]any{
					"field_index": k,
					"field_value": field,
					"header":      header,
				},
			}
		}
	}

	return record, nil
}

func (dr *{{ .StructName }}) SQLColumns() []string {
	return []string{
		"am_id",
		"am_created_at",
		"am_updated_at",
		"am_meta",
{{- range .Fields }}
		"{{ .PostgresName }}",
{{- end }}
	}
}

func (dr *{{ .StructName }}) SQLTable() string {
	return "{{ .TableName}}"
}

func (dr *{{ .StructName }}) SQLValues() ([]any, error) {
	if dr.AMId == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return nil, &errors.Object{
				Id:     "{{ uuidv4 }}",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to generate UUID.",
				Cause:  err.Error(),
			}
		}
		dr.AMId = u
	}

	now := time.Now()

	if dr.AMCreatedAt.IsZero() {
		dr.AMCreatedAt = now
	}

	values := []any{
		dr.AMId,
		dr.AMCreatedAt,
		now,
		dr.AMMeta,
{{- range .Fields }}
		dr.{{ .GoName }},
{{- end }}
	}

	return values, nil
}
`

const sqlFileTemplate = `
create table {{ .TableName }} (
	am_id         uuid not null,
	am_created_at timestamp with time zone not null,
	am_updated_at timestamp with time zone not null,
	am_meta       jsonb,
{{ range .Fields }}
	{{ .PostgresName }}{{ .PostgresNamePadding }}{{ .PostgresType }}{{ if not .GoPtr }} not null{{ end }},
{{- end }}

	primary key (am_id)
);
`
