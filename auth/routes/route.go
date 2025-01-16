package routes

import (
	"auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	user := api.Group("/user")

	user.Post("/login", controllers.Login)
	user.Post("/register", controllers.Register)
	user.Put("/upload/:email", controllers.UploadUserImage)
	user.Put("/:email", controllers.Update)
	user.Delete("/:email", controllers.Delete)
}