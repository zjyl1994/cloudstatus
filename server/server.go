package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/zjyl1994/cloudstatus/infra/define"
	"github.com/zjyl1994/cloudstatus/infra/vars"
	"github.com/zjyl1994/cloudstatus/service/record"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Server(cmd *cobra.Command, args []string) {
	var err error
	vars.DebugMode, err = cmd.Flags().GetBool("debug")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	if vars.DebugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	vars.Listen, err = cmd.Flags().GetString("listen")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}

	vars.NodeAliveTimeout, err = cmd.Flags().GetInt("alive")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}

	// load config
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}

	bConf, err := os.ReadFile(configFile)
	if err != nil {
		slog.Error("Load config", slog.String("err", err.Error()))
		return
	}

	var cfg define.ServerConfig
	err = json.Unmarshal(bConf, &cfg)
	if err != nil {
		slog.Error("Unmarshal config", slog.String("err", err.Error()))
		return
	}
	vars.Config = cfg

	// init db
	dbFile, err := cmd.Flags().GetString("db")
	if err != nil {
		slog.Error("Error", slog.String("err", err.Error()))
		return
	}
	gormLogger := slogGorm.New()
	vars.DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		slog.Error("Open database", slog.String("err", err.Error()))
		return
	}
	err = vars.DB.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		slog.Error("Switch to WAL", slog.String("err", err.Error()))
		return
	}
	err = vars.DB.AutoMigrate(&define.MeasureRecord{})
	if err != nil {
		slog.Error("Database migrate", slog.String("err", err.Error()))
		return
	}
	// clean data
	cleanDataFn := func() {
		if err = record.CleanRecord(); err != nil {
			slog.Error("Measure data clean", slog.String("err", err.Error()))
		}
	}
	cronInstance := cron.New()
	cronInstance.AddFunc("@daily", cleanDataFn)
	cronInstance.Start()
	cleanDataFn()
	// run web server
	webErrCh := make(chan error, 1)
	go func(ch chan error) {
		slog.Info("Starting server", slog.String("listen", vars.Listen))
		if err = Start(vars.Listen); err != nil {
			webErrCh <- err
		}
	}(webErrCh)
	// register ctrl+c handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-webErrCh:
		if err != nil {
			slog.Error("Web server faild", slog.String("err", err.Error()))
			return
		}
	case sig := <-sigChan:
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			slog.Info("Signal receive", slog.String("singal", sig.String()))

			cronInstance.Stop()

			if vars.App != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				vars.App.ShutdownWithContext(ctx)
			}
			if vars.DB != nil {
				rawDB, _ := vars.DB.DB()
				if rawDB != nil {
					rawDB.Close()
				}
			}
		}
	}
}
