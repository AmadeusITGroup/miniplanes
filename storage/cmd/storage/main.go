
package main

import (
	"fmt"
	"os"

	loads "github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/amadeusitgroup/miniapp/storage/cmd/config"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi"
	"github.com/amadeusitgroup/miniapp/storage/pkg/gen/restapi/operations"
)

const (
	mongoHostParamName = "mongo-host"
	mongoHostDefault   = "mongo"
	mongoPortParamName = "mongo-port"
	mongoPortDefault   = 27017

	mongoDBNameParamName = "mongo-db-name"
	mongoDBNameDefault   = "miniapp"
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

		title := "miniapp storage"
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
