package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/db"

	"github.com/amadeusitgroup/miniapp/storage/pkg/backend/mongo"
	"github.com/globalsign/mgo"
)

const (
	dbName              = "miniapp"
	coursesCollection   = "courses"
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
	schedulesCollection = "schedules"
)

type MongoDB struct {
	mongoHost string
	mongoPort string
}

func NewMongoDB(mongoHost string, mongoPort int) *MongoDB {
	return &MongoDB{
		mongoHost: mongoHost,
		mongoPort: strconv.Itoa(mongoPort),
	}
}

func (m *MongoDB) dialString() string {
	return strings.Join([]string{m.mongoHost, m.mongoPort}, ":")
}

func (m *MongoDB) GetSchedules() ([]*mongo.Schedule, error) {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	defer db.Close() // clean up when we’re done
	var schedules []*mongo.Schedule
	if err := db.DB(dbName).C(schedulesCollection).Find(nil).All(&schedules); err != nil {
		return []*mongo.Schedule{}, err
	}

	// write it out
	return schedules, nil
}

func (m *MongoDB) GetCourses() ([]*mongo.Course, error) {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	//db := context.Get(r, "database").(*mgo.Session)
	var courses []*mongo.Course
	if err := db.DB(dbName).C(coursesCollection).Find(nil).Sort("-when").Limit(100).All(&courses); err != nil {
		return []*mongo.Course{}, err
	}
	return courses, nil
}

func (m *MongoDB) GetAirports() ([]*mongo.Airport, error) {
	var (
		dbAirports []*mongo.Airport
	)
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		return dbAirports, err
	}
	defer db.Close() // clean up when we’re done
	err = db.DB(dbName).C(airportsCollection).Find(nil).All(&dbAirports)
	return dbAirports, err
}

func (m *MongoDB) GetAirlines() ([]*mongo.Airline, error) {
	var (
		dbAirlines []*mongo.Airline
	)
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		return dbAirlines, err
	}
	defer db.Close() // clean up when we’re done
	err = db.DB(dbName).C(airlinesCollection).Find(nil).All(&dbAirlines)
	return dbAirlines, err
}

func computeDepartureTime() string {
	return "0800"
}

func computeArrivalTime(origin, destination *mongo.Airport, departureTime string) string {
	if origin == nil || destination == nil {
		return "2259"
	}

	distance := db.Distance(origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude)
	distanceKm := distance / 1000
	flightTimeInHours := distanceKm / 800
	fmt.Printf("%s -> %s: distance km= %f, flight-time hours=%f \n", origin.City, destination.City, distanceKm, flightTimeInHours)

	return "1200"
}

func main() {
	m := NewMongoDB("localhost", 9999) // 9999 taken by default
	/*airlines, err := m.GetAirlines()
	if err != nil {
		log.Fatalf("Unable to load airlines: %v", err)
	}
	for i := range airlines {
		fmt.Printf("airline: %#v\n", airlines[i])
	}
	*/

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

	flightNumberPerAirline := map[string]int16{}
	for _, course := range courses {
		flightNumberPerAirline[course.Airline] = flightNumberPerAirline[course.Airline] + 1

		schedule := mongo.Schedule{
			Origin:           course.SourceAirportID,
			Destination:      course.DestinationAirportID,
			FlightNumber:     strings.Join([]string{course.Airline, fmt.Sprintf("%03d", flightNumberPerAirline[course.Airline])}, ""),
			OperatingCarrier: course.Airline,
			DaysOperated:     "1234567",
			Departure:        computeDepartureTime(),
			Arrival:          computeArrivalTime(ID2Airports[course.SourceAirportID], ID2Airports[course.DestinationAirportID], computeDepartureTime()),
		}
		fmt.Printf("%#v\n", schedule)
	}

	os.Exit(0)
}
