package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
	utils "github.com/mappichat/go-api.git/src/utils"
)

var PORT = os.Getenv("PORT")

func main() {
	utils.ConfigureEnv()

	if err := database.Initialize(strings.Split(utils.Env.DB_HOST, ",")); err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})

	api := app.Group("/api")

	api.Get("/posts", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")
		payload := struct {
			Level          int8    `json:"level"`
			Latitude       float32 `json:"latitude"`
			Longitude      float32 `json:"longitude"`
			LatitudeDelta  float32 `json:"latitude_delta"`
			LongitudeDelta float32 `json:"longitude_delta"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		posts, err := database.ReadPosts(
			payload.Level, payload.Latitude,
			payload.Longitude, payload.LatitudeDelta,
			payload.LongitudeDelta,
		)

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
			ID:         uuid.NewString(),
			Title:      payload.Title,
			Body:       payload.Body,
			UserHandle: "@test_handle",
			Timestamp:  time.Now(),
			Latitude:   payload.Latitude,
			Longitude:  payload.Longitude,
			Level:      0,
			ReplyCount: 0,
			UpVotes:    0,
			DownVotes:  0,
		}

		if err := database.InsertPost(&newPost); err != nil {
			return err
		}

		return c.JSON(newPost)
	})

	log.Fatal(app.Listen(":" + utils.Env.PORT))
}
