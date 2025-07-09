package gconf

import (
	"encoding/base64"
	"maps"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"abodemine/lib/errors"
)

type ConfigValues struct {
	Bool    map[string]string `json:"bool,omitempty" yaml:"bool,omitempty"`
	Bytes   map[string]string `json:"bytes,omitempty" yaml:"bytes,omitempty"`
	Decimal map[string]string `json:"decimal,omitempty" yaml:"decimal,omitempty"`
	Float64 map[string]string `json:"float64,omitempty" yaml:"float64,omitempty"`
	Int     map[string]string `json:"int,omitempty" yaml:"int,omitempty"`
	Int32   map[string]string `json:"int32,omitempty" yaml:"int32,omitempty"`
	Int64   map[string]string `json:"int64,omitempty" yaml:"int64,omitempty"`
	String  map[string]string `json:"string,omitempty" yaml:"string,omitempty"`
	Time    map[string]string `json:"time,omitempty" yaml:"time,omitempty"`
	Uint    map[string]string `json:"uint,omitempty" yaml:"uint,omitempty"`
	Uint32  map[string]string `json:"uint32,omitempty" yaml:"uint32,omitempty"`
	Uint64  map[string]string `json:"uint64,omitempty" yaml:"uint64,omitempty"`
}

// Values is a concurrency-safe struct that holds
// the parsed values from ConfigValues.
// It always returns a copy of the value.
type Values struct {
	_bool    map[string]bool
	_bytes   map[string][]byte
	_decimal map[string]decimal.Decimal
	_float64 map[string]float64
	_int     map[string]int
	_int32   map[string]int32
	_int64   map[string]int64
	_string  map[string]string
	_time    map[string]time.Time
	_uint    map[string]uint
	_uint32  map[string]uint32
	_uint64  map[string]uint64
}

func (v *Values) Bool(key string) (bool, error) {
	val, ok := v._bool[key]

	if !ok {
		return false, &errors.Object{
			Id:     "f05e06fe-972a-4e40-9eba-5b0ff87a5cdd",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Bytes(key string) ([]byte, error) {
	val, ok := v._bytes[key]

	if !ok {
		return nil, &errors.Object{
			Id:     "01299479-2cf0-462b-bea7-7360e1530016",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	// Return a copy of the value.
	return val[:], nil
}

func (v *Values) Decimal(key string) (decimal.Decimal, error) {
	val, ok := v._decimal[key]

	if !ok {
		return decimal.Decimal{}, &errors.Object{
			Id:     "f67e95e6-d741-477e-8f12-d249c00c480a",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Float64(key string) (float64, error) {
	val, ok := v._float64[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "7edefa50-ae71-444d-8204-0e9678f8e9a4",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Int(key string) (int, error) {
	val, ok := v._int[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "698916de-63c0-45e6-bd14-b96949444acc",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Int32(key string) (int32, error) {
	val, ok := v._int32[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "55601033-e9dc-43c1-b41b-4c1fa5a857fc",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Int64(key string) (int64, error) {
	val, ok := v._int64[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "f0df0ebd-4c20-4fee-a33a-b0fb795eb2a8",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) String(key string) (string, error) {
	val, ok := v._string[key]

	if !ok {
		return "", &errors.Object{
			Id:     "a235629e-3f27-4744-83f0-9e4433bb68c8",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Time(key string) (time.Time, error) {
	val, ok := v._time[key]

	if !ok {
		return time.Time{}, &errors.Object{
			Id:     "f2d861f5-1f6d-4661-bc8e-4a0bce46bbe7",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Uint(key string) (uint, error) {
	val, ok := v._uint[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "764b042f-f062-4fc5-822a-0ad0df0260ce",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Uint32(key string) (uint32, error) {
	val, ok := v._uint32[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "2f9c0b17-4242-49db-a645-bfee6fd63c3a",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func (v *Values) Uint64(key string) (uint64, error) {
	val, ok := v._uint64[key]

	if !ok {
		return 0, &errors.Object{
			Id:     "7fa8b943-6831-4e51-9b85-db028b422cf5",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Key not found.",
		}
	}

	return val, nil
}

func LoadConfigValues(v *ConfigValues) (*Values, error) {
	if v == nil {
		return nil, nil
	}

	// Bool.
	bools := make(map[string]bool, len(v.Bool))
	for k, v := range v.Bool {
		val, err := strconv.ParseBool(v)
		if err != nil {
			return nil, &errors.Object{
				Id:     "5788d35e-cf30-49a3-bf82-0b0f3531d777",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse bool.",
				Cause:  err.Error(),
			}
		}
		bools[k] = val
	}

	// Bytes.
	bytes := make(map[string][]byte, len(v.Bool))
	for k, v := range v.Bool {
		val, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, &errors.Object{
				Id:     "f8a0b1c4-2d3e-4c5b-9f6d-7a0e5f3f2b8c",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to decode base64.",
				Cause:  err.Error(),
			}
		}
		bytes[k] = val
	}

	// Decimal.
	decimals := make(map[string]decimal.Decimal, len(v.Decimal))
	for k, v := range v.Decimal {
		val, err := decimal.NewFromString(v)
		if err != nil {
			return nil, &errors.Object{
				Id:     "1d20f6ff-68a1-4f80-8154-6f1223c6a477",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse decimal.",
				Cause:  err.Error(),
			}
		}
		decimals[k] = val
	}

	// Float64.
	floats := make(map[string]float64, len(v.Float64))
	for k, v := range v.Float64 {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, &errors.Object{
				Id:     "066f09b9-039e-4fca-9edf-64f328996d8f",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse float64.",
				Cause:  err.Error(),
			}
		}
		floats[k] = val
	}

	// Int.
	ints := make(map[string]int, len(v.Int))
	for k, v := range v.Int {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, &errors.Object{
				Id:     "67b2bc15-d3c4-4eb1-9eff-6c77d1571c28",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse int.",
				Cause:  err.Error(),
			}
		}
		ints[k] = val
	}

	// Int32.
	int32s := make(map[string]int32, len(v.Int32))
	for k, v := range v.Int32 {
		val, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, &errors.Object{
				Id:     "c91bca2b-c342-4549-9bf6-c2190c3da492",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse int32.",
				Cause:  err.Error(),
			}
		}
		int32s[k] = int32(val)
	}

	// Int64.
	int64s := make(map[string]int64, len(v.Int64))
	for k, v := range v.Int64 {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, &errors.Object{
				Id:     "a119420f-9ac2-4c33-ba71-99bb045bebef",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse int64.",
				Cause:  err.Error(),
			}
		}
		int64s[k] = val
	}

	// String.
	strings := make(map[string]string, len(v.String))
	maps.Copy(strings, v.String)

	// Time.
	times := make(map[string]time.Time, len(v.Time))
	for k, v := range v.Time {
		val, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, &errors.Object{
				Id:     "3070ba5c-9376-40dc-9854-f00723447d9e",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse time.",
				Cause:  err.Error(),
			}
		}
		times[k] = val
	}

	// Uint.
	uints := make(map[string]uint, len(v.Uint))
	for k, v := range v.Uint {
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, &errors.Object{
				Id:     "d6474ab3-44fd-4c4d-b8c0-5ddc77a32a5d",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse uint.",
				Cause:  err.Error(),
			}
		}
		uints[k] = uint(val)
	}

	// Uint32.
	uint32s := make(map[string]uint32, len(v.Uint32))
	for k, v := range v.Uint32 {
		val, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, &errors.Object{
				Id:     "7ff2b7fc-c837-4a3c-8602-74c4f9824519",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse uint32.",
				Cause:  err.Error(),
			}
		}
		uint32s[k] = uint32(val)
	}

	// Uint64.
	uint64s := make(map[string]uint64, len(v.Uint64))
	for k, v := range v.Uint64 {
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, &errors.Object{
				Id:     "16f0b0e7-ebc2-4052-81dd-92466c0464f1",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse uint64.",
				Cause:  err.Error(),
			}
		}
		uint64s[k] = val
	}

	return &Values{
		_bool:    bools,
		_bytes:   bytes,
		_decimal: decimals,
		_float64: floats,
		_int:     ints,
		_int32:   int32s,
		_int64:   int64s,
		_string:  strings,
		_time:    times,
		_uint:    uints,
		_uint32:  uint32s,
		_uint64:  uint64s,
	}, nil
}

func NewEmptyValues() *Values {
	return &Values{
		_bool:    make(map[string]bool),
		_bytes:   make(map[string][]byte),
		_decimal: make(map[string]decimal.Decimal),
		_float64: make(map[string]float64),
		_int:     make(map[string]int),
		_int32:   make(map[string]int32),
		_int64:   make(map[string]int64),
		_string:  make(map[string]string),
		_time:    make(map[string]time.Time),
		_uint:    make(map[string]uint),
		_uint32:  make(map[string]uint32),
		_uint64:  make(map[string]uint64),
	}
}
