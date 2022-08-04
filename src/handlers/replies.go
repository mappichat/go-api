package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
)

func HandleReplies(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID string `json:"id" validate:"required" db:"id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		replies := []database.Reply{}
		database.Sqldb.Select(
			&replies,
			"SELECT * FROM replies WHERE id=$1",
			payload.ID,
		)

		return c.JSON(replies)
	})

	router.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID    string  `json:"post_id" validate:"required" db:"post_id"`
			Body      string  `json:"body" validate:"required" db:"body"`
			Latitude  float64 `json:"latitude" validate:"required" db:"latitude"`
			Longitude float64 `json:"longitude" validate:"required" db:"longitude"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		newReply := database.Reply{
			ID:        uuid.NewString(),
			PostID:    payload.PostID,
			AccountId: c.Locals("account_id").(string),
			Body:      payload.Body,
			Latitude:  payload.Latitude,
			Longitude: payload.Longitude,
			TimeStamp: time.Now().Round(time.Microsecond),
		}

		if _, err := database.Sqldb.NamedExec(
			"INSERT INTO replies (id, post_id, account_id, body, latitude, longitude, time_stamp)",
			newReply,
		); err != nil {
			return err
		}

		return c.JSON(newReply)
	})

	router.Patch("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID         string `json:"id" validate:"required"`
			UpdateBody struct {
				Body string `json:"body" db:"body"`
			} `json:"update_body" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		setStmt := ""
		if payload.UpdateBody.Body != "" {
			setStmt += "body=:body"
		}

		if _, err := database.Sqldb.NamedExec(
			fmt.Sprintf("UPDATE replies SET %s WHERE id=%s AND account_id=%s", setStmt, payload.ID, c.Locals("account_id").(string)),
			payload.UpdateBody,
		); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	router.Delete("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID string `json:"id" validate:"required" db:"id"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		if _, err := database.Sqldb.Exec(
			"DELETE FROM replies WHERE id=$1 AND account_id=$2",
			payload.ID, c.Locals("account_id").(string),
		); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
