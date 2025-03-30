package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/cloudstatus/infra/vars"
)

func Start(listen string) error {
	app := fiber.New()
	vars.App = app

	return app.Listen(listen)
}
