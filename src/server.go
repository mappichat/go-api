package main

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mappichat/go-api.git/src/database"
	"github.com/mappichat/go-api.git/src/handlers"
	utils "github.com/mappichat/go-api.git/src/utils"
)

func main() {
	utils.ConfigureEnv()
	_, err := database.SqlInitialize(utils.Env.DB_CONNECTION_STRING)
	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})

	webhooks := app.Group("/webhooks")

	handlers.HandleWebhooks(webhooks)

	api := app.Group("/api")

	keyRefreshDuration := 5 * time.Minute
	api.Use(jwtware.New(jwtware.Config{
		SigningMethod:       "RS256",
		KeySetURL:           utils.Env.AUTH_JWKS_URI,
		KeyRefreshInterval:  &keyRefreshDuration,
		KeyRefreshRateLimit: &keyRefreshDuration,
	}))

	api.Use(func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		accountID, ok := claims["account_id"].(string)
		if !ok {
			return errors.New("jwt has no account_id claim")
		}
		c.Locals("account_id", accountID)
		return c.Next()
	})

	handlers.HandlePosts(api.Group("/posts"))
	handlers.HandleReplies(api.Group("/replies"))
	handlers.HandleVotes(api.Group("/votes"))

	log.Fatal(app.Listen(":" + utils.Env.PORT))
}
