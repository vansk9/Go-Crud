package database

import (
	"go-fiber-api/models"
)

func Migrate() {
    DB.AutoMigrate(&models.User{})
    DB.AutoMigrate(&models.CartItem{})
    DB.AutoMigrate(&models.Product{})
}
