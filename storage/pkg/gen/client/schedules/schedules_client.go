// Code generated by go-swagger; DO NOT EDIT.

package schedules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new schedules API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for schedules API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetSchedules get schedules API
*/
func (a *Client) GetSchedules(params *GetSchedulesParams) (*GetSchedulesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetSchedulesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetSchedules",
		Method:             "GET",
		PathPattern:        "/schedules",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetSchedulesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetSchedulesOK), nil

}

/*
AddSchedule Creates a new schedule. Duplicates are not allowed
*/
func (a *Client) AddSchedule(params *AddScheduleParams) (*AddScheduleCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddScheduleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "addSchedule",
		Method:             "POST",
		PathPattern:        "/schedules",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &AddScheduleReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*AddScheduleCreated), nil

}

/*
DeleteSchedule Delete an existant schedules
*/
func (a *Client) DeleteSchedule(params *DeleteScheduleParams) (*DeleteScheduleNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteScheduleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteSchedule",
		Method:             "DELETE",
		PathPattern:        "/schedules/{id}",
		ProducesMediaTypes: []string{""},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteScheduleReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteScheduleNoContent), nil

}

/*
GetSchedule get schedule API
*/
func (a *Client) GetSchedule(params *GetScheduleParams) (*GetScheduleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetScheduleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getSchedule",
		Method:             "GET",
		PathPattern:        "/schedules/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetScheduleReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetScheduleOK), nil

}

/*
UpdateSchedule Updates an existant Schedule
*/
func (a *Client) UpdateSchedule(params *UpdateScheduleParams) (*UpdateScheduleCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateScheduleParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updateSchedule",
		Method:             "PUT",
		PathPattern:        "/schedules/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateScheduleReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateScheduleCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
