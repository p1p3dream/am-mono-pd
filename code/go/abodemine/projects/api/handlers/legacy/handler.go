package legacy

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/projects/api/domains/auth"
	"abodemine/projects/api/domains/legacy"
	"abodemine/projects/api/domains/legacy/models"
)

type Handler interface {
	Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type handler struct {
	ArcDomain    arc.Domain
	AuthDomain   auth.Domain
	LegacyDomain legacy.Domain
}

type NewHandlerInput struct {
	ArcDomain    arc.Domain
	AuthDomain   auth.Domain
	LegacyDomain legacy.Domain
}

func NewHandler(in *NewHandlerInput) *handler {
	return &handler{
		ArcDomain:    in.ArcDomain,
		AuthDomain:   in.AuthDomain,
		LegacyDomain: in.LegacyDomain,
	}
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{
		AuthorizationHeader: r.Header["Authorization"],
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "4cd05b07-fcdb-4734-878e-0dc80f2666a8"))
		return
	}

	arcRequest := authOut.Request

	var request []models.PropertySearchRequests
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", &errors.Object{
			Id:     "857a50a4-b2b7-453f-b3c2-ae233e7cefdd",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to read request body.",
			Cause:  err.Error(),
		})
		return
	}

	response, err := h.LegacyDomain.Search(arcRequest, request)
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "aac5e2af-66df-4cae-8c97-f291c6bd7267"))
		return
	}

	arc.HttpApiDataResponse(h.ArcDomain, w, http.StatusOK, response)
}
