package engine

import (
	"fmt"
	"reflect"
	"testing"

	itinerarymodels "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	"github.com/davecgh/go-spew/spew"
)

/*
type graphMock struct {
	innerGraph     *simple.WeightedDirectedGraph
	airportID2Node map[int32]graph.Node
	node2AirportID map[graph.Node]int32
}

func (l *graphMock) AirportID2Node(id int32) graph.Node {
	return l.airportID2Node[id]
}

func (l *graphMock) Node2AirportID(n graph.Node) int32 {
	return l.node2AirportID[n]
}

var (
	C, D, E, F, G, H int32 = 1, 2, 3, 4, 5, 6
)

// C->1
// D->2
// F->3
// G->4
// H->5
func (l *graphMock) InnerGraph() *simple.WeightedDirectedGraph {
	l.innerGraph = simple.NewWeightedDirectedGraph(0, math.Inf(1))
	l.airportID2Node = make(map[int32]graph.Node, 0)
	l.airportID2Node[C] = simple.Node('C')
	l.airportID2Node[D] = simple.Node('D')
	l.airportID2Node[F] = simple.Node('F')
	l.airportID2Node[G] = simple.Node('G')
	l.airportID2Node[H] = simple.Node('H')

	l.node2AirportID = make(map[graph.Node]int32, 0)
	l.node2AirportID[simple.Node('C')] = C
	l.node2AirportID[simple.Node('D')] = D
	l.node2AirportID[simple.Node('E')] = E
	l.node2AirportID[simple.Node('F')] = F
	l.node2AirportID[simple.Node('G')] = G
	l.node2AirportID[simple.Node('H')] = H

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
*/

func NewString(s string) *string {
	return &s
}

func NewInt(i int32) *int32 {
	return &i
}

func NewBool(b bool) *bool {
	return &b
}

func TestCompute(t *testing.T) {
	type args struct {
		airports      []*storagemodels.Airport
		schedules     []*storagemodels.Schedule
		From          string
		DepartureDate string
		DepartureTime string
		To            string
		numberOfPaths int
	}
	tests := []struct {
		name             string
		args             args
		wantErrorMessage string
		want             []*itinerarymodels.Itinerary
	}{
		{
			name: "no ID for From ",
			args: args{
				From: "",
			},
			wantErrorMessage: `can't find airportID for `,
			want:             []*itinerarymodels.Itinerary{},
		},
		{
			name: "no ID for To",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
				},
				From: "NCE",
				To:   "",
			},
			wantErrorMessage: `can't find airportID for `,
		},
		{
			name: "no schedules no itinerary",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "JFK",
					},
				},
				From: "NCE",
				To:   "JFK",
			},
			wantErrorMessage: "",                             // no schedules no error
			want:             []*itinerarymodels.Itinerary{}, // no itineraries
		},
		{
			name: "one valid schedule",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "JFK",
					},
				},
				schedules: []*storagemodels.Schedule{
					{
						ArrivalTime:      NewString("1300"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("AF01"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "JFK",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "", // no schedules no error
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-JFK",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1300",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "JFK",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "two valid schedules",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "JFK",
					},
				},
				schedules: []*storagemodels.Schedule{
					{
						ArrivalTime:      NewString("1300"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("AF01"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("1305"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1005"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("BA01"),
						OperatingCarrier: NewString("BA"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "JFK",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-JFK",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1300",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "JFK",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
					},
				},
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-JFK",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1305",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1005",
							Destination:      "JFK",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "NCE",
							SegmentID:        0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGraph(tt.args.airports, tt.args.schedules)
			if err != nil {
				panic(err) // should not happen. Stop here.
			}
			got, err := g.Compute(tt.args.From, tt.args.DepartureDate, tt.args.DepartureTime, tt.args.To, tt.args.numberOfPaths)
			if err != nil {
				if err.Error() != tt.wantErrorMessage {
					t.Errorf("Expected error was %q: got %q", tt.wantErrorMessage, err.Error())
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("\nExpected\n")
				spew.Dump(tt.want)
				fmt.Printf("\nGot:\n")
				spew.Dump(got)
				t.Error()
			}
		})
	}
}
