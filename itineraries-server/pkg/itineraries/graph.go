package itineraries

import (
	"fmt"
	"math"

	storageclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

// Segment it's an itinerary segment
type Segment struct {
	FromAirport  IATACode
	Departure    string
	ToAirport    IATACode
	Arrival      string
	FlightNumber string
}

// Itinerary it's an array of Segment (https://aviation.stackexchange.com/questions/14567/what-is-the-difference-between-slice-segment-and-leg)
type Itinerary []*Segment

// Graph represents operations we can do on itinerary graph
type Graph interface {
	Compute(from, to string, numberOfItineraries int) ([]Itinerary, error)
}

type realGraph struct {
	innerGraph     *simple.WeightedDirectedGraph
	airportID2Node map[int64]graph.Node
	node2AirportID map[graph.Node]int64
}

func (r *realGraph) AirportID2Node(id int64) (graph.Node, error) {
	n, ok := r.airportID2Node[id]
	if !ok {
		return nil, fmt.Errorf("unable to find inner graph node for Airport ID %d", id)
	}
	return n, nil
}

func (r *realGraph) edgeAdder() graph.WeightedEdgeAdder {
	return r.innerGraph
}

func (r *realGraph) nodeAdder() graph.NodeAdder {
	return r.innerGraph
}

func (r *realGraph) graph() graph.Graph {
	return r.innerGraph
}

// NewGraph creates the itinerary graph
func NewGraph() (Graph, error) {
	itineraryGraph := &realGraph{
		innerGraph:     simple.NewWeightedDirectedGraph(0, math.Inf(1)),
		airportID2Node: make(map[int64]graph.Node, 0),
		node2AirportID: make(map[graph.Node]int64, 0),
	}
	airportsOK, err := storageclient.Default.Airports.GetAirports(nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve airports: %v", err)
	}
	airports := airportsOK.Payload
	if err != nil {
		return nil, fmt.Errorf("Unable to get airports: %v", err)
	}
	for _, a := range airports {
		newNode := itineraryGraph.nodeAdder().NewNode()
		itineraryGraph.airportID2Node[a.AirportID] = newNode
		itineraryGraph.node2AirportID[newNode] = a.AirportID
	}
	schedulesOK, err := storageclient.Default.Schedules.GetSchedules(nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve schedules: %v", err)
	}
	schedules := schedulesOK.Payload
	if err != nil {
		fmt.Printf("No schedules found...\n")
		return nil, fmt.Errorf("unable to fetch schedules: %v", err)
	}
	for _, s := range schedules {
		fmt.Printf("Adding edge %d->%d\n", s.Origin, s.Destination)
		e := simple.WeightedEdge{simple.Node(s.Origin), simple.Node(s.Destination), 1}
		itineraryGraph.edgeAdder().SetWeightedEdge(e)
	}
	return itineraryGraph, nil
}

// Compute computes itineraries
func (r *realGraph) Compute(from, to string, numberOfPaths int) ([]Itinerary, error) {
	solutions := []Itinerary{}
	fromAirportID, err := airportIDFromIATA(from)
	if err != nil {
		return solutions, fmt.Errorf("can't find airportID for %q: %v", from, err)
	}
	toAirportID, err := airportIDFromIATA(to)
	if err != nil {
		return solutions, fmt.Errorf("can't find airportID for %q: %v", to, err)
	}
	fromNode, err := r.AirportID2Node(fromAirportID)
	if err != nil {
		return solutions, err
	}

	toNode, err := r.AirportID2Node(toAirportID)
	if err != nil {
		return solutions, err
	}
	paths := path.YenKShortestPaths(r.graph().(graph.Graph), numberOfPaths, fromNode, toNode)
	for i, p := range paths {
		if p == nil {
			continue
		}
		solutions[i] = make([]*Segment, len(p))
		for _, s := range p {
			airportID := s.ID()
			fmt.Println("AirportID -> %d\n", airportID)
			solutions[i] = append(solutions[i], &Segment{
				FromAirport:  "NCE",
				Departure:    "0900",
				ToAirport:    "HTW",
				Arrival:      "1127",
				FlightNumber: "BA111",
			})
		}
	}
	return solutions, nil
}
