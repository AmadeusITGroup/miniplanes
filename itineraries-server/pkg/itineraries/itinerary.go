package itinerary

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

var itineraryGraph *realGraph

// AirportID is a string
type AirportID string

// Compute computes itineraries
func Compute(itineraryGraph Graph, From, To AirportID, numberOfPaths int) [][]AirportID {
	g := itineraryGraph.InnerGraph()

	paths := path.YenKShortestPaths(g, numberOfPaths, itineraryGraph.AirportID2Node(From), itineraryGraph.AirportID2Node(To))
	if paths == nil {
		return nil
	}
	solutions := make([][]AirportID, len(paths))
	for i, p := range paths {
		if p == nil {
			continue
		}
		solutions[i] = make([]AirportID, len(p))
		for j, n := range p {
			solutions[i][j] = itineraryGraph.Node2AirportID(n)
		}
	}
	return solutions
}

// Graph represents operations we can do on itinerary graph
type Graph interface {
	InnerGraph() *simple.WeightedDirectedGraph
	AirportID2Node(AirportID) graph.Node
	Node2AirportID(graph.Node) AirportID
}

// GetGraph creates the itinerary graph
func GetGraph() Graph {
	if itineraryGraph == nil {
		itineraryGraph = &realGraph{
			innerGraph:     simple.NewWeightedDirectedGraph(0, math.Inf(1)),
			airportID2Node: make(map[AirportID]graph.Node, 0),
			node2AirportID: make(map[graph.Node]AirportID, 0),
		}
		itineraryGraph.innerGraph = simple.NewWeightedDirectedGraph(0, math.Inf(1))
	}
	return itineraryGraph
}

type realGraph struct {
	innerGraph     *simple.WeightedDirectedGraph
	airportID2Node map[AirportID]graph.Node
	node2AirportID map[graph.Node]AirportID
}

func (r *realGraph) Load() *simple.WeightedDirectedGraph {
	r.innerGraph = simple.NewWeightedDirectedGraph(0, math.Inf(1))
	//...
	r.airportID2Node = make(map[AirportID]graph.Node, 0)
	r.node2AirportID = make(map[graph.Node]AirportID, 0)
	return r.innerGraph
}

func (r *realGraph) AirportID2Node(id AirportID) graph.Node {
	return nil
}

func (r *realGraph) Node2AirportID(n graph.Node) AirportID {
	return ""
}

func (r *realGraph) InnerGraph() *simple.WeightedDirectedGraph {
	return r.innerGraph
}

func (r *realGraph) update() error {

}
