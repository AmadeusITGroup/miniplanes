// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteScheduleParams creates a new DeleteScheduleParams object
// with the default values initialized.
func NewDeleteScheduleParams() *DeleteScheduleParams {
	var ()
	return &DeleteScheduleParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteScheduleParamsWithTimeout creates a new DeleteScheduleParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteScheduleParamsWithTimeout(timeout time.Duration) *DeleteScheduleParams {
	var ()
	return &DeleteScheduleParams{

		timeout: timeout,
	}
}

// NewDeleteScheduleParamsWithContext creates a new DeleteScheduleParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteScheduleParamsWithContext(ctx context.Context) *DeleteScheduleParams {
	var ()
	return &DeleteScheduleParams{

		Context: ctx,
	}
}

// NewDeleteScheduleParamsWithHTTPClient creates a new DeleteScheduleParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteScheduleParamsWithHTTPClient(client *http.Client) *DeleteScheduleParams {
	var ()
	return &DeleteScheduleParams{
		HTTPClient: client,
	}
}

/*DeleteScheduleParams contains all the parameters to send to the API endpoint
for the delete schedule operation typically these are written to a http.Request
*/
type DeleteScheduleParams struct {

	/*ID
	  The id of the item

	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete schedule params
func (o *DeleteScheduleParams) WithTimeout(timeout time.Duration) *DeleteScheduleParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete schedule params
func (o *DeleteScheduleParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete schedule params
func (o *DeleteScheduleParams) WithContext(ctx context.Context) *DeleteScheduleParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete schedule params
func (o *DeleteScheduleParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete schedule params
func (o *DeleteScheduleParams) WithHTTPClient(client *http.Client) *DeleteScheduleParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete schedule params
func (o *DeleteScheduleParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete schedule params
func (o *DeleteScheduleParams) WithID(id int64) *DeleteScheduleParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete schedule params
func (o *DeleteScheduleParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteScheduleParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
