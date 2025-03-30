package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/client"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Cloudstatus report agent",
	Run:   client.Client,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("report", "", "Remote report url")
	clientCmd.Flags().String("node", "", "Node ID")
	clientCmd.Flags().String("token", "", "Node token")
	clientCmd.Flags().Int("interval", 3, "Report interval in seconds")
	clientCmd.Flags().Bool("sensors", true, "Load tempature use lm-sensors")
}
