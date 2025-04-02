package client

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/fxamacker/cbor/v2"
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
	if reportUrl == "" { // remote url not set，print data
		slog.Error("report url not set")
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

	// init data for filling server status
	report(nodeId, reportUrl, token, time.Second, sensors)

	for {
		err := report(nodeId, reportUrl, token, time.Duration(interval)*time.Second, sensors)
		if err != nil {
			slog.Error("Measure error", slog.String("err", err.Error()))
		}
	}
}

func report(nodeId, reportUrl, token string, interval time.Duration, useSensors bool) error {
	samples, err := measure.Measure(interval, useSensors)
	if err != nil {
		slog.Error("Measure error", slog.String("err", err.Error()))
		return err
	}

	if nodeId == "" {
		samples.NodeID = samples.Host.Hostname
	} else {
		samples.NodeID = nodeId
	}

	bCbor, err := cbor.Marshal(samples)
	if err != nil {
		slog.Error("Marshal error", slog.String("err", err.Error()))
		return err
	}
	slog.Debug("Measure", slog.Int("len", len(bCbor)), slog.Any("data", samples))

	hReq, err := http.NewRequest(http.MethodPost, reportUrl, bytes.NewReader(bCbor))
	if err != nil {
		slog.Error("New report error", slog.String("err", err.Error()))
		return err
	}
	hReq.Header.Set("Content-Type", "application/cbor")
	hReq.Header.Set("Authorization", "Bearer "+token)

	hc := http.Client{Timeout: interval}
	resp, err := hc.Do(hReq)
	if err != nil {
		slog.Error("Report send error", slog.String("err", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Report send error",
			slog.Int("status", resp.StatusCode),
			slog.String("body", string(body)))
		return fmt.Errorf("bad server response code %d", resp.StatusCode)
	}

	return nil
}
