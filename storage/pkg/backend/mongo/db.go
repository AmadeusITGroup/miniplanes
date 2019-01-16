package mongo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
)

var mongoHost string

const (
	dbName              = "miniapp"
	coursesCollection   = "courses"
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
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

func (m *MongoDB) GetCourses() ([]*Schedule, error) {
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

func (m *MongoDB) Getcourses(w http.ResponseWriter, r *http.Request) {
	db, err := mgo.Dial(m.dialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	//db := context.Get(r, "database").(*mgo.Session)
	var courses []*Course
	if err := db.DB(dbName).C(coursesCollection).
		Find(nil).Sort("-when").Limit(100).All(&courses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

// HandleInsert
func HandleInsert(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)
	// decode the request body
	var c Airline
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// give the comment a unique ID and set the time
	//c.ID = bson.NewObjectId()
	//c.When = time.Now()
	// insert it into the database
	if err := db.DB(dbName).C(airlinesCollection).Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// redirect to it
	http.Redirect(w, r, "/airlines/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func findFlights(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	source := params["source"]
	destination := params["destination"]
	//db := context.Get(r, "database").(*mgo.Session)

	fmt.Printf("find flights... %q -> %q\n", source, destination)
	return
}
