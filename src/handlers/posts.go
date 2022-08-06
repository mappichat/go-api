package handlers

import (
	"fmt"
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
			Level        int8    `query:"level" json:"level"`
			MinLatitude  float32 `query:"min_latitude" json:"min_latitude" validate:"required"`
			MaxLatitude  float32 `query:"max_latitude" json:"max_latitude" validate:"required"`
			MinLongitude float32 `query:"min_longitude" json:"min_longitude" validate:"required"`
			MaxLongitude float32 `query:"max_longitude" json:"max_longitude" validate:"required"`
		}{}

		if err := c.QueryParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		posts := []database.Post{}
		if err := database.Sqldb.Select(
			&posts,
			"SELECT * FROM posts WHERE post_level=$1 AND latitude >= $2 AND latitude <= $3 AND longitude >= $4 AND longitude <= $5",
			payload.Level, payload.MinLatitude, payload.MaxLatitude, payload.MinLongitude, payload.MaxLongitude,
		); err != nil {
			return err
		}

		return c.JSON(posts)
	})

	router.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")
		payload := struct {
			Title     string  `json:"title" validate:"required"`
			Body      string  `json:"body" validate:"required"`
			Latitude  float64 `json:"latitude" validate:"required"`
			Longitude float64 `json:"longitude" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		newPost := database.Post{
			ID:        uuid.NewString(),
			AccountId: c.Locals("account_id").(string),
			Title:     payload.Title,
			Body:      payload.Body,
			Level:     0,
			Latitude:  payload.Latitude,
			Longitude: payload.Longitude,
			TimeStamp: time.Now().Round(time.Microsecond),
		}

		if _, err := database.Sqldb.NamedExec(
			"INSERT INTO posts (id, account_id, title, body, post_level, latitude, longitude, time_stamp) VALUES (:id,:account_id,:title,:body,:post_level,:latitude,:longitude,:time_stamp)",
			newPost,
		); err != nil {
			return err
		}

		return c.JSON(newPost)
	})

	router.Patch("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			ID         string `json:"id" validate:"required"`
			UpdateBody struct {
				Title string `json:"title" db:"title"`
				Body  string `json:"body" db:"body"`
			} `json:"update_body" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		setStmt := ""
		if payload.UpdateBody.Title != "" {
			setStmt += "title=:title,"
		}
		if payload.UpdateBody.Body != "" {
			setStmt += "body=:body"
		}

		if _, err := database.Sqldb.NamedExec(
			fmt.Sprintf("UPDATE posts SET %s WHERE id='%s' AND account_id='%s'", setStmt, payload.ID, c.Locals("account_id").(string)),
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
			ID string `query:"id" json:"id" validate:"required" db:"id"`
		}{}

		if err := c.QueryParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		if _, err := database.Sqldb.Exec(
			"DELETE FROM posts WHERE id=$1 AND account_id=$2",
			payload.ID, c.Locals("account_id").(string),
		); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
