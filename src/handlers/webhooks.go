package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/google/uuid"
	"github.com/mappichat/go-api.git/src/database"
	"github.com/mappichat/go-api.git/src/utils"
)

func HandleWebhooks(webhooks fiber.Router) {
	auth0 := webhooks.Group("/auth0")

	auth0.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.Env.AUTH_WEBHOOK_JWT_SECRET),
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

		email := []string{}
		if err := database.Sqldb.Select(
			&email,
			"SELECT email FROM accounts WHERE email=$1",
			payload.Email,
		); err != nil {
			return err
		}

		if len(email) == 0 {
			accountID := uuid.NewString()
			userHandle := fmt.Sprintf("@%s", strings.ToLower(payload.Username))

			handles := []string{}
			if err := database.Sqldb.Select(
				&handles,
				"SELECT user_handle FROM accounts WHERE user_handle LIKE $1",
				userHandle+"%",
			); err != nil {
				return err
			}

			if len(handles) > 0 {
				handleSet := map[string]bool{}
				for _, handle := range handles {
					handleSet[handle] = true
				}
				i := 1
				for i <= 9999999 {
					if _, ok := handleSet[userHandle+strconv.Itoa(i)]; !ok {
						userHandle = userHandle + strconv.Itoa(i)
						break
					}
					i++
				}
				if i == 9999999 {
					return errors.New("too many usernames with prefix " + userHandle)
				}
			}

			if _, err := database.Sqldb.Exec(
				"INSERT INTO accounts (id, email, user_handle) VALUES ($1,$2,$3)",
				accountID, payload.Email, userHandle,
			); err != nil {
				return err
			}

			return c.JSON(map[string]interface{}{"account_id": accountID, "user_handle": userHandle})
		} else {
			return c.JSON("account already exists")
		}

	})
}
