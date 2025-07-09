package auth

import (
	"net/url"
	"time"

	"abodemine/domains/arc"
	"abodemine/domains/token"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/projects/api/entities"
	"abodemine/projects/saas/domains/user"
)

type TokenValidateInput struct {
	DownstreamHost    string
	TokenExchangeBody *entities.TokenExchangeBody
}

type TokenValidateOutput struct {
	SessionToken string
	RedirectUri  string
	Expire       time.Duration
}

func (dom *domain) TokenValidate(r *arc.Request, in *TokenValidateInput) (*TokenValidateOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "35103cf4-0750-4048-bc0b-f4f12631f494",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	tokn := in.TokenExchangeBody

	if tokn == nil {
		return nil, &errors.Object{
			Id:     "1760a2a2-1af8-4545-8ead-032d159b97fe",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing token exchange body.",
		}
	}

	redirectUrl, err := url.Parse(tokn.RedirectUri)
	if err != nil {
		return nil, &errors.Object{
			Id:     "7275c8fc-6827-4696-8146-a756c08802e0",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse redirect URL.",
		}
	}

	if redirectUrl.Host != in.DownstreamHost {
		return nil, &errors.Object{
			Id:     "1f2ad831-bbd1-4330-bf17-a1a8ac54fcdd",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Redirect URL host does not match referrer.",
			Meta: map[string]any{
				"downstream_host": in.DownstreamHost,
				"redirect_host":   redirectUrl.Host,
			},
		}
	}

	selectUserOut, err := dom.UserDomain.SelectUser(r, &user.SelectUserInput{
		OrganizationId: tokn.OrganizationId,
		ExternalId:     tokn.ExternalId,
	})
	if err != nil {
		lastErr := errors.Last(err)

		if lastErr.Code != errors.Code_NOT_FOUND {
			return nil, errors.Forward(err, "94766c72-6ee0-423f-a4c0-57804c91a3e4")
		}

		insertUserOut, err := dom.UserDomain.InsertUser(r, &user.InsertUserInput{
			OrganizationId: tokn.OrganizationId,
			User: &user.User{
				OrganizationId: tokn.OrganizationId,
				RoleName:       consts.RoleSaasWhitelabelUser,
				ExternalId:     tokn.ExternalId,
			},
		})
		if err != nil {
			return nil, errors.Forward(err, "0b3145a0-4780-4525-8ef9-94ba961778c9")
		}

		selectUserOut = &user.SelectUserOutput{
			User: insertUserOut.User,
		}
	}

	user := selectUserOut.User

	var ttl time.Duration

	if tokn.Expire > 0 {
		ttl = time.Duration(tokn.Expire) * time.Second
	} else {
		defaultExpire, err := r.Dom().SelectDuration(consts.ConfigKeyDurationSaasSession)
		if err != nil {
			return nil, errors.Forward(err, "71b9060a-a05a-499e-87a5-c9c45efd8058")
		}

		ttl = defaultExpire
	}

	session, err := dom.ArcDomain.CreateServerSession(r, &arc.CreateServerSessionInput{
		OrganizationId: user.OrganizationId,
		UserId:         user.Id,
		RoleName:       user.RoleName,
		SessionType:    arc.SessionTypeSaasServer,
		Timezone:       "UTC",
		Username:       user.Username,
		TTL:            ttl,
	})
	if err != nil {
		return nil, errors.Forward(err, "934cef8c-2c8d-43db-b3ac-cdf774412ff1")
	}

	encodeTokenOut, err := dom.TokenDomain.EncodeToken(r, &token.EncodeTokenInput{
		Id:               session.Id(),
		OrganizationId:   user.OrganizationId,
		TokenType:        token.TypeSaasServerSession,
		Expire:           ttl,
		Version:          token.VersionPasetoV4Local,
		PasetoConfigName: consts.ConfigKeyPasetoSession,
	})
	if err != nil {
		return nil, errors.Forward(err, "f08a949c-aa54-433b-971f-b2f6fdd87258")
	}

	request, err := dom.ArcDomain.CreateRequest(&arc.CreateRequestInput{
		Context: r.Context(),
	})
	if err != nil {
		return nil, errors.Forward(err, "ef6aeca1-e4a4-4af0-850f-317f17050937")
	}

	request.SetSession(session)

	out := &TokenValidateOutput{
		SessionToken: encodeTokenOut.Value,
		RedirectUri:  tokn.RedirectUri,
		Expire:       ttl,
	}

	return out, nil
}
