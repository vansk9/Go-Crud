package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"

	"github.com/gofiber/fiber/v2"
)



func SetupProductRoutes(app *fiber.App) {
    productGroup := app.Group("/products")

    productGroup.Get("/", controllers.GetProducts)
    productGroup.Get("/:id", controllers.GetProduct)

    productGroup.Post("/", middleware.AdminOnly, controllers.CreateProduct)
    productGroup.Put("/:id", middleware.AdminOnly, controllers.UpdateProduct)
    productGroup.Delete("/:id", middleware.AdminOnly, controllers.DeleteProduct)
}

