package main

import (
	"log"

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
	// initialize new Comment struct
	cmt := new(Comment)

	// parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Printf("error parsing cmt %v: ", err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
}
