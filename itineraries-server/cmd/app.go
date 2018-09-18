package cmd

import (
	"github.com/amadeusitgroup/miniapp/itinerary/pkg/itinerary"
	"github.com/spf13/cobra"
)

var Port string

func init() {
	RootCmd.PersistentFlags().StringVarP(&Port, "port", "p", "8080", "port")
	RootCmd.PersistentFlags().StringVarP(&Mongo, "mongo", "m", "localhost", "mongo host")
}

var RootCmd = &cobra.Command{
	Use: "itinerary",
	Run: func(cmd *cobra.Command, args []string) {
		a := itinerary.NewApplication(Port)
		a.Run()
	},
}
