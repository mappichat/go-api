package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
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

	api := app.Group("/api")
	handlers.HandlePosts(api)
	handlers.HandleReplies(api)

	log.Fatal(app.Listen(":" + utils.Env.PORT))
}
