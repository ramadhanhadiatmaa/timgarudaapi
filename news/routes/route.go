package routes

import (
	"news/controllers"
	"news/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1", middlewares.Auth)
	api2 := app.Group("/v1")

	cat := api.Group("/category")
	cat.Get("/", controllers.ShowCat)
	cat.Get("/:id", controllers.IndexCat)
	cat.Post("/", controllers.CreateCat)
	cat.Put("/:id", controllers.UpdateCat)
	cat.Delete("/:id", controllers.DeleteCat)

	news := api.Group("/news")
	news.Get("/:id", controllers.IndexNews)
	news.Post("/", controllers.CreateNews)
	news.Put("/upload/:id", controllers.UploadNewsImage)
	news.Put("/:id", controllers.UpdateNews)
	news.Delete("/:id", controllers.DeleteNews)

	news2 := api2.Group("/news")
	news2.Get("/", controllers.ShowNews)

	comment := api.Group("/comment")
	comment.Get("/", controllers.ShowNewsCom)
	comment.Get("/:id", controllers.IndexNewsCom)
	comment.Post("/", controllers.CreateNewsCom)
	comment.Put("/:id", controllers.UpdateNewsCom)
	comment.Delete("/:id", controllers.DeleteNewsCom)

	like := api.Group("/like")
	like.Get("/", controllers.ShowNewsLike)
	like.Get("/:id", controllers.IndexNewsLike)
	like.Post("/", controllers.CreateNewsLike)
	like.Delete("/:id", controllers.DeleteNewsLike)

}
