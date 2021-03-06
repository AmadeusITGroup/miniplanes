// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

// AddScheduleCreatedCode is the HTTP code returned for type AddScheduleCreated
const AddScheduleCreatedCode int = 201

/*AddScheduleCreated Created

swagger:response addScheduleCreated
*/
type AddScheduleCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Schedule `json:"body,omitempty"`
}

// NewAddScheduleCreated creates AddScheduleCreated with default headers values
func NewAddScheduleCreated() *AddScheduleCreated {

	return &AddScheduleCreated{}
}

// WithPayload adds the payload to the add schedule created response
func (o *AddScheduleCreated) WithPayload(payload *models.Schedule) *AddScheduleCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add schedule created response
func (o *AddScheduleCreated) SetPayload(payload *models.Schedule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddScheduleCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddScheduleDefault unexpected error

swagger:response addScheduleDefault
*/
type AddScheduleDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddScheduleDefault creates AddScheduleDefault with default headers values
func NewAddScheduleDefault(code int) *AddScheduleDefault {
	if code <= 0 {
		code = 500
	}

	return &AddScheduleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add schedule default response
func (o *AddScheduleDefault) WithStatusCode(code int) *AddScheduleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add schedule default response
func (o *AddScheduleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the add schedule default response
func (o *AddScheduleDefault) WithPayload(payload *models.Error) *AddScheduleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add schedule default response
func (o *AddScheduleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddScheduleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
