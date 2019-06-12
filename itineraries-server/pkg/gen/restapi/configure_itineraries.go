/*

MIT License

Copyright (c) 2019 Amadeus s.a.s.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"time"

	"net/http"

	"github.com/amadeusitgroup/miniplanes/itineraries-server/cmd/config"
	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/engine"
	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/restapi/operations/readiness"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"

	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/restapi/operations"
	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/restapi/operations/itineraries"
	airportsclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airports"
	schedulesclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/schedules"
	storage_models "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	httptransport "github.com/go-openapi/runtime/client"
)

//go:generate swagger generate server --target .. --name itineraries --spec ../swagger/swagger.yaml

var (
	airports  []*storage_models.Airport
	schedules []*storage_models.Schedule
)

func configureFlags(api *operations.ItinerariesAPI) {
}

func makeDescription(from, to string) string {
	return fmt.Sprintf("Itinerary: %s - %s", from, to)
}

func refreshSchedulesIfNeeded() {
	log.Trace("Refreshhing schedules")
	storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
	client := schedulesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
	Ok, err := client.GetSchedules(schedulesclient.NewGetSchedulesParams())
	if err != nil {
		log.Errorf("Current number of schedule %d. Unable to refresh Schedules: %v", len(schedules), err)
		return
	}
	schedules = Ok.Payload
}

func refreshAirportsIfNeeded() {
	log.Trace("Refreshing Airports")
	storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
	client := airportsclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
	Ok, err := client.GetAirports(airportsclient.NewGetAirportsParams())
	if err != nil {
		log.Errorf("Current number of airports %d. Unable to refresh airports: %v", len(airports), err)
		return
	}
	airports = Ok.Payload
}

func makeItineraryID() string {
	return "this_is_my_itinerary_ID"
}

func getItineraries(from, to, departureDate, departureTime *string) ([]*models.Itinerary, error) {
	log.Debugf("getItineraries from:%s, to:%s, departureDate:%s, departureTime:%s", *from, *to, *departureDate, *departureTime)
	refreshSchedulesIfNeeded()
	refreshAirportsIfNeeded()
	itineraryGraph, err := engine.NewGraph(airports, schedules)
	if err != nil {
		log.Errorf("unable to instantiate itineraries-server engine: %v", err)
		return []*models.Itinerary{}, fmt.Errorf("unable to instantiate itineraries-server engine: %v", err)
	}

	departureDateTime := fmt.Sprintf("%sT%s", *departureDate, *departureTime)
	d, err := time.Parse("2006-01-02T15:04", departureDateTime)
	if err != nil {
		return []*models.Itinerary{}, fmt.Errorf("Unable to parse departure date time: %v", err)
	}
	dDate := fmt.Sprintf("%02d%02d", int(d.Day()), int(d.Month()))
	dTime := fmt.Sprintf("%02d%02d", int(d.Hour()), int(d.Minute()))
	maxNumberOfPaths := 5
	return itineraryGraph.Compute(*from, dDate, dTime, *to, maxNumberOfPaths)
}

func isAlive() bool {
	return true // TODO: fill it
}

func isReady() bool {
	return true // TODO: fill it
}

func configureAPI(api *operations.ItinerariesAPI) http.Handler {
	api.ServeError = errors.ServeError
	api.Logger = log.Printf
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.LivenessGetLiveHandler = liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
		var r middleware.Responder
		r = liveness.NewGetLiveServiceUnavailable()
		if isAlive() {
			r = liveness.NewGetLiveOK()
		}
		return r
	})

	api.ReadinessGetReadyHandler = readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
		var r middleware.Responder
		r = readiness.NewGetReadyServiceUnavailable()
		if isReady() {
			r = readiness.NewGetReadyOK()
		}
		return r
	})

	api.ItinerariesGetItinerariesHandler = itineraries.GetItinerariesHandlerFunc(func(params itineraries.GetItinerariesParams) middleware.Responder {
		mergedParams := itineraries.NewGetItinerariesParams()
		if params.From == nil || len(*params.From) == 0 {
			message := fmt.Sprintf("No Origin:")
			log.Error(message)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		mergedParams.From = params.From
		log.Tracef("From: %s", *mergedParams.From)

		if params.To == nil || len(*params.To) == 0 {
			message := fmt.Sprintf("No Destination")
			log.Error(message)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		mergedParams.To = params.To
		log.Tracef("To: %s", *params.To)

		if params.DepartureDate == nil || len(*params.DepartureDate) == 0 {
			message := fmt.Sprintf("No DepartureDate")
			log.Error(message)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		mergedParams.DepartureDate = params.DepartureDate
		log.Tracef("DepartureDate: %s", *mergedParams.DepartureDate)

		if params.DepartureTime != nil && len(*params.DepartureTime) != 0 {
			mergedParams.DepartureTime = params.DepartureTime
		} // otherwise DepartureTime is defaulted from swagger
		log.Tracef("DepartureTime: %s", *params.DepartureTime)

		if params.ReturnDate == nil || len(*params.ReturnDate) == 0 {
			message := fmt.Sprintf("No ReturnDate")
			log.Error(message)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		mergedParams.ReturnDate = params.ReturnDate
		log.Tracef("ReturnDate: %s", *params.ReturnDate)

		if params.ReturnTime != nil && len(*params.ReturnTime) != 0 {
			mergedParams.ReturnTime = params.ReturnTime
		} // otherwise ReturnTime is defaulted from swagger
		log.Tracef("ReturnTime: %s", *mergedParams.ReturnTime)

		its, err := getItineraries(mergedParams.From, mergedParams.To, mergedParams.DepartureDate, mergedParams.DepartureTime)
		if err != nil {
			message := fmt.Sprintf("Couldn't get itineraries: %v", err)
			log.Error(message)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &message})
		}
		if len(its) == 0 {
			message := fmt.Sprintf("No itineraries found")
			log.Warn(message)
			return itineraries.NewGetItinerariesNotFound().WithPayload(&models.Error{Code: http.StatusNotFound, Message: &message})
		}
		return itineraries.NewGetItinerariesOK().WithPayload(its)
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
