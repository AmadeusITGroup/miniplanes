// Code generated by go-swagger; DO NOT EDIT.

package liveness

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/amadeusitgroup/miniapp/storage/pkg/models"
)

// GetLiveOKCode is the HTTP code returned for type GetLiveOK
const GetLiveOKCode int = 200

/*GetLiveOK liveness probe

swagger:response getLiveOK
*/
type GetLiveOK struct {
}

// NewGetLiveOK creates GetLiveOK with default headers values
func NewGetLiveOK() *GetLiveOK {

	return &GetLiveOK{}
}

// WriteResponse to the client
func (o *GetLiveOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// GetLiveServiceUnavailableCode is the HTTP code returned for type GetLiveServiceUnavailable
const GetLiveServiceUnavailableCode int = 503

/*GetLiveServiceUnavailable if not alive

swagger:response getLiveServiceUnavailable
*/
type GetLiveServiceUnavailable struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLiveServiceUnavailable creates GetLiveServiceUnavailable with default headers values
func NewGetLiveServiceUnavailable() *GetLiveServiceUnavailable {

	return &GetLiveServiceUnavailable{}
}

// WithPayload adds the payload to the get live service unavailable response
func (o *GetLiveServiceUnavailable) WithPayload(payload *models.Error) *GetLiveServiceUnavailable {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get live service unavailable response
func (o *GetLiveServiceUnavailable) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLiveServiceUnavailable) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(503)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetLiveDefault generic error response

swagger:response getLiveDefault
*/
type GetLiveDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLiveDefault creates GetLiveDefault with default headers values
func NewGetLiveDefault(code int) *GetLiveDefault {
	if code <= 0 {
		code = 500
	}

	return &GetLiveDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get live default response
func (o *GetLiveDefault) WithStatusCode(code int) *GetLiveDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get live default response
func (o *GetLiveDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get live default response
func (o *GetLiveDefault) WithPayload(payload *models.Error) *GetLiveDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get live default response
func (o *GetLiveDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLiveDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}