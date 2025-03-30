package define

type StatExchangeFormat struct {
	Percent struct {
		CPU  float64 `json:"cpu"`
		Mem  float64 `json:"mem"`
		Swap float64 `json:"swap"`
		Disk float64 `json:"disk"`
	} `json:"percent"`
	Load struct {
		Load1  float64 `json:"load1"`
		Load5  float64 `json:"load5"`
		Load15 float64 `json:"load15"`
	} `json:"load"`
	Memory UsageStat `json:"memory"`
	Swap   UsageStat `json:"swap"`
	Disk   struct {
		UsageStat
		Read  uint64 `json:"read"`
		Write uint64 `json:"write"`
	} `json:"disk"`
	Network struct {
		Rx   uint64 `json:"rx"`
		Tx   uint64 `json:"tx"`
		Send uint64 `json:"send"`
		Recv uint64 `json:"recv"`
	} `json:"network"`
	Uptime      uint64             `json:"uptime"`
	Hostname    string             `json:"hostname"`
	Interval    uint64             `json:"interval"`
	ReportTime  int64              `json:"report"`
	Temperature map[string]float64 `json:"temperature"`
}

type UsageStat struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}
