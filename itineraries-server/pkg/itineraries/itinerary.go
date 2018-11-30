package itineraries

import (
	"fmt"

	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/models"

	storageclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client"
	//httptransport "github.com/go-openapi/runtime/client"
)

type IATACode string

var (
	airportID2Airports map[int64]*models.Airport
	iata2Airports      map[string]*models.Airport
)

func airportIDFromIATA(IATA string) (int64, error) {
	//transport := httptransport.New(os.Getenv("TODOLIST_HOST"), "", nil)
	resp, err := storageclient.Default.Airports.GetAirports(nil)
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve airports: %v", err)
	}
	airports := resp.Payload
	if err != nil {
		return 0, fmt.Errorf("unable to get fetch airports from DB: %v", err)
	}
	for _, a := range airports {
		if IATA == a.IATA {
			return a.AirportID, nil
		}
	}
	return 0, fmt.Errorf("couldn't find Airport ID for airport with IATA code %q", IATA)
}

func airportIDFromName(name string) (int64, error) {
	resp, err := storageclient.Default.Airports.GetAirports(nil)
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve airports: %v", err)
	}
	airports := resp.Payload
	if err != nil {
		return 0, fmt.Errorf("unable to get fetch airports from DB: %v", err)
	}
	for _, a := range airports {
		if name == a.Name {
			return a.AirportID, nil
		}
	}
	return 0, fmt.Errorf("couldn't find Airport ID for airport with name %q", name)
}
