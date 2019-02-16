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
	"strconv"
	"strings"

	"github.com/jinzhu/copier"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

func (m *MongoDB) GetAirlines() ([]*models.Airline, error) {
	var airlines []*models.Airline
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return airlines, err
	}
	defer mgoDB.Close()
	var dbAirlines []*Airline
	if err = mgoDB.DB(m.dbName).C(airlinesCollection).Find(nil).All(&dbAirlines); err != nil {
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
	if err := mgoDB.DB(m.dbName).C(coursesCollection).Find(nil).Sort("-when").Limit(100).All(&dbCourses); err != nil {
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

func (m *MongoDB) InsertSchedule(s *models.Schedule) (*models.Schedule, error) {
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
	mgoDB, err := mgo.Dial(m.DialString())
	if err != nil {
		log.Errorf("Cannot connect to  MongoDB: %v", err)
		return err
	}
	err = mgoDB.DB(m.dbName).C(schedulesCollection).Remove(bson.M{"ID": id})
	return err
}
