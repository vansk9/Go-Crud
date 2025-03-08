package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Register & Login
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	

	// Protected Route (Butuh Login)
	api.Get("/protected", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Anda berhasil mengakses protected route"} )
	})
}
