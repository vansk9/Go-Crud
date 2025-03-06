package database

import (
    "go-fiber-api/models"
)

func Migrate() {
    DB.AutoMigrate(&models.User{}) // Ini akan otomatis menambahkan kolom baru jika belum ada
    DB.AutoMigrate(&models.CartItem{})
}
