package listings

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/domains/listings"
	"abodemine/lib/errors"
	"abodemine/projects/api/domains/auth"
)

type Handler interface {
	ExchangeToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type handler struct {
	ArcDomain      arc.Domain
	AuthDomain     auth.Domain
	ListingsDomain listings.Domain
}

type NewHandlerInput struct {
	ArcDomain      arc.Domain
	AuthDomain     auth.Domain
	ListingsDomain listings.Domain
}

func NewHandler(in *NewHandlerInput) *handler {
	return &handler{
		ArcDomain:      in.ArcDomain,
		AuthDomain:     in.AuthDomain,
		ListingsDomain: in.ListingsDomain,
	}
}

func (h *handler) GetListings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{
		AuthorizationHeader: r.Header["Authorization"],
	})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "597e64eb-43aa-4d21-9740-9ebfa8f5cfa0"))
		return
	}

	arcRequest := authOut.Request

	searchListingsInput := &SearchListingsInput{}
	if err := json.NewDecoder(r.Body).Decode(&searchListingsInput); err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), &errors.Object{
			Id:     "7cb7020b-d014-42d1-9552-9af9ba25f515",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid request body.",
		})
		return
	}

	searchListingsInputDomain := searchListingsInput.ToDomainModel()
	out, err := h.ListingsDomain.SearchListings(arcRequest, &searchListingsInputDomain)
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, arcRequest.Id().String(), errors.Forward(err, "54591cb9-62ba-435d-a690-cf067cfe53cc"))
		return
	}

	arc.HttpApiDataResponse(h.ArcDomain, w, http.StatusOK, out)
}
