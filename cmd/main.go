package main

import (
	"log"
	"net/http"
	"os"

	"go-fiber-api/database"

	// User
	userController "go-fiber-api/internal/app/user/controller"
	userRepo "go-fiber-api/internal/app/user/repository"
	userService "go-fiber-api/internal/app/user/service"

	// Cart
	cartController "go-fiber-api/internal/app/cart/controller"
	cartRepo "go-fiber-api/internal/app/cart/repository"
	cartService "go-fiber-api/internal/app/cart/service"

	// Product
	productController "go-fiber-api/internal/app/product/controller"
	productRepo "go-fiber-api/internal/app/product/repository"
	productService "go-fiber-api/internal/app/product/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file tidak ditemukan, menggunakan default environment")
	}

	database.ConnectDB()

	// Inisialisasi repository dan service
	userRepo := userRepo.NewUserRepository(database.DB)
	userService := userService.NewUserService(userRepo)

	productRepo := productRepo.NewProductRepository(database.DB)
	productService := productService.NewProductService(productRepo)

	cartRepo := cartRepo.NewCartRepository(database.DB)
	cartService := cartService.NewCartService(cartRepo, productRepo) // ✅ Inject productRepo

	// Setup HTTP multiplexer
	mux := http.NewServeMux()

	// Inisialisasi controller
	userController.NewUserController(mux, userService)
	productController.NewProductController(mux, productService)
	cartController.NewCartController(mux, cartService)

	// Jalankan server
	port := os.Getenv("PORT")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

	log.Printf("🚀 Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
