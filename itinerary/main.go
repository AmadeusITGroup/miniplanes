package main

import (
	"log"
	"os"

	"github.com/amadeusitgroup/miniapp/itinerary/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
