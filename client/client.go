package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/service/measure"
)

func Client(cmd *cobra.Command, args []string) {
	debugMode, err := cmd.Flags().GetBool("debug")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	if debugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	reportUrl, err := cmd.Flags().GetString("report")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	nodeId, err := cmd.Flags().GetString("node")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	interval, err := cmd.Flags().GetInt("interval")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	sensors, err := cmd.Flags().GetBool("sensors")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}

	for {
		intdur := time.Duration(interval) * time.Second
		samples, err := measure.Measure(intdur, sensors)
		if err != nil {
			slog.Error("Measure error", slog.String("err", err.Error()))
			continue
		}

		if nodeId == "" {
			samples.NodeID = samples.Host.Hostname
		} else {
			samples.NodeID = nodeId
		}

		bJson, err := json.Marshal(samples)
		if err != nil {
			slog.Error("Marshal error", slog.String("err", err.Error()))
			continue
		}

		slog.Debug("Measure", slog.Int("len", len(bJson)), slog.String("body", string(bJson)))

		if reportUrl == "" { // remote url not set，print data
			slog.Error("report url not set")
			continue
		}

		hReq, err := http.NewRequest(http.MethodPost, reportUrl, bytes.NewReader(bJson))
		if err != nil {
			slog.Error("New report error", slog.String("err", err.Error()))
			continue
		}
		hReq.Header.Set("Content-Type", "application/json")
		hReq.Header.Set("Authorization", "Bearer "+token)

		hc := http.Client{Timeout: intdur}
		resp, err := hc.Do(hReq)
		if err != nil {
			slog.Error("Report send error", slog.String("err", err.Error()))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			slog.Error("Report send error",
				slog.Int("status", resp.StatusCode),
				slog.String("body", string(body)))
		}

		resp.Body.Close()
	}
}
