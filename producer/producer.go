package main

import (
	"github.com/gofiber/fiber/v2"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	app := fiber.New()
	api := app.Group("/api/vi")

	api.POST("/comments", createComment)

	app.Listen(":3000")
}

func createComment(c *fiber.Ctx) error {
	cmt := new(Comment)
}
