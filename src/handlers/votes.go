package handlers

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mappichat/go-api.git/src/database"
	"github.com/mappichat/go-api.git/src/utils"
	"github.com/uber/h3-go/v3"
)

func HandleVotes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID string `json:"post_id" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		votes := []database.Vote{}
		if err := database.Sqldb.Select(
			&votes,
			"SELECT * FROM votes WHERE post_id=$1",
			payload.PostID,
		); err != nil {
			return err
		}

		return c.JSON(votes)
	})

	router.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID     string  `json:"post_id" validate:"required"`
			VoteWeight float32 `json:"vote_weight" validate:"required"`
			Level      int8    `json:"level"`
			Latitude   float64 `json:"latitude" validate:"required"`
			Longitude  float64 `json:"longitude" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		dest := struct {
			Latitude  float64 `db:"latitude"`
			Longitude float64 `db:"longitude"`
		}{}
		if err := database.Sqldb.Get(
			&dest,
			"SELECT latitude, longitude FROM posts WHERE id=$1 and account_id=$2",
			payload.PostID, c.Locals("account_id").(string),
		); err != nil {
			return err
		}

		distance := h3.DistanceBetween(
			h3.FromGeo(h3.GeoCoord{Latitude: payload.Latitude, Longitude: payload.Longitude}, utils.Env.MAX_RESOLUTION),
			h3.FromGeo(h3.GeoCoord{Latitude: dest.Latitude, Longitude: dest.Longitude}, utils.Env.MAX_RESOLUTION),
		)

		weight := math.Pow(utils.Env.VOTE_DISTANCE_MULTIPLIER, float64(distance))
		if payload.VoteWeight < 0 {
			weight = -weight
		}

		newVote := database.Vote{
			PostID:     payload.PostID,
			AccountId:  c.Locals("account_id").(string),
			VoteWeight: weight,
			Level:      payload.Level,
			Latitude:   payload.Latitude,
			Longitude:  payload.Longitude,
			TimeStamp:  time.Now().Round(time.Microsecond),
		}

		if _, err := database.Sqldb.NamedExec(
			"INSERT INTO votes (post_id, account_id, vote_weight, vote_level, latitude, longitude, time_stamp) VALUES (:post_id,:account_id,:vote_weight,:vote_level,:latitude,:longitude,:time_stamp)",
			newVote,
		); err != nil {
			return err
		}

		return c.JSON(payload)
	})

	router.Patch("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		// This endpoint just reverses the vote weight
		payload := struct {
			PostID string `json:"post_id" validate:"required"`
			Level  int8   `json:"level"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		if _, err := database.Sqldb.Exec(
			"UPDATE votes SET vote_weight=-vote_weight WHERE post_id=$1 AND account_id=$2 AND vote_level=$3",
			payload.PostID, c.Locals("account_id").(string), payload.Level,
		); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	router.Delete("/", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			PostID string `json:"post_id" validate:"required" db:"post_id"`
			Level  int8   `json:"level" db:"vote_level"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		if _, err := database.Sqldb.Exec(
			"DELETE FROM votes WHERE post_id=$1 AND account_id=$2 AND vote_level=$3",
			payload.PostID, c.Locals("account_id").(string), payload.Level,
		); err != nil {
			return err
		}

		return c.SendStatus(200)
	})
}
