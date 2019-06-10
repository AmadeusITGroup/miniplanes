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
package engine

import (
	"fmt"
	"sort"
	"time"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"

	"github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/db"
	itinerarymodels "github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
)

var (
	AverageSpeedKmH = float64(700)
	//HalfHourOverhead = float64(.5)
	FlightOverhead = float64(.75)
)

func ComputeArrivalDateTime(year int, departureDate, departureTime string, origin, destination *models.Airport) (string, string, error) {

	arrivalDate := ""
	arrivalTime := ""

	var months, days, hours, minutes int
	fmt.Sscanf(departureTime, "%02d%02d", &hours, &minutes)
	fmt.Sscanf(departureDate, "%02d%02d", &days, &months)
	distance := db.Distance(origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude)
	distanceKm := float64(distance / 1000)

	formattedHourDuration := fmt.Sprintf("%fh", (FlightOverhead + (distanceKm / AverageSpeedKmH)))
	flightDuration, err := time.ParseDuration(formattedHourDuration)
	if err != nil {
		return arrivalDate, arrivalTime, fmt.Errorf("unable to parse duration: %s", formattedHourDuration)
	}

	originLocation, err := time.LoadLocation(origin.TZ)
	if err != nil {
		return arrivalDate, arrivalTime, fmt.Errorf("unknown origin timezone: %s: %v", origin.TZ, err)
	}
	localDepartureTime := time.Date(year, time.Month(months), days, hours, minutes, int(0), int(0), originLocation)
	utcLocation, _ := time.LoadLocation("UTC") // no error check here since we hardcode "UTC"
	utcDepartureTime := localDepartureTime.In(utcLocation)
	utcArrivalTime := utcDepartureTime.Add(flightDuration)
	arrivalTimeLocation, err := time.LoadLocation(destination.TZ)
	if err != nil {
		return arrivalDate, arrivalTime, fmt.Errorf("unknown TZ: %s: %v", destination.TZ, err)
	}

	localArrivalTime := utcArrivalTime.In(arrivalTimeLocation)
	arrivalDate = fmt.Sprintf("%02d%02d", localArrivalTime.Day(), localArrivalTime.Month())
	arrivalTime = fmt.Sprintf("%02d%02d", localArrivalTime.Hour(), localArrivalTime.Minute())
	return arrivalDate, arrivalTime, nil
}

func (r *realGraph) lookupAirport(id int32) *models.Airport {
	for i := range r.airports {
		if r.airports[i].AirportID == id {
			return r.airports[i]
		}
	}
	return nil
}

func (r *realGraph) buildSegmentFromSchedule(s *storagemodels.Schedule, departureDate string) (*itinerarymodels.Segment, error) {
	origin := r.lookupAirport(*s.Origin)
	destination := r.lookupAirport(*s.Destination)
	arrivalDate, arrivalTime, err := ComputeArrivalDateTime(2019, departureDate, *s.DepartureTime, origin, destination)
	if err != nil {
		return nil, fmt.Errorf("Unable to compute arrival date and arrival time: %v", err)
	}
	seg := &itinerarymodels.Segment{
		ArrivalDate:      arrivalDate, //arrivalDate,
		ArrivalTime:      arrivalTime,
		DepartureDate:    departureDate, //*s.DepartureDate,
		DepartureTime:    *s.DepartureTime,
		Destination:      r.airportFromID[*s.Destination].IATA,
		FlightNumber:     *s.FlightNumber,
		OperatingCarrier: *s.OperatingCarrier,
		Origin:           r.airportFromID[*s.Origin].IATA,
		SegmentID:        0,
	}
	return seg, nil
}

func checkScheduleIntegrity(s *storagemodels.Schedule) error {
	if *s.Origin == 0 {
		return fmt.Errorf("no origin airport found for schedule")
	}
	if *s.Destination == 0 {
		return fmt.Errorf("no destination airport found for schedule")
	}
	return nil
}

func (r *realGraph) computeAllSegments(from, departure, departureDate, to string, separationDegree int) ([][]*itinerarymodels.Segment, error) {
	segments := [][]*itinerarymodels.Segment{}
	log.Debugf("computeAllSegments %s->%s", from, to)

	fromAirport, found := r.airportFromIATA[from]
	if !found {
		return segments, fmt.Errorf("can't find airportID for %s", from)
	}

	toAirport, found := r.airportFromIATA[to]
	if !found {
		return segments, fmt.Errorf("can't find airportID for %s", to)
	}

	if fromAirport.AirportID == toAirport.AirportID {
		log.Debugf("same from and to (%d,%d) ", fromAirport.AirportID, toAirport.AirportID)
		return segments, nil
	}
	if separationDegree == 0 {
		log.Debugf("Separation degree %d, no segments found", separationDegree)
		return segments, nil //fmt.Errorf("no segments found")
	}

	schedules, ok := r.originToSchedules[fromAirport.AirportID]
	if !ok {
		log.Warnf("Cannot find outbound schedules from %s", fromAirport.IATA)
		return segments, nil //fmt.Errorf("No segments from %s", fromAirport.IATA)
	}
	for _, s := range schedules {
		if !oKtoBeTaken(departure, *s.DepartureTime) { // schedule is too early
			log.Debugf("Schedule too early to be taken, departure=%s, schedule departure time=%s", departure, *s.DepartureTime)
			break // schedules are sorted for departure time as soon the current schedule departure time is too late for the request we can skip
		}

		currentSegment, err := r.buildSegmentFromSchedule(s, departureDate)
		destinationAirport := r.airportFromID[*s.Destination]

		currentSegmentFanout, err := r.computeAllSegments(destinationAirport.IATA, departure, departureDate, to, separationDegree-1)
		if err != nil {
			log.Errorf("Unable to compute all segments from %s to %s: %v", destinationAirport.IATA, to, err)
			return [][]*itinerarymodels.Segment{}, err
		}
		ss := []*itinerarymodels.Segment{}
		s := new(itinerarymodels.Segment)
		copier.Copy(s, currentSegment)
		ss = append(ss, s)
		if len(currentSegmentFanout) == 0 {
			if s.Destination == to {
				segments = append(segments, ss)
			}
		}
		for i := range currentSegmentFanout {
			if currentSegmentFanout[i][len(currentSegmentFanout[i])-1].Destination != to {
				continue // only fanout with right destination is kept
			}
			ss2 := append(ss, currentSegmentFanout[i]...)
			segments = append(segments, ss2)
		}
	}
	return segments, nil
}

func makeDescription(from, to, departureDate, departureTime string) string {
	return fmt.Sprintf("%s:%s - %s-%s", departureDate, departureTime, from, to)
}

func makeItineraryID() string {
	return "MY ID"
}

// ComputeItineraries computes itineraries
func (r *realGraph) Compute(from, departureDate, departureTime, to string, numberOfPaths int) ([]*itinerarymodels.Itinerary, error) {
	log.Debugf("realGraph.Compute")
	log.Debugf("Compute itineraries: %s->%s, departureDate: %s, departureTime: %s", from, to, departureDate, departureTime)
	solutions := []*itinerarymodels.Itinerary{}
	maxDegreefSeparation := 4

	var departureLocation *time.Location
	var err error
	for _, a := range r.airports { // TODO: fix maps to avoid linear scan
		if a.IATA == from {
			if departureLocation, err = time.LoadLocation(a.TZ); err != nil {
				return solutions, fmt.Errorf("unable to determine departure timezone: %v", err)
			}
			break
		}
	}
	if departureLocation == nil {
		return solutions, fmt.Errorf("unable to determine timezone for %s", from)
	}

	segments, err := r.computeAllSegments(from, departureTime, departureDate, to, maxDegreefSeparation)
	if err != nil {
		return solutions, err
	}

	for i := range segments {
		itinerary := &itinerarymodels.Itinerary{
			Description: makeDescription(from, to, departureDate, departureTime),
			ItineraryID: "MY ID",
			Segments:    segments[i],
		}
		solutions = append(solutions, itinerary)
	}
	log.Debugf("Computed %d solutions", len(solutions))
	return solutions, nil
}

// Graph represents operations we can do on itinerary graph
type Graph interface {
	Compute(from, departureDate, departureTime, to string, numberOfPaths int) ([]*itinerarymodels.Itinerary, error)
}

type realGraph struct {
	airports        []*storagemodels.Airport
	schedules       []*storagemodels.Schedule
	airportFromIATA map[string]*storagemodels.Airport
	airportFromID   map[int32]*storagemodels.Airport

	originToSchedules map[int32][]*storagemodels.Schedule
}

func splitHourMinutes(t string) (int32, int32, error) {
	var h, m int32
	_, err := fmt.Sscanf(t, "%02d%02d", &h, &m)
	return h, m, err
}

// it means t2 is later on than t1
func oKtoBeTaken(t1, t2 string) bool {
	return t2 > t1
}

// NewGraph creates the itinerary graph
func NewGraph(airports []*storagemodels.Airport, schedules []*storagemodels.Schedule) (Graph, error) {
	itineraryGraph := &realGraph{
		schedules:       []*storagemodels.Schedule{},
		airports:        airports,
		airportFromIATA: make(map[string]*models.Airport, 0),
		airportFromID:   make(map[int32]*models.Airport, 0),
		//airportIDFromIATA: make(map[string]int32, 0),
		//IATAFromAirportID: make(map[int32]string, 0),
		originToSchedules: make(map[int32][]*storagemodels.Schedule, 0),
	}
	for i := range schedules {
		if err := checkScheduleIntegrity(schedules[i]); err != nil {
			log.Warnf("bad schedule found: %v", err)
			continue
		}
		itineraryGraph.schedules = append(itineraryGraph.schedules, schedules[i])
		if s, ok := itineraryGraph.originToSchedules[*schedules[i].Origin]; ok {
			itineraryGraph.originToSchedules[*schedules[i].Origin] = insertSort(s, schedules[i])
		} else {
			s := []*storagemodels.Schedule{schedules[i]}
			itineraryGraph.originToSchedules[*schedules[i].Origin] = s
		}
	}

	for _, airport := range airports {
		itineraryGraph.airportFromIATA[airport.IATA] = airport
		itineraryGraph.airportFromID[airport.AirportID] = airport
		//itineraryGraph.airportIDFromIATA[airport.IATA] = airport.AirportID
		//itineraryGraph.IATAFromAirportID[airport.AirportID] = airport.IATA
	}
	return itineraryGraph, nil
}

func insertSort(data []*storagemodels.Schedule, el *storagemodels.Schedule) []*storagemodels.Schedule {
	index := sort.Search(len(data), func(i int) bool { return *data[i].DepartureTime <= *el.DepartureTime })
	data = append(data, &storagemodels.Schedule{})
	copy(data[index+1:], data[index:])
	data[index] = el
	return data
}
