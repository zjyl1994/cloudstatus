package record

import (
	"github.com/zjyl1994/cloudstatus/infra/define"
	"github.com/zjyl1994/cloudstatus/infra/vars"
	"gorm.io/gorm"
)

func WriteRecord(def *define.StatExchangeFormat) error {
	return vars.DB.Transaction(func(tx *gorm.DB) error {
		var measure define.MeasureRecord
		measure.NodeID = def.NodeID
		measure.Timestamp = def.ReportTime
		measure.CPU = def.Percent.CPU
		measure.Memory = def.Percent.Mem
		measure.Swap = def.Percent.Swap
		measure.Disk = def.Percent.Disk
		measure.Load1 = def.Load.Load1
		measure.Load5 = def.Load.Load5
		measure.Load15 = def.Load.Load15
		measure.DiskRx = def.Disk.Rx
		measure.DiskWx = def.Disk.Wx
		measure.NetRx = def.Network.Rx
		measure.NetTx = def.Network.Tx
		measure.NetSend = def.Network.Send
		measure.NetRecv = def.Network.Recv
		err := tx.Create(&measure).Error
		if err != nil {
			return err
		}

		ts := make([]define.TemperatureRecord, 0, len(def.Temperature))
		for name, temp := range def.Temperature {
			ts = append(ts, define.TemperatureRecord{
				NodeID:      def.NodeID,
				Timestamp:   def.ReportTime,
				Name:        name,
				Temperature: temp,
			})
		}
		return tx.CreateInBatches(ts, len(ts)).Error
	})
}

func GetNetTraffic() ([]define.TrafficCalcResult, error) {
	var results []define.TrafficCalcResult
	err := vars.DB.Model(&define.MeasureRecord{}).
		Select("node_id,SUM(net_send) AS net_send, SUM(net_recv) AS net_recv").
		Group("node_id").
		Find(&results).Error
	return results, err
}
