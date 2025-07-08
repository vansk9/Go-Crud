package main

import (
	"log"

	"go-fiber-api/database"
	"go-fiber-api/models"
	"go-fiber-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file tidak ditemukan, menggunakan default environment.")
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", 
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	database.ConnectDB()
	err = database.DB.AutoMigrate(&models.CartItem{})
	if err != nil {
		log.Fatalf("❌ Gagal migrasi CartItem: %v", err)
	}

	routes.SetupAuthRoutes(app)
	routes.SetupProductRoutes(app)
	routes.SetupCartRoutes(app)
	log.Println("🚀 Fiber running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
