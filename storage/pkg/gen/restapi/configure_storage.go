// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/airlines"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/airports"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/courses"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/readiness"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations/schedules"
)

//go:generate swagger generate server --target ../../pkg/gen --name storage --spec ../swagger.yaml --exclude-main

func configureFlags(api *operations.StorageAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
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

	api.AirlinesGetAirlinesHandler = airlines.GetAirlinesHandlerFunc(func(params airlines.GetAirlinesParams) middleware.Responder {
		return middleware.NotImplemented("operation airlines.GetAirlines has not yet been implemented")
	})
	api.AirportsGetAirportsHandler = airports.GetAirportsHandlerFunc(func(params airports.GetAirportsParams) middleware.Responder {
		return middleware.NotImplemented("operation airports.GetAirports has not yet been implemented")
	})
	api.CoursesGetCoursesHandler = courses.GetCoursesHandlerFunc(func(params courses.GetCoursesParams) middleware.Responder {
		return middleware.NotImplemented("operation courses.GetCourses has not yet been implemented")
	})
	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		return middleware.NotImplemented("operation liveness.GetLive has not yet been implemented")
	})
	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		return middleware.NotImplemented("operation readiness.GetReady has not yet been implemented")
	})
	api.SchedulesGetSchedulesHandler = schedules.GetSchedulesHandlerFunc(func(params schedules.GetSchedulesParams) middleware.Responder {
		return middleware.NotImplemented("operation schedules.GetSchedules has not yet been implemented")
	})
	api.SchedulesAddScheduleHandler = schedules.AddScheduleHandlerFunc(func(params schedules.AddScheduleParams) middleware.Responder {
		return middleware.NotImplemented("operation schedules.AddSchedule has not yet been implemented")
	})
	api.SchedulesDeleteScheduleHandler = schedules.DeleteScheduleHandlerFunc(func(params schedules.DeleteScheduleParams) middleware.Responder {
		return middleware.NotImplemented("operation schedules.DeleteSchedule has not yet been implemented")
	})
	api.SchedulesGetScheduleHandler = schedules.GetScheduleHandlerFunc(func(params schedules.GetScheduleParams) middleware.Responder {
		return middleware.NotImplemented("operation schedules.GetSchedule has not yet been implemented")
	})
	api.SchedulesUpdateScheduleHandler = schedules.UpdateScheduleHandlerFunc(func(params schedules.UpdateScheduleParams) middleware.Responder {
		return middleware.NotImplemented("operation schedules.UpdateSchedule has not yet been implemented")
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
