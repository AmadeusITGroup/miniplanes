/*

MIT License

Copyright (c) 2019 Amadeus s.a.s.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package mongo

import (
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/jinzhu/copier"

	"gopkg.in/mgo.v2/bson"
)

type Airline struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	AirlineID int64         `json:"airlineID" bson:"airlineID"` // Unique OpenFlights identifier for this airline.
	Name      string        `json:"name" bson:"name"`           // Name of the airline.
	Alias     string        `json:"alias" bson:"alias"`         //Alias of the airline. For example, All Nippon Airways is commonly known as "ANA".
	IATA      string        `json:"IATA" bson:"IATA"`           //2-letter IATA code, if available.
	ICAO      string        `json:"ICAO" bson:"ICAO"`           //3-letter ICAO code, if available.
	Callsign  string        `json:"callsign" bson:"callsign"`   //Airline callsign.
	Country   string        `json:"country" bson:"country"`     //Country or territory where airline is incorporated.
	Active    string        `json:"active" bson:"active"`       //"Y" if the airline is or has until recently been operational, "N" if it is defunct. This field is not reliable: in particular, major airlines that stopped flying long ago, but have not had their IATA code reassigned (eg. Ansett/AN), will incorrectly show as "Y".

}

func (a Airline) ToModel() (*models.Airline, error) {
	v := &models.Airline{}
	err := copier.Copy(v, &a)
	return v, err
}

type Airport struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	AirportID int32         `json:"airportID" bson:"airportID"` // Unique OpenFlights identifier for this airport.
	Name      string        `json:"name" bson:"name"`           //Name of airport. May or may not contain the City name.
	City      string        `json:"city" bson:"city"`           // Main city served by airport. May be spelled differently from Name.
	Country   string        `json:"country" bson:"country"`     // Country or territory where airport is located. See countries.dat to cross-reference to ISO 3166-1 codes.
	IATA      string        `json:"IATA" bson:"IATA"`           //3-letter IATA code. Null if not assigned/unknown.
	ICAO      string        `json:"ICAO" bson:"ICAO"`           //4-letter ICAO code. Null if not assigned.
	Latitude  float64       `json:"latitude" bson:"latitude"`   //Decimal degrees, usually to six significant digits. Negative is South, positive is North.
	Longitude float64       `json:"longitude" bson:"longitude"` //Decimal degrees, usually to six significant digits. Negative is West, positive is East.
	Altitude  string        `json:"altitude" bson:"altitude"`   //In feet.
	Timezone  string        `json:"timezone" bson:"timezone"`   //Hours offset from UTC. Fractional hours are expressed as decimals, eg. India is 5.5.
	DST       string        `json:"DST" bson:"DST"`             //Daylight savings time. One of E (Europe), A (US/Canada), S (South America), O (Australia), Z (New Zealand), N (None) or U (Unknown). See also: Help: Time
	TZ        string        `json:"TZ" bson:"TZ"`               //database time zone	Timezone in "tz" (Olson) format, eg. "America/Los_Angeles".
	Type      string        `json:"type" bson:"type"`           //	Type of the airport. Value "airport" for air terminals, "station" for train stations, "port" for ferry terminals and "unknown" if not known. In airports.csv, only type=airport is included.
	Source    string        `json:"source" bson:"source"`       //Source of this data. "OurAirports" for data sourced from OurAirports, "Legacy" for old data not matched to OurAirports (mostly DAFIF), "User" for unverified user contributions. In airports.csv, only source=OurAirports is included.

}

func (a Airport) ToModel() (*models.Airport, error) {
	v := &models.Airport{}
	err := copier.Copy(v, &a)
	return v, err
}

type Schedule struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	ScheduleID       int64         `json:"scheduleID" bson:"scheduleID"`
	Origin           int32         `json:"origin" bson:"origin"`
	Destination      int32         `json:"destination" bson:"destination"`
	FlightNumber     string        `json:"flightNumber" bson:"flightNumber"`
	OperatingCarrier string        `json:"operatingCarrier" bson:"operatingCarrier"`
	DaysOperated     string        `json:"daysOperated" bson:"daysOperated"`
	DepartureTime    string        `json:"departureTime" bson:"departureTime"`
}

func (s Schedule) ToModel() (*models.Schedule, error) {
	v := &models.Schedule{}
	err := copier.Copy(v, &s)
	return v, err
}
