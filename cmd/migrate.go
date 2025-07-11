package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env jika belum production
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "local" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Ambil env variabel
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Validasi
	if host == "" || port == "" || user == "" || name == "" || sslmode == "" {
		log.Fatal("❌ Environment DB config tidak lengkap")
	}

	// Format DB URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, sslmode)

	// Inisialisasi migrasi
	m, err := migrate.New(
		"file://migrations", // folder migrasi kamu
		dbURL,
	)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi migrasi: %v", err)
	}

	// Jalankan migrasi
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("❌ Gagal migrasi: %v", err)
	}

	log.Println("✅ Migrasi database berhasil")
}
