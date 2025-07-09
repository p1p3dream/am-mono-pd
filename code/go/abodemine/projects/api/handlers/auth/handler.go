package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/projects/api/domains/auth"
)

type Handler interface {
	ExchangeToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
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

type TokenExchangeInput struct {
	ClientId    uuid.UUID `json:"client_id,omitempty"`
	Expire      int       `json:"expire,omitempty"`
	ExternalId  string    `json:"external_id,omitempty"`
	RedirectUri string    `json:"redirect_uri,omitempty"`
}

type TokenExchangeOutput struct {
	Token string `json:"token,omitempty"`
}

func (h *handler) TokenExchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{
		AuthorizationHeader: r.Header["Authorization"],
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "3a218eb4-375c-432a-aac0-50f9d94d8e38"))
		return
	}

	arcRequest := authOut.Request

	inB, err := io.ReadAll(r.Body)
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), &errors.Object{
			Id:     "3ac3e3b9-c7a6-4bbe-b2f5-591413a7faab",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to read request body.",
			Cause:  err.Error(),
		})
		return
	}

	tokenExchangeInput := new(TokenExchangeInput)

	if err := json.Unmarshal(inB, tokenExchangeInput); err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), &errors.Object{
			Id:     "0197d777-5880-4582-ac8e-7553a9abe81a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid request body.",
			Cause:  err.Error(),
		})
		return
	}

	tokenExchangeOut, err := h.AuthDomain.TokenExchange(arcRequest, &auth.TokenExchangeInput{
		ClientId:    tokenExchangeInput.ClientId,
		Expire:      tokenExchangeInput.Expire,
		ExternalId:  tokenExchangeInput.ExternalId,
		RedirectUri: tokenExchangeInput.RedirectUri,
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), errors.Forward(err, "1e85ecb3-9642-46d8-9df5-dcbd4a92b325"))
		return
	}

	out := &TokenExchangeOutput{
		Token: tokenExchangeOut.Token,
	}

	arc.HttpApiDataResponse(h.ArcDomain, w, http.StatusOK, out)
}
