package handlers

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	listings_domain "abodemine/domains/listings"
	"abodemine/domains/token"
	"abodemine/lib/app"
	"abodemine/projects/saas/conf"
	auth_domain "abodemine/projects/saas/domains/auth"
	"abodemine/projects/saas/domains/user"
	auth_handler "abodemine/projects/saas/handlers/auth"
	listings_handler "abodemine/projects/saas/handlers/listings"
	"abodemine/projects/saas/handlers/tests"
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
	})

	// addressDomain := legacy_address_domain.NewDomain(&legacy_address_domain.NewDomainInput{})

	tokenDomain := token.NewDomain(&token.NewDomainInput{})
	userDomain := user.NewDomain(&user.NewDomainInput{})

	authDomain := auth_domain.NewDomain(&auth_domain.NewDomainInput{
		ArcDomain:   arcDomain,
		TokenDomain: tokenDomain,
		UserDomain:  userDomain,
	})

	authHandler := auth_handler.NewHandler(&auth_handler.NewHandlerInput{
		ArcDomain:  arcDomain,
		AuthDomain: authDomain,
	})

	testsHandler := tests.NewHandler(&tests.NewHandlerInput{
		ArcDomain: arcDomain,
	})

	listingsDomain := listings_domain.NewDomain(&listings_domain.NewDomainInput{
		ArcDomain: arcDomain,
		// AddressDomain: addressDomain,
	})

	listingsHandler := listings_handler.NewHandler(&listings_handler.NewHandlerInput{
		ArcDomain:      arcDomain,
		AuthDomain:     authDomain,
		ListingsDomain: listingsDomain,
	})

	router := httprouter.New()

	// Health check.
	var healthCheckLog sync.Once

	////////////////////////////////////////////////////////////////////////////
	// Api routes.
	////////////////////////////////////////////////////////////////////////////

	apiPrefix := "/api"

	router.GET(
		apiPrefix+"/685f26df-dc42-4009-8784-63bd162d0255",
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

	router.POST(
		apiPrefix+"/auth/load",
		authHandler.Load,
	)

	router.GET(
		apiPrefix+"/auth/token/validate/:token",
		authHandler.TokenValidate,
	)

	router.POST(
		apiPrefix+"/listings",
		listingsHandler.GetListings,
	)

	if c.File.DeploymentEnvironment < app.DeploymentEnvironment_TESTING {
		return router
	}

	////////////////////////////////////////////////////////////////////////////
	// Testing routes.
	////////////////////////////////////////////////////////////////////////////

	testsPrefix := "/zxkvkpkzhpah"

	router.GET(
		testsPrefix+"/token/exchange",
		testsHandler.TokenExchangeView,
	)

	router.POST(
		testsPrefix+"/token/exchange",
		testsHandler.TokenExchange,
	)

	return router
}
