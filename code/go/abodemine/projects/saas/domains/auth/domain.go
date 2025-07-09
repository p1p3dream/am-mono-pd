package auth

import (
	"context"
	"net/http"

	"github.com/fxamacker/cbor/v2"

	"abodemine/domains/arc"
	"abodemine/domains/token"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/projects/api/entities"
	"abodemine/projects/saas/domains/user"
)

type Domain interface {
	Authenticate(ctx context.Context, in *AuthenticateInput) (*AuthenticateOutput, error)

	Load(r *arc.Request) (*LoadOutput, error)
	TokenValidate(r *arc.Request, in *TokenValidateInput) (*TokenValidateOutput, error)
}

type domain struct {
	repository Repository

	ArcDomain   arc.Domain
	TokenDomain token.Domain
	UserDomain  user.Domain
}

type NewDomainInput struct {
	Repository Repository

	ArcDomain   arc.Domain
	TokenDomain token.Domain
	UserDomain  user.Domain
}

func NewDomain(in *NewDomainInput) *domain {
	rep := in.Repository

	if rep == nil {
		rep = &repository{}
	}

	return &domain{
		repository:  rep,
		ArcDomain:   in.ArcDomain,
		TokenDomain: in.TokenDomain,
		UserDomain:  in.UserDomain,
	}
}

const (
	AuthMethodCookie = iota
	AuthMethodToken
)

type AuthenticateInput struct {
	AuthMethod  int
	HttpRequest *http.Request
	Token       string
}

type AuthenticateOutput struct {
	Request           *arc.Request
	TokenExchangeBody *entities.TokenExchangeBody
}

func (dom *domain) Authenticate(ctx context.Context, in *AuthenticateInput) (*AuthenticateOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "097fae72-0459-4dd6-b69d-95d0ca815920",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	var pasetoConfigName string
	var tokenStr string
	var tokenType uint16

	switch in.AuthMethod {
	case AuthMethodCookie:
		if in.HttpRequest == nil {
			return nil, &errors.Object{
				Id:     "022454bb-eb9c-44f6-b438-111b5fe59ff8",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Missing http request.",
			}
		}

		cookie, err := in.HttpRequest.Cookie(consts.CookieAbodeMineSaasWebSession)
		if err != nil {
			return nil, &errors.Object{
				Id:     "5915278c-3ad1-476d-b26e-206037283409",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Missing cookie",
			}
		}

		pasetoConfigName = consts.ConfigKeyPasetoSession
		tokenStr = cookie.Value
		tokenType = token.TypeSaasServerSession
	case AuthMethodToken:
		pasetoConfigName = consts.ConfigKeyPasetoTokenExchange
		tokenStr = in.Token
		tokenType = token.TypeTokenExchange
	default:
		return nil, &errors.Object{
			Id:     "06e2db52-631b-4a6d-8425-91f05f30de1e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing auth.",
		}
	}

	systemRequest, err := dom.ArcDomain.CreateRequest(&arc.CreateRequestInput{
		Context: ctx,
	})
	if err != nil {
		return nil, errors.Forward(err, "5327cce6-fd70-43eb-93ac-03eef1e2eeca")
	}

	systemSession, err := dom.ArcDomain.CreateServerSession(
		systemRequest,
		&arc.CreateServerSessionInput{
			OrganizationId: consts.AbodeMineOrganizationId(),
			UserId:         consts.AbodeMineBotUserId(),
			RoleName:       consts.RoleSystemAuthCheckUser,
			SessionType:    arc.SessionTypeSystem,
			TTL:            1,
			DoNotSave:      true,
		})
	if err != nil {
		return nil, errors.Forward(err, "ef9574c7-89be-4fa4-ac28-bef2aafa7a4d")
	}

	systemRequest.SetSession(systemSession)

	decodeTokenOut, err := dom.TokenDomain.DecodeToken(systemRequest, &token.DecodeTokenInput{
		Value:            tokenStr,
		TokenType:        tokenType,
		PasetoConfigName: pasetoConfigName,
	})
	if err != nil {
		return nil, errors.Forward(err, "d395767a-d638-44e2-90da-18c2b5d4086f")
	}

	switch tokenType {
	case token.TypeSaasServerSession:
		session, err := dom.ArcDomain.SelectServerSession(systemRequest, &arc.SelectServerSessionInput{
			OrganizationId: decodeTokenOut.OrganizationId,
			Id:             decodeTokenOut.Id,
			SessionType:    arc.SessionTypeSaasServer,
		})
		if err != nil {
			return nil, errors.Forward(err, "7cfd8da9-8438-404c-8283-04485d7fc853")
		}

		request, err := dom.ArcDomain.CreateRequest(&arc.CreateRequestInput{
			Context: ctx,
		})
		if err != nil {
			return nil, errors.Forward(err, "ad578635-a525-49f0-ba3e-5f1974f3b852")
		}

		request.SetSession(session)

		return &AuthenticateOutput{
			Request: request,
		}, nil
	case token.TypeTokenExchange:
		selectTokenOut, err := dom.TokenDomain.SelectToken(systemRequest, &token.SelectTokenInput{
			OrganizationId:  decodeTokenOut.OrganizationId,
			Id:              decodeTokenOut.Id,
			TokenType:       token.TypeTokenExchange,
			QuotaDecreaseBy: 1,
		})
		if err != nil {
			return nil, errors.Forward(err, "db298b9e-6f25-4d24-ad9e-6e8012ceaaf7")
		}

		tokenExchangeBody := new(entities.TokenExchangeBody)

		if err := cbor.Unmarshal(selectTokenOut.Body(), tokenExchangeBody); err != nil {
			return nil, &errors.Object{
				Id:     "2e85140b-1333-4c57-951f-c9f15db11ac2",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to unmarshal token body.",
				Cause:  err.Error(),
			}
		}

		return &AuthenticateOutput{
			Request:           systemRequest,
			TokenExchangeBody: tokenExchangeBody,
		}, nil
	}

	return nil, &errors.Object{
		Id:     "3a8a1280-4a86-4f20-9bb2-40f627033759",
		Code:   errors.Code_INTERNAL,
		Detail: "Unsupported token type.",
	}
}

type LoadOutput struct {
	WebSession *arc.WebClientSession `json:"web_session,omitempty"`
}

func (dom *domain) Load(r *arc.Request) (*LoadOutput, error) {
	session := r.Session()

	out := &LoadOutput{
		WebSession: &arc.WebClientSession{
			IsLoading:       true,
			IsAuthenticated: true,
			SessionType:     arc.SessionTypeSaasWebClient,
			Timezone:        session.Timezone(),
		},
	}

	return out, nil
}
