package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type CVSFileSchedule struct {
	Origin      string
	Destination string
	//Via              string
	FlightNumber     string
	OperatingCarrier string
	DaysOperated     string // 12 means operated Monday, Tuesday... unused for the moment
	Departure        string
	Arrival          string
	//EffFrom          string // date in format 'dd Mon YYYY'
	//EffTill          string // date in format dd-Mon-YY
}

func (s *CVSFileSchedule) ToStrings() []string {
	return []string{
		s.Origin,
		s.Destination,
		s.FlightNumber,
		s.OperatingCarrier,
		s.DaysOperated,
		s.Departure,
		s.Arrival,
	}
}

type airport struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	AirportID string        `json:"airportID" bson:"airportID"` // Unique OpenFlights identifier for this airport.
	Name      string        `json:"name" bson:"name"`           //Name of airport. May or may not contain the City name.
	City      string        `json:"city" bson:"city"`           // Main city served by airport. May be spelled differently from Name.
	Country   string        `json:"country" bson:"country"`     // Country or territory where airport is located. See countries.dat to cross-reference to ISO 3166-1 codes.
	IATA      string        `json:"IATA" bson:"IATA"`           //3-letter IATA code. Null if not assigned/unknown.
	ICAO      string        `json:"ICAO" bson:"ICAO"`           //4-letter ICAO code. Null if not assigned.
	Latitude  string        `json:"latitude" bson:"latitude"`   //Decimal degrees, usually to six significant digits. Negative is South, positive is North.
	Longitude string        `json:"longitude" bson:"longitude"` //Decimal degrees, usually to six significant digits. Negative is West, positive is East.
	Altitude  string        `json:"altitude" bson:"altitude"`   //In feet.
	Timezone  string        `json:"timezone" bson:"timezone"`   //Hours offset from UTC. Fractional hours are expressed as decimals, eg. India is 5.5.
	DST       string        `json:"DST" bson:"DST"`             //Daylight savings time. One of E (Europe), A (US/Canada), S (South America), O (Australia), Z (New Zealand), N (None) or U (Unknown). See also: Help: Time
	TZ        string        `json:"TZ" bson:"TZ"`               //database time zone	Timezone in "tz" (Olson) format, eg. "America/Los_Angeles".
	Type      string        `json:"type" bson:"type"`           //	Type of the airport. Value "airport" for air terminals, "station" for train stations, "port" for ferry terminals and "unknown" if not known. In airports.csv, only type=airport is included.
	Source    string        `json:"source" bson:"source"`       //Source of this data. "OurAirports" for data sourced from OurAirports, "Legacy" for old data not matched to OurAirports (mostly DAFIF), "User" for unverified user contributions. In airports.csv, only source=OurAirports is included.
}

type airportMap struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	PseudoName string        `json:"pseudoName" bson:"pseudoName"`
	Name       string        `json:"name" bson:"name"`
	IATA       string        `json:"IATA" bson:"IATA"`
}

const (
	dbName                = "miniapp"
	routesCollection      = "routes"
	airportsCollection    = "airports"
	airlinesCollection    = "airlines"
	airportsMapCollection = "airportsmap"
)

func findAirportID(airportName string, airports []*airport) (string, error) {
	for i := range airports {
		if airportName == airports[i].Name {
			return airports[i].AirportID, nil
		}
	}
	return "", fmt.Errorf("couldn't find airport %q", airportName)
}

func getIATACodeFromName(airports []*airport, airportsMap []*airportMap, airportName string) (string, error) {
	if len(airportName) == 0 {
		return "", fmt.Errorf("Empty airport name given")
	}
	for _, a := range airports {
		if airportName == a.Name {
			if len(a.IATA) == 0 {
				log.Fatalf("Empty IATA code for %q", a.Name)
			}
			return a.IATA, nil
		}
	}
	for _, a := range airportsMap {
		if airportName == a.Name {
			if len(a.IATA) == 0 {
				log.Fatalf("Empty IATA code for in airportMap %q", a.Name)
			}
			return a.IATA, nil
		}
	}

	// Now best effort :)
	withAirport := airportName + " Airport"
	for _, a := range airports {
		if withAirport == a.Name {
			if len(a.IATA) == 0 {
				log.Fatalf("Empty IATA code for in airportMap %q", a.Name)
			}
			return a.IATA, nil
		}
	}

	// with City
	for _, a := range airports {
		if airportName == a.City {
			if len(a.IATA) == 0 {
				log.Fatalf("Empty IATA code for %q", a.Name)
			}
		}
		return a.IATA, nil
	}

	return "", fmt.Errorf("can't find IATA code for %q", airportName)
}

type byCityName []*airport

func (b byCityName) Len() int {
	return len(b)
}

func (b byCityName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func main() {
	csvFile, _ := ioutil.ReadFile("JetAirwaysFlightSchedules.csv")

	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	var airports []*airport

	if err := db.DB(dbName).C(airportsCollection).Find(nil).All(&airports); err != nil {
		log.Fatalf("Cannot get airpts: %v\n", err)
	}

	var airportsMap []*airportMap
	if err := db.DB(dbName).C(airportsMapCollection).Find(nil).All(&airportsMap); err != nil {
		log.Fatalf("Cannot get airportsMap: %v\n", err)
	}

	s := string(csvFile)
	reader := csv.NewReader(strings.NewReader(s))
	writer := csv.NewWriter(os.Stdout)

	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if r[7] == "Daily" {
			r[7] = "1234567"
		}

		origin, err := getIATACodeFromName(airports, airportsMap, r[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Origin. Cannot obtain IATA code for %q\n", r[0])
			continue
		}

		destination, err := getIATACodeFromName(airports, airportsMap, r[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Destination. Cannot obtain IATA code for %q -> %#v\n", r[2], r)
			continue
		}

		if origin == destination {
			continue
		}

		s := CVSFileSchedule{
			Origin:      origin,
			Destination: destination,
			//Via:              r[3],
			FlightNumber:     r[4],
			OperatingCarrier: r[5],
			DaysOperated:     r[7],
			Departure:        r[8],
			Arrival:          r[10],
		}

		writer.Write(s.ToStrings())
		writer.Flush()
	}
	os.Exit(0)
}
