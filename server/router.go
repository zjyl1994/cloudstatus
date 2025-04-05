package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	cloudstatusfe "github.com/zjyl1994/cloudstatus/cloudstatus-fe"
	"github.com/zjyl1994/cloudstatus/infra/vars"
)

func Start(listen string) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	vars.App = app
	app.Use(cors.New())

	apiG := app.Group("/api")
	{
		apiG.Post("/report", handleAPIReport)
		apiG.Get("/overview", handleOverview)
		apiG.Get("/charts", handleCharts)
		apiG.Get("/nodes", handleNodes)
	}

	app.Use(filesystem.New(filesystem.Config{
		Root:         http.FS(cloudstatusfe.FrontendAssets),
		PathPrefix:   "build/client",
		Index:        "index.html",
		NotFoundFile: "build/client/index.html",
	}))
	return app.Listen(listen)
}
