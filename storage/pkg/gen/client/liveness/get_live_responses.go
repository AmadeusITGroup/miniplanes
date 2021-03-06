// Code generated by go-swagger; DO NOT EDIT.

package liveness

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

// GetLiveReader is a Reader for the GetLive structure.
type GetLiveReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetLiveReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetLiveOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 503:
		result := NewGetLiveServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewGetLiveDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetLiveOK creates a GetLiveOK with default headers values
func NewGetLiveOK() *GetLiveOK {
	return &GetLiveOK{}
}

/*GetLiveOK handles this case with default header values.

liveness probe
*/
type GetLiveOK struct {
}

func (o *GetLiveOK) Error() string {
	return fmt.Sprintf("[GET /live][%d] getLiveOK ", 200)
}

func (o *GetLiveOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetLiveServiceUnavailable creates a GetLiveServiceUnavailable with default headers values
func NewGetLiveServiceUnavailable() *GetLiveServiceUnavailable {
	return &GetLiveServiceUnavailable{}
}

/*GetLiveServiceUnavailable handles this case with default header values.

if not alive
*/
type GetLiveServiceUnavailable struct {
	Payload *models.Error
}

func (o *GetLiveServiceUnavailable) Error() string {
	return fmt.Sprintf("[GET /live][%d] getLiveServiceUnavailable  %+v", 503, o.Payload)
}

func (o *GetLiveServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLiveDefault creates a GetLiveDefault with default headers values
func NewGetLiveDefault(code int) *GetLiveDefault {
	return &GetLiveDefault{
		_statusCode: code,
	}
}

/*GetLiveDefault handles this case with default header values.

generic error response
*/
type GetLiveDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get live default response
func (o *GetLiveDefault) Code() int {
	return o._statusCode
}

func (o *GetLiveDefault) Error() string {
	return fmt.Sprintf("[GET /live][%d] GetLive default  %+v", o._statusCode, o.Payload)
}

func (o *GetLiveDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
