package search_handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/models"
	"abodemine/projects/api/domains/auth"
	"abodemine/projects/api/domains/search"
)

type Handler interface {
	SearchProperty(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type handler struct {
	authDomain   auth.Domain
	searchDomain search.Domain
	arcDomain    arc.Domain
}

func NewHandler(authDomain auth.Domain, searchDomain search.Domain, arcDomain arc.Domain) Handler {
	return &handler{
		authDomain:   authDomain,
		searchDomain: searchDomain,
		arcDomain:    arcDomain,
	}
}

type SearchPropertyInput struct {
	Aupid   uuid.UUID `json:"aupid"`
	Layouts []string  `json:"layouts"`

	models.ApiSearchAddress
}

func (h *handler) SearchProperty(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Authenticate request.
	authOut, err := h.authDomain.Authenticate(r.Context(), &auth.AuthenticateInput{
		AuthorizationHeader: r.Header["Authorization"],
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.arcDomain, w, "", err)
		return
	}
	arcRequest := authOut.Request

	var input SearchPropertyInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		arc.HttpApiErrorResponse(h.arcDomain, w, arcRequest.Id().String(), &errors.Object{
			Id:     "6244d443-5de2-4ee0-aa93-e10da4d935bc",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid request body.",
		})
		return
	}

	searchAuroraOut, err := h.searchDomain.SearchProperty(arcRequest, &search.SearchPropertyInput{
		Aupid: val.Ternary(
			input.Aupid != uuid.Nil,
			&input.Aupid,
			nil,
		),
		Layouts:          input.Layouts,
		ApiSearchAddress: &input.ApiSearchAddress,
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.arcDomain, w, arcRequest.Id().String(), err)
		return
	}

	properties := searchAuroraOut.PropertyEntities

	// Ensure we return an emtpy slice instead of null.
	if properties == nil {
		properties = []*entities.Property{}
	}

	// Write response
	arc.HttpApiDataResponse(h.arcDomain, w, http.StatusOK, properties)
}
