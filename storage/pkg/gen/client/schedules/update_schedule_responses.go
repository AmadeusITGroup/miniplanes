// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
)

// UpdateScheduleReader is a Reader for the UpdateSchedule structure.
type UpdateScheduleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateScheduleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewUpdateScheduleCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewUpdateScheduleBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewUpdateScheduleNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateScheduleCreated creates a UpdateScheduleCreated with default headers values
func NewUpdateScheduleCreated() *UpdateScheduleCreated {
	return &UpdateScheduleCreated{}
}

/*UpdateScheduleCreated handles this case with default header values.

Updated Succesfully
*/
type UpdateScheduleCreated struct {
	Payload *models.Schedule
}

func (o *UpdateScheduleCreated) Error() string {
	return fmt.Sprintf("[PUT /schedules/{id}][%d] updateScheduleCreated  %+v", 201, o.Payload)
}

func (o *UpdateScheduleCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Schedule)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateScheduleBadRequest creates a UpdateScheduleBadRequest with default headers values
func NewUpdateScheduleBadRequest() *UpdateScheduleBadRequest {
	return &UpdateScheduleBadRequest{}
}

/*UpdateScheduleBadRequest handles this case with default header values.

invalid id
*/
type UpdateScheduleBadRequest struct {
}

func (o *UpdateScheduleBadRequest) Error() string {
	return fmt.Sprintf("[PUT /schedules/{id}][%d] updateScheduleBadRequest ", 400)
}

func (o *UpdateScheduleBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateScheduleNotFound creates a UpdateScheduleNotFound with default headers values
func NewUpdateScheduleNotFound() *UpdateScheduleNotFound {
	return &UpdateScheduleNotFound{}
}

/*UpdateScheduleNotFound handles this case with default header values.

schedule not found
*/
type UpdateScheduleNotFound struct {
}

func (o *UpdateScheduleNotFound) Error() string {
	return fmt.Sprintf("[PUT /schedules/{id}][%d] updateScheduleNotFound ", 404)
}

func (o *UpdateScheduleNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}