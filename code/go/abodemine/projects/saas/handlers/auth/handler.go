package auth

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/projects/saas/domains/auth"
)

type Handler interface {
	Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

	TokenValidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type handler struct {
	ArcDomain  arc.Domain
	AuthDomain auth.Domain
}

type NewHandlerInput struct {
	ArcDomain  arc.Domain
	AuthDomain auth.Domain
}

func NewHandler(in *NewHandlerInput) *handler {
	return &handler{
		ArcDomain:  in.ArcDomain,
		AuthDomain: in.AuthDomain,
	}
}

func (h *handler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{HttpRequest: r})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Unauthenticated(err, "ffb8bc98-6f92-4f31-a120-786d1596be2d"))
		return
	}

	arcRequest := authOut.Request

	loadOut, err := h.AuthDomain.Load(authOut.Request)
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), errors.Forward(err, "9b97079b-589c-4356-abf2-36aa55435004"))
		return
	}

	arc.HttpApiDataResponse(h.ArcDomain, w, http.StatusOK, loadOut)
}

type TokenValidateInput struct {
	Token string `json:"token,omitempty"`
}

func (h *handler) TokenValidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{
		AuthMethod: auth.AuthMethodToken,
		Token:      ps.ByName("token"),
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "a8fe15c2-cfb1-4214-a715-7ba69a0c8bf7"))
		return
	}

	arcRequest := authOut.Request
	tokenExchangeBody := authOut.TokenExchangeBody

	validateTokenOut, err := h.AuthDomain.TokenValidate(arcRequest, &auth.TokenValidateInput{
		DownstreamHost:    r.Host,
		TokenExchangeBody: tokenExchangeBody,
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), errors.Forward(err, "a7a8f225-7f3c-4b42-9014-dfd00561826b"))
		return
	}

	cookie := &http.Cookie{
		Name:     consts.CookieAbodeMineSaasWebSession,
		Value:    validateTokenOut.SessionToken,
		Path:     "/",
		Expires:  time.Now().Add(validateTokenOut.Expire),
		MaxAge:   int(validateTokenOut.Expire.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, validateTokenOut.RedirectUri, http.StatusFound)
}
