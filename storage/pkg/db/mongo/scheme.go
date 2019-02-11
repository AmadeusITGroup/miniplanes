package mongo

import (
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	"github.com/jinzhu/copier"

	"gopkg.in/mgo.v2/bson"
)

type Airline struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	AirlineID string        `json:"airlineID" bson:"airlineID"` // Unique OpenFlights identifier for this airline.
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

type Course struct {
	ID                   bson.ObjectId `json:"id" bson:"_id"`
	Airline              string        `json:"airline" bson:"airline"`                           //2-letter (IATA) or 3-letter (ICAO) code of the airline
	AirlineID            string        `json:"airlineID" bson:"airlineID"`                       // Unique OpenFlights identifier for this airline.
	SourceAirport        string        `json:"sourceAirport" bson:"sourceAirport"`               //3-letter (IATA) or 4-letter (ICAO) code of the source airport.
	SourceAirportID      int32         `json:"sourceAirportID" bson:"sourceAirportID"`           //Unique OpenFlights identifier for source airport (see Airport)
	DestinationAirport   string        `json:"destinationAirport" bson:"destinationAirport"`     //3-letter (IATA) or 4-letter (ICAO) code of the destination airport.
	DestinationAirportID int32         `json:"destinationAirportID" bson:"destinationAirportID"` //Unique OpenFlights identifier for destination airport (see Airport)
	Codeshare            string        `json:"codeshare" bson:"codeshare"`                       //"Y" if this flight is a codeshare (that is, not operated by Airline, but another carrier), empty otherwise.
	Stops                string        `json:"stops" bson:"stops"`                               //Number of stops on this flight ("0" for direct)
	Equipment            string        `json:"equipment" bson:"equipment"`                       //3-letter codes for plane type(s) generally used on this flight, separated by spaces

}

func (c Course) ToModel() (*models.Course, error) {
	v := &models.Course{}
	err := copier.Copy(v, &c)
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
	Origin           int32         `json:"origin" bson:"origin"`
	Destination      int32         `json:"destination" bson:"destination"`
	FlightNumber     string        `json:"flightNumber" bson:"flightNumber"`
	OperatingCarrier string        `json:"operatingCarrier" bson:"operatingCarrier"`
	DaysOperated     string        `json:"daysOperated" bson:"daysOperated"`
	Departure        string        `json:"departure" bson:"departure"`
	Arrival          string        `json:"arrival" bson:"arrival"`
	ArriveNextDay    bool          `json:"arriveNextDay" bson:"arriveNextDay"`
}

func (s Schedule) ToModel() (*models.Schedule, error) {
	v := &models.Schedule{}
	err := copier.Copy(v, &s)
	return v, err
}
