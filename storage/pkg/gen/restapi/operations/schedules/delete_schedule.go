// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteScheduleHandlerFunc turns a function with the right signature into a delete schedule handler
type DeleteScheduleHandlerFunc func(DeleteScheduleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteScheduleHandlerFunc) Handle(params DeleteScheduleParams) middleware.Responder {
	return fn(params)
}

// DeleteScheduleHandler interface for that can handle valid delete schedule params
type DeleteScheduleHandler interface {
	Handle(DeleteScheduleParams) middleware.Responder
}

// NewDeleteSchedule creates a new http.Handler for the delete schedule operation
func NewDeleteSchedule(ctx *middleware.Context, handler DeleteScheduleHandler) *DeleteSchedule {
	return &DeleteSchedule{Context: ctx, Handler: handler}
}

/*DeleteSchedule swagger:route DELETE /schedules/{id} schedules deleteSchedule

Delete an existant schedules

*/
type DeleteSchedule struct {
	Context *middleware.Context
	Handler DeleteScheduleHandler
}

func (o *DeleteSchedule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteScheduleParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
