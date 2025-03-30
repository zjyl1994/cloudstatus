package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/service/sensors"
)

// sensorsCmd represents the sensors command
var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Check lm-sensors output",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := sensors.Get()
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		fmt.Println("Sensors Output:")
		for sensorName, temp := range result {
			fmt.Printf("%s: %.2f\n", sensorName, temp)
		}
	},
}

func init() {
	rootCmd.AddCommand(sensorsCmd)
}
