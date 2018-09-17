// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/amadeusitgroup/miniapp/itinerary/models"
	"github.com/amadeusitgroup/miniapp/itinerary/restapi/operations"
	"github.com/amadeusitgroup/miniapp/itinerary/restapi/operations/itineraries"
)

//go:generate swagger generate server --target .. --name itineraries --spec ../swagger/swagger.yaml

func configureFlags(api *operations.ItinerariesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func getItineraries(from, to *string) []*models.Itinerary {
	var itineraries []*models.Itinerary

	return itineraries
}

func configureAPI(api *operations.ItinerariesAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ItinerariesGetItinerariesHandler = itineraries.GetItinerariesHandlerFunc(func(params itineraries.GetItinerariesParams) middleware.Responder {
		mergedParam := itineraries.NewGetItinerariesParams()
		if params.From != nil {
			mergedParam.From = params.From
		}
		if params.To != nil {
			mergedParam.To = params.To
		}

		if mergedParam.From == nil || mergedParam.To == nil {

		}

		its := getItineraries(mergedParam.From, mergedParam.To)

		fmt.Printf("Itinerary request: %q -> %q\n", *mergedParam.From, *mergedParam.To)
		return itineraries.NewGetItinerariesOK().WithPayload(its)
		//return middleware.NotImplemented("operation itineraries.GetItineraries has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
