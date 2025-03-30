package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/cloudstatus/infra/vars"
)

func Start(listen string) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	vars.App = app

	apiG := app.Group("/api")
	{
		apiG.Post("/report", handleAPIReport)
		apiG.Get("/overview", handleOverview)
		apiG.Get("/detail", handleDetail)
	}
	return app.Listen(listen)
}
