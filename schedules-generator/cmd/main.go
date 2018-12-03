package main

import (
	"log"
	"os"

	storageclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client"
)

const (
	dbName             = "miniapp"
	routesCollection   = "routes"
	airportsCollection = "airports"
	airlinesCollection = "airlines"
)

func main() {
	// airports
	airportsResp, err := storageclient.Default.Airports.GetAirports(nil)
	if err != nil {
		log.Fatalf("couldn't retrieve airports: %v", err)
	}
	airports := airportsResp.Payload
	if err != nil {
		log.Fatalf("unable to get fetch airports from DB: %v", err)
	}
	// airlines
	airlinesResp, err := storageclient.Default.Airlines.GetAirlines(nil)
	if err != nil {
		log.Fatalf("couldn't retrieve airlines: %v", err)
	}
	airlines := airlinesResp.Payload
	if err != nil {
		log.Fatalf("unable to get fetch airports from DB: %v", err)
	}

	// routes
	routesResp, err := storageclient.Default.Routes.GetRoutes(nil)
	if err :=

	// now create the schedules

	os.Exit(0)
}
