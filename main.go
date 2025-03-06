package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-fiber-api/database"
	"go-fiber-api/routes"
	"go-fiber-api/models"
)

func main() {
	app := fiber.New()

	// Tambahkan CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Izinkan semua origin
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Koneksi Database
	database.ConnectDB()
	database.DB.AutoMigrate(&models.CartItem{})

	// Setup Routes
	routes.SetupAuthRoutes(app)
	routes.SetupProductRoutes(app)
	routes.SetupCartRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
