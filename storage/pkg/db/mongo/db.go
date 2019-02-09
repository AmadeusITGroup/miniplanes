/*
Copyright 2018 Amadeus SaS All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mongo

import (
	"strconv"
	"strings"

	"github.com/jinzhu/copier"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/amadeusitgroup/miniapp/storage/cmd/config"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	log "github.com/sirupsen/logrus"
)

var mongoHost string

const (
	coursesCollection   = "courses"
	airportsCollection  = "airports"
	airlinesCollection  = "airlines"
	airlinesCourses     = "courses"
	schedulesCollection = "schedules"
)

// MongoDB implements miniapp storage interface for MongoDB
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

func (m *MongoDB) Ping() error {
	_, err := mgo.Dial(m.DialString())
	return err
}

func (m *MongoDB) DialString() string {
	return strings.Join([]string{m.mongoHost, m.mongoPort}, ":")
}

func (m *MongoDB) GetAirlines() ([]*models.Airline, error) {
	var airlines []*models.Airline
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return airlines, err
	}
	defer mgoDB.Close()
	var dbAirlines []*Airline
	if err = mgoDB.DB(config.MongoDBName).C(airlinesCollection).Find(nil).All(&dbAirlines); err != nil {
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

func (m *MongoDB) GetCourses() ([]*models.Course, error) {
	var courses []*models.Course
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return courses, err
	}
	defer mgoDB.Close()
	var dbCourses []*Course
	if err := mgoDB.DB(config.MongoDBName).C(coursesCollection).Find(nil).Sort("-when").Limit(100).All(&dbCourses); err != nil {
		return courses, err
	}
	for i := range dbCourses {
		c, err := dbCourses[i].ToModel()
		if err != nil {
			log.Errorf("Unable to remap course mongo scheme to Course model: %v - v", err, dbCourses[i])
			continue
		}
		courses = append(courses, c)
	}
	return courses, nil // TODO: don't swallow errors
}

func (m *MongoDB) GetAirports() ([]*models.Airport, error) {
	var airports []*models.Airport
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return airports, err
	}
	defer mgoDB.Close()
	dbAirports := []*Airport{}
	if err := mgoDB.DB(config.MongoDBName).C(airportsCollection).Find(nil).All(&dbAirports); err != nil {
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

func (m *MongoDB) InsertSchedule(s *models.Schedule) (*models.Schedule, error) {
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	defer mgoDB.Close()
	schedule := new(Schedule)
	copier.Copy(schedule, s)
	schedule.ID = bson.NewObjectId()
	err = mgoDB.DB(config.MongoDBName).C(schedulesCollection).Insert(schedule)
	if err != nil {
		return nil, err
	}
	var v models.Schedule
	err = mgoDB.DB(config.MongoDBName).C(schedulesCollection).FindId(schedule.ID).One(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *MongoDB) GetSchedules() ([]*models.Schedule, error) {
	var schedules []*models.Schedule
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return schedules, err
	}
	defer mgoDB.Close()

	var dbSchedules []*Schedule
	if err := mgoDB.DB(config.MongoDBName).C(schedulesCollection).Find(nil).All(&schedules); err != nil {
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

func (m *MongoDB) DeleteSchedule(id int32) error {
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return err
	}
	err = mgoDB.DB(config.MongoDBName).C(schedulesCollection).Remove(bson.M{"ID": id})
	return err
}
