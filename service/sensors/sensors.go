package sensors

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
)

var tempRegex = regexp.MustCompile(`^temp\d+_input$`)

func Get() (map[string]float64, error) {
	cmd := exec.Command("sensors", "-j")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("call sensors: %w", err)
	}

	return parseSensorsOutput(output)
}

func parseSensorsOutput(output []byte) (map[string]float64, error) {
	var result map[string]any
	err := json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	temperatureMap := make(map[string]float64)

SENSOR_LOOP:
	for sensorName, sensorDataRaw := range result {
		sensorData, ok := sensorDataRaw.(map[string]any)
		if !ok {
			continue SENSOR_LOOP
		}

		for _, fieldValueRaw := range sensorData {
			fieldValue, ok := fieldValueRaw.(map[string]any)
			if !ok {
				continue
			}

			for tempKey, tempValue := range fieldValue {
				if tempRegex.MatchString(tempKey) {
					tempValueFloat, ok := tempValue.(float64)
					if !ok {
						continue
					}

					temperatureMap[sensorName] = tempValueFloat

					continue SENSOR_LOOP
				}
			}
		}
	}
	return temperatureMap, nil
}
