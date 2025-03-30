package server

import (
	"log/slog"
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

		traffic, err := record.GetNetTraffic()
		if err != nil {
			return nil, err
		}
		tm := make(map[string]define.TrafficCalcResult, len(traffic))
		for _, t := range traffic {
			tm[t.NodeId] = t
		}

		result := make([]define.StatExchangeFormat, 0, len(vars.Nodes))
		for _, node := range vars.Nodes {
			stat, ok := statCache.Get(node.ID)
			if !ok {
				continue
			}

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

func handleDetail(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
