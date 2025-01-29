package routes

import (
	"news/controllers"
	"news/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	api := app.Group("/v1")

	cat := api.Group("/category", middlewares.Auth)
	cat.Get("/", controllers.ShowCat)
	cat.Get("/:id", controllers.IndexCat)
	cat.Post("/", controllers.CreateCat)
	cat.Put("/:id", controllers.UpdateCat)
	cat.Delete("/:id", controllers.DeleteCat)

	news := api.Group("/news", middlewares.Auth)
	news.Get("/:id", controllers.IndexNews)
	news.Post("/", controllers.CreateNews)
	news.Put("/upload/:id", controllers.UploadNewsImage)
	news.Put("/:id", controllers.UpdateNews)
	news.Delete("/:id", controllers.DeleteNews)

	publicNews := api.Group("/public/news")
	publicNews.Get("/", controllers.ShowNews)

	comment := api.Group("/comment", middlewares.Auth)
	comment.Get("/", controllers.ShowNewsCom)
	comment.Get("/:id", controllers.IndexNewsCom)
	comment.Post("/", controllers.CreateNewsCom)
	comment.Put("/:id", controllers.UpdateNewsCom)
	comment.Delete("/:id", controllers.DeleteNewsCom)

	like := api.Group("/like", middlewares.Auth)
	like.Get("/", controllers.ShowNewsLike)
	like.Get("/:id", controllers.IndexNewsLike)
	like.Post("/", controllers.CreateNewsLike)
	like.Delete("/:id", controllers.DeleteNewsLike)

}
