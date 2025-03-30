package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/cloudstatus/infra/define"
	"github.com/zjyl1994/cloudstatus/infra/rwmap"
	"github.com/zjyl1994/cloudstatus/infra/vars"
	"github.com/zjyl1994/cloudstatus/service/record"
	"golang.org/x/sync/singleflight"
)

var (
	statCache  = new(rwmap.Map[string, define.StatExchangeFormat])
	overviewSf singleflight.Group
	chartsSf   singleflight.Group
)

func handleAPIReport(c *fiber.Ctx) error {
	// check token
	authHeader := c.Get(fiber.HeaderAuthorization)
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	authHeader = strings.TrimSpace(authHeader)
	if authHeader != vars.Token {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	// parse data
	var data define.StatExchangeFormat
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if vars.DebugMode {
		slog.Debug("Receive data", slog.Any("data", data))
	}
	// save data
	statCache.Set(data.NodeID, data)
	err = record.WriteRecord(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

type overviewResponse struct {
	UpdateAt int64                       `json:"update_at"`
	Nodes    []define.StatExchangeFormat `json:"nodes"`
}

func handleOverview(c *fiber.Ctx) error {
	ret, err, _ := overviewSf.Do("overview", func() (interface{}, error) {
		now := time.Now().Unix()

		// calcaute monthly traffic
		traffic, err := record.GetNetTraffic()
		if err != nil {
			return nil, err
		}
		tm := make(map[string]define.TrafficCalcResult, len(traffic))
		for _, t := range traffic {
			tm[t.NodeId] = t
		}
		// get all node stat from cache
		result := make([]define.StatExchangeFormat, 0, len(vars.Nodes))
		for _, node := range vars.Nodes {
			stat, ok := statCache.Get(node.ID)
			if !ok {
				result = append(result, define.StatExchangeFormat{
					NodeID:    node.ID,
					NodeName:  node.Name,
					NodeAlive: false,
				})
				continue
			}

			stat.NodeName = node.Name
			stat.NodeAlive = (time.Now().Unix() - stat.ReportTime) < int64(vars.NodeAliveTimeout)

			// set monthly traffic data
			if td, ok := tm[node.ID]; ok {
				stat.Network.Send = td.NetSend
				stat.Network.Recv = td.NetRecv
			}

			result = append(result, stat)
		}

		return overviewResponse{UpdateAt: now, Nodes: result}, nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(ret)
}

type ChartsResponse struct {
	CPU         []ChartsPercentItem            `json:"cpu"`
	Memory      []ChartsPercentItem            `json:"memory"`
	Swap        []ChartsPercentItem            `json:"swap"`
	DiskSpeed   []ChartsSpeedItem              `json:"disk_speed"`
	NetSpeed    []ChartsSpeedItem              `json:"net_speed"`
	Load        []ChartsLoadItem               `json:"load"`
	Temperature map[string][]ChartsPercentItem `json:"temperature"`
}

type ChartsPercentItem struct {
	DateTime string  `json:"time"`
	Value    float64 `json:"value"`
}

type ChartsSpeedItem struct {
	DateTime string `json:"time"`
	Rx       int64  `json:"rx"`
	Tx       int64  `json:"tx"`
}

type ChartsLoadItem struct {
	DateTime string  `json:"time"`
	Load1    float64 `json:"load1"`
	Load5    float64 `json:"load5"`
	Load15   float64 `json:"load15"`
}

func handleCharts(c *fiber.Ctx) error {
	now := time.Now().Unix()
	nodeId := c.Query("id")
	if nodeId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing node id")
	}
	endTime := c.QueryInt("end")
	if endTime == 0 {
		endTime = int(now)
	}
	startTime := c.QueryInt("start")
	if startTime == 0 {
		startTime = endTime - 3600
	}
	if startTime > endTime {
		return c.Status(fiber.StatusBadRequest).SendString("Start time must be less than end time")
	}
	sresp, err, _ := chartsSf.Do(fmt.Sprintf("charts-%s-%d-%d", nodeId, startTime, endTime), func() (interface{}, error) {
		// load data
		mrList, err := record.LoadRecord(nodeId, int64(startTime), int64(endTime))
		if err != nil {
			return nil, err
		}
		sort.Slice(mrList, func(i, j int) bool {
			return mrList[i].Timestamp < mrList[j].Timestamp
		})
		// convert to resp
		resp := ChartsResponse{
			CPU:         make([]ChartsPercentItem, 0, len(mrList)),
			Memory:      make([]ChartsPercentItem, 0, len(mrList)),
			Swap:        make([]ChartsPercentItem, 0, len(mrList)),
			DiskSpeed:   make([]ChartsSpeedItem, 0, len(mrList)),
			NetSpeed:    make([]ChartsSpeedItem, 0, len(mrList)),
			Load:        make([]ChartsLoadItem, 0, len(mrList)),
			Temperature: make(map[string][]ChartsPercentItem),
		}
		for _, mr := range mrList {
			dateTime := time.Unix(mr.Timestamp, 0).Format(time.DateTime)
			resp.CPU = append(resp.CPU, ChartsPercentItem{
				DateTime: dateTime,
				Value:    mr.CPU,
			})
			resp.Memory = append(resp.Memory, ChartsPercentItem{
				DateTime: dateTime,
				Value:    mr.Memory,
			})
			resp.Swap = append(resp.Swap, ChartsPercentItem{
				DateTime: dateTime,
				Value:    mr.Swap,
			})
			resp.DiskSpeed = append(resp.DiskSpeed, ChartsSpeedItem{
				DateTime: dateTime,
				Rx:       int64(mr.DiskRx),
				Tx:       int64(mr.DiskWx),
			})
			resp.NetSpeed = append(resp.NetSpeed, ChartsSpeedItem{
				DateTime: dateTime,
				Rx:       int64(mr.NetRx),
				Tx:       int64(mr.NetTx),
			})
			resp.Load = append(resp.Load, ChartsLoadItem{
				DateTime: dateTime,
				Load1:    mr.Load1,
				Load5:    mr.Load5,
				Load15:   mr.Load15,
			})
			var tempMap map[string]float64
			err = json.Unmarshal([]byte(mr.Temperature), &tempMap)
			if err == nil {
				for k, v := range tempMap {
					resp.Temperature[k] = append(resp.Temperature[k], ChartsPercentItem{
						DateTime: dateTime,
						Value:    v,
					})
				}
			}
		}
		return resp, nil
	})

	if err != nil {
		return err
	}

	return c.JSON(sresp)
}
