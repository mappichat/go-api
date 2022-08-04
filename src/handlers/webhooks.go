package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
	"github.com/mappichat/go-api.git/src/utils"
)

func HandleWebhooks(webhooks fiber.Router) {
	auth0 := webhooks.Group("/auth0")

	auth0.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.Env.AUTH0_EXTENSION_SECRET),
	}))

	auth0.Post("/post-user-registration", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")
		c.Accepts("application/json")

		payload := struct {
			Username string `json:"username" validate:"required"`
			Email    string `json:"email" validate:"required"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		if err := validate.Struct(payload); err != nil {
			return err
		}

		accountID := uuid.NewString()
		userHandle := fmt.Sprintf("@%s", payload.Username)

		if _, err := database.Sqldb.Exec(
			"INSERT INTO accounts (id, email, user_handle) VALUES ($1,$2,$3)",
			accountID, payload.Email, userHandle,
		); err != nil {
			return err
		}

		return c.JSON(map[string]interface{}{"account_id": accountID, "user_handle": userHandle})
	})
}
