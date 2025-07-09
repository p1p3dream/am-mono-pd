package arc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type SessionHeader byte

const (
	SessionHeaderZstdCompression SessionHeader = 0
)

const (
	// System.
	SessionTypeSystem uint16 = 1

	// Api.
	SessionTypeApiServer uint16 = 100

	// Saas.
	SessionTypeSaasServer    uint16 = 200
	SessionTypeSaasWebClient uint16 = 201

	// ESO.
	SessionTypeESOWhiteLabelServer    uint16 = 300
	SessionTypeESOWhiteLabelWebClient uint16 = 301
)

var validSessionTypes = map[uint16]struct{}{
	SessionTypeSystem:        {},
	SessionTypeApiServer:     {},
	SessionTypeSaasServer:    {},
	SessionTypeSaasWebClient: {},
}

type ServerSession interface {
	Flags() []string
	Id() uuid.UUID
	KeyHash() string
	KeyId() uuid.UUID
	KeyType() int16
	OrganizationId() uuid.UUID
	RoleName() string
	SessionType() uint16
	Timezone() string
	Username() string
	UserId() uuid.UUID
	Valid() bool

	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
}

// DO NOT change or reuse the cbor tags.
// They are used to serialize and deserialize the struct
// and could lead to data corruption if changed.
type serverSessionPayload struct {
	Id uuid.UUID `cbor:"1,keyasint,omitempty"`

	OrganizationId uuid.UUID `cbor:"2,keyasint,omitempty"`
	UserId         uuid.UUID `cbor:"3,keyasint,omitempty"`
	KeyId          uuid.UUID `cbor:"4,keyasint,omitempty"`

	RoleName    string `cbor:"5,keyasint,omitempty"`
	SessionType uint16 `cbor:"6,keyasint,omitempty"`
	Timezone    string `cbor:"7,keyasint,omitempty"`
	Username    string `cbor:"8,keyasint,omitempty"`

	Flags []string `cbor:"9,keyasint,omitempty"`

	keyHash string
	keyType int16
}

// We use a container to prevent recursive serialization of the payload.
type serverSessionContainer struct {
	payload serverSessionPayload
}

type NewServerSessionFromBytesInput struct {
	Data    []byte
	KeyHash string
	KeyType int16
}

func NewServerSessionFromBytes(in *NewServerSessionFromBytesInput) (*serverSessionContainer, error) {
	sess := new(serverSessionContainer)

	if err := sess.UnmarshalBinary(in.Data); err != nil {
		return nil, errors.Forward(err, "b82bcfb4-f819-429a-b0fd-3dc13b96c39a")
	}

	sess.payload.keyHash = in.KeyHash
	sess.payload.keyType = in.KeyType

	return sess, nil
}

func (c *serverSessionContainer) Flags() []string {
	return c.payload.Flags
}

func (c *serverSessionContainer) Id() uuid.UUID {
	return c.payload.Id
}

func (c *serverSessionContainer) KeyHash() string {
	return c.payload.keyHash
}

func (c *serverSessionContainer) KeyId() uuid.UUID {
	return c.payload.KeyId
}

func (c *serverSessionContainer) KeyType() int16 {
	return c.payload.keyType
}

func (c *serverSessionContainer) OrganizationId() uuid.UUID {
	return c.payload.OrganizationId
}

func (c *serverSessionContainer) RoleName() string {
	return c.payload.RoleName
}

func (c *serverSessionContainer) SessionType() uint16 {
	return c.payload.SessionType
}

func (c *serverSessionContainer) Timezone() string {
	return c.payload.Timezone
}

func (c *serverSessionContainer) UserId() uuid.UUID {
	return c.payload.UserId
}

func (c *serverSessionContainer) Username() string {
	return c.payload.Username
}

func (c *serverSessionContainer) Valid() bool {
	return c.payload.Id != uuid.Nil
}

func (c *serverSessionContainer) MarshalBinary() ([]byte, error) {
	b, err := cbor.Marshal(c.payload)
	if err != nil {
		return nil, &errors.Object{
			Id:     "4adefb6c-3587-400e-9bcc-99951841e822",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to marshal.",
			Cause:  err.Error(),
		}
	}

	var body []byte
	var header byte

	if len(b) > 1024 {
		header |= 1 << SessionHeaderZstdCompression
		buf := new(bytes.Buffer)

		enc, err := zstd.NewWriter(buf)
		if err != nil {
			return nil, &errors.Object{
				Id:     "a9890862-9947-4e21-95da-77b5c69d3d07",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to create zstd writer.",
				Cause:  err.Error(),
			}
		}

		if _, err := io.Copy(enc, bytes.NewReader(b)); err != nil {
			if err := enc.Close(); err != nil {
				log.Error().
					Err(err).
					Str("id", "e1faf7f9-b63e-47a0-b0bf-bc8fb2895683").
					Msg("Failed to close zstd writer.")
			}

			return nil, &errors.Object{
				Id:     "db633267-d38a-4117-8930-68e80813ed78",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to write to zstd writer.",
				Cause:  err.Error(),
			}
		}

		if err := enc.Close(); err != nil {
			return nil, &errors.Object{
				Id:     "debda895-e4c1-461d-89e5-982e91d6fd15",
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

func (c *serverSessionContainer) UnmarshalBinary(data []byte) error {
	header := data[0]

	var body []byte

	switch {
	case header == 0:
		body = data[1:]
	case header&(1<<SessionHeaderZstdCompression) != 0:
		dec, err := zstd.NewReader(bytes.NewReader(data[1:]))
		if err != nil {
			return &errors.Object{
				Id:     "5ac8ad38-7693-4558-b189-b4e6d730d97c",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to create zstd reader.",
				Cause:  err.Error(),
			}
		}

		defer dec.Close()
		buf := new(bytes.Buffer)

		if _, err := io.Copy(buf, dec); err != nil {
			return &errors.Object{
				Id:     "63a5c819-7b66-4512-a085-390e6a572f1a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to read from zstd reader.",
				Cause:  err.Error(),
			}
		}

		body = buf.Bytes()
	}

	if err := cbor.Unmarshal(body, &c.payload); err != nil {
		return &errors.Object{
			Id:     "c7a0b695-9739-47a1-b6ef-49f41a316563",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to unmarshal server session.",
			Cause:  err.Error(),
		}
	}

	return nil
}

type CreateServerSessionInput struct {
	OrganizationId uuid.UUID
	UserId         uuid.UUID
	KeyHash        string
	KeyId          uuid.UUID
	KeyType        int16
	RoleName       string
	SessionType    uint16
	Timezone       string
	Username       string
	TTL            time.Duration
	Flags          []string

	DoNotSave bool
}

func (dom *domain) CreateServerSession(r *Request, in *CreateServerSessionInput) (ServerSession, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "51e73ab3-7ec5-4f93-82db-314f38f7deb7",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.OrganizationId == uuid.Nil {
		return nil, &errors.Object{
			Id:     "c2a27e96-04df-40d8-818c-d8febdb3fb48",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing OrganizationId.",
		}
	}

	if in.UserId == uuid.Nil && in.KeyId == uuid.Nil {
		return nil, &errors.Object{
			Id:     "84d2ad91-2b61-4823-a683-114520dfa7c4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "One of UserId or ApiKeyId must be provided.",
		}
	}

	if _, ok := validSessionTypes[in.SessionType]; !ok {
		return nil, &errors.Object{
			Id:     "a18fdfa3-8c83-4f0c-9436-8455f42459ba",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported SessionType.",
			Meta: map[string]any{
				"SessionType": in.SessionType,
			},
		}
	}

	if in.TTL == 0 {
		return nil, &errors.Object{
			Id:     "961054c6-9364-40c3-ba03-f9062fd4c27f",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "TTL is required.",
		}
	}

	id, err := val.NewUUID4()
	if err != nil {
		return nil, errors.Forward(err, "14496afb-3c34-4691-b196-26acc229ae93")
	}

	sess := &serverSessionContainer{
		payload: serverSessionPayload{
			Id:             id,
			OrganizationId: in.OrganizationId,
			UserId:         in.UserId,
			KeyId:          in.KeyId,
			RoleName:       in.RoleName,
			SessionType:    in.SessionType,
			Timezone:       in.Timezone,
			Username:       in.Username,
			Flags:          in.Flags,
			keyHash:        in.KeyHash,
			keyType:        in.KeyType,
		},
	}

	if in.DoNotSave {
		return sess, nil
	}

	payload, err := sess.MarshalBinary()
	if err != nil {
		return nil, errors.Forward(err, "98f33f6e-05b5-44f8-8a78-8cc8b2f66cd0")
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeySession)
	if err != nil {
		return nil, errors.Forward(err, "876afded-8403-47c0-91f2-33316e8ba730")
	}

	createScript, err := r.Dom().SelectValkeyScript("create-session")
	if err != nil {
		return nil, errors.Forward(err, "d722b1fe-1054-40c2-bcb3-4e8be2dbfb6f")
	}

	key := fmt.Sprintf(
		"{%s}:sess:%s:%s",
		base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		val.Uint16ToRawBase64(in.SessionType),
		base64.StdEncoding.EncodeToString(id[:])[:22],
	)

	scriptOut := createScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key},
		[]string{
			strconv.Itoa(int(in.TTL.Seconds())),
			string(payload),
		},
	)

	if scriptOut.Error() != nil {
		return nil, &errors.Object{
			Id:     "1ed76331-0390-4cc7-b3d0-66ff32f04869",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute script.",
			Cause:  scriptOut.Error().Error(),
		}
	}

	scriptOutStr, err := scriptOut.ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "9027bdf0-b9b9-4c29-a7f8-88acca8df98d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output.",
			Cause:  err.Error(),
		}
	}

	if scriptOutStr != "OK" {
		return nil, &errors.Object{
			Id:     "365a1cd8-1b55-48fe-aa6f-a531361593ca",
			Code:   errors.Code_INTERNAL,
			Detail: "Failed to create session.",
		}
	}

	return sess, nil
}

type SelectServerSessionInput struct {
	OrganizationId uuid.UUID
	Id             uuid.UUID
	SessionType    uint16
}

func (dom *domain) SelectServerSession(r *Request, in *SelectServerSessionInput) (ServerSession, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "1f1a2c75-8088-4f7b-9a60-c9ad1afc9306",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.OrganizationId == uuid.Nil {
		return nil, &errors.Object{
			Id:     "5603dbee-573a-4362-b58b-88e54ea61d99",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing OrganizationId.",
		}
	}

	if in.Id == uuid.Nil {
		return nil, &errors.Object{
			Id:     "7f5ee5ce-afeb-4a24-827e-2c6ecba61e1c",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Id.",
		}
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeySession)
	if err != nil {
		return nil, errors.Forward(err, "31d29d7b-f772-4ba0-b390-adfb3f96fd44")
	}

	selectSessionScript, err := r.Dom().SelectValkeyScript("select-session")
	if err != nil {
		return nil, errors.Forward(err, "effb77d0-d491-427a-857e-ade0c5353c1b")
	}

	key := fmt.Sprintf(
		"{%s}:sess:%s:%s",
		base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		val.Uint16ToRawBase64(in.SessionType),
		base64.StdEncoding.EncodeToString(in.Id[:])[:22],
	)

	scriptOut := selectSessionScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key},
		nil,
	)

	if scriptOut.Error() != nil {
		switch scriptOut.Error().Error() {
		case "KEY_NOT_FOUND":
			return nil, &errors.Object{
				Id:     "e3e4ded9-9b1a-483c-a23b-60e743f95ab5",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Session not found.",
			}
		default:
			return nil, &errors.Object{
				Id:     "a3f082a7-a3a2-4a6b-b01b-cbccbf4c9cb3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to execute script.",
				Cause:  scriptOut.Error().Error(),
			}
		}
	}

	scriptOutArray, err := scriptOut.ToArray()
	if err != nil {
		return nil, &errors.Object{
			Id:     "aeec2365-abb7-4c12-95b0-8a46025cbb23",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output.",
			Cause:  err.Error(),
		}
	}

	if len(scriptOutArray) < 1 {
		return nil, &errors.Object{
			Id:     "00e05b52-95c9-44eb-a7f1-f2314410f204",
			Code:   errors.Code_INTERNAL,
			Detail: "Script output length mismatch.",
		}
	}

	payloadStr, err := scriptOutArray[0].ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "4fc36f39-8ef3-4541-bb7c-16c7a634ee32",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output[0].",
			Cause:  err.Error(),
		}
	}

	sess := new(serverSessionContainer)

	if err := sess.UnmarshalBinary([]byte(payloadStr)); err != nil {
		return nil, errors.Forward(err, "ca9e4c63-8fc3-487f-bcaf-1bc5c9bf8992")
	}

	return sess, nil
}

type DeleteServerSessionInput struct {
	OrganizationId uuid.UUID
	Id             uuid.UUID
	SessionType    uint16
}

func (dom *domain) DeleteServerSession(r *Request, in *DeleteServerSessionInput) error {
	if in == nil {
		return &errors.Object{
			Id:     "c2976214-464e-4f96-98c8-eba4394fdeb4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.OrganizationId == uuid.Nil {
		return &errors.Object{
			Id:     "95145a17-80fe-4deb-a4ae-26425d34d7c0",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing OrganizationId.",
		}
	}

	if in.Id == uuid.Nil {
		return &errors.Object{
			Id:     "0dad30a3-4c5e-404a-a000-098f8b293a3e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Id.",
		}
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeySession)
	if err != nil {
		return errors.Forward(err, "e95c1e72-0e96-462d-ac20-535b3734fa6e")
	}

	selectSessionScript, err := r.Dom().SelectValkeyScript("delete-session")
	if err != nil {
		return errors.Forward(err, "5916ebe2-6e9c-4005-8a91-5a850b23265b")
	}

	key := fmt.Sprintf(
		"{%s}:sess:%s:%s",
		base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		val.Uint16ToRawBase64(in.SessionType),
		base64.StdEncoding.EncodeToString(in.Id[:])[:22],
	)

	scriptOut := selectSessionScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key},
		nil,
	)

	if scriptOut.Error() != nil {
		return &errors.Object{
			Id:     "a54cb4c3-3926-4bc9-b55b-e742ae5963af",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute script.",
			Cause:  scriptOut.Error().Error(),
		}
	}

	return nil
}

type WebClientSession struct {
	DisplayName     string `json:"display_name,omitempty"`
	IsLoading       bool   `json:"is_loading,omitempty"`
	IsAuthenticated bool   `json:"is_authenticated,omitempty"`
	IsOffline       bool   `json:"is_offline,omitempty"`
	SessionToken    string `json:"session_token,omitempty"`
	SessionType     uint16 `json:"session_type,omitempty"`
	Timezone        string `json:"timezone,omitempty"`
}
