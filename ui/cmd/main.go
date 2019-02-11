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
