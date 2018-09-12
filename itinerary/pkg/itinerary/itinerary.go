package itinerary

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

// App contains the App string
type App struct {
	port string
}

// NewApplication creates and initiliazes a backend App
func NewApplication(port string) *App {
	p := fmt.Sprintf(":%s", port)
	return &App{
		port: p,
	}
}

// Run runs the backend application
func (a *App) Run() {
	r := mux.NewRouter()
	if err := http.ListenAndServe(a.port, r); err != nil {
		log.Fatal(err)
	}
}

type AirportID string

// Compute computes ittinerariies
func Compute(itineraryGraph Graph, From, To AirportID, numberOfPaths int) [][]AirportID {
	g := itineraryGraph.Load()

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

type Graph interface {
	Load() *simple.WeightedDirectedGraph
	AirportID2Node(AirportID) graph.Node
	Node2AirportID(graph.Node) AirportID
}

func NewGraph() Graph {
	return &realGraph{}
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
