package itinerary

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type graphMock struct {
	innerGraph     *simple.WeightedDirectedGraph
	airportID2Node map[AirportID]graph.Node
	node2AirportID map[graph.Node]AirportID
}

func (l *graphMock) AirportID2Node(id AirportID) graph.Node {
	return l.airportID2Node[id]
}

func (l *graphMock) Node2AirportID(n graph.Node) AirportID {
	return l.node2AirportID[n]
}

func (l *graphMock) InnerGraph() *simple.WeightedDirectedGraph {
	l.innerGraph = simple.NewWeightedDirectedGraph(0, math.Inf(1))
	l.airportID2Node = make(map[AirportID]graph.Node, 0)
	l.airportID2Node["C"] = simple.Node('C')
	l.airportID2Node["D"] = simple.Node('D')
	l.airportID2Node["F"] = simple.Node('F')
	l.airportID2Node["G"] = simple.Node('G')
	l.airportID2Node["H"] = simple.Node('H')

	l.node2AirportID = make(map[graph.Node]AirportID, 0)
	l.node2AirportID[simple.Node('C')] = "C"
	l.node2AirportID[simple.Node('D')] = "D"
	l.node2AirportID[simple.Node('E')] = "E"
	l.node2AirportID[simple.Node('F')] = "F"
	l.node2AirportID[simple.Node('G')] = "G"
	l.node2AirportID[simple.Node('H')] = "H"

	edges := []simple.WeightedEdge{
		{F: simple.Node('C'), T: simple.Node('D'), W: 3},
		{F: simple.Node('C'), T: simple.Node('E'), W: 2},
		{F: simple.Node('E'), T: simple.Node('D'), W: 1},
		{F: simple.Node('D'), T: simple.Node('F'), W: 4},
		{F: simple.Node('E'), T: simple.Node('F'), W: 2},
		{F: simple.Node('E'), T: simple.Node('G'), W: 3},
		{F: simple.Node('F'), T: simple.Node('G'), W: 2},
		{F: simple.Node('F'), T: simple.Node('H'), W: 1},
		{F: simple.Node('G'), T: simple.Node('H'), W: 2},
	}
	for _, e := range edges {
		l.innerGraph.SetWeightedEdge(e)
	}
	return l.innerGraph
}

func TestCompute(t *testing.T) {
	type args struct {
		From          AirportID
		To            AirportID
		numberOfPaths int
	}
	tests := []struct {
		name string
		args args
		want [][]AirportID
	}{
		{
			name: "wikipedia example",
			args: args{
				From:          "C",
				To:            "H",
				numberOfPaths: 3,
			},
			want: [][]AirportID{
				{"C", "E", "F", "H"},
				{"C", "E", "G", "H"},
				{"C", "D", "F", "H"},
			},
		},
	}

	loader := &graphMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compute(loader, tt.args.From, tt.args.To, tt.args.numberOfPaths); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compute() = %v, want %v", got, tt.want)
			}
		})
	}
}
