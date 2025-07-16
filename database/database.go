package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	env := os.Getenv("ENV")
	dsn := buildDSN(env)

	// 🔍 Log semua environment variable
	fmt.Println("=== ENV CHECK ===")
	fmt.Println("ENV:", env)
	fmt.Println("DSN:", dsn)
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_SSLMODE:", os.Getenv("DB_SSLMODE"))

	// 🚫 Validasi wajib
	required := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
	}
	for k, v := range required {
		if v == "" {
			log.Fatalf("❌ ENV %s is empty", k)
		}
	}

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database: %v", err)
	}

	log.Println("✅ Berhasil konek ke database")
}

func buildDSN(env string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}
