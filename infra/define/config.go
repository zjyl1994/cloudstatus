package define

type ServerConfig struct {
	Token string       `json:"token"`
	Title string       `json:"title"`
	Nodes []ServerNode `json:"nodes"`
}

type ServerNode struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Location string `json:"location"`
	ResetDay int    `json:"reset_day"`
}
