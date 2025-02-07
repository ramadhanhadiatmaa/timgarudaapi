package main

import (
	"data/models"
	"data/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	models.ConnectDatabase()

	app := fiber.New()

	app.Use(cors.New())

	routes.Route(app)

	app.Listen(":6411")
}