package main

import (
	"log"
	"net/http"

	"go-fiber-api/database"
	userController "go-fiber-api/internal/app/user/controller"
	userRepo "go-fiber-api/internal/app/user/repository"
	userService "go-fiber-api/internal/app/user/service"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDB()

	repo := userRepo.NewUserRepository(database.DB)
	service := userService.NewUserService(repo)

	mux := http.NewServeMux()
	userController.NewUserController(mux, service)

	log.Println("🚀 Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
