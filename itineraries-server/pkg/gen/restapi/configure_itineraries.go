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
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/itineraries"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/restapi/operations/readiness"
	storageclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/airports"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/schedules"
)

var (
	airportsId2Name = make(map[int64]string)
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
		storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
		transport := storageclient.DefaultTransportConfig().WithHost(storageURL)
		client := storageclient.NewHTTPClientWithConfig(nil, transport)
		//client := storageclient.New(transport, strfmt.Default)
		modItineraries := []*models.Itinerary{}

		if len(airportsId2Name) == 0 {
			getParams := &airports.GetAirportsParams{Context: context.Background()}
			resp, err := client.Airports.GetAirports(getParams)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				//return itineraries.NewGetItinerariesOK().WithPayload(modItineraries)
				return itineraries.NewGetItinerariesNotFound()
			}
			for i := range resp.Payload {
				a := resp.Payload[i]
				airportsId2Name[a.AirportID] = a.IATA
			}
		}

		getParams := &schedules.GetSchedulesParams{Context: context.Background()}
		resp, err := client.Schedules.GetSchedules(getParams)
		if err != nil {
			return itineraries.NewGetItinerariesBadRequest()
		}

		for i := range resp.Payload {
			schedule := resp.Payload[i]
			segment := &models.Segment{
				//ArrivalDate: *schedule.Arrival, // obtained from OperatedDay
				ArrivalTime:      *schedule.ArrivalTime,
				DepartureTime:    *schedule.DepartureTime,
				ArriveNextDay:    *schedule.ArriveNextDay,
				Destination:      airportsId2Name[*schedule.Destination],
				FlightNumber:     *schedule.FlightNumber,
				OperatingCarrier: *schedule.OperatingCarrier,
				Origin:           airportsId2Name[*schedule.Origin],
				SegmentID:        0,
			}
			itinerary := &models.Itinerary{
				Description: "my awesome itinerary",
				ItineraryID: "",
				Segments:    []*models.Segment{segment},
			}
			modItineraries = append(modItineraries, itinerary)
			break
		}
		return itineraries.NewGetItinerariesOK().WithPayload(modItineraries)
	})
	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		return middleware.NotImplemented("operation liveness.GetLive has not yet been implemented")
	})
	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		return middleware.NotImplemented("operation readiness.GetReady has not yet been implemented")
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
