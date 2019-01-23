package engine

import (
	"fmt"

	"github.com/jinzhu/copier"

	itinerarymodels "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
)

func (r *realGraph) buildSegmentFromSchedule(s *storagemodels.Schedule, departureDate string) *itinerarymodels.Segment {
	arrivalDate := departureDate //TOOD fix bug here arrival date could be different
	return &itinerarymodels.Segment{
		ArrivalDate:      arrivalDate,
		ArrivalTime:      *s.ArrivalTime,                      //    string    `json:"ArrivalTime,omitempty"`
		ArriveNextDay:    *s.ArriveNextDay,                    //    bool    `json:"ArriveNextDay,omitempty"`
		DepartureDate:    departureDate,                       //    string    `json:"DepartureDate,omitempty"`
		DepartureTime:    *s.DepartureTime,                    //    string    `json:"DepartureTime,omitempty"`
		Destination:      r.IATAFromAirportID[*s.Destination], //    string    `json:"Destination,omitempty"`
		FlightNumber:     *s.FlightNumber,                     //    string    `json:"FlightNumber,omitempty"`
		OperatingCarrier: *s.OperatingCarrier,                 //    string    `json:"OperatingCarrier,omitempty"`
		Origin:           r.IATAFromAirportID[*s.Origin],      //    string    `json:"Origin,omitempty"`
		SegmentID:        0,                                   //    int64    `json:"SegmentID,omitempty"`
	}
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
	if fromAirportID == toAirportID || separationDegree == 0 {
		// log.Infof("Found destination airport")
		return segments, nil
	}
	for _, s := range r.schedules { // for each schedules...
		if *s.Origin != fromAirportID {
			continue
		}
		if !oKtoBeTaken(departureTime, *s.DepartureTime) { // not good schedules or too early
			// log.Infof("Cannot be taken...")
			continue
		}
		currentSegment := r.buildSegmentFromSchedule(s, departureDate)

		currentSegmentFanout, err := r.computeAllSegments(r.IATAFromAirportID[*s.Destination], departureDate, *s.ArrivalTime, to, separationDegree-1)
		if err != nil {
			return [][]*itinerarymodels.Segment{}, err
		}
		ss := []*itinerarymodels.Segment{}
		s := new(itinerarymodels.Segment)
		copier.Copy(s, currentSegment)
		ss = append(ss, s)
		if len(currentSegmentFanout) == 0 {
			segments = append(segments, ss)
		}
		for i := range currentSegmentFanout {
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
