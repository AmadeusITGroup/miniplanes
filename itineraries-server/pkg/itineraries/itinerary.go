package itineraries

import (
	"fmt"

	"github.com/amadeusitgroup/miniapp/itineraries-server/pkg/db"
)

type IATACode string

var (
	airportID2Airports map[int64]*db.Airport
	iata2Airports      map[string]*db.Airport
)

func airportIDFromIATA(IATA string) (int64, error) {
	airports, err := db.GetAirports()
	if err != nil {
		return 0, fmt.Errorf("unable to get fetch airports from DB: %v", err)
	}
	for _, a := range airports {
		if IATA == a.IATA {
			return a.ID, nil
		}
	}
	return 0, fmt.Errorf("couldn't find Airport ID for airport with IATA code %q", IATA)
}

func airportIDFromName(name string) (int64, error) {
	airports, err := db.GetAirports() // TODO: REST call
	if err != nil {
		return 0, fmt.Errorf("unable to get fetch airports from DB: %v", err)
	}
	for _, a := range airports {
		if name == a.Name {
			return a.ID, nil
		}
	}
	return 0, fmt.Errorf("couldn't find Airport ID for airport with name %q", name)
}
