package routes

import (
	"schedule/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1" /* , middlewares.Auth */)

	garuda := api.Group("/garuda")
	garuda.Post("/", controllers.Create)
	garuda.Get("/:id", controllers.Index)
	garuda.Put("/:id", controllers.Update)
	garuda.Delete("/:id", controllers.Delete)

	team := api.Group("/team")
	team.Post("/", controllers.CreateTeam)
	team.Get("/:id", controllers.IndexTeam)
	team.Put("/:id", controllers.UpdateTeam)
	team.Delete("/:id", controllers.DeleteTeam)
}
