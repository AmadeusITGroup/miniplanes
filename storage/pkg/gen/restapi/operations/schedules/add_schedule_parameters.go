// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	models "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

// NewAddScheduleParams creates a new AddScheduleParams object
// no default values defined in spec.
func NewAddScheduleParams() AddScheduleParams {

	return AddScheduleParams{}
}

// AddScheduleParams contains all the bound params for the add schedule operation
// typically these are obtained from a http.Request
//
// swagger:parameters addSchedule
type AddScheduleParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Schedule
	  Required: true
	  In: body
	*/
	Schedule *models.Schedule
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewAddScheduleParams() beforehand.
func (o *AddScheduleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Schedule
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("schedule", "body"))
			} else {
				res = append(res, errors.NewParseError("schedule", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Schedule = &body
			}
		}
	} else {
		res = append(res, errors.Required("schedule", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
