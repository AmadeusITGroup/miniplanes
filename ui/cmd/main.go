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
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/amadeusitgroup/miniapp/ui/cmd/config"
	"github.com/amadeusitgroup/miniapp/ui/pkg/www"
)

const (
	portParamName                  = "port"
	portDefault                    = 8080
	storageHostParamName           = "storage-host"
	storageHostDefault             = "storage"
	storagePortParamName           = "storage-port"
	storagePortDefault             = 12345
	itinerariesServerPortParamName = "itineraries-server-port"
	itinerariesServerPortDefault   = 54321
	itinerariesServerHostParamName = "itineraries-server-host"
	itinerariesServerHostDefault   = "itineraries-server"
)

func main() {

	flag.IntVar(&config.Port, portParamName, portDefault, "defines the ui port")
	flag.StringVar(&config.StorageHost, storageHostParamName, storageHostDefault, "defines the storage server endpoint")
	flag.IntVar(&config.StoragePort, storagePortParamName, storagePortDefault, "defines the storage service port")
	flag.StringVar(&config.ItinerariesServerHost, itinerariesServerHostParamName, itinerariesServerHostDefault, "defines the itineraries server endpoint")
	flag.IntVar(&config.ItinerariesServerPort, itinerariesServerPortParamName, itinerariesServerPortDefault, "defines the itineraries service port")

	flag.Parse()

	serverCfg := www.Config{
		Host:         "localhost:" + strconv.Itoa(config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	htmlServer := www.Start(serverCfg)
	defer htmlServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("Thanks for running ui. Hope it was OK")
}
