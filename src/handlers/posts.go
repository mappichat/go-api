package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
)

func HandlePosts(router fiber.Router) {
	router.Get("/posts", func(c *fiber.Ctx) error {
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

	router.Post("/posts", func(c *fiber.Ctx) error {
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

	router.Patch("/posts", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID         string                 `json:"id"`
			UpdateBody map[string]interface{} `json:"update_body"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.UpdatePost(payload.ID, payload.UpdateBody); err != nil {
			return err
		}

		post, err := database.ReadPost(payload.ID)
		if err != nil {
			return err
		}

		return c.JSON(*post)
	})

	router.Delete("/posts", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID string `json:"id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.DeletePost(payload.ID); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
