package define

type MeasureRecord struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;not null"`
	NodeID    string `gorm:"index:ix_mr_node_time"`
	Timestamp int64  `gorm:"index:ix_mr_node_time"`
	CPU       float64
	Memory    float64
	Swap      float64
	Disk      float64
	Load1     float64
	Load5     float64
	Load15    float64
	DiskRx    uint64
	DiskWx    uint64
	NetRx     uint64
	NetTx     uint64
	NetSend   uint64 `gorm:"index:ix_mr_node_time"`
	NetRecv   uint64 `gorm:"index:ix_mr_node_time"`
}

type TemperatureRecord struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;not null"`
	NodeID      string `gorm:"index:ix_tr_node_time"`
	Timestamp   int64  `gorm:"index:ix_tr_node_time"`
	Name        string
	Temperature float64
}

type TrafficCalcResult struct {
	NodeId  string `gorm:"column:node_id"`
	NetSend uint64 `gorm:"column:net_send"`
	NetRecv uint64 `gorm:"column:net_recv"`
}
