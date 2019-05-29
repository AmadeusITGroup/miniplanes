// Code generated by go-swagger; DO NOT EDIT.

package airports

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

// NewAddAirportParams creates a new AddAirportParams object
// with the default values initialized.
func NewAddAirportParams() *AddAirportParams {
	var ()
	return &AddAirportParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAddAirportParamsWithTimeout creates a new AddAirportParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAddAirportParamsWithTimeout(timeout time.Duration) *AddAirportParams {
	var ()
	return &AddAirportParams{

		timeout: timeout,
	}
}

// NewAddAirportParamsWithContext creates a new AddAirportParams object
// with the default values initialized, and the ability to set a context for a request
func NewAddAirportParamsWithContext(ctx context.Context) *AddAirportParams {
	var ()
	return &AddAirportParams{

		Context: ctx,
	}
}

// NewAddAirportParamsWithHTTPClient creates a new AddAirportParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAddAirportParamsWithHTTPClient(client *http.Client) *AddAirportParams {
	var ()
	return &AddAirportParams{
		HTTPClient: client,
	}
}

/*AddAirportParams contains all the parameters to send to the API endpoint
for the add airport operation typically these are written to a http.Request
*/
type AddAirportParams struct {

	/*Airport
	  Airport

	*/
	Airport *models.Airport

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the add airport params
func (o *AddAirportParams) WithTimeout(timeout time.Duration) *AddAirportParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add airport params
func (o *AddAirportParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add airport params
func (o *AddAirportParams) WithContext(ctx context.Context) *AddAirportParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add airport params
func (o *AddAirportParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add airport params
func (o *AddAirportParams) WithHTTPClient(client *http.Client) *AddAirportParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add airport params
func (o *AddAirportParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAirport adds the airport to the add airport params
func (o *AddAirportParams) WithAirport(airport *models.Airport) *AddAirportParams {
	o.SetAirport(airport)
	return o
}

// SetAirport adds the airport to the add airport params
func (o *AddAirportParams) SetAirport(airport *models.Airport) {
	o.Airport = airport
}

// WriteToRequest writes these params to a swagger request
func (o *AddAirportParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Airport != nil {
		if err := r.SetBodyParam(o.Airport); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}