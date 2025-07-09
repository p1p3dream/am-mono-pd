package handlers

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/domains/assessor"
	"abodemine/domains/avm"
	listings "abodemine/domains/listings"
	"abodemine/domains/property"
	"abodemine/domains/recorder"
	"abodemine/domains/token"
	"abodemine/middleware"
	"abodemine/projects/api/conf"
	auth "abodemine/projects/api/domains/auth"
	search "abodemine/projects/api/domains/search"
	auth_handler "abodemine/projects/api/handlers/auth"
	listings_handler "abodemine/projects/api/handlers/listings"
	search_handler "abodemine/projects/api/handlers/search"
)

func Router(c *conf.Config) *httprouter.Router {
	arcDomain := arc.NewDomain(&arc.NewDomainInput{
		DeploymentEnvironment: c.File.DeploymentEnvironment,
		Casbin:                c.Casbin,
		Duration:              c.Duration,
		Flags:                 c.File.Flags,
		OpenSearch:            c.OpenSearch,
		Paseto:                c.Paseto,
		PgxPool:               c.PGxPool,
		Valkey:                c.Valkey,
		ValkeyScript:          c.ValkeyScript,
		Values:                c.Values,
	})

	tokenDomain := token.NewDomain(&token.NewDomainInput{})

	authDomain := auth.NewDomain(&auth.NewDomainInput{
		ArcDomain:   arcDomain,
		TokenDomain: tokenDomain,
	})

	// legacyDomain := legacy.NewDomain(&legacy.NewDomainInput{
	// 	Config:     c,
	// 	AuthDomain: authDomain,
	// })

	addressDomain := address.NewDomain(&address.NewDomainInput{})
	assessorDomain := assessor.NewDomain(&assessor.NewDomainInput{})
	avmDomain := avm.NewDomain(&avm.NewDomainInput{})
	recorderDomain := recorder.NewDomain(&recorder.NewDomainInput{})

	listingsDomain := listings.NewDomain(&listings.NewDomainInput{
		ArcDomain:     arcDomain,
		AuthDomain:    authDomain,
		AddressDomain: addressDomain,
	})

	propertyDomain := property.NewDomain(&property.NewDomainInput{
		AddressDomain:  addressDomain,
		AssessorDomain: assessorDomain,
		AvmDomain:      avmDomain,
		ListingDomain:  listingsDomain,
		RecorderDomain: recorderDomain,
	})

	searchDomain := search.NewDomain(&search.NewDomainInput{
		AddressDomain:  addressDomain,
		AuthDomain:     authDomain,
		PropertyDomain: propertyDomain,
	})

	authHandler := auth_handler.NewHandler(&auth_handler.NewHandlerInput{
		ArcDomain:  arcDomain,
		AuthDomain: authDomain,
	})

	// legacyHandler := legacy_handler.NewHandler(&legacy_handler.NewHandlerInput{
	// 	ArcDomain:    arcDomain,
	// 	AuthDomain:   authDomain,
	// 	LegacyDomain: legacyDomain,
	// })

	listingsHandler := listings_handler.NewHandler(&listings_handler.NewHandlerInput{
		ArcDomain:      arcDomain,
		AuthDomain:     authDomain,
		ListingsDomain: listingsDomain,
	})

	searchHandler := search_handler.NewHandler(authDomain, searchDomain, arcDomain)

	// httprouter doesn't support subrouting, so we have to prefix all routes.
	router := httprouter.New()

	// Health check.
	var healthCheckLog sync.Once

	router.GET(
		"/d5aecc76-03e9-4920-8ee5-30f3c390bc52",
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			healthCheckLog.Do(func() {
				log.Info().
					Str("host", r.Host).
					Str("remote_addr", r.RemoteAddr).
					Msg("Health check. This event is logged only once.")
			})

			w.WriteHeader(http.StatusOK)
		},
	)

	////////////////////////////////////////////////////////////////////////////
	// v2 routes.
	////////////////////////////////////////////////////////////////////////////

	// v2 was deprecated in favor of v3, and all resources were removed as of
	// June 10, 2025.

	// v2Prefix := "/api/v2"

	// router.POST(v2Prefix+"/search", middleware.GzipHandler(middleware.SentryMiddlewareHandler(legacyHandler.Search)))

	////////////////////////////////////////////////////////////////////////////
	// v3 routes.
	////////////////////////////////////////////////////////////////////////////

	v3Prefix := "/api/v3"

	router.POST(
		v3Prefix+"/auth/token/exchange",
		middleware.GzipHandler(middleware.SentryMiddlewareHandler(authHandler.TokenExchange)),
	)

	router.POST(
		v3Prefix+"/listings",
		middleware.GzipHandler(middleware.SentryMiddlewareHandler(listingsHandler.GetListings)),
	)

	router.POST(
		v3Prefix+"/search",
		middleware.GzipHandler(middleware.SentryMiddlewareHandler(searchHandler.SearchProperty)),
	)

	return router
}
