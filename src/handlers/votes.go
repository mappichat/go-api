package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mappichat/go-api.git/src/database"
)

func HandleVotes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID string `json:"post_id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		votes, err := database.ReadVotes(payload.PostID)
		if err != nil {
			return err
		}

		return c.JSON(votes)
	})

	router.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := database.Vote{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		payload.TimeStamp = time.Now()

		if err := database.InsertVote(&payload); err != nil {
			return err
		}

		return c.JSON(payload)
	})

	router.Patch("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID     string                 `json:"post_id"`
			UserHandle string                 `json:"user_handle"`
			Level      int8                   `json:"level"`
			UpdateBody map[string]interface{} `json:"update_body"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.UpdateVote(payload.PostID, payload.UserHandle, payload.Level, payload.UpdateBody); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	router.Delete("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID     string `json:"post_id"`
			UserHandle string `json:"user_handle"`
			Level      int8   `json:"level"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.DeleteVote(payload.PostID, payload.UserHandle, payload.Level); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
