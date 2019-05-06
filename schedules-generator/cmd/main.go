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

	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/db"
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

type departureArrivalTime struct {
	departure     string
	arrival       string
	arriveNextDay bool
}

type pnrTime struct {
	h, m int8
	l    *time.Location
}

func NewPNRTime(h, m int8, timezone *time.Location) (*pnrTime, error) {
	if timezone == nil {
		timezone, _ = time.LoadLocation("UTC")
	}
	return &pnrTime{
		h: h,
		m: m,
		l: timezone,
	}, nil
}

func (p *pnrTime) String() string {
	return fmt.Sprintf("%02d%02d", p.h, p.h)
}

func pnrTimeToTime(pnrTime string, location *time.Location) time.Time {
	now := time.Now() // to get year, month, day
	var h, m int
	fmt.Sscanf(pnrTime, "%02d%02d", &h, &m)
	return time.Date(now.Year(), now.Month(), now.Day(), h, m, int(0), int(0), location)
}

func timeToPnrTime(t *time.Time) string {
	return fmt.Sprintf("%02d%02d", t.Hour(), t.Minute())
}

func computeDepartureArrivalTimes(origin, destination *models.Airport) ([]*departureArrivalTime, error) {
	depArrtimes := []*departureArrivalTime{}
	averageSpeedKmH := float64(875)
	halfHourOverhead := float64(.5)
	arriveNextDay := false

	if origin == nil || destination == nil {
		return depArrtimes, fmt.Errorf("missing origin or destination airport")
	}
	if origin == destination {
		return depArrtimes, fmt.Errorf("origin and destination, same airport")
	}
	if _, err := time.LoadLocation(destination.TZ); err != nil {
		return depArrtimes, fmt.Errorf("bad TZ for destination airport %q: %v", destination.Name, err)
	}
	if _, err := time.LoadLocation(origin.TZ); err != nil {
		return depArrtimes, fmt.Errorf("bad TZ for origin airport %q: %v", origin.Name, err)
	}

	distance := db.Distance(origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude)
	distanceKm := float64(distance / 1000)
	formattedHourDuration := fmt.Sprintf("%fh", (halfHourOverhead + (distanceKm / averageSpeedKmH)))
	flightDuration, err := time.ParseDuration(formattedHourDuration)
	if err != nil {
		return depArrtimes, fmt.Errorf("unable to parse duration: %s", formattedHourDuration)
	}

	log.Infof("Flight duration %q->%q:%v\n", origin.City, destination.City, flightDuration)
	now := time.Now() // to get year, month, day

	departureTimes := []string{"0800", "1100", "1500", "1700"}

	for _, departureTime := range departureTimes {
		var h, m int
		fmt.Sscanf(departureTime, "%02d%02d", &h, &m)
		originLocation, err := time.LoadLocation(origin.TZ)
		if err != nil {
			return depArrtimes, fmt.Errorf("unknown TZ: %s: %v", origin.TZ, err)
		}

		localDepartureTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, int(0), int(0), originLocation)
		utcLocation, _ := time.LoadLocation("UTC") // no error check here since we hardcode "UTC"
		utcDepartureTime := localDepartureTime.In(utcLocation)

		utcArrivalTime := utcDepartureTime.Add(flightDuration)
		arrivalTimeLocation, _ := time.LoadLocation(destination.TZ) // already checked error
		localArrivalTime := utcArrivalTime.In(arrivalTimeLocation)

		if localArrivalTime.Hour() < 7 || localArrivalTime.Hour() > 22 {
			continue // silly workaround to avoid crazy arrival time
		}

		if localArrivalTime.Day() != localDepartureTime.Day() {
			arriveNextDay = true
		}
		depArrtimes = append(depArrtimes,
			&departureArrivalTime{departure: departureTime,
				arrival:       fmt.Sprintf("%02d%02d", localArrivalTime.Hour(), localArrivalTime.Minute()),
				arriveNextDay: arriveNextDay})
		if arriveNextDay {
			break // don't arrive next day twice
		}
	}
	return depArrtimes, nil
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
			continue
		}

		destinationAirportID, err := strconv.ParseInt(line[5], 10, 32)
		if err != nil {
			continue
		}

		sourceAirportID, err := strconv.ParseInt(line[3], 10, 32)
		if err != nil {
			continue
		}
		stops, err := strconv.ParseInt(line[7], 10, 64)
		if err != nil {
			continue
		}
		//		airline
		//		airlineID
		//		sourceAirport
		//		sourceAirportID
		//		destinationAirport
		//		destinationAirportID
		//		codeshare
		//		stops
		//		equipment

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
		log.Fatalf("Unable to load courses: %v", err)
	}

	//schedules := []*models.Schedule{}
	flightNumberPerAirline := map[string]int16{}
	sclient := schedulesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
	var scheduleID int64
	for i, course := range routes {
		scheduleID = int64(i)
		flightNumberPerAirline[course.Airline] = flightNumberPerAirline[course.Airline] + 1

		departureArrivalTimes, err := computeDepartureArrivalTimes(ID2Airports[course.SourceAirportID], ID2Airports[course.DestinationAirportID])
		if err != nil {
			log.Infof("Cannot compute arrival time... %v", err)
			continue
		}
		for _, d := range departureArrivalTimes {
			log.Infof("DepartureTime: %s, ArrivalTime: %s", d.departure, d.arrival)
			flightNumber := strings.Join([]string{course.Airline, fmt.Sprintf("%03d", flightNumberPerAirline[course.Airline])}, "")
			daysOperated := "1234567"
			addSchedulePayload := schedulesclient.NewAddScheduleParams()
			addSchedulePayload.Schedule = &models.Schedule{
				ScheduleID:       scheduleID,
				Origin:           course.SourceAirportID,
				Destination:      course.DestinationAirportID,
				FlightNumber:     flightNumber,
				OperatingCarrier: course.Airline,
				DaysOperated:     daysOperated,
				DepartureTime:    d.departure,
				ArrivalTime:      d.arrival,
				ArriveNextDay:    d.arriveNextDay,
			}
			//sclient := schedulesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
			_, err := sclient.AddSchedule(addSchedulePayload)
			if err != nil {
				log.Errorf("Unable to add schedule: %v\n", err)
			}
		}
	}
	os.Exit(0)
}
