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

	loads "github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/amadeusitgroup/miniplanes/storage/cmd/config"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/restapi/operations"
)

const (
	mongoHostParamName = "mongo-host"
	mongoHostDefault   = "mongo"
	mongoPortParamName = "mongo-port"
	mongoPortDefault   = 27017

	mongoDBNameParamName = "mongo-db-name"
	mongoDBNameDefault   = "miniplanes"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	var server *restapi.Server
	flag.StringVar(&config.MongoHost, mongoHostParamName, mongoHostDefault, "the mongo service name")
	flag.IntVar(&config.MongoPort, mongoPortParamName, mongoPortDefault, "the port of the mongo service")
	flag.StringVar(&config.MongoDBName, mongoDBNameParamName, mongoDBNameDefault, "name of the Mongo DB")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  storage-server [OPTIONS]\n\n")

		title := "miniplanes storage"
		fmt.Fprint(os.Stderr, title+"\n\n")
		desc := "needs to add a description"
		if desc != "" {
			fmt.Fprintf(os.Stderr, desc+"\n\n")
		}
		fmt.Fprintln(os.Stderr, flag.CommandLine.FlagUsages())
	}

	flag.Parse()

	log.Infof("Running storage version: %s", config.Version)
	log.Infof("Running storage with %s: %s", mongoHostParamName, config.MongoHost)
	log.Infof("Running storage with %s: %d", mongoPortParamName, config.MongoPort)
	log.Infof("Running storage with %s: %s", mongoDBNameParamName, config.MongoDBName)

	api := operations.NewStorageAPI(swaggerSpec)
	server = restapi.NewServer(api)
	defer server.Shutdown()

	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
	log.Info("Thanks for running storage. Hope it was OK")
}
