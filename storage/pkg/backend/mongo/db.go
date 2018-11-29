package mongo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
)

var mongoHost string

const (
	dbName              = "miniapp"
	routesCollection    = "routes"
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
	schedulesCollection = "schedules"
)

// MongoDB implemtns miniapp storage interface for MongoDB
type MongoDB struct {
	mongoEndPoint string
}

func NewMongoDB(endPoint string) *MongoDB {
	return &MongoDB{mongoEndPoint: endPoint}
}

func (m *MongoDB) GetSchedules() ([]*Schedule, error) {
	db, err := mgo.Dial("localhost:27017")
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
	db, err := mgo.Dial("localhost")
	if err != nil {
		return dbAirlines, err
	}
	defer db.Close() // clean up when we’re done
	err = db.DB(dbName).C(airlinesCollection).Find(nil).All(&dbAirlines)
	return dbAirlines, err
}

func (m *MongoDB) GetRoutes(w http.ResponseWriter, r *http.Request) {
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer db.Close() // clean up when we’re done
	//db := context.Get(r, "database").(*mgo.Session)
	var routes []*Route
	if err := db.DB(dbName).C(routesCollection).
		Find(nil).Sort("-when").Limit(100).All(&routes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(routes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (m *MongoDB) GetAirports() ([]*Airport, error) {
	var (
		dbAirports []*Airport
	)
	db, err := mgo.Dial("localhost")
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
