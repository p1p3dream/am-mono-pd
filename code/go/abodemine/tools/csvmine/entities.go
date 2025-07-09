package main

import (
	"encoding/json"

	"github.com/shopspring/decimal"

	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type MinMaxValue struct {
	DecimalValue *decimal.Decimal `json:"decimal_value,omitempty"`
	StringValue  *string          `json:"string_value,omitempty"`
}

func (m *MinMaxValue) UnmarshalJSON(data []byte) error {
	var decimalValue decimal.Decimal
	if err := json.Unmarshal(data, &decimalValue); err == nil {
		m.DecimalValue = &decimalValue
		m.StringValue = val.PtrRef(string(data))
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		m.StringValue = &stringValue
		return nil
	}

	return &errors.Object{
		Id:     "b1400349-337a-476b-ba9b-2ce108d11994",
		Code:   errors.Code_FAILED_PRECONDITION,
		Detail: "Compatible type not found.",
	}
}

type ColumnStat struct {
	ColumnID     int             `json:"column_id"`
	ColumnName   string          `json:"column_name"`
	Type         string          `json:"type"`
	Nulls        bool            `json:"nulls"`
	NonNulls     int             `json:"nonnulls"`
	Unique       int             `json:"unique"`
	Min          MinMaxValue     `json:"min"`
	Max          MinMaxValue     `json:"max"`
	Sum          decimal.Decimal `json:"sum"`
	Mean         decimal.Decimal `json:"mean"`
	Median       decimal.Decimal `json:"median"`
	Stdev        decimal.Decimal `json:"stdev"`
	MaxPrecision int             `json:"maxprecision"`
}

type DataFileTemplateInput struct {
	PackageName string
	StructName  string
	TableName   string

	Fields []*DataFileTemplateField
}

const (
	GoValidationBoolPtrFromStringIfNonZero    = "BOOL_PTR_FROM_STRING_IF_NON_ZERO"
	GoValidationDecimalPtrFromStringIfNonZero = "DECIMAL_PTR_FROM_STRING_IF_NON_ZERO"
	GoValidationFloat64PtrFromStringIfNonZero = "FLOAT64_PTR_FROM_STRING_IF_NON_ZERO"
	GoValidationIntPtrFromStringIfNonZero     = "INT_PTR_FROM_STRING_IF_NON_ZERO"
	GoValidationInt64PtrFromStringIfNonZero   = "INT64_PTR_FROM_STRING_IF_NON_ZERO"
	GoValidationMustBool                      = "MUST_BOOL"
	GoValidationMustDecimal                   = "MUST_DECIMAL"
	GoValidationMustFloat64                   = "MUST_FLOAT64"
	GoValidationMustInt                       = "MUST_INT"
	GoValidationMustInt64                     = "MUST_INT64"
	GoValidationMustString                    = "MUST_STRING"
	GoValidationMustTime                      = "MUST_TIME"
	GoValidationStringPtrIfNonZero            = "STRING_PTR_IF_NON_ZERO"
	GoValidationTimePtrFromStringIfNonZero    = "TIME_PTR_FROM_STRING_IF_NON_ZERO"
)

type DataFileTemplateField struct {
	ColumnName string

	GoName           string
	GoType           string
	GoPtr            bool
	GoValidationType string
	GoTimeFormat     string

	PostgresName        string
	PostgresNamePadding string
	PostgresType        string
}
