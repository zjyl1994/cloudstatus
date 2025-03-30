package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Cloudstatus report and web panel server",
	Run:   server.Server,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("config", "config.json", "Config file")
	serverCmd.Flags().String("listen", "127.0.0.1:10567", "Server listen address")
	serverCmd.Flags().String("db", "cloudstatus.db", "Database file")
	serverCmd.Flags().Int("alive", 60, "Alive time for nodes")
}
