package vars

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/cloudstatus/infra/define"
	"gorm.io/gorm"
)

var (
	DebugMode bool
	App       *fiber.App
	DB        *gorm.DB
	Token     string
	Listen    string
	Nodes     []define.ServerNode
)
