package server

// formatFloat 将浮点数格式化为保留两位小数
func formatFloat(value float64) float64 {
	return float64(int64(value*100)) / 100
}

// formatFloatMap 将map中的浮点数值格式化为保留两位小数
func formatFloatMap(data map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(data))
	for k, v := range data {
		result[k] = formatFloat(v)
	}
	return result
}