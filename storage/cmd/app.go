package cmd

import (
	"github.com/amadeusitgroup/miniapp/storage"
	"github.com/spf13/cobra"
)

var port, host string

func init() {
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "port")
	RootCmd.PersistentFlags().StringVarP(&host, "host", "s", "localhost", "mongo host")
}

var RootCmd = &cobra.Command{
	Use: "storage",
	Run: func(cmd *cobra.Command, args []string) {
		a := storage.NewApplication(host, port)
		a.Run()
	},
}
