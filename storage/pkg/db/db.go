package db

import "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"

// DB defines db interface
type DB interface {
	Ping() error
	DialString() string
	GetSchedules() ([]*models.Schedule, error)
	GetAirlines() ([]*models.Airline, error)
	GetCourses() ([]*models.Course, error)
	GetAirports() ([]*models.Airport, error)
	InsertSchedule(s *models.Schedule) error
}
