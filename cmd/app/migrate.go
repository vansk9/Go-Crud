package app

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	Run:   runMigrations,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrations(cmd *cobra.Command, args []string) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || name == "" || sslmode == "" {
		log.Fatal("❌ Environment DB config tidak lengkap")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, sslmode)

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("❌ Gagal migrasi: %v", err)
	}

	log.Println("✅ Migrasi database berhasil")
}
