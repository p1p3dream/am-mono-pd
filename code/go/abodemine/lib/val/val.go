// package val provides utilities for working with values.
package val

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"abodemine/lib/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Boolean.
////////////////////////////////////////////////////////////////////////////////

// Ternary returns one of two values based on a condition, similar to the ternary
// operator (?:) in other languages.
// Returns:
//   - ifTrue: if condition is true
//   - ifFalse: if condition is false
func Ternary[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

////////////////////////////////////////////////////////////////////////////////
// Coalesce.
////////////////////////////////////////////////////////////////////////////////

// Coalesce returns the first non-zero value from the given arguments.
func Coalesce[T comparable](values ...T) T {
	return CoalesceSlice(values)
}

// CoalesceSlice returns the first non-zero value from the given slice.
func CoalesceSlice[T comparable](values []T) T {
	var zero T

	for _, v := range values {
		if v != zero {
			return v
		}
	}

	return zero
}

////////////////////////////////////////////////////////////////////////////////
// Date.
////////////////////////////////////////////////////////////////////////////////

func IntegerDate(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

////////////////////////////////////////////////////////////////////////////////
// Pointers.
////////////////////////////////////////////////////////////////////////////////

// PtrRef returns a pointer to the given value.
func PtrRef[T any](v T) *T {
	return &v
}

// PtrDeref returns the value of the given pointer
// or a zero value if the pointer is nil.
func PtrDeref[T any](v *T) T {
	if v == nil {
		return *new(T)
	}
	return *v
}

// PtrEnsure returns a pointer to a new instance of T if v is nil.
func PtrEnsure[T any](v *T) *T {
	if v == nil {
		return new(T)
	}
	return v
}

////////////////////////////////////////////////////////////////////////////////
// Pointers from concrete values if non-zero.
// Per-type functions for performance.
////////////////////////////////////////////////////////////////////////////////

func BoolPtrFromStringIfNonZero(s string) (*bool, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, &errors.Object{
			Id:    "4973e13e-5de0-4fb6-8a19-a572c9d441bc",
			Code:  errors.Code_INVALID_ARGUMENT,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}

func BoolPtrFromYNString(s string) (*bool, error) {
	switch s {
	case "Y":
		return PtrRef(true), nil
	case "N":
		return PtrRef(false), nil
	case "":
		return nil, nil
	}

	return nil, &errors.Object{
		Id:   "58ec8280-b1cd-4f93-a2ad-a484160171c4",
		Code: errors.Code_INVALID_ARGUMENT,
		Meta: map[string]any{
			"value": s,
		},
	}
}

func DecimalPtrFromString(s string) (*decimal.Decimal, error) {
	v, err := decimal.NewFromString(s)
	if err != nil {
		return nil, &errors.Object{
			Id:    "d827753c-f482-4b33-b96e-3371bcd1a32f",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}

func DecimalPtrFromStringIfNonZero(s string) (*decimal.Decimal, error) {
	if s == "" {
		return nil, nil
	}

	return DecimalPtrFromString(s)
}

func Float64PtrFromStringIfNonZero(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, &errors.Object{
			Id:    "ea504258-bec7-4f95-b069-2061f6eec546",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}

func IntPtrFromStringIfNonZero(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, &errors.Object{
			Id:    "226f2cc9-1a37-4e2b-94a6-1c18f4837622",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}

func IntPtrFromFloat64StringIfNonZero(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, &errors.Object{
			Id:    "dbbba877-e40f-4510-a424-2c63d4e2ff3e",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	intValue := int(v)

	return &intValue, nil
}

func Int64PtrFromString(s string) (*int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, &errors.Object{
			Id:    "22561116-04dd-43b8-9c96-09537c4ebe97",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}

func Int64PtrFromStringIfNonZero(s string) (*int64, error) {
	if s == "" {
		return nil, nil
	}

	return Int64PtrFromString(s)
}

func StringPtrIfNonZero(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func TimePtrFromString(layout, s string) (*time.Time, error) {
	v, err := time.Parse(layout, s)
	if err != nil {
		return nil, &errors.Object{
			Id:    "e9b05d9f-0fb1-40f0-849e-630086eb234c",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"layout": layout,
				"value":  s,
			},
		}
	}

	return &v, nil
}

func TimePtrFromStringIfNonZero(layout, s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	return TimePtrFromString(layout, s)
}

func UUIDPtrFromStringIfNonZero(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}

	v, err := uuid.Parse(s)
	if err != nil {
		return nil, &errors.Object{
			Id:    "e5cf1b2c-712d-41ae-8c75-8a3c5d73782c",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	return &v, nil
}
