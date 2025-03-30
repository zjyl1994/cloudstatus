package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	}
	return app.Listen(listen)
}
