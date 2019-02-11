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
	"log"
	"net/http"

	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/restapi/operations/readiness"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	itinerary_utils "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/itineraries"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/models"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/restapi/operations"
	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/restapi/operations/itineraries"
)

//go:generate swagger generate server --target .. --name itineraries --spec ../swagger/swagger.yaml

func configureFlags(api *operations.ItinerariesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func makeDescription(from, to string) string {
	return fmt.Sprintf("Itinerary: %s - %s", from, to)
}

func getItineraries(from, to *string) ([]*models.Itinerary, error) {
	itineraryGraph, err := itinerary_utils.NewGraph()
	if err != nil {
	}
	maxNumberOfPaths := 10
	From := "NCE"
	To := "JFK"
	solutions, err := itineraryGraph.Compute(From, To, maxNumberOfPaths)
	var itineraries []*models.Itinerary
	if err != nil {
		return itineraries, err
	}
	for _, solution := range solutions {
		var steps []*models.ItineraryStep
		for _, segment := range solution {
			s := &models.ItineraryStep{
				From: string(segment.FromAirport),
				To:   string(segment.ToAirport),
			}
			steps = append(steps, s)
		}
		if len(steps) != 0 {
			i := &models.Itinerary{
				Description: makeDescription(*from, *to),
				Steps:       steps,
				Distance:    1,
			}
			itineraries = append(itineraries, i)
		}
	}
	return itineraries, nil
}

func isAlive() bool {
	return true // TODO: fill it
}

func isReady() bool {
	return true // TODO: fill it
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
		mergedParam := itineraries.NewGetItinerariesParams()
		if params.From != nil {
			mergedParam.From = params.From
		}
		if params.To != nil {
			mergedParam.To = params.To
		}

		if mergedParam.From == nil || len(*mergedParam.From) == 0 {
			errorMessage := "Missing `from` parameter"
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &errorMessage})
		}
		if mergedParam.To == nil || len(*mergedParam.To) == 0 {
			errorMessage := "Missing 'to' parameter"
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &errorMessage})
		}

		its, err := getItineraries(mergedParam.From, mergedParam.To)
		if err != nil {
			errorMessage := fmt.Sprintf("Cannot get itineraries: %v", err)
			return itineraries.NewGetItinerariesBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: &errorMessage})
		}
		if len(its) == 0 {
			errorMessage := fmt.Sprintf("No itineraries found")
			return itineraries.NewGetItinerariesNotFound().WithPayload(&models.Error{Code: http.StatusNotFound, Message: &errorMessage})
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
