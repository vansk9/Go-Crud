package database

import (
	"go-fiber-api/internal/app/user/model"
)

func Migrate() {
	DB.AutoMigrate(&model.User{})
	// DB.AutoMigrate(&model.CartItem{})
	// DB.AutoMigrate(&model.Product{})
}
