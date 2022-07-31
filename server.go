package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/database"
	utils "github.com/mappichat/go-api.git/utils"
)

var PORT = os.Getenv("PORT")

func main() {
	utils.ConfigureEnv()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})

	api := app.Group("/api")

	api.Get("/posts", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")

		posts, err := database.ReadPosts()
		if err != nil {
			return err
		}
		return c.JSON(posts)
	})

	api.Post("/posts", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")
		payload := struct {
			Title     string  `json:"title"`
			Body      string  `json:"body"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		newPost := database.Post{
			Id:         uuid.NewString(),
			Title:      payload.Title,
			Body:       payload.Body,
			UserHandle: "@test_handle",
			Timestamp:  time.Now(),
			Latitude:   payload.Latitude,
			Longitude:  payload.Longitude,
			Level:      0,
			Replies:    []database.Reply{},
			UpVotes:    []string{},
			DownVotes:  []string{},
		}

		if err := database.InsertPost(&newPost); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":" + utils.Env.PORT))
}
