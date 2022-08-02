package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
)

func HandleReplies(router fiber.Router) {
	router.Get("/replies", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID string `json:"post_id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		replies, err := database.ReadReplies(payload.PostID)
		if err != nil {
			return err
		}

		return c.JSON(replies)
	})

	router.Post("/replies", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID     string `json:"post_id"`
			UserHandle string `json:"user_handle"`
			Body       string `json:"body"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		newReply := database.Reply{
			PostID:     payload.PostID,
			ID:         uuid.NewString(),
			UserHandle: payload.UserHandle,
			Body:       payload.Body,
			Timestamp:  time.Now(),
		}

		if err := database.InsertReply(&newReply); err != nil {
			return err
		}

		return c.JSON(newReply)
	})

	router.Patch("/replies", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID     string                 `json:"post_id"`
			ID         string                 `json:"id"`
			UpdateBody map[string]interface{} `json:"update_body"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.UpdateReply(payload.PostID, payload.ID, payload.UpdateBody); err != nil {
			return err
		}

		reply, err := database.ReadReply(payload.PostID, payload.ID)
		if err != nil {
			return err
		}

		return c.JSON(*reply)
	})

	router.Delete("/replies", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID string `json:"post_id"`
			ID     string `json:"id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		if err := database.DeleteReply(payload.PostID, payload.ID); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
