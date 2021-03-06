// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// DeleteScheduleNoContentCode is the HTTP code returned for type DeleteScheduleNoContent
const DeleteScheduleNoContentCode int = 204

/*DeleteScheduleNoContent Deleted successfully

swagger:response deleteScheduleNoContent
*/
type DeleteScheduleNoContent struct {
}

// NewDeleteScheduleNoContent creates DeleteScheduleNoContent with default headers values
func NewDeleteScheduleNoContent() *DeleteScheduleNoContent {

	return &DeleteScheduleNoContent{}
}

// WriteResponse to the client
func (o *DeleteScheduleNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteScheduleBadRequestCode is the HTTP code returned for type DeleteScheduleBadRequest
const DeleteScheduleBadRequestCode int = 400

/*DeleteScheduleBadRequest Invalid ID

swagger:response deleteScheduleBadRequest
*/
type DeleteScheduleBadRequest struct {
}

// NewDeleteScheduleBadRequest creates DeleteScheduleBadRequest with default headers values
func NewDeleteScheduleBadRequest() *DeleteScheduleBadRequest {

	return &DeleteScheduleBadRequest{}
}

// WriteResponse to the client
func (o *DeleteScheduleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// DeleteScheduleNotFoundCode is the HTTP code returned for type DeleteScheduleNotFound
const DeleteScheduleNotFoundCode int = 404

/*DeleteScheduleNotFound Schedule not found

swagger:response deleteScheduleNotFound
*/
type DeleteScheduleNotFound struct {
}

// NewDeleteScheduleNotFound creates DeleteScheduleNotFound with default headers values
func NewDeleteScheduleNotFound() *DeleteScheduleNotFound {

	return &DeleteScheduleNotFound{}
}

// WriteResponse to the client
func (o *DeleteScheduleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}
