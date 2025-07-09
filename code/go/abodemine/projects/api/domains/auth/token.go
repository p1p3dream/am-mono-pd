package auth

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/domains/token"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/api/entities"
)

type TokenExchangeInput struct {
	ClientId    uuid.UUID
	Expire      int
	ExternalId  string
	RedirectUri string
}

type TokenExchangeOutput struct {
	Token string
}

func (dom *domain) TokenExchange(r *arc.Request, in *TokenExchangeInput) (*TokenExchangeOutput, error) {
	if err := r.CasbinEnforce(
		consts.ConfigKeyCasbinApiDefault,
		"/auth/token/exchange",
		"write",
	); err != nil {
		return nil, errors.Forward(err, "7daf3d95-8211-44b0-917b-64f0f667f2ce")
	}

	if in == nil {
		return nil, &errors.Object{
			Id:     "f24a26f5-5e27-493b-a9f7-9e65e22fa601",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if in.ClientId == uuid.Nil {
		return nil, &errors.Object{
			Id:     "e4f9fc98-39eb-47f8-88d4-bac53610ee5a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing client_id.",
		}
	}

	if in.Expire < 0 {
		return nil, &errors.Object{
			Id:     "eb205831-fde3-4431-ac52-60194a3daa9b",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Expire must be nonnegative.",
		}
	}

	maxExpire, err := r.Dom().SelectDuration(consts.ConfigKeyDurationSaasWhitelabelSessionMax)
	if err != nil {
		return nil, errors.Forward(err, "06b39bf3-49d8-4b1b-af6b-adb40bbb71df")
	}

	if in.Expire > int(maxExpire.Seconds()) {
		return nil, &errors.Object{
			Id:     "bfe399c0-e208-44d8-8a2d-d949d3075aa3",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: fmt.Sprintf("Expire exceeds maximum of %ds.", int(maxExpire.Seconds())),
		}
	}

	externalId := strings.TrimSpace(in.ExternalId)

	if externalId == "" {
		return nil, &errors.Object{
			Id:     "ffbf320d-2633-42a0-8d3a-16593ad68837",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing external_id.",
		}
	}

	selectClientRecordOut, err := dom.repository.SelectClientRecord(r, &SelectClientRecordInput{
		Id:             in.ClientId,
		OrganizationId: r.Session().OrganizationId(),
	})
	if err != nil {
		return nil, errors.Forward(err, "78cda612-28a0-41e9-a506-5092e855f904")
	}

	client := selectClientRecordOut.Record

	if client == nil {
		return nil, &errors.Object{
			Id:     "58cbecaa-23cd-48b1-9cfd-e5e75584c728",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Client not found.",
		}
	}

	if in.RedirectUri == "" && client.RedirectUri == "" {
		return nil, &errors.Object{
			Id:     "ab919728-e9df-4b9f-b6a8-736b53cb2628",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing redirect_uri.",
		}
	}

	defaultExpire, err := r.Dom().SelectDuration(consts.ConfigKeyDurationSaasWhitelabelSessionDefault)
	if err != nil {
		return nil, errors.Forward(err, "6156e30d-d3aa-4279-89a9-7789475b6b35")
	}

	tokenExchangeBody := &entities.TokenExchangeBody{
		OrganizationId: r.Session().OrganizationId(),
		Expire:         val.Ternary(in.Expire == 0, int(defaultExpire.Seconds()), in.Expire),
		ExternalId:     externalId,
		RedirectUri:    val.Coalesce(in.RedirectUri, client.RedirectUri),
	}

	tokenExchangeTtl, err := r.Dom().SelectDuration(consts.ConfigKeyDurationApiTokenExchangeTtl)
	if err != nil {
		return nil, errors.Forward(err, "eda9bae6-db45-407e-b35d-fcd68d35c69e")
	}

	createOut, err := dom.TokenDomain.CreateToken(r, &token.CreateTokenInput{
		OrganizationId:         r.Session().OrganizationId(),
		TokenType:              token.TypeTokenExchange,
		Value:                  tokenExchangeBody,
		Quota:                  1,
		TTL:                    tokenExchangeTtl,
		EncodeToken:            true,
		EncodeTokenVersion:     token.VersionPasetoV4Local,
		EncodePasetoConfigName: consts.ConfigKeyPasetoTokenExchange,
	})
	if err != nil {
		return nil, errors.Forward(err, "11496789-4942-4d6c-8787-9a7d970e2a1c")
	}

	out := &TokenExchangeOutput{
		Token: createOut.EncodedToken,
	}

	return out, nil
}
