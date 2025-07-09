package must

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"math"
	"strconv"
)

func GenRandomInt32() int32 {
	r := GenRandomUint32()

	if r > math.MaxInt32 {
		// Should be a fair distribution.
		return -int32(r / 2)
	}

	return int32(r)
}

func GenRandomInt64() int64 {
	r := GenRandomUint64()

	if r > math.MaxInt64 {
		// Should be a fair distribution.
		return -int64(r / 2)
	}

	return int64(r)
}

func GenRandomInt32Range(begin, end int32) int32 {
	if begin >= end {
		panic("begin must be less than end")
	}

	// Calculate the range size
	rangeSize := uint32(end - begin)

	// Generate a random int32 and take modulo to fit within the range.
	// This gives a value between 0 and rangeSize-1.
	r := GenRandomUint32() % rangeSize

	// Add the beginning offset to get the final result in the desired range.
	return begin + int32(r)
}

func GenRandomInt64Range(begin, end int64) int64 {
	if begin >= end {
		panic("begin must be less than end")
	}

	// Calculate the range size
	rangeSize := uint64(end - begin)

	// Generate a random int64 and take modulo to fit within the range.
	// This gives a value between 0 and rangeSize-1.
	r := GenRandomUint64() % rangeSize

	// Add the beginning offset to get the final result in the desired range.
	return begin + int64(r)
}

func GenRandomUint32() uint32 {
	var bytes [4]byte

	_, err := rand.Read(bytes[:])
	if err != nil {
		panic(err)
	}

	return binary.BigEndian.Uint32(bytes[:])
}

func GenRandomUint64() uint64 {
	var bytes [8]byte

	_, err := rand.Read(bytes[:])
	if err != nil {
		panic(err)
	}

	return binary.BigEndian.Uint64(bytes[:])
}

func MarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func MarshalJSONIndent(v any, prefix, indent string) []byte {
	b, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}

	return b
}

func ParseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return v
}
