package token

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type Domain interface {
	CreateToken(r *arc.Request, in *CreateTokenInput) (*CreateTokenOutput, error)
	SelectToken(r *arc.Request, in *SelectTokenInput) (Token, error)
	DeleteToken(r *arc.Request, in *DeleteTokenInput) error

	EncodeToken(r *arc.Request, in *EncodeTokenInput) (*EncodeTokenOutput, error)
	DecodeToken(r *arc.Request, in *DecodeTokenInput) (*DecodeTokenOutput, error)
}

type domain struct{}

type NewDomainInput struct{}

func NewDomain(in *NewDomainInput) Domain {
	return &domain{}
}

type CreateTokenInput struct {
	OrganizationId uuid.UUID
	TokenType      uint16
	Body           []byte
	Value          any
	Quota          uint32

	TTL time.Duration

	EncodeToken            bool
	EncodeTokenVersion     uint16
	EncodePasetoConfigName string
}

type CreateTokenOutput struct {
	EncodedToken string
	Token        Token
}

func (dom *domain) CreateToken(r *arc.Request, in *CreateTokenInput) (*CreateTokenOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "4bc74386-db1a-4925-b859-8791b777eea4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if _, ok := validTokenTypes[in.TokenType]; !ok {
		return nil, &errors.Object{
			Id:     "d833ed02-ce0f-4d4c-a59b-e6d8d6c0768f",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported TokenType.",
			Meta: map[string]any{
				"TokenType": in.TokenType,
			},
		}
	}

	if in.TTL == 0 {
		return nil, &errors.Object{
			Id:     "e8660186-fe03-40f4-a46f-5621bc20f1be",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "TTL is required.",
		}
	}

	id, err := val.NewUUID4()
	if err != nil {
		return nil, errors.Forward(err, "4d059d57-a2cc-40c8-8654-72b24d1120fb")
	}

	var body []byte

	if in.Value == nil {
		body = in.Body
	} else {
		body, err = cbor.Marshal(in.Value)
		if err != nil {
			return nil, &errors.Object{
				Id:     "bd124fe1-b211-4756-94d8-606591960392",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to marshal token value",
				Cause:  err.Error(),
			}
		}
	}

	tokn := &tokenContainer{tokenPayload{
		Id:             id,
		OrganizationId: in.OrganizationId,
		TokenType:      in.TokenType,
		Body:           body,
		quota:          in.Quota,
	}}

	// Encode token before Redis ops to ensure we won't have
	// useless data there in case of failure.

	var encodedToken string

	if in.EncodeToken {
		encodeTokenOut, err := dom.EncodeToken(r, &EncodeTokenInput{
			Id:             tokn.payload.Id,
			OrganizationId: tokn.payload.OrganizationId,
			TokenType:      tokn.payload.TokenType,
			Expire:         in.TTL,

			Version:          in.EncodeTokenVersion,
			PasetoConfigName: in.EncodePasetoConfigName,
		})
		if err != nil {
			return nil, errors.Forward(err, "7e568cd5-cfc2-48df-8a75-305123c1c62d")
		}

		encodedToken = encodeTokenOut.Value
	}

	payload, err := tokn.MarshalBinary()
	if err != nil {
		return nil, errors.Forward(err, "09ee2536-5e1a-465a-9360-74f4e22c2579")
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeyToken)
	if err != nil {
		return nil, errors.Forward(err, "1000ec7b-3e84-4e4b-a5f1-e0ac38cdeb85")
	}

	createScript, err := r.Dom().SelectValkeyScript("create-token")
	if err != nil {
		return nil, errors.Forward(err, "992359e6-68a4-459f-a14d-ff7d5d4e6e32")
	}

	key := new(strings.Builder)

	if in.OrganizationId != uuid.Nil {
		key.WriteString(fmt.Sprintf(
			"{%s}:",
			base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		))
	}

	key.WriteString(fmt.Sprintf(
		"tokn:%s:%s",
		val.Uint16ToRawBase64(in.TokenType),
		base64.StdEncoding.EncodeToString(id[:])[:22],
	))

	scriptOut := createScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key.String()},
		[]string{
			strconv.Itoa(int(in.TTL.Seconds())),
			string(payload),
			strconv.FormatInt(int64(in.Quota), 10),
		},
	)

	if scriptOut.Error() != nil {
		return nil, &errors.Object{
			Id:     "5897e73e-098e-4710-8dd1-eb012b527e7d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute script.",
			Cause:  scriptOut.Error().Error(),
		}
	}

	scriptOutStr, err := scriptOut.ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "7ffb1cf2-1b35-441d-a569-b981de02284e",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output.",
			Cause:  err.Error(),
		}
	}

	if scriptOutStr != "OK" {
		return nil, &errors.Object{
			Id:     "50df55ec-68c5-43a3-881d-ebb5ecec9177",
			Code:   errors.Code_INTERNAL,
			Detail: "Failed to create session.",
		}
	}

	out := &CreateTokenOutput{
		EncodedToken: encodedToken,
		Token:        tokn,
	}

	return out, nil
}

type SelectTokenInput struct {
	OrganizationId  uuid.UUID
	Id              uuid.UUID
	TokenType       uint16
	QuotaDecreaseBy uint32
	ReturnQuota     bool
}

func (dom *domain) SelectToken(r *arc.Request, in *SelectTokenInput) (Token, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "83823559-7879-49b8-a6c7-2c553bf248bc",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.Id == uuid.Nil {
		return nil, &errors.Object{
			Id:     "a8a80671-7163-4acb-a1b8-357740fffbf4",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Id.",
		}
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeyToken)
	if err != nil {
		return nil, errors.Forward(err, "b06302d4-7991-4b5f-8ae3-7c8a6e86456d")
	}

	selectScript, err := r.Dom().SelectValkeyScript("select-token")
	if err != nil {
		return nil, errors.Forward(err, "e49e7938-5b23-4614-b9e9-f3f94e996481")
	}

	key := new(strings.Builder)

	if in.OrganizationId != uuid.Nil {
		key.WriteString(fmt.Sprintf(
			"{%s}:",
			base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		))
	}

	key.WriteString(fmt.Sprintf(
		"tokn:%s:%s",
		val.Uint16ToRawBase64(in.TokenType),
		base64.StdEncoding.EncodeToString(in.Id[:])[:22],
	))

	scriptOut := selectScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key.String()},
		[]string{
			strconv.FormatInt(int64(in.QuotaDecreaseBy), 10),
			strconv.FormatBool(in.ReturnQuota),
		},
	)

	if scriptOut.Error() != nil {
		switch scriptOut.Error().Error() {
		case "KEY_NOT_FOUND":
			return nil, &errors.Object{
				Id:     "b52821f0-3127-4595-9927-70147f482472",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Token not found.",
			}
		case "INSUFFICIENT_QUOTA":
			return nil, &errors.Object{
				Id:     "2c4dc6dd-6844-4110-b99e-49ba1cf14e5c",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Detail: "Insufficient quota.",
			}
		default:
			return nil, &errors.Object{
				Id:     "4e0395be-ee74-4af3-85bb-17583f9d8d53",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to execute script.",
				Cause:  scriptOut.Error().Error(),
			}
		}
	}

	scriptOutArray, err := scriptOut.ToArray()
	if err != nil {
		return nil, &errors.Object{
			Id:     "aa019848-686d-4e9d-ba56-24693a897bb7",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output.",
			Cause:  err.Error(),
		}
	}

	if len(scriptOutArray) < 2 {
		return nil, &errors.Object{
			Id:     "cdcbbc8c-acf6-4818-8301-3eecd4e4589c",
			Code:   errors.Code_INTERNAL,
			Detail: "Script output length mismatch.",
		}
	}

	payloadStr, err := scriptOutArray[0].ToString()
	if err != nil {
		return nil, &errors.Object{
			Id:     "c2ab919d-0c69-4267-91c1-2629bf406ba5",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse script output[0].",
			Cause:  err.Error(),
		}
	}

	var quota uint32

	if in.ReturnQuota {
		quotaInt64, err := scriptOutArray[1].ToInt64()
		if err != nil {
			return nil, &errors.Object{
				Id:     "188f0129-412d-4797-a39e-8734900d5e9e",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to parse script output[1].",
				Cause:  err.Error(),
			}
		}

		quota = uint32(quotaInt64)
	}

	tokn := &tokenContainer{tokenPayload{quota: quota}}

	if err := tokn.UnmarshalBinary([]byte(payloadStr)); err != nil {
		return nil, errors.Forward(err, "bd85340d-cab5-44bf-ae30-d3591725df24")
	}

	return tokn, nil
}

type DeleteTokenInput struct {
	OrganizationId uuid.UUID
	Id             uuid.UUID
	TokenType      uint16
}

func (dom *domain) DeleteToken(r *arc.Request, in *DeleteTokenInput) error {
	if in == nil {
		return &errors.Object{
			Id:     "641d5822-bd75-4289-b5f2-e8ca88f94363",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.Id == uuid.Nil {
		return &errors.Object{
			Id:     "4f2fb8d2-e499-414c-a87f-687d9a98607d",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Id.",
		}
	}

	valkeyCli, err := r.Dom().SelectValkey(consts.ConfigKeyValkeyToken)
	if err != nil {
		return errors.Forward(err, "fd229ab0-4cc7-49a0-af02-4ab94c394d4b")
	}

	deleteScript, err := r.Dom().SelectValkeyScript("delete-token")
	if err != nil {
		return errors.Forward(err, "50070df8-fc4a-4a5a-9b6e-58f756df59a2")
	}

	key := new(strings.Builder)

	if in.OrganizationId != uuid.Nil {
		key.WriteString(fmt.Sprintf(
			"{%s}:",
			base64.StdEncoding.EncodeToString(in.OrganizationId[:])[:22],
		))
	}

	key.WriteString(fmt.Sprintf(
		"tokn:%s:%s",
		val.Uint16ToRawBase64(in.TokenType),
		base64.StdEncoding.EncodeToString(in.Id[:])[:22],
	))

	scriptOut := deleteScript.Exec(
		r.Context(),
		valkeyCli,
		[]string{key.String()},
		nil,
	)

	if scriptOut.Error() != nil {
		return &errors.Object{
			Id:     "bc9ff9a3-5b12-4905-b8bd-c58e15caf02b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute script.",
			Cause:  scriptOut.Error().Error(),
		}
	}

	return nil
}

type WebClientToken struct {
	DisplayName     string
	IsLoading       bool
	IsAuthenticated bool
	IsOffline       bool
	TokenToken      string
	TokenType       uint16
	Timezone        string
}

type EncodeTokenInput struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
	TokenType      uint16
	Expire         time.Duration

	Version          uint16
	PasetoConfigName string
}

type EncodeTokenOutput struct {
	Value     string
	ExpiresAt time.Time
}

func (dom *domain) EncodeToken(r *arc.Request, in *EncodeTokenInput) (*EncodeTokenOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "3c06d800-5ff5-4db2-bd32-34eebb94ff17",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.Expire == 0 {
		return nil, &errors.Object{
			Id:     "447b138e-ce27-4926-b820-f930533ee1d5",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Expire is required.",
		}
	}

	encoded := new(strings.Builder)

	// Insert our version header.
	encoded.WriteString(val.Uint16ToRawUrlSafeBase64(in.Version))

	var expiresAt time.Time

	switch in.Version {
	case VersionPasetoV4Local:
		pasetoConfig, err := r.Dom().SelectPaseto(in.PasetoConfigName)
		if err != nil {
			return nil, errors.Forward(err, "c11a9c9b-91cc-4b0a-98b3-44692eb298e1")
		}

		now := time.Now()
		tokenTypeStr := val.Uint16ToRawBase64(in.TokenType)

		token := paseto.NewToken()
		token.SetIssuedAt(now)
		token.SetNotBefore(now)
		token.SetExpiration(now.Add(in.Expire))
		token.SetString(
			val.Uint16ToRawBase64(TypeKey),
			tokenTypeStr,
		)

		body := &tokenIdentifier{
			Id:             in.Id,
			OrganizationId: in.OrganizationId,
			TokenType:      in.TokenType,
		}

		if err := token.Set(tokenTypeStr, body); err != nil {
			return nil, &errors.Object{
				Id:     "51018313-7dcf-41fa-a5c9-75e0664ca20c",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to set token key.",
				Meta: map[string]any{
					"key": tokenTypeStr,
				},
			}
		}

		// Remove paseto header.
		encoded.WriteString(strings.TrimPrefix(
			token.V4Encrypt(pasetoConfig.V4SymmetricKey, nil),
			"v4.local.",
		))
	default:
		return nil, &errors.Object{
			Id:     "fd6964fd-5170-469c-a2c4-698c212b3ffe",
			Code:   errors.Code_INTERNAL,
			Detail: "Unsupported token version.",
		}
	}

	out := &EncodeTokenOutput{
		Value:     encoded.String(),
		ExpiresAt: expiresAt,
	}

	return out, nil
}

type DecodeTokenInput struct {
	Value            string
	TokenType        uint16
	PasetoConfigName string
}

type DecodeTokenOutput struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
	TokenType      uint16

	Version uint16
}

func (dom *domain) DecodeToken(r *arc.Request, in *DecodeTokenInput) (*DecodeTokenOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "0ae73580-53af-47fa-8860-5f3ff99ebeaa",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if len(in.Value) < 3 {
		return nil, &errors.Object{
			Id:     "02939888-9ac3-47ca-ac0e-7f934da25f93",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing value.",
		}
	}

	version, err := val.Uint16FromRawUrlSafeBase64(in.Value[:3])
	if err != nil {
		return nil, &errors.Object{
			Id:     "42497ebf-76bf-42d0-a78a-ca5d3f041a93",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid version.",
			Cause:  err.Error(),
		}
	}

	switch version {
	case VersionPasetoV4Local:
		pasetoConfig, err := r.Dom().SelectPaseto(in.PasetoConfigName)
		if err != nil {
			return nil, errors.Forward(err, "5cfdf1b7-4cd0-423c-b54d-7d12b6012cea")
		}

		tokenStr := "v4.local." + in.Value[3:]

		token, err := pasetoConfig.Parser.ParseV4Local(pasetoConfig.V4SymmetricKey, tokenStr, nil)
		if err != nil {
			return nil, &errors.Object{
				Id:     "c126b63e-bec6-4ea9-9ead-800fb556577a",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid auth token.",
				Cause:  err.Error(),
			}
		}

		tokenTypeStr, err := token.GetString(val.Uint16ToRawBase64(TypeKey))
		if err != nil {
			return nil, &errors.Object{
				Id:     "d0c7ae33-e8dd-4603-88ce-d83a821e4b36",
				Code:   errors.Code_INTERNAL,
				Detail: "Token type not found.",
				Cause:  err.Error(),
			}
		}

		tokenType, err := val.Uint16FromRawBase64(tokenTypeStr)
		if err != nil {
			return nil, errors.Forward(err, "739822b7-660a-478d-a9ca-0feca6f29955")
		}

		if tokenType != in.TokenType {
			return nil, &errors.Object{
				Id:     "a6b3c4c6-5f1f-4f7a-bb8f-2f6d4e9b3f6b",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Token type mismatch.",
			}
		}

		body := new(tokenIdentifier)

		if err := token.Get(tokenTypeStr, body); err != nil {
			return nil, &errors.Object{
				Id:     "964feb15-034e-46f6-b155-6fb6e98322fb",
				Code:   errors.Code_INTERNAL,
				Detail: "Token body not found.",
				Cause:  err.Error(),
			}
		}

		out := &DecodeTokenOutput{
			Id:             body.Id,
			OrganizationId: body.OrganizationId,
			TokenType:      body.TokenType,
			Version:        version,
		}

		return out, nil
	}

	return nil, &errors.Object{
		Id:     "8f7ecbb6-dd86-4844-a23a-f2a473d209be",
		Code:   errors.Code_INVALID_ARGUMENT,
		Detail: "Unsupported token version.",
	}
}
