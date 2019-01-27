// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/amadeusitgroup/miniapp/itineraries-server/cmd/config"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/engine"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/itineraries"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/readiness"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/version"
	storageclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/airports"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/schedules"
	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
)

var (
	airportID2Airport = make(map[int32]*storagemodels.Airport)
	airportsIATA2ID   = make(map[string]int32)
)

//go:generate swagger generate server --target ../../pkg/gen --name itineraries --spec ../swagger.yaml --exclude-main

func configureFlags(api *operations.ItinerariesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ItinerariesAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ItinerariesGetItinerariesHandler = itineraries.GetItinerariesHandlerFunc(func(params itineraries.GetItinerariesParams) middleware.Responder {

		from := params.From
		fmt.Printf("FROM: %s\n", *from)

		to := params.To
		fmt.Printf("TO: %s\n", *to)

		departureDate := params.DepartureDate
		fmt.Printf("Departure Date: %s\n", *departureDate)

		departureTime := "0800"
		fmt.Printf("Departure Time: %s\n", departureTime)

		//returnDate := params.ReturnDate
		//fmt.Printf("Return Date: %s\n", *returnDate)

		storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
		transport := storageclient.DefaultTransportConfig().WithHost(storageURL)
		client := storageclient.NewHTTPClientWithConfig(nil, transport)

		airportsParams := &airports.GetAirportsParams{Context: context.Background()}
		airportsResp, err := client.Airports.GetAirports(airportsParams)
		if err != nil {
			return itineraries.NewGetItinerariesNotFound() // TODO: try later set error
		}

		schedulesParams := &schedules.GetSchedulesParams{Context: context.Background()}
		schedulesResp, err := client.Schedules.GetSchedules(schedulesParams)
		if err != nil {
			return itineraries.NewGetItinerariesNotFound() // TODO: try later set error
		}

		itineraryGraph, err := engine.NewGraph(airportsResp.Payload, schedulesResp.Payload)
		if err != nil {
			return itineraries.NewGetItinerariesNotFound() // TODO: try later set error
		}

		modItineraries, err := itineraryGraph.Compute(*from, *departureDate, departureTime, *to, 5)
		if err != nil {
			return itineraries.NewGetItinerariesNotFound() // TODO: try later set error
		}

		return itineraries.NewGetItinerariesOK().WithPayload(modItineraries)
	})
	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		// liveness.NewGetLiveServiceUnavailable()
		return liveness.NewGetLiveOK()
	})
	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		//readiness.NewGetReadyServiceUnavailable()
		return readiness.NewGetReadyOK()
	})

	api.VersionGetVersionHandler = version.GetVersionHandlerFunc(func(params version.GetVersionParams) middleware.Responder {
		tmp := &models.Version{
			Version: config.Version,
		}
		return version.NewGetVersionOK().WithPayload(tmp)
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
