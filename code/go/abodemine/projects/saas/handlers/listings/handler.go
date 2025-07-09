package listings

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"abodemine/domains/arc"
	"abodemine/domains/listings"
	"abodemine/lib/errors"
	"abodemine/projects/saas/domains/auth"
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
	authOut, err := h.AuthDomain.Authenticate(r.Context(), &auth.AuthenticateInput{HttpRequest: r})
	if err != nil {
		arc.HttpApiErrorResponse(h.ArcDomain, w, "", errors.Forward(err, "b7609198-3b15-4422-a9ab-9abde2e8e67e"))
		return
	}

	arcRequest := authOut.Request

	searchListingsInput := &SearchListingsInput{}
	if err := json.NewDecoder(r.Body).Decode(&searchListingsInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	searchListingsInputDomain := searchListingsInput.ToDomainModel()
	out, err := h.ListingsDomain.SearchListings(arcRequest, &searchListingsInputDomain)
	if err != nil {
		http.Error(w, "Error searching listings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	arc.HttpApiDataResponse(h.ArcDomain, w, http.StatusOK, out)
}
