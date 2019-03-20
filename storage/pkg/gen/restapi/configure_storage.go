// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/amadeusitgroup/miniplanes/storage/cmd/config"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/db/mongo"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/airlines"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/airports"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/courses"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/readiness"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/schedules"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/version"
)

//go:generate swagger generate server --target ../../pkg/gen --name storage --spec ../swagger.yaml --exclude-main

func configureFlags(*operations.StorageAPI) {
}

func configureAPI(api *operations.StorageAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.Logger = log.Infof
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// GET Airlines
	api.AirlinesGetAirlinesHandler = airlines.GetAirlinesHandlerFunc(func(params airlines.GetAirlinesParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		modAirlines, err := db.GetAirlines()
		if err != nil {
			message := fmt.Sprintf("unable to retrieve airlines: %v", err)
			log.Warn(message)
			return airlines.NewGetAirlinesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		return airlines.NewGetAirlinesOK().WithPayload(modAirlines)
	})

	// GET Airports
	api.AirportsGetAirportsHandler = airports.GetAirportsHandlerFunc(func(params airports.GetAirportsParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		modAirports, err := db.GetAirports()
		if err != nil {
			message := fmt.Sprintf("unable to retrieve airports: %v", err)
			log.Warn(message)
			return airports.NewGetAirportsBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		return airports.NewGetAirportsOK().WithPayload(modAirports)
	})

	// GET Courses
	api.CoursesGetCoursesHandler = courses.GetCoursesHandlerFunc(func(params courses.GetCoursesParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		modCourses, err := db.GetCourses()
		if err != nil {
			message := fmt.Sprintf("unable to retrieve courses: %v", err)
			log.Warn(message)
			return courses.NewGetCoursesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		return courses.NewGetCoursesOK().WithPayload(modCourses)
	})

	// GET Liveness
	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		if err := db.Ping(); err != nil {
			message := fmt.Sprintf("unable to ping DB %s: %v", db.DialString(), err)
			log.Warnf(message)
			return liveness.NewGetLiveServiceUnavailable().WithPayload(&models.Error{Code: http.StatusServiceUnavailable, Message: &message})
		}
		log.Debug("Storage is alive!")
		return liveness.NewGetLiveOK()
	})

	// GET Readiness
	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		if err := db.Ping(); err != nil {
			message := fmt.Sprintf("unable to ping DB %s: %v", db.DialString(), err)
			log.Warnf(message)
			return readiness.NewGetReadyServiceUnavailable().WithPayload(&models.Error{Code: http.StatusServiceUnavailable, Message: &message})
		}
		log.Trace("Storage is ready!")
		return readiness.NewGetReadyOK()
	})

	// GET Schedules
	api.SchedulesGetSchedulesHandler = schedules.GetSchedulesHandlerFunc(func(params schedules.GetSchedulesParams) middleware.Responder {
		log.Trace("Serving Schedules...")
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		modSchedules, err := db.GetSchedules()
		if err != nil {
			log.Errorf("Could't get schedules: %v\n", err)
			message := fmt.Sprintf("unable to retrieve airports: %v", err)
			return airports.NewGetAirportsBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		return schedules.NewGetSchedulesOK().WithPayload(modSchedules)
	})

	// AddSchedule
	api.SchedulesAddScheduleHandler = schedules.AddScheduleHandlerFunc(func(params schedules.AddScheduleParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		modSchedule, err := db.InsertSchedule(params.Schedule)
		if err != nil {
			return schedules.NewAddScheduleDefault(422) // todo Add 422 Unprocessable entity, 409 conflict (even if already exists)
		}
		return schedules.NewAddScheduleCreated().WithPayload(modSchedule)
	})

	// DELETE Schedule
	api.SchedulesDeleteScheduleHandler = schedules.DeleteScheduleHandlerFunc(func(params schedules.DeleteScheduleParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		err := db.DeleteSchedule(params.ID)
		if err != nil {
			return schedules.NewDeleteScheduleBadRequest()
		}
		return schedules.NewDeleteScheduleNoContent()
	})

	// GET Schedule<ID>
	api.SchedulesGetScheduleHandler = schedules.GetScheduleHandlerFunc(func(params schedules.GetScheduleParams) middleware.Responder {
		//return middleware.NotImplemented("operation schedules.GetSchedule has not yet been implemented")
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		schedule, err := db.GetSchedule(params.ID)
		if err != nil {
			return schedules.NewGetSchedulesBadRequest()
		}
		return schedules.NewGetScheduleOK().WithPayload(schedule)
	})

	// PUT Schedule
	api.SchedulesUpdateScheduleHandler = schedules.UpdateScheduleHandlerFunc(func(params schedules.UpdateScheduleParams) middleware.Responder {
		db := mongo.NewMongoDB(config.MongoHost, config.MongoPort, config.MongoDBName)
		schedule, err := db.UpdateSchedule(params.ID, params.Schedule)
		if err != nil {
			return schedules.NewUpdateScheduleBadRequest()
		}
		return schedules.NewUpdateScheduleCreated().WithPayload(schedule)
	})

	// GetVersion
	api.VersionGetVersionHandler = version.GetVersionHandlerFunc(func(params version.GetVersionParams) middleware.Responder {
		log.Tracef("Serving Version: %s", config.Version)
		return version.NewGetVersionOK().WithPayload(&models.Version{
			Version: config.Version,
		})
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
