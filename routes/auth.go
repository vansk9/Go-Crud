package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/protected", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Anda berhasil mengakses protected route"} )
	})
}
