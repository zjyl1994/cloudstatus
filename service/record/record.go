package record

import (
	"time"

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

func CleanRecord() error {
	validNodeMap := make(map[string]struct{})
	for _, node := range vars.Nodes {
		validNodeMap[node.ID] = struct{}{}
	}
	validNodes := make([]string, 0, len(validNodeMap))
	for node := range validNodeMap {
		validNodes = append(validNodes, node)
	}
	currentDayInMonth := time.Now().Day()
	return vars.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("node_id NOT IN ?", validNodes).Delete(&define.MeasureRecord{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("node_id NOT IN ?", validNodes).Delete(&define.TemperatureRecord{}).Error
		if err != nil {
			return err
		}
		for _, node := range vars.Nodes {
			if node.ResetDay == currentDayInMonth {
				err = tx.Where("node_id =?", node.ID).Delete(&define.MeasureRecord{}).Error
				if err != nil {
					return err
				}
				err = tx.Where("node_id = ?", node.ID).Delete(&define.TemperatureRecord{}).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
