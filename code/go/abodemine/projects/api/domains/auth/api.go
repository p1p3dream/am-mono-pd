package auth

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
)

type ApiKeyType int16

const (
	ApiKeyTypeLegacy ApiKeyType = 100
)

// Base64 returns the base64 representation of ApiKeyType without padding.
func (t ApiKeyType) Base64() string {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(t))
	return base64.RawStdEncoding.EncodeToString(buf)
}

var validApiKeyTypes = map[ApiKeyType]struct{}{
	ApiKeyTypeLegacy: {},
}

type ApiKeyStatus int16

const (
	ApiKeyStatusActive  ApiKeyStatus = 100
	ApiKeyStatusExpired ApiKeyStatus = 200
	ApiKeyStatusRevoked ApiKeyStatus = 300
)

type ApiKey struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]any

	OrganizationId uuid.UUID
	UserId         uuid.UUID
	RoleId         uuid.UUID
	// RoleName here to prevent a database lookup on api key validation.
	RoleName    string
	KeyType     ApiKeyType
	KeyStatus   ApiKeyStatus
	ExpiresAt   time.Time
	LastUsedAt  time.Time
	RevokedAt   time.Time
	RevokedBy   uuid.UUID
	KeyHash     string
	Name        string
	Description string
}

type SelectApiSessionInput struct {
	OrganizationId uuid.UUID
	KeyType        ApiKeyType
	KeyHash        string
}

type SelectApiSessionOutput struct {
	Invalid        bool
	QuotaExhausted string
	Session        arc.ServerSession
}

func (dom *domain) SelectApiSession(r *arc.Request, in *SelectApiSessionInput) (*SelectApiSessionOutput, error) {
	if err := r.CasbinEnforce(
		consts.ConfigKeyCasbinApiDefault,
		"/auth/session",
		"read",
	); err != nil {
		return nil, errors.Forward(err, "89676e26-b66b-46e5-86ef-601338babc67")
	}

	if in == nil {
		return nil, &errors.Object{
			Id:     "9ad2907f-63e6-4fc6-8549-1afd79ca5739",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.KeyHash == "" {
		return nil, &errors.Object{
			Id:     "3c8bc274-21c0-4524-bd59-7d49f6beb627",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing key hash.",
		}
	}

	if _, ok := validApiKeyTypes[in.KeyType]; !ok {
		return nil, &errors.Object{
			Id:     "2a57d955-4083-4f2a-b3e0-7adffd1940df",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported ApiKeyType.",
			Meta: map[string]any{
				"ApiKeyType": in.KeyType,
			},
		}
	}

	valkeyCli, err := dom.ArcDomain.SelectValkey(consts.ConfigKeyValkeyApi)
	if err != nil {
		return nil, errors.Forward(err, "5ee26f89-ea29-41be-891b-f93924fdd9c8")
	}

	selectScript, err := dom.ArcDomain.SelectValkeyScript("select-api-session")
	if err != nil {
		return nil, errors.Forward(err, "ba8b0b15-464b-4f9c-ac36-a94a580fdc02")
	}

	key := fmt.Sprintf(
		"apik:%s:%s",
		in.KeyType.Base64(),
		in.KeyHash,
	)

	scriptOut := selectScript.Exec(
		context.Background(),
		valkeyCli,
		[]string{key},
		nil,
	)

	if scriptOut.Error() != nil {
		switch scriptOut.Error().Error() {
		case "KEY_NOT_FOUND":
			return nil, &errors.Object{
				Id:     "c0a66e0c-025a-4426-9419-b45195c6945a",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Token not found.",
			}
		case "INSUFFICIENT_QUOTA":
			return nil, &errors.Object{
				Id:     "a9aeef93-b3da-472f-befe-245eeec795c4",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Detail: "Insufficient quota.",
			}
		default:
			return nil, &errors.Object{
				Id:     "e5c5fa0a-f9ae-41f8-a045-14d7db3addc3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to execute script.",
				Cause:  scriptOut.Error().Error(),
			}
		}
	}

	scriptOutArray, err := scriptOut.ToArray()
	if err != nil {
		return nil, &errors.Object{
			Id:     "8ed19877-a746-40be-adb0-f52725e8f365",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output.",
			Cause:  err.Error(),
		}
	}

	if len(scriptOutArray) < 3 {
		return nil, &errors.Object{
			Id:     "93aee2ec-887a-49ee-9be7-07594b0e675e",
			Code:   errors.Code_INTERNAL,
			Detail: "Script output length mismatch.",
		}
	}

	invalidStr, err := scriptOutArray[0].ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "05396b5d-6b66-4f00-8621-f6e272f8d0be",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output[0].",
			Cause:  err.Error(),
		}
	}

	if invalidStr == "true" {
		return &SelectApiSessionOutput{
			Invalid: true,
		}, nil
	}

	quotaExhaustedStr, err := scriptOutArray[1].ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "9b313bcb-7a51-4cc9-b6ff-702ada9e607d",
			Code:   errors.Code_RESOURCE_EXHAUSTED,
			Detail: "Failed to parse script output[1].",
			Cause:  err.Error(),
		}
	}

	if quotaExhaustedStr != "" {
		return &SelectApiSessionOutput{
			QuotaExhausted: quotaExhaustedStr,
		}, nil
	}

	sessionStr, err := scriptOutArray[2].ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "90e5955e-dbdd-4714-adc8-f3fad60f65f2",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output[2].",
			Cause:  err.Error(),
		}
	}

	session, err := arc.NewServerSessionFromBytes(&arc.NewServerSessionFromBytesInput{
		Data:    []byte(sessionStr),
		KeyHash: in.KeyHash,
		KeyType: int16(in.KeyType),
	})
	if err != nil {
		return nil, errors.Forward(err, "cde6075a-f7a3-442d-9633-2828ab92b8a1")
	}

	out := &SelectApiSessionOutput{
		Session: session,
	}

	return out, nil
}

type UpdateApiSessionInput struct {
	Invalid        bool
	QuotaExhausted string
	KeyType        ApiKeyType
	KeyHash        string
	Session        arc.ServerSession
	TTL            int32
}

type UpdateApiSessionOutput struct{}

func (dom *domain) UpdateApiSession(r *arc.Request, in *UpdateApiSessionInput) (*UpdateApiSessionOutput, error) {
	if err := r.CasbinEnforce(
		consts.ConfigKeyCasbinApiDefault,
		"/auth/session",
		"write",
	); err != nil {
		return nil, errors.Forward(err, "8d9d4cb2-39bd-4746-b4c8-dde573b2f0dd")
	}

	if in == nil {
		return nil, &errors.Object{
			Id:     "d54f30a4-7063-4f6b-aee0-d3675de4b10b",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.KeyHash == "" {
		return nil, &errors.Object{
			Id:     "bfa0fccc-6d5e-4ebf-aeff-656338fd781a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing key hash.",
		}
	}

	if _, ok := validApiKeyTypes[in.KeyType]; !ok {
		return nil, &errors.Object{
			Id:     "c3d4617b-345d-42d7-b85a-9a8b6c5630a5",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported ApiKeyType.",
			Meta: map[string]any{
				"ApiKeyType": in.KeyType,
			},
		}
	}

	valkeyCli, err := dom.ArcDomain.SelectValkey(consts.ConfigKeyValkeyApi)
	if err != nil {
		return nil, errors.Forward(err, "0ef0d125-c1aa-45f0-9b18-c974b9146f4d")
	}

	updateScript, err := dom.ArcDomain.SelectValkeyScript("update-api-session")
	if err != nil {
		return nil, errors.Forward(err, "832d8693-8c6b-45dc-be7d-431635559fcd")
	}

	key := fmt.Sprintf(
		"apik:%s:%s",
		in.KeyType.Base64(),
		in.KeyHash,
	)

	var sessionBytes []byte

	if !in.Invalid && in.QuotaExhausted == "" {
		sessionBytes, err = in.Session.MarshalBinary()
		if err != nil {
			return nil, errors.Forward(err, "09cc4c8f-3764-45fe-8798-90d423479109")
		}
	}

	scriptOut := updateScript.Exec(
		context.Background(),
		valkeyCli,
		[]string{key},
		[]string{
			strconv.FormatBool(in.Invalid),
			in.QuotaExhausted,
			string(sessionBytes),
			strconv.FormatInt(int64(in.TTL), 10),
		},
	)

	if scriptOut.Error() != nil {
		return nil, &errors.Object{
			Id:     "51e5e4f0-450c-42f3-8b34-1e7b8e8c4fc0",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute script.",
			Cause:  scriptOut.Error().Error(),
		}
	}

	out := &UpdateApiSessionOutput{}

	return out, nil
}
