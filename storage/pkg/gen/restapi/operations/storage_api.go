// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	security "github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/airlines"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/airports"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/liveness"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/readiness"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/schedules"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations/version"
)

// NewStorageAPI creates a new Storage instance
func NewStorageAPI(spec *loads.Document) *StorageAPI {
	return &StorageAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		AirlinesGetAirlinesHandler: airlines.GetAirlinesHandlerFunc(func(params airlines.GetAirlinesParams) middleware.Responder {
			return middleware.NotImplemented("operation AirlinesGetAirlines has not yet been implemented")
		}),
		AirportsGetAirportsHandler: airports.GetAirportsHandlerFunc(func(params airports.GetAirportsParams) middleware.Responder {
			return middleware.NotImplemented("operation AirportsGetAirports has not yet been implemented")
		}),
		LivenessGetLiveHandler: liveness.GetLiveHandlerFunc(func(params liveness.GetLiveParams) middleware.Responder {
			return middleware.NotImplemented("operation LivenessGetLive has not yet been implemented")
		}),
		ReadinessGetReadyHandler: readiness.GetReadyHandlerFunc(func(params readiness.GetReadyParams) middleware.Responder {
			return middleware.NotImplemented("operation ReadinessGetReady has not yet been implemented")
		}),
		SchedulesGetSchedulesHandler: schedules.GetSchedulesHandlerFunc(func(params schedules.GetSchedulesParams) middleware.Responder {
			return middleware.NotImplemented("operation SchedulesGetSchedules has not yet been implemented")
		}),
		VersionGetVersionHandler: version.GetVersionHandlerFunc(func(params version.GetVersionParams) middleware.Responder {
			return middleware.NotImplemented("operation VersionGetVersion has not yet been implemented")
		}),
		AirlinesAddAirlineHandler: airlines.AddAirlineHandlerFunc(func(params airlines.AddAirlineParams) middleware.Responder {
			return middleware.NotImplemented("operation AirlinesAddAirline has not yet been implemented")
		}),
		AirportsAddAirportHandler: airports.AddAirportHandlerFunc(func(params airports.AddAirportParams) middleware.Responder {
			return middleware.NotImplemented("operation AirportsAddAirport has not yet been implemented")
		}),
		SchedulesAddScheduleHandler: schedules.AddScheduleHandlerFunc(func(params schedules.AddScheduleParams) middleware.Responder {
			return middleware.NotImplemented("operation SchedulesAddSchedule has not yet been implemented")
		}),
		SchedulesDeleteScheduleHandler: schedules.DeleteScheduleHandlerFunc(func(params schedules.DeleteScheduleParams) middleware.Responder {
			return middleware.NotImplemented("operation SchedulesDeleteSchedule has not yet been implemented")
		}),
		SchedulesGetScheduleHandler: schedules.GetScheduleHandlerFunc(func(params schedules.GetScheduleParams) middleware.Responder {
			return middleware.NotImplemented("operation SchedulesGetSchedule has not yet been implemented")
		}),
		SchedulesUpdateScheduleHandler: schedules.UpdateScheduleHandlerFunc(func(params schedules.UpdateScheduleParams) middleware.Responder {
			return middleware.NotImplemented("operation SchedulesUpdateSchedule has not yet been implemented")
		}),
	}
}

/*StorageAPI needs to add a description */
type StorageAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer runtime.Producer

	// AirlinesGetAirlinesHandler sets the operation handler for the get airlines operation
	AirlinesGetAirlinesHandler airlines.GetAirlinesHandler
	// AirportsGetAirportsHandler sets the operation handler for the get airports operation
	AirportsGetAirportsHandler airports.GetAirportsHandler
	// LivenessGetLiveHandler sets the operation handler for the get live operation
	LivenessGetLiveHandler liveness.GetLiveHandler
	// ReadinessGetReadyHandler sets the operation handler for the get ready operation
	ReadinessGetReadyHandler readiness.GetReadyHandler
	// SchedulesGetSchedulesHandler sets the operation handler for the get schedules operation
	SchedulesGetSchedulesHandler schedules.GetSchedulesHandler
	// VersionGetVersionHandler sets the operation handler for the get version operation
	VersionGetVersionHandler version.GetVersionHandler
	// AirlinesAddAirlineHandler sets the operation handler for the add airline operation
	AirlinesAddAirlineHandler airlines.AddAirlineHandler
	// AirportsAddAirportHandler sets the operation handler for the add airport operation
	AirportsAddAirportHandler airports.AddAirportHandler
	// SchedulesAddScheduleHandler sets the operation handler for the add schedule operation
	SchedulesAddScheduleHandler schedules.AddScheduleHandler
	// SchedulesDeleteScheduleHandler sets the operation handler for the delete schedule operation
	SchedulesDeleteScheduleHandler schedules.DeleteScheduleHandler
	// SchedulesGetScheduleHandler sets the operation handler for the get schedule operation
	SchedulesGetScheduleHandler schedules.GetScheduleHandler
	// SchedulesUpdateScheduleHandler sets the operation handler for the update schedule operation
	SchedulesUpdateScheduleHandler schedules.UpdateScheduleHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *StorageAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *StorageAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *StorageAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *StorageAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *StorageAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *StorageAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *StorageAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the StorageAPI
func (o *StorageAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.AirlinesGetAirlinesHandler == nil {
		unregistered = append(unregistered, "airlines.GetAirlinesHandler")
	}

	if o.AirportsGetAirportsHandler == nil {
		unregistered = append(unregistered, "airports.GetAirportsHandler")
	}

	if o.LivenessGetLiveHandler == nil {
		unregistered = append(unregistered, "liveness.GetLiveHandler")
	}

	if o.ReadinessGetReadyHandler == nil {
		unregistered = append(unregistered, "readiness.GetReadyHandler")
	}

	if o.SchedulesGetSchedulesHandler == nil {
		unregistered = append(unregistered, "schedules.GetSchedulesHandler")
	}

	if o.VersionGetVersionHandler == nil {
		unregistered = append(unregistered, "version.GetVersionHandler")
	}

	if o.AirlinesAddAirlineHandler == nil {
		unregistered = append(unregistered, "airlines.AddAirlineHandler")
	}

	if o.AirportsAddAirportHandler == nil {
		unregistered = append(unregistered, "airports.AddAirportHandler")
	}

	if o.SchedulesAddScheduleHandler == nil {
		unregistered = append(unregistered, "schedules.AddScheduleHandler")
	}

	if o.SchedulesDeleteScheduleHandler == nil {
		unregistered = append(unregistered, "schedules.DeleteScheduleHandler")
	}

	if o.SchedulesGetScheduleHandler == nil {
		unregistered = append(unregistered, "schedules.GetScheduleHandler")
	}

	if o.SchedulesUpdateScheduleHandler == nil {
		unregistered = append(unregistered, "schedules.UpdateScheduleHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *StorageAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *StorageAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	return nil

}

// Authorizer returns the registered authorizer
func (o *StorageAPI) Authorizer() runtime.Authorizer {

	return nil

}

// ConsumersFor gets the consumers for the specified media types
func (o *StorageAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *StorageAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *StorageAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the storage API
func (o *StorageAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *StorageAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/airlines"] = airlines.NewGetAirlines(o.context, o.AirlinesGetAirlinesHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/airports"] = airports.NewGetAirports(o.context, o.AirportsGetAirportsHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/live"] = liveness.NewGetLive(o.context, o.LivenessGetLiveHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/ready"] = readiness.NewGetReady(o.context, o.ReadinessGetReadyHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/schedules"] = schedules.NewGetSchedules(o.context, o.SchedulesGetSchedulesHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/version"] = version.NewGetVersion(o.context, o.VersionGetVersionHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/airlines"] = airlines.NewAddAirline(o.context, o.AirlinesAddAirlineHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/airports"] = airports.NewAddAirport(o.context, o.AirportsAddAirportHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/schedules"] = schedules.NewAddSchedule(o.context, o.SchedulesAddScheduleHandler)

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/schedules/{id}"] = schedules.NewDeleteSchedule(o.context, o.SchedulesDeleteScheduleHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/schedules/{id}"] = schedules.NewGetSchedule(o.context, o.SchedulesGetScheduleHandler)

	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/schedules/{id}"] = schedules.NewUpdateSchedule(o.context, o.SchedulesUpdateScheduleHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *StorageAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *StorageAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *StorageAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *StorageAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}
