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
		Rx uint64 `json:"rx"`
		Wx uint64 `json:"wx"`
	} `json:"disk"`
	Network struct {
		Rx   uint64 `json:"rx"`
		Tx   uint64 `json:"tx"`
		Send uint64 `json:"sb"`
		Recv uint64 `json:"rb"`
	} `json:"network"`
	Uptime      uint64             `json:"uptime"`
	Hostname    string             `json:"hostname"`
	NodeID      string             `json:"node_id"`
	Interval    uint64             `json:"interval"`
	ReportTime  int64              `json:"report"`
	Temperature map[string]float64 `json:"temperature"`
	Metadata    ServerNode         `json:"metadata"`
	NodeAlive   bool               `json:"node_alive"`
}

type UsageStat struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}
