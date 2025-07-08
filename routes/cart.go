package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCartRoutes(app *fiber.App) {
	cartGroup := app.Group("/cart", middleware.AuthMiddleware)

	cartGroup.Get("/", controllers.GetCartItems) // Lihat semua item di keranjang
	cartGroup.Get("/:id", controllers.GetCartItem) // Lihat satu item di keranjang
	cartGroup.Post("/", controllers.AddCartItem) // Tambah satu item ke keranjang
	cartGroup.Post("/bulk", controllers.AddMultipleCartItems) // Tambah beberapa item ke keranjang sekaligus
	cartGroup.Put("/:id", controllers.UpdateCartItem) // Update item di keranjang
	cartGroup.Delete("/:id", controllers.DeleteCartItem) // Hapus satu item dari keranjang
	cartGroup.Delete("/bulk", controllers.DeleteMultipleCartItems) // Hapus beberapa item sekaligus
}
