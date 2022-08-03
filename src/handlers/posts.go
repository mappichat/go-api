package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
)

func HandlePosts(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")
		payload := struct {
			Level int8     `json:"level"`
			Tiles []string `json:"tiles"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		posts, err := database.ReadPosts(payload.Level, payload.Tiles)

		if err != nil {
			return err
		}

		return c.JSON(posts)
	})

	router.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")
		payload := struct {
			Tile      string  `json:"tile"`
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
			Tile:       payload.Tile,
			Title:      payload.Title,
			Body:       payload.Body,
			AccountId:  c.Locals("account_id").(string),
			UserHandle: c.Locals("user_handle").(string),
			TimeStamp:  time.Now(),
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

	router.Patch("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			Level      int8                   `json:"level"`
			Tile       string                 `json:"tile"`
			ID         string                 `json:"id"`
			UpdateBody map[string]interface{} `json:"update_body"`
		}{}

		newMap := map[string]interface{}{}
		if _, ok := payload.UpdateBody["title"]; ok {
			newMap["title"] = payload.UpdateBody["title"]
		}
		if _, ok := payload.UpdateBody["body"]; ok {
			newMap["body"] = payload.UpdateBody["body"]
		}
		payload.UpdateBody = newMap

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.UpdatePost(payload.ID, payload.Tile, payload.Level, c.Locals("account_id").(string), payload.UpdateBody); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	router.Delete("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID    string `json:"id"`
			Tile  string `json:"tile"`
			Level int8   `json:"level"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.DeletePost(payload.ID, payload.Tile, payload.Level, c.Locals("account_id").(string)); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
