package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/service/sensors"
)

// sensorsCmd represents the sensors command
var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Check sensors output",
	Long:  `Get sensors output and check it works.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sensorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sensorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
