package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/db"
	"github.com/amadeusitgroup/miniapp/storage/pkg/backend/mongo"
)

func computeDepartureTime() string {
	return "1520" // TODO: having something more random here...
}

func computeArrivalTime(origin, destination *mongo.Airport, departureTime string) (string, bool, error) {
	averageSpeedKmH := float64(800)
	arriveNextDay := false
	if origin == nil || destination == nil {
		return "2259", arriveNextDay, nil
	}

	distance := db.Distance(origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude)
	distanceKm := float64(distance / 1000)
	formattedHourDuration := fmt.Sprintf("%fh", (distanceKm / averageSpeedKmH))
	flightDuration, err := time.ParseDuration(formattedHourDuration)
	if err != nil {
		return "", arriveNextDay, fmt.Errorf("unable to parse duration: %s", formattedHourDuration)
	}

	now := time.Now() // to get year, month, day
	var h, m int
	fmt.Sscanf(departureTime, "%02d%02d", &h, &m)
	originLocation, err := time.LoadLocation(origin.TZ)
	if err != nil {
		return "", arriveNextDay, fmt.Errorf("unknown TZ: %s: %v", origin.TZ, err)
	}
	localDepartureTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, int(0), int(0), originLocation)
	utcLocation, _ := time.LoadLocation("UTC") // no error check here since we hardcode "UTC"
	utcDepartureTime := localDepartureTime.In(utcLocation)

	utcArrivalTime := utcDepartureTime.Add(flightDuration)
	arrivalTimeLocation, err := time.LoadLocation(destination.TZ)
	if err != nil {
		return "", arriveNextDay, fmt.Errorf("unknown TZ: %s: %v", destination.TZ, err)
	}

	localArrivalTime := utcArrivalTime.In(arrivalTimeLocation)
	if localArrivalTime.Day() != localDepartureTime.Day() {
		arriveNextDay = true
	}
	return fmt.Sprintf("%02d%02d", localArrivalTime.Hour(), localArrivalTime.Minute()), arriveNextDay, nil
}

const (
	defaultMongoIP   = "localhost"
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

	fmt.Printf("%s %s %d\n", csvFileName, mongoIP, mongoPort)

	m := mongo.NewMongoDB(mongoIP, mongoPort) // 9999 take by default
	ID2Airports := map[int32]*mongo.Airport{}
	airports, err := m.GetAirports()
	if err != nil {
		log.Fatalf("Unable to load airports: %v", err)
	}
	for i := range airports {
		ID2Airports[airports[i].AirportID] = airports[i]
	}

	courses, err := m.GetCourses()
	if err != nil {
		log.Fatalf("Unable to load courses: %v", err)
	}

	var file io.Writer
	if generateCSV {
		file, err = os.Create(csvFileName)
		if err != nil {
			log.Fatalf("Couldn't open file\n")
		}
	}

	schedules := []*mongo.Schedule{}
	flightNumberPerAirline := map[string]int16{}
	for _, course := range courses {
		flightNumberPerAirline[course.Airline] = flightNumberPerAirline[course.Airline] + 1
		depTime := computeDepartureTime()
		arrTime, arriveNextDay, err := computeArrivalTime(ID2Airports[course.SourceAirportID], ID2Airports[course.DestinationAirportID], depTime)
		if err != nil {
			fmt.Printf("Cannot compute arrival time... %v", err)
			continue
		}

		schedules = append(schedules, &mongo.Schedule{
			Origin:           course.SourceAirportID,
			Destination:      course.DestinationAirportID,
			FlightNumber:     strings.Join([]string{course.Airline, fmt.Sprintf("%03d", flightNumberPerAirline[course.Airline])}, ""),
			OperatingCarrier: course.Airline,
			DaysOperated:     "1234567",
			Departure:        depTime,
			Arrival:          arrTime,
			ArriveNextDay:    arriveNextDay,
		})
	}

	if generateCSV {
		file, err = os.Create(csvFileName)
		if err != nil {
			log.Fatalf("Couldn't open file\n")
		}
		writer := csv.NewWriter(file)
		for _, s := range schedules {
			data := []string{
				fmt.Sprint(s.Origin),
				fmt.Sprint(s.Destination),
				s.FlightNumber,
				s.OperatingCarrier,
				s.DaysOperated,
				s.Departure,
				s.Arrival,
				strconv.FormatBool(s.ArriveNextDay),
			}
			err := writer.Write(data)
			if err != nil {
				fmt.Printf("Something went wrong writing csv file: %v\n", err)
				continue
			}
			writer.Flush()
		}
	} else {
		for _, s := range schedules {
			if err := m.InsertSchedule(s); err != nil {
				fmt.Printf("Couldn't insert schedule %#v: %v\n", s, err)
			}
		}
	}
	os.Exit(0)
}
