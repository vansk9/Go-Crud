package database

import (
	"fmt"
	"go-fiber-api/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslmode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database: %v", err)
	}

	fmt.Println("✅ Berhasil konek ke PostgreSQL")

	fmt.Println("📦 Running migrations...")
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("❌ Gagal migrasi database: %v", err)
	}
	fmt.Println("✅ Migrations completed!")
}
