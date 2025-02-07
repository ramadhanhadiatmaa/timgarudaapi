package routes

import (
	"schedule/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1" /* , middlewares.Auth */)

	schedule := api.Group("/schedule")
	schedule.Post("/", controllers.Create)
	schedule.Get("/", controllers.Show)
	schedule.Get("/:id", controllers.Index)
	schedule.Put("/:id", controllers.Update)
	schedule.Delete("/:id", controllers.Delete) 

	team := api.Group("/team")
	team.Post("/", controllers.CreateTeam)
	team.Get("/", controllers.ShowTeam)
	team.Get("/:id", controllers.IndexTeam)
	team.Put("/:id", controllers.UpdateTeam)
	team.Delete("/:id", controllers.DeleteTeam)
}
