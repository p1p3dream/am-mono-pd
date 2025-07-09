package token

import (
	"bytes"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
)

type Token interface {
	Id() uuid.UUID
	OrganizationId() uuid.UUID
	TokenType() uint16
	Body() []byte

	// Quota returns the remaining quota of the token.
	// Although Valkey supports int64 for most int operations,
	// it doesn't work well on scripts because the number type
	// only supports 53 bits of precision.
	// Additionally, for this use case we're only interested in
	// positive integers with a high enough size.
	Quota() uint32
}

type tokenIdentifier struct {
	Id             uuid.UUID `cbor:"1,keyasint,omitempty" json:"id,omitempty"`
	OrganizationId uuid.UUID `cbor:"2,keyasint,omitempty" json:"organization_id,omitempty"`
	TokenType      uint16    `cbor:"4,keyasint,omitempty" json:"token_type,omitempty"`
}

// DO NOT change or reuse the cbor tags.
// They are used to serialize and deserialize the struct
// and could lead to data corruption if changed.
type tokenPayload struct {
	Id             uuid.UUID `cbor:"1,keyasint,omitempty" json:"id,omitempty"`
	OrganizationId uuid.UUID `cbor:"2,keyasint,omitempty" json:"organization_id,omitempty"`

	TokenType uint16 `cbor:"4,keyasint,omitempty" json:"token_type,omitempty"`
	Body      []byte `cbor:"5,keyasint,omitempty" json:"body,omitempty"`

	// No need to serialize the quota.
	quota uint32 `cbor:"-" json:"-"`
}

type tokenContainer struct {
	payload tokenPayload
}

func (c *tokenContainer) Id() uuid.UUID {
	return c.payload.Id
}

func (c *tokenContainer) OrganizationId() uuid.UUID {
	return c.payload.OrganizationId
}

func (c *tokenContainer) TokenType() uint16 {
	return c.payload.TokenType
}

func (c *tokenContainer) Body() []byte {
	return c.payload.Body
}

func (c *tokenContainer) Quota() uint32 {
	return c.payload.quota
}

func (c *tokenContainer) MarshalBinary() ([]byte, error) {
	b, err := cbor.Marshal(c.payload)
	if err != nil {
		return nil, &errors.Object{
			Id:     "80705d81-b497-44cd-8123-da4430452a37",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to marshal.",
			Cause:  err.Error(),
		}
	}

	var body []byte
	var header byte

	if len(b) > 1024 {
		header |= 1 << HeaderZstdCompression
		buf := new(bytes.Buffer)

		enc, err := zstd.NewWriter(buf)
		if err != nil {
			return nil, &errors.Object{
				Id:     "2a193d4f-f30c-4a5d-a854-01f1962e93d2",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to create zstd writer.",
				Cause:  err.Error(),
			}
		}

		if _, err := io.Copy(enc, bytes.NewReader(b)); err != nil {
			if err := enc.Close(); err != nil {
				log.Error().
					Err(err).
					Str("id", "acd6d0ab-f820-4783-aaf6-5643c5e73729").
					Msg("Failed to close zstd writer.")
			}

			return nil, &errors.Object{
				Id:     "f155b2f8-0852-4857-a1c1-be5b8b869892",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to write to zstd writer.",
				Cause:  err.Error(),
			}
		}

		if err := enc.Close(); err != nil {
			return nil, &errors.Object{
				Id:     "07061ccc-4aa6-4845-819e-087840cff2e6",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to close zstd writer.",
				Cause:  err.Error(),
			}
		}

		body = buf.Bytes()
	}

	if header == 0 {
		body = b
	}

	payload := append([]byte{header}, body...)

	return payload, nil
}

func (c *tokenContainer) UnmarshalBinary(data []byte) error {
	header := data[0]

	var body []byte

	switch {
	case header == 0:
		body = data[1:]
	case header&(1<<HeaderZstdCompression) != 0:
		dec, err := zstd.NewReader(bytes.NewReader(data[1:]))
		if err != nil {
			return &errors.Object{
				Id:     "e60f94eb-abde-4f3e-a98f-b641ab37c81e",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to create zstd reader.",
				Cause:  err.Error(),
			}
		}

		defer dec.Close()
		buf := new(bytes.Buffer)

		if _, err := io.Copy(buf, dec); err != nil {
			return &errors.Object{
				Id:     "8b6b0ff4-36a6-4e17-9509-df944b8b5140",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to read from zstd reader.",
				Cause:  err.Error(),
			}
		}

		body = buf.Bytes()
	}

	if err := cbor.Unmarshal(body, &c.payload); err != nil {
		return &errors.Object{
			Id:     "e9fcfdfa-e462-4451-80cc-7535d0c20df3",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to unmarshal server session.",
			Cause:  err.Error(),
		}
	}

	return nil
}
