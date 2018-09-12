package cmd

import (
	"github.com/amadeusitgroup/miniapp/backend/pkg/backend"
	"github.com/spf13/cobra"
)

var Mongo string
var Port string

func init() {
	RootCmd.PersistentFlags().StringVarP(&Mongo, "mongo", "m", "localhost", "mongo host")
	RootCmd.PersistentFlags().StringVarP(&Port, "port", "p", "8080", "port")
}

var RootCmd = &cobra.Command{
	Use: "backend",
	Run: func(cmd *cobra.Command, args []string) {
		a := backend.NewApplication(Port, Mongo)
		a.Run()
	},
}
