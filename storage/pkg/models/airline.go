// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Airline airline
// swagger:model airline
type Airline struct {

	// active
	Active bool `json:"Active,omitempty"`

	// i a t a
	IATA string `json:"IATA,omitempty"`

	// name
	Name string `json:"Name,omitempty"`
}

// Validate validates this airline
func (m *Airline) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Airline) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Airline) UnmarshalBinary(b []byte) error {
	var res Airline
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
