package mongo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongoHost string

const (
	dbName              = "miniapp"
	coursesCollection   = "courses"
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
	airlinesCourses     = "courses"
	schedulesCollection = "schedules"
)

// MongoDB implemtns miniapp storage interface for MongoDB
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

func (m *MongoDB) GetSchedules() ([]*Schedule, error) {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	defer db.Close() // clean up when we’re done
	//db := context.Get(r, "database").(*mgo.Session)
	var schedules []*Schedule
	if err := db.DB(dbName).C(schedulesCollection).Find(nil).All(&schedules); err != nil {
		return []*Schedule{}, err
	}

	// write it out
	return schedules, nil
}

func (m *MongoDB) GetAirlines() ([]*Airline, error) {
	var (
		dbAirlines []*Airline
	)
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		return dbAirlines, err
	}
	defer db.Close() // clean up when we’re done
	err = db.DB(dbName).C(airlinesCollection).Find(nil).All(&dbAirlines)
	return dbAirlines, err
}

func (m *MongoDB) GetCourses() ([]*Course, error) {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	//db := context.Get(r, "database").(*mgo.Session)
	var courses []*Course
	if err := db.DB(dbName).C(coursesCollection).Find(nil).Sort("-when").Limit(100).All(&courses); err != nil {
		return []*Course{}, err
	}
	return courses, nil
}

func (m *MongoDB) GetAirports() ([]*Airport, error) {
	var (
		dbAirports []*Airport
	)
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		return dbAirports, err
	}
	defer db.Close() // clean up when we’re done
	err = db.DB(dbName).C(airportsCollection).Find(nil).All(&dbAirports)
	return dbAirports, err
}

func (m *MongoDB) InsertSchedule(s *Schedule) error {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	s.ID = bson.NewObjectId()
	return db.DB(dbName).C(schedulesCollection).Insert(s)
}

func findFlights(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	source := params["source"]
	destination := params["destination"]
	//db := context.Get(r, "database").(*mgo.Session)

	fmt.Printf("find flights... %q -> %q\n", source, destination)
	return
}
