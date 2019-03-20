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
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/db"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/db/mongo"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

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
	defaultMongoIP   = "127.0.0.1"
	defaultMongoPort = 27017
)

var (
	csvFileName string
	mongoIP     string
	mongoPort   int
)

func init() {

	flag.StringVar(&csvFileName, "csv-file-name", "", "csv file to be generated. If no csv-file-name is supplied schedules will be inserted in mongo")
	flag.StringVar(&mongoIP, "mongo-host", defaultMongoIP, "mongo endpoint")
	flag.IntVar(&mongoPort, "mongoPort", defaultMongoPort, "mongo port")
}

func main() {

	flag.Parse()
	generateCSV := false
	if len(csvFileName) > 0 {
		generateCSV = true
	}

	log.Infof("%s %s %d\n", csvFileName, mongoIP, mongoPort)

	m := mongo.NewMongoDB(mongoIP, mongoPort, "miniplanes")
	ID2Airports := map[int32]*models.Airport{}
	airports, err := m.GetAirports()
	if err != nil || len(airports) == 0 {
		log.Fatalf("Unable to load airports %s : %v", m.DialString(), err)
	}
	for i := range airports {
		ID2Airports[airports[i].AirportID] = airports[i]
	}

	courses, err := m.GetCourses()
	if err != nil || len(courses) == 0 {
		log.Fatalf("Unable to load courses: %v", err)
	}

	var file io.Writer
	if generateCSV {
		file, err = os.Create(csvFileName)
		if err != nil {
			log.Fatalf("Couldn't open file\n")
		}
	}
	schedules := []*models.Schedule{}
	flightNumberPerAirline := map[string]int16{}
	var scheduleID int64
	for i, course := range courses {
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
			schedules = append(schedules, &models.Schedule{
				ScheduleID:       scheduleID,
				Origin:           course.SourceAirportID,
				Destination:      course.DestinationAirportID,
				FlightNumber:     flightNumber,
				OperatingCarrier: course.Airline,
				DaysOperated:     daysOperated,
				DepartureTime:    d.departure,
				ArrivalTime:      d.arrival,
				ArriveNextDay:    d.arriveNextDay,
			})
		}
	}

	if generateCSV {
		file, err = os.Create(csvFileName)
		if err != nil {
			log.Fatalf("Couldn't open file\n")
		}
		writer := csv.NewWriter(file)
		for _, s := range schedules {
			data := []string{
				fmt.Sprintf("%d", s.ScheduleID),
				fmt.Sprint(s.Origin),
				fmt.Sprint(s.Destination),
				s.FlightNumber,
				s.OperatingCarrier,
				s.DaysOperated,
				s.DepartureTime,
				s.ArrivalTime,
				strconv.FormatBool(s.ArriveNextDay),
			}
			err := writer.Write(data)
			if err != nil {
				log.Errorf("Something went wrong writing csv file: %v\n", err)
				continue
			}
			writer.Flush()
		}
	} else {
		for _, s := range schedules {
			if _, err := m.InsertSchedule(s); err != nil {
				log.Errorf("Couldn't insert schedule %#v: %v\n", s, err)
			}
		}
	}
	os.Exit(0)
}
