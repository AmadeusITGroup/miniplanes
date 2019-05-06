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
package mongo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

var mongoHost string

const (
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
	airlinesCourses     = "courses"
	schedulesCollection = "schedules"
)

// MongoDB implements miniplanes storage interface for MongoDB
type MongoDB struct {
	mongoHost string
	mongoPort string
	dbName    string
}

func NewMongoDB(mongoHost string, mongoPort int, dbName string) *MongoDB {
	return &MongoDB{
		mongoHost: mongoHost,
		mongoPort: strconv.Itoa(mongoPort),
		dbName:    dbName,
	}
}

func (m *MongoDB) Ping() error {
	_, err := mgo.Dial(m.DialString())
	return err
}

func (m *MongoDB) DialString() string {
	return strings.Join([]string{m.mongoHost, m.mongoPort}, ":")
}

// GetAirlines returns all airlines stored in mongo db
func (m *MongoDB) GetAirlines() ([]*models.Airline, error) {
	log.Debugf("MongoDB.GetAirlines")
	var airlines []*models.Airline
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return airlines, err
	}
	defer mgoDB.Close()
	var dbAirlines []*Airline
	if err = mgoDB.DB(m.dbName).C(airlinesCollection).Find(nil).All(&dbAirlines); err != nil {
		log.Errorf("Unable to get airlines: %v", err)
		return airlines, err
	}
	for i := range dbAirlines {
		a, err := dbAirlines[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap airline mongo scheme to Airline model: %v - %v", err, dbAirlines[i])
			continue
		}
		airlines = append(airlines, a)
	}
	return airlines, nil //  TODO: don't swallow errors
}

// InsertAirline insert airline if is not already stored in mongo db
func (m *MongoDB) InsertAirline(a *models.Airline) (*models.Airline, error) {
	log.Debugf("MongoDB.InsertAirlines")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to MongoDB: %v", err)
		return nil, fmt.Errorf("cannot dial mongo: %v", err)
	}
	defer mgoDB.Close()

	var dbAirlines []*Airline
	if err = mgoDB.DB(m.dbName).C(airlinesCollection).Find(nil).All(&dbAirlines); err != nil {
		return nil, err
	}
	for i := range dbAirlines {
		modAirline, err := dbAirlines[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap airline mongo scheme to Airline model: %v - %v", err, dbAirlines[i])
			continue
		}
		if modAirline.AirlineID == a.AirlineID {
			log.Errorf("Airline with ID %d already exists in DB", modAirline.AirlineID)
			return nil, new(ConflictError)
		}
	}
	airline := new(Airline)
	err = copier.Copy(airline, a)
	if err != nil {
		log.Errorf("Cannot convert model airline to mongo airline: %v", err)
		return nil, new(UnprocessableError)
	}
	airline.ID = bson.NewObjectId()
	err = mgoDB.DB(m.dbName).C(airlinesCollection).Insert(airline)
	if err != nil {
		log.Errorf("Cannot insert airline in mongo DB: %v", err)
		return nil, err
	}
	log.Infof("airline #%v inserted", a)
	return a, nil
}

func (m *MongoDB) GetAirports() ([]*models.Airport, error) {
	log.Debugf("MongoDB.GetAirports")
	var airports []*models.Airport
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return airports, err
	}
	defer mgoDB.Close()
	dbAirports := []*Airport{}
	if err := mgoDB.DB(m.dbName).C(airportsCollection).Find(nil).All(&dbAirports); err != nil {
		return airports, err
	}
	for i := range dbAirports {
		a, err := dbAirports[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap airport mongo scheme to Airport model: %v - %v", err, dbAirports[i])
			continue
		}
		airports = append(airports, a)
	}
	return airports, nil // TODO: don't swallow errors
}

// InsertAirport insert airport if is not already stored in mongo db
func (m *MongoDB) InsertAirport(a *models.Airport) (*models.Airport, error) {
	log.Debugf("MongoDB.InsertAirport")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to MongoDB: %v", err)
		return nil, fmt.Errorf("cannot dial mongo: %v", err)
	}
	defer mgoDB.Close()

	dbAirports := []*Airport{}
	if err := mgoDB.DB(m.dbName).C(airportsCollection).Find(nil).All(&dbAirports); err != nil {
		return nil, err
	}
	for i := range dbAirports {
		modAiport, err := dbAirports[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap airline mongo scheme to Airport model: %v - %v", err, dbAirports[i])
			continue
		}
		if modAiport.AirportID == a.AirportID {
			log.Errorf("Airport with ID %d already exists in DB", modAiport.AirportID)
			return nil, new(ConflictError)
		}
	}
	airport := new(Airport)
	err = copier.Copy(airport, a)
	if err != nil {
		log.Errorf("Cannot convert model airline to mongo airline: %v", err)
		return nil, new(UnprocessableError)
	}
	airport.ID = bson.NewObjectId()
	err = mgoDB.DB(m.dbName).C(airportsCollection).Insert(airport)
	if err != nil {
		log.Errorf("Cannot insert airport in mongo DB: %v", err)
		return nil, err
	}
	log.Infof("airport #%v inserted", a)
	return a, nil
}

func (m *MongoDB) InsertSchedule(s *models.Schedule) (*models.Schedule, error) {
	log.Debugf("MongoDB.InsertSchedule")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer mgoDB.Close()
	schedule := new(Schedule)
	copier.Copy(schedule, s)
	schedule.ID = bson.NewObjectId()
	err = mgoDB.DB(m.dbName).C(schedulesCollection).Insert(schedule)
	if err != nil {
		return nil, err
	}
	var v models.Schedule
	err = mgoDB.DB(m.dbName).C(schedulesCollection).FindId(schedule.ID).One(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *MongoDB) GetSchedules() ([]*models.Schedule, error) {
	log.Debugf("MongoDB.GetSchedules")
	var schedules []*models.Schedule
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return schedules, err
	}
	defer mgoDB.Close()

	var dbSchedules []*Schedule
	if err := mgoDB.DB(m.dbName).C(schedulesCollection).Find(nil).All(&dbSchedules); err != nil {
		return schedules, err
	}
	for i := range dbSchedules {
		s, err := dbSchedules[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap schedule mongo scheme to Schedule model: %v - %v", err, schedules[i])
			continue
		}
		schedules = append(schedules, s)
	}
	return schedules, nil //  TODO: don't swallow errors
}

func (m *MongoDB) GetSchedule(id int64) (*models.Schedule, error) {
	log.Debugf("MongoDB.GetSchedule")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return nil, err
	}
	var (
		modSchedule *models.Schedule
		schedule    Schedule
	)

	err = mgoDB.DB(m.dbName).C(schedulesCollection).Find(bson.M{"ID": id}).All(&schedule)
	if err != nil {
		return nil, err
	}
	copier.Copy(modSchedule, schedule)
	return modSchedule, nil
}

func (m *MongoDB) UpdateSchedule(id int64, schedule *models.Schedule) (*models.Schedule, error) {
	log.Debugf("MongoDB.UpdateSchedule")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return nil, err
	}
	var s Schedule
	copier.Copy(&s, schedule)
	err = mgoDB.DB(m.dbName).C(schedulesCollection).Update(bson.M{"ID": id}, s)
	if err != nil {
		return nil, err
	}
	// TODO: UPDATE REF
	return schedule, nil
}

func (m *MongoDB) DeleteSchedule(id int64) error {
	log.Debugf("MongoDB.DeleteSchedule")
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return err
	}
	err = mgoDB.DB(m.dbName).C(schedulesCollection).Remove(bson.M{"ID": id})
	return err
}
