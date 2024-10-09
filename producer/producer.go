package main

import (
	"encoding/json"
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
	// convert body into bytes and send it to kafka
	cmtIntBytes, err := json.Marshal(cmt)
	PushCommentToQueue("comments", cmtIntBytes)

	// Return Comment in JSON format
	c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfuly",
		"comment": cmt,
	})
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		return err
	}
	return err
}
