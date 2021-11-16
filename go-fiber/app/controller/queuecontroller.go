package controller

import (
	"github.com/fusesuphasin/go-fiber/app/infrastructure"
	"github.com/fusesuphasin/go-fiber/app/interfaces"
	"github.com/gofiber/fiber/v2"
)

// A UserController belong to the interface layer.
type QueueController struct {
	Logger interfaces.Logger
}

func NewQueueController(logger interfaces.Logger) *QueueController {
	return &QueueController{
		Logger: logger,
	}
}

func (controller *QueueController) TestQueue(c *fiber.Ctx) error {
	query := "Hello World!!"

	infrastructure.SendQueue(query, "TestQueue")

	return c.JSON(&fiber.Map{
		"success": "ok",
	})
}

func (controller *QueueController) TestGetFromQueue(c *fiber.Ctx) error {

	infrastructure.GetQueue("TestQueue")

	return c.JSON(&fiber.Map{
		"success": "ok1",
	})
}