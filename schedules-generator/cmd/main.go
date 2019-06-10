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
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"

	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/engine"
	airportsclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airports"
	schedulesclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/schedules"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	httptransport "github.com/go-openapi/runtime/client"
)

type Route struct {
	// airline
	Airline string
	// airline ID
	AirlineID int64
	// code share
	CodeShare string
	// destination airport
	DestinationAirport string
	// destination airport ID
	DestinationAirportID int32
	// equipment
	Equipment string
	// source airport
	SourceAirport string
	// source airport ID
	SourceAirportID int32
	// stops
	Stops int64
}

func computeDepartureArrivals(origin, destination *models.Airport) ([]string, error) {
	arrivalTimes := []string{}

	if origin == nil || destination == nil {
		return arrivalTimes, fmt.Errorf("missing origin or destination airport")
	}
	if origin == destination {
		return arrivalTimes, fmt.Errorf("origin and destination, same airport")
	}
	if _, err := time.LoadLocation(destination.TZ); err != nil {
		return arrivalTimes, fmt.Errorf("bad TZ for destination airport %q: %v", destination.Name, err)
	}
	if _, err := time.LoadLocation(origin.TZ); err != nil {
		return arrivalTimes, fmt.Errorf("bad TZ for origin airport %q: %v", origin.Name, err)
	}

	/*distance := db.Distance(origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude)
	distanceKm := float64(distance / 1000)
	formattedHourDuration := fmt.Sprintf("%fh", (engine.FlightOverhead + (distanceKm / engine.AverageSpeedKmH)))
	flightDuration, err := time.ParseDuration(formattedHourDuration)
	if err != nil {
		return arrivalTimes, fmt.Errorf("unable to parse duration: %s", formattedHourDuration)
	}

	log.Infof("Flight duration %q->%q:%v\n", origin.City, destination.City, flightDuration)
	now := time.Now() // to get year, month, day
	*/

	departureTimes := []string{"0800", "1100", "1500", "1700"}

	for _, departureTime := range departureTimes {
		_, arrivalTime, err := engine.ComputeArrivalDateTime(2019, "2412", departureTime, origin, destination)
		if err != nil {
			log.Errorf("Unable to compute Arrival Time for %s -> %s, departureTime = %s: %v", origin.IATA, destination.IATA, departureTime, err)
			continue
		}
		var hours, minutes int
		fmt.Sscanf(arrivalTime, "%02d%02d", &hours, &minutes)

		if hours < 6 || hours > 23 {
			continue
		}
		arrivalTimes = append(arrivalTimes, departureTime)
	}
	return arrivalTimes, nil
}

const (
	storageHostParamName = "storage-host"
	storageHostDefault   = "storage"
	storagePortParamName = "storage-port"
	storagePortDefault   = 12345
	routesFileParamName  = "routes-file"
	routesFileDefault    = "routes.dat"
)

var (
	storagePort int
	storageHost string
	routesFile  string
)

func init() {
	flag.IntVar(&storagePort, storagePortParamName, storagePortDefault, "the port of storage service")
	flag.StringVar(&storageHost, storageHostParamName, storageHostDefault, "the name of the storage service")
	flag.StringVar(&routesFile, routesFileParamName, routesFileDefault, "routes file")
}

func main() {

	flag.Parse()

	log.Infof("%s %d\n", storageHost, storagePort)
	storageURL := fmt.Sprintf("%s:%d", storageHost, storagePort)
	client := airportsclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
	OK, err := client.GetAirports(airportsclient.NewGetAirportsParams())
	if err != nil {
		log.Fatalf("Unable to load airports from storage %s : %v", storageURL, err)
	}
	ID2Airports := map[int32]*models.Airport{}
	for i := range OK.Payload {
		ID2Airports[OK.Payload[i].AirportID] = OK.Payload[i]
	}

	rFile, err := os.Open(routesFile)
	if err != nil {
		log.Fatalf("Unable to import routes: %v", err)
	}
	reader := csv.NewReader(bufio.NewReader(rFile))
	var routes []Route
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		airlineID, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			log.Errorf("Unable to convert airlinedID for line %s: %v", line, err)
			continue
		}

		destinationAirportID, err := strconv.ParseInt(line[5], 10, 32)
		if err != nil {
			log.Errorf("Unable to convert destinationAirportID for line %s: %v", line, err)
			continue
		}

		sourceAirportID, err := strconv.ParseInt(line[3], 10, 32)
		if err != nil {
			log.Errorf("Unable to convert sourceAirportID for line %s: %v", line, err)
			continue
		}
		stops, err := strconv.ParseInt(line[7], 10, 64)
		if err != nil {
			log.Errorf("Unable to convert stops for line %s: %v", line, err)
			continue
		}

		r := Route{
			Airline:              line[0],
			AirlineID:            airlineID,
			SourceAirport:        line[3],
			SourceAirportID:      int32(sourceAirportID),
			DestinationAirport:   line[3],
			DestinationAirportID: int32(destinationAirportID),
			CodeShare:            line[6],
			Stops:                stops,
			Equipment:            line[8]}
		routes = append(routes, r)
	}
	if err != nil || len(routes) == 0 {
		log.Fatalf("Unable to load routes: %v", err)
	}

	flightNumberPerAirline := map[string]int16{}
	sclient := schedulesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
	var scheduleID int64
	for i, route := range routes {
		scheduleID = int64(i)
		flightNumberPerAirline[route.Airline] = flightNumberPerAirline[route.Airline] + 1

		departureArrivals, err := computeDepartureArrivals(ID2Airports[route.SourceAirportID], ID2Airports[route.DestinationAirportID])
		if err != nil {
			log.Infof("Cannot compute arrival time... %v", err)
			continue
		}
		for _, d := range departureArrivals {
			log.Infof("DepartureDate: %s", d)
			flightNumber := strings.Join([]string{route.Airline, fmt.Sprintf("%03d", flightNumberPerAirline[route.Airline])}, "")
			daysOperated := "1234567"
			addSchedulePayload := schedulesclient.NewAddScheduleParams()
			addSchedulePayload.Schedule = &models.Schedule{
				ScheduleID:       &scheduleID,
				Origin:           &route.SourceAirportID,
				Destination:      &route.DestinationAirportID,
				FlightNumber:     &flightNumber,
				OperatingCarrier: &route.Airline,
				DaysOperated:     &daysOperated,
				DepartureTime:    &d,
			}
			log.Infof("Adding Schedule %+v", addSchedulePayload.Schedule)
			if _, err := sclient.AddSchedule(addSchedulePayload); err != nil {
				log.Errorf("Unable to add schedule: %v\n", err)
			}
		}
	}
	os.Exit(0)
}
