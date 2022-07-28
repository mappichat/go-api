package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var PORT = os.Getenv("PORT")

func main() {
	utils.configure()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
