package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/mappichat/go-api.git/src/database"
	"github.com/mappichat/go-api.git/src/handlers"
	utils "github.com/mappichat/go-api.git/src/utils"
)

func main() {
	utils.ConfigureEnv()

	if err := database.Initialize(strings.Split(utils.Env.DB_HOST, ",")); err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.Env.JWT_SECRET),
	}))

	api := app.Group("/api")
	handlers.HandlePosts(api.Group("/posts"))
	handlers.HandleReplies(api.Group("/replies"))
	handlers.HandleVotes(api.Group("/votes"))

	log.Fatal(app.Listen(":" + utils.Env.PORT))
}
