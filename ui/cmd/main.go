package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/amadeusitgroup/miniapp/ui/pkg/www"
)

var (
	Port                   int
	StorageServicePort     int
	StorageServiceName     string
	ItinerariesServicePort int
	ItinerariesServiceName string
)

func main() {

	flag.IntVar(&Port, "port", 8080, "defines the ui port")
	flag.StringVar(&StorageServiceName, "storage-service", "storage", "defines the storage server endpoint")
	flag.IntVar(&StorageServicePort, "storage-port", 8090, "defines the storage service port")
	flag.StringVar(&ItinerariesServiceName, "itineraries-service", "itineraries", "defines the itineraries server endpoint")
	flag.IntVar(&ItinerariesServicePort, "itineraries-port", 8100, "defines the itineraries service port")

	flag.Parse()

	serverCfg := www.Config{
		Host:         "localhost:" + strconv.Itoa(Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	htmlServer := www.Start(serverCfg)
	defer htmlServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("main : shutting down")
}
