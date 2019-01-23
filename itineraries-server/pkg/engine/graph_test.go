package engine

import (
	"fmt"
	"reflect"
	"testing"

	itinerarymodels "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	"github.com/davecgh/go-spew/spew"
)

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
						IATA:      "LHR",
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
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "", // no schedules no error
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1300",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "LHR",
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
						IATA:      "LHR",
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
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1300",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "LHR",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
					},
				},
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1305",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1005",
							Destination:      "LHR",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "NCE",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "one valid schedule, one too early",

			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "LHR",
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
						ArrivalTime:      NewString("1005"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("0705"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("BA01"),
						OperatingCarrier: NewString("BA"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1300",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "LHR",
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
			name: "two segments flight ",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "LHR",
					},
					{
						AirportID: 3,
						IATA:      "CDG",
					},
				},
				schedules: []*storagemodels.Schedule{
					{
						ArrivalTime:      NewString("1205"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(3),
						FlightNumber:     NewString("AF01"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("2105"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1305"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("BA01"),
						OperatingCarrier: NewString("BA"),
						Origin:           NewInt(3),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1205",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "CDG",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "2105",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1305",
							Destination:      "LHR",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "CDG",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "2xtwo segments flights ",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "LHR",
					},
					{
						AirportID: 3,
						IATA:      "CDG",
					},
					{
						AirportID: 4,
						IATA:      "CPH",
					},
				},
				schedules: []*storagemodels.Schedule{
					{
						ArrivalTime:      NewString("1205"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(3),
						FlightNumber:     NewString("AF01"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("1405"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1305"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("BA01"),
						OperatingCarrier: NewString("BA"),
						Origin:           NewInt(3),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("1120"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(4),
						FlightNumber:     NewString("AF02"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("1205"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1310"),
						Destination:      NewInt(2),
						FlightNumber:     NewString("AF18"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(4),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1205",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "CDG",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1405",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1305",
							Destination:      "LHR",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "CDG",
							SegmentID:        0,
						},
					},
				},
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1120",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1000",
							Destination:      "CPH",
							FlightNumber:     "AF02",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1205",
							ArriveNextDay:    false,
							DepartureDate:    "2412",
							DepartureTime:    "1310",
							Destination:      "LHR",
							FlightNumber:     "AF18",
							OperatingCarrier: "AF",
							Origin:           "CPH",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "no segments",
			//enabled: true,
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
					},
					{
						AirportID: 2,
						IATA:      "LHR",
					},
					{
						AirportID: 3,
						IATA:      "CDG",
					},
					{
						AirportID: 4,
						IATA:      "CPH",
					},
				},
				schedules: []*storagemodels.Schedule{
					{
						ArrivalTime:      NewString("1205"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(3),
						FlightNumber:     NewString("AF01"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
					{
						ArrivalTime:      NewString("1120"),
						ArriveNextDay:    NewBool(false),
						DaysOperated:     NewString("1234567"),
						DepartureTime:    NewString("1000"),
						Destination:      NewInt(4),
						FlightNumber:     NewString("AF02"),
						OperatingCarrier: NewString("AF"),
						Origin:           NewInt(1),
						//ScheduleID    *int64
					},
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want:             []*itinerarymodels.Itinerary{},
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
					t.Errorf("%s - Expected error was %q: got %q", tt.name, tt.wantErrorMessage, err.Error())
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("\nExpected\n")
				spew.Dump(tt.want)
				fmt.Printf("\nGot:\n")
				spew.Dump(got)
				t.Errorf("%s", tt.name)
			}
		})
	}
}
