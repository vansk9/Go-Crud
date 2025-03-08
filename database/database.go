package database

import (
	"fmt"
	"log"
	"os"
	"go-fiber-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Ambil koneksi database dari environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=miftah87 dbname=gol port=5432 sslmode=disable"
	}

	// Koneksi ke PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek ke database:", err)
	}

	fmt.Println("✅ Berhasil koneksi ke PostgreSQL")

	// Migrasi tabel
	fmt.Println("📦 Running migrations...")
	err = DB.AutoMigrate(&models.User{}, &models.CartItem{}, &models.Product{})
	if err != nil {
		log.Fatal("❌ Gagal migrasi database:", err)
	}
	fmt.Println("✅ Migrations completed!")
}


