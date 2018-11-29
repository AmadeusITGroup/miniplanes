// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/amadeusitgroup/miniapp/storage/pkg/backend/mongo"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/airlines"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/airports"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/readiness"
)

//go:generate swagger generate server --target ../../pkg/gen --name storage --spec ../swagger.yaml --exclude-main

var (
	MongoHost string
	MongoPort int
)

func configureFlags(api *operations.StorageAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func isAlive() bool {
	return true
}

func isReady() bool {
	return true
}

func configureAPI(api *operations.StorageAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Airlines
	api.AirlinesGetAirlinesHandler = airlines.GetAirlinesHandlerFunc(func(params airlines.GetAirlinesParams) middleware.Responder {
		db := mongo.NewMongoDB("")
		dbAirlines, err := db.GetAirlines()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			message := fmt.Sprintf("unable to retrieve airlines: %v", err)
			return airlines.NewGetAirlinesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		modAirlines := []*models.Airline{}
		for _, a := range dbAirlines {
			tmp := &models.Airline{}
			copier.Copy(tmp, a)
			modAirlines = append(modAirlines, tmp)
		}
		return airlines.NewGetAirlinesOK().WithPayload(modAirlines)
	})

	// Airports
	api.AirportsGetAirportsHandler = airports.GetAirportsHandlerFunc(func(params airports.GetAirportsParams) middleware.Responder {
		db := mongo.NewMongoDB("")
		dbAirports, err := db.GetAirports()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			message := fmt.Sprintf("unable to retrieve airports: %v", err)
			return airports.NewGetAirportsBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}

		modAirports := []*models.Airport{}
		for _, a := range dbAirports {
			tmp := &models.Airport{}
			copier.Copy(tmp, a)
			modAirports = append(modAirports, tmp)
		}

		//return airports, nil
		return airports.NewGetAirportsOK().WithPayload(modAirports)
	})

	// Routes

	// Schedules

	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) (r middleware.Responder) {
		r = liveness.NewGetLiveServiceUnavailable()
		if isAlive() {
			r = liveness.NewGetLiveOK()
		}
		return r
	})
	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) (r middleware.Responder) {
		r = readiness.NewGetReadyServiceUnavailable()
		if isReady() {
			r = readiness.NewGetReadyOK()
		}
		return r
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
