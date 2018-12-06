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
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/schedules"
)

//go:generate swagger generate server --target ../../pkg/gen --name storage --spec ../swagger.yaml --exclude-main

var (
	MongoHost string
	MongoPort int
)

func configureFlags(*operations.StorageAPI) {
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
		db := mongo.NewMongoDB(MongoHost, MongoPort)
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
		db := mongo.NewMongoDB(MongoHost, MongoPort)
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
		return airports.NewGetAirportsOK().WithPayload(modAirports)
	})

	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		return liveness.NewGetLiveOK()
	})

	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		return readiness.NewGetReadyOK()
	})

	/*api.SchedulesGetRoutesHandler = schedules.GetRoutesHandlerFunc(func(params schedules.GetRoutesParams) middleware.Responder {
		db := mongo.NewMongoDB(MongoHost, MongoPort)
		dbRoutes, err := db.GetRoutes()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			message := fmt.Sprintf("unable to retrieve routes: %v", err)
			return routes.NewGetRoutesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		return middleware.NotImplemented("operation schedules.GetRoutes has not yet been implemented")
	})*/

	api.SchedulesGetSchedulesHandler = schedules.GetSchedulesHandlerFunc(func(params schedules.GetSchedulesParams) middleware.Responder {
		db := mongo.NewMongoDB(MongoHost, MongoPort)
		dbSchedules, err := db.GetSchedules()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			message := fmt.Sprintf("unable to retrieve airports: %v", err)
			return airports.NewGetAirportsBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		modSchedules := []*models.Schedule{}
		for _, a := range dbSchedules {
			tmp := &models.Schedule{}
			copier.Copy(tmp, a)
			modSchedules = append(modSchedules, tmp)
		}
		return schedules.NewGetSchedulesOK().WithPayload(modSchedules)
	})
	api.AddScheduleHandler = operations.AddScheduleHandlerFunc(func(params operations.AddScheduleParams) middleware.Responder {
		return middleware.NotImplemented("operation .AddSchedule has not yet been implemented")
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