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
	"reflect"
	"testing"
	"time"

	itinerarymodels "github.com/amadeusitgroup/miniplanes/itineraries-server/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	storagemodels "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/davecgh/go-spew/spew"
)

var copenhagenTimezone *time.Location

var (
	NCE = &storagemodels.Airport{
		AirportID: 1,
		IATA:      "NCE",
		TZ:        "Europe/Paris",
		Latitude:  43.6584014893,
		Longitude: 7.215869903560001,
	}
	LHR = &storagemodels.Airport{
		AirportID: 2,
		IATA:      "LHR",
		TZ:        "Europe/London",
		Latitude:  51.4706,
		Longitude: -0.461941,
	}
	JFK = &storagemodels.Airport{
		AirportID: 3,
		IATA:      "JFK",
		TZ:        "America/New_York",
		Latitude:  40.63980103,
		Longitude: -73.77890015,
	}
	CDG = &storagemodels.Airport{
		AirportID: 4,
		IATA:      "CDG",
		TZ:        "Europe/Paris",
		Latitude:  49.0127983093,
		Longitude: 2.54999995232,
	}
	CPH = &storagemodels.Airport{
		AirportID: 5,
		IATA:      "CPH",
		TZ:        "Europe/Copenhagen",
		Latitude:  55.617900848389,
		Longitude: 12.656000137329,
	}
	SFO = &storagemodels.Airport{
		AirportID: 6,
		IATA:      "SFO",
		TZ:        "America/Los_Angeles",
		Latitude:  37.61899948120117,
		Longitude: -122.375,
	}
	SEA = &storagemodels.Airport{
		AirportID: 7,
		IATA:      "SEA",
		TZ:        "America/Los_Angeles",
		Latitude:  47.44900131225586,
		Longitude: -122.30899810791016,
	}
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

func NewSchedule(origin int32, departureTime string, destination int32, operatingCarrier, daysOperated, flightNumber string) *storagemodels.Schedule {
	return &storagemodels.Schedule{
		DaysOperated:     &daysOperated,
		DepartureTime:    &departureTime,
		Destination:      &destination,
		FlightNumber:     &flightNumber,
		OperatingCarrier: &operatingCarrier,
		Origin:           &origin,
	}
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
			wantErrorMessage: `unable to determine timezone for `,
			want:             []*itinerarymodels.Itinerary{},
		},
		{
			name: "no ID for To",
			args: args{
				airports: []*storagemodels.Airport{
					{
						AirportID: 1,
						IATA:      "NCE",
						TZ:        "Europe/Paris",
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
				airports: []*storagemodels.Airport{NCE, LHR, JFK},
				From:     "NCE",
				To:       "JFK",
			},
			wantErrorMessage: "No segments from NCE",         // no schedules no error
			want:             []*itinerarymodels.Itinerary{}, // no itineraries
		},
		{
			name: "one valid schedule",
			args: args{
				airports: []*storagemodels.Airport{NCE, LHR},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "1000", LHR.AirportID, "AF", "1234567", "AF01"),
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
							ArrivalTime:      "1114",
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
				airports: []*storagemodels.Airport{NCE, LHR},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "1000", LHR.AirportID, "AF", "1234567", "AF01"),

					NewSchedule(NCE.AirportID, "1005", LHR.AirportID, "BA", "1234567", "BA01"),
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
							ArrivalTime:      "1119",
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
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1114",
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
			name: "two valid schedules, one too early",
			args: args{
				airports: []*storagemodels.Airport{NCE, LHR},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "1000", LHR.AirportID, "AF", "1234567", "AF01"),
					NewSchedule(NCE.AirportID, "0830", LHR.AirportID, "BA", "1234567", "BA01"),
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0900",
				DepartureDate: "2412",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "2412:0900 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1114",
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
			name: "over night flight",
			args: args{
				airports: []*storagemodels.Airport{LHR, JFK},
				schedules: []*storagemodels.Schedule{
					NewSchedule(LHR.AirportID, "2205", JFK.AirportID, "BA", "1234567", "BA01"),
				},
				From:          "LHR",
				To:            "JFK",
				DepartureTime: "2100",
				DepartureDate: "0707",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "0707:2100 - LHR-JFK",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "0807",
							ArrivalTime:      "0145",
							DepartureDate:    "0707",
							DepartureTime:    "2205",
							Destination:      "JFK",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "LHR",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "over night flight to SFO",
			args: args{
				airports: []*storagemodels.Airport{LHR, JFK, SFO},
				schedules: []*storagemodels.Schedule{
					NewSchedule(LHR.AirportID, "2305", JFK.AirportID, "BA", "1234567", "BA01"),
					NewSchedule(JFK.AirportID, "0405", SFO.AirportID, "AA", "1234567", "AA01"),
					NewSchedule(CDG.AirportID, "0800", SEA.AirportID, "AA", "1234567", "AA02"),
				},
				From:          "LHR",
				To:            "SFO",
				DepartureTime: "2200",
				DepartureDate: "0707",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "0707:2200 - LHR-SFO",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "0807",
							ArrivalTime:      "0245",
							DepartureDate:    "0707",
							DepartureTime:    "2305",
							Destination:      "JFK",
							FlightNumber:     "BA01",
							OperatingCarrier: "BA",
							Origin:           "LHR",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0807",
							ArrivalTime:      "0746",
							DepartureDate:    "0807",
							DepartureTime:    "0405",
							Destination:      "SFO",
							FlightNumber:     "AA01",
							OperatingCarrier: "AA",
							Origin:           "JFK",
							SegmentID:        0,
						},
					},
				},
			},
		},
		{
			name: "2xtwo segments flights ",
			args: args{
				airports: []*storagemodels.Airport{NCE, LHR, CDG, CPH},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "1000", CDG.AirportID, "AF", "1234567", "AF01"),
					NewSchedule(CDG.AirportID, "1305", LHR.AirportID, "BA", "1234567", "BA01"),
					NewSchedule(NCE.AirportID, "1000", CPH.AirportID, "AF", "1234567", "AF02"),
					NewSchedule(CPH.AirportID, "1310", LHR.AirportID, "AF", "1234567", "AF18"),
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
							ArrivalTime:      "1243",
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
							ArrivalTime:      "1419",
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
				&itinerarymodels.Itinerary{
					Description: "2412:0800 - NCE-LHR",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "2412",
							ArrivalTime:      "1144",
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
							ArrivalTime:      "1319",
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
			name: "no segments",
			args: args{
				airports: []*storagemodels.Airport{NCE, LHR, CDG, CPH},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "1000", CDG.AirportID, "AF", "1234567", "AF01"),
					NewSchedule(NCE.AirportID, "1000", CPH.AirportID, "AF", "1234567", "AF02"),
				},
				From:          "NCE",
				To:            "LHR",
				DepartureTime: "0800",
				DepartureDate: "2412",
			},
			wantErrorMessage: "No segments from CPH",
			want:             []*itinerarymodels.Itinerary{},
		},
		{
			name: "multiple routes",
			args: args{
				airports: []*storagemodels.Airport{NCE, CDG, LHR, JFK, SEA},
				schedules: []*storagemodels.Schedule{
					NewSchedule(NCE.AirportID, "0840", LHR.AirportID, "AF", "1234567", "AF01"),
					NewSchedule(NCE.AirportID, "0900", CDG.AirportID, "AF", "1234567", "AF02"),
					NewSchedule(LHR.AirportID, "1300", JFK.AirportID, "AA", "1234567", "AA02"),
					NewSchedule(JFK.AirportID, "1700", SEA.AirportID, "AA", "1234567", "AA01"),
					NewSchedule(JFK.AirportID, "1700", CDG.AirportID, "AA", "1234567", "AA02"),
					NewSchedule(CDG.AirportID, "2000", SEA.AirportID, "AA", "1234567", "AA01"),
				},
				From:          "NCE",
				To:            "SEA",
				DepartureTime: "0800",
				DepartureDate: "0707",
			},
			wantErrorMessage: "",
			want: []*itinerarymodels.Itinerary{
				&itinerarymodels.Itinerary{
					Description: "0707:0800 - NCE-SEA",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "1044",
							DepartureDate:    "0707",
							DepartureTime:    "0900",
							Destination:      "CDG",
							FlightNumber:     "AF02",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "2315",
							DepartureDate:    "0707",
							DepartureTime:    "2000",
							Destination:      "SEA",
							FlightNumber:     "AA01",
							OperatingCarrier: "AA",
							Origin:           "CDG",
							SegmentID:        0,
						},
					},
				},
				&itinerarymodels.Itinerary{
					Description: "0707:0800 - NCE-SEA",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "0954",
							DepartureDate:    "0707",
							DepartureTime:    "0840",
							Destination:      "LHR",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "1640",
							DepartureDate:    "0707",
							DepartureTime:    "1300",
							Destination:      "JFK",
							FlightNumber:     "AA02",
							OperatingCarrier: "AA",
							Origin:           "LHR",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0807",
							ArrivalTime:      "0805",
							DepartureDate:    "0707",
							DepartureTime:    "1700",
							Destination:      "CDG",
							FlightNumber:     "AA02",
							OperatingCarrier: "AA",
							Origin:           "JFK",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0807",
							ArrivalTime:      "2315",
							DepartureDate:    "0807",
							DepartureTime:    "2000",
							Destination:      "SEA",
							FlightNumber:     "AA01",
							OperatingCarrier: "AA",
							Origin:           "CDG",
							SegmentID:        0,
						},
					},
				},
				&itinerarymodels.Itinerary{
					Description: "0707:0800 - NCE-SEA",
					ItineraryID: "MY ID",
					Segments: []*itinerarymodels.Segment{
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "0954",
							DepartureDate:    "0707",
							DepartureTime:    "0840",
							Destination:      "LHR",
							FlightNumber:     "AF01",
							OperatingCarrier: "AF",
							Origin:           "NCE",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "1640",
							DepartureDate:    "0707",
							DepartureTime:    "1300",
							Destination:      "JFK",
							FlightNumber:     "AA02",
							OperatingCarrier: "AA",
							Origin:           "LHR",
							SegmentID:        0,
						},
						&itinerarymodels.Segment{
							ArrivalDate:      "0707",
							ArrivalTime:      "2018",
							DepartureDate:    "0707",
							DepartureTime:    "1700",
							Destination:      "SEA",
							FlightNumber:     "AA01",
							OperatingCarrier: "AA",
							Origin:           "JFK",
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

func TestComputeArrivalDateTime(t *testing.T) { // TODO
	type args struct {
		year          int
		departureTime string
		departureDate string
		origin        *models.Airport
		destination   *models.Airport
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ComputeArrivalDateTime(tt.args.year, tt.args.departureTime, tt.args.departureDate, tt.args.origin, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeArrivalDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComputeArrivalDateTime() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ComputeArrivalDateTime() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
