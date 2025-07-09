package val

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/google/uuid"

	"abodemine/lib/errors"
)

func ByteArray16ToSlice(v [16]byte) []byte {
	return v[:]
}

func NewUUID4() (uuid.UUID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, &errors.Object{
			Id:     "eb2e058d-5c81-4814-b03d-eb02d9cf8290",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to generate UUIDv4.",
			Cause:  err.Error(),
		}
	}

	return u, nil
}

func NewUUID7() (uuid.UUID, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, &errors.Object{
			Id:     "f00c9540-5388-4899-9394-5bc2e03713e3",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to generate UUIDv7.",
			Cause:  err.Error(),
		}
	}

	return u, nil
}

func Uint16FromRawBase64(s string) (uint16, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return 0, &errors.Object{
			Id:     "5e162577-e4da-4ff4-afa0-f0087d9dc839",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to decode base64 string.",
			Cause:  err.Error(),
		}
	}

	return binary.BigEndian.Uint16(b), nil
}

func Uint16FromRawUrlSafeBase64(s string) (uint16, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, &errors.Object{
			Id:     "5e869d41-6088-4295-b5ee-88e285eae2b4",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to decode base64 string.",
			Cause:  err.Error(),
		}
	}

	return binary.BigEndian.Uint16(b), nil
}

// Always 3 characters long.
func Uint16ToRawBase64(v uint16) string {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return base64.RawStdEncoding.EncodeToString(buf)
}

// Always 3 characters long.
func Uint16ToRawUrlSafeBase64(v uint16) string {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return base64.RawURLEncoding.EncodeToString(buf)
}

func UUIDFromBytes(b []byte) (uuid.UUID, error) {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil, &errors.Object{
			Id:     "33f3f4ab-e652-4629-8eb3-2aebc823109b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to generate UUID from bytes.",
			Cause:  err.Error(),
		}
	}

	return u, nil
}

func UUIDFromString(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, &errors.Object{
			Id:     "5581ea9a-4ab2-4877-8d29-bfdf4efd0c56",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to generate UUID from string.",
			Cause:  err.Error(),
		}
	}

	return u, nil
}
