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

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"

	itinerarymodels "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
)

func (r *realGraph) buildSegmentFromSchedule(s *storagemodels.Schedule, departureDate string) *itinerarymodels.Segment {
	arrivalDate := departureDate //TOOD fix bug here arrival date could be different
	return &itinerarymodels.Segment{
		ArrivalDate:      arrivalDate,
		ArrivalTime:      s.ArrivalTime,
		ArriveNextDay:    s.ArriveNextDay,
		DepartureDate:    departureDate,
		DepartureTime:    s.DepartureTime,
		Destination:      r.IATAFromAirportID[s.Destination],
		FlightNumber:     s.FlightNumber,
		OperatingCarrier: s.OperatingCarrier,
		Origin:           r.IATAFromAirportID[s.Origin],
		SegmentID:        0,
	}
}

func checkScheduleIntegrity(s *storagemodels.Schedule) error {
	if s.Origin == 0 {
		return fmt.Errorf("no origin airport found for schedule")
	}
	if s.Destination == 0 {
		return fmt.Errorf("no destination airport found for schedule")
	}
	if len(s.DepartureTime) == 0 {
		return fmt.Errorf("no departureTime found for schedule")
	}
	return nil
}

func (r *realGraph) computeAllSegments(from, departureDate, departureTime, to string, separationDegree int) ([][]*itinerarymodels.Segment, error) {
	segments := [][]*itinerarymodels.Segment{}
	fromAirportID, found := r.airportIDFromIATA[from]
	if !found {
		return segments, fmt.Errorf("can't find airportID for %s", from)
	}
	toAirportID, found := r.airportIDFromIATA[to]
	if !found {
		return segments, fmt.Errorf("can't find airportID for %s", to)
	}
	if fromAirportID == toAirportID {
		return segments, nil
	}
	if separationDegree == 0 {
		return segments, fmt.Errorf("no segments found")
	}
	for _, s := range r.schedules { // for each schedules...
		if err := checkScheduleIntegrity(s); err != nil {
			log.Warnf("bad schedule found: %v", err)
			continue
		}
		if s.Origin != fromAirportID {
			continue
		}
		if !oKtoBeTaken(departureTime, s.DepartureTime) { // not good schedules or too early
			log.Tracef("Impossible to take flight with: requested departure %s - flight departure %s", departureDate, s.DepartureTime)
			continue
		}
		currentSegment := r.buildSegmentFromSchedule(s, departureDate)
		currentSegmentFanout, err := r.computeAllSegments(r.IATAFromAirportID[s.Destination], departureDate, s.ArrivalTime, to, separationDegree-1)
		if err != nil {
			return [][]*itinerarymodels.Segment{}, err
		}
		ss := []*itinerarymodels.Segment{}
		s := new(itinerarymodels.Segment)
		copier.Copy(s, currentSegment)
		ss = append(ss, s)
		if len(currentSegmentFanout) == 0 {
			if s.Destination == to {
				segments = append(segments, ss) // if no fanout check destination
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
	solutions := []*itinerarymodels.Itinerary{}
	maxDegreefSeparation := 4
	segments, err := r.computeAllSegments(from, departureDate, departureTime, to, maxDegreefSeparation)
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
	return solutions, nil
}

// Graph represents operations we can do on itinerary graph
type Graph interface {
	Compute(from, departureDate, departureTime, to string, numberOfPaths int) ([]*itinerarymodels.Itinerary, error)
	//InnerGraph() graph.Graph
}

type realGraph struct {
	airports          []*storagemodels.Airport
	schedules         []*storagemodels.Schedule
	airportIDFromIATA map[string]int32
	IATAFromAirportID map[int32]string
}

func splitHourMinutes(t string) (int32, int32, error) {
	var h, m int32
	_, err := fmt.Sscanf(t, "%02d%02d", &h, &m)
	return h, m, err
}

// it means t2 is later on than t1
func oKtoBeTaken(t1, t2 string) bool {
	h1, _, err := splitHourMinutes(t1)
	if err != nil {
		return false
	}
	h2, _, err := splitHourMinutes(t2)
	if err != nil {
		return false
	}
	if h2 > h1 { // must be in same day
		return true
	}
	return false
}

// NewGraph creates the itinerary graph
func NewGraph(airports []*storagemodels.Airport, schedules []*storagemodels.Schedule) (Graph, error) {
	itineraryGraph := &realGraph{
		schedules:         schedules,
		airports:          airports,
		airportIDFromIATA: make(map[string]int32, 0),
		IATAFromAirportID: make(map[int32]string, 0),
	}

	for _, airport := range airports {
		itineraryGraph.airportIDFromIATA[airport.IATA] = airport.AirportID
		itineraryGraph.IATAFromAirportID[airport.AirportID] = airport.IATA
	}
	return itineraryGraph, nil
}
