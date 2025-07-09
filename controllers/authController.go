package controllers

import (
	"go-fiber-api/database"
	"go-fiber-api/middleware"
	"go-fiber-api/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func Register(c *fiber.Ctx) error {
	// Gunakan struct untuk parsing yang lebih aman
	type RegisterInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Validasi input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "All fields are required"})
	}

	// Cek email sudah terdaftar
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email already registered"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to encrypt password"})
	}

	// Buat user baru dengan name
	user := models.User{
		Name:       input.Name, // Tambahkan name
		Email:      input.Email,
		Password:   string(hashedPassword),
		Permission: "user",
		// Pin: 0, // Bisa dihapus jika tidak digunakan
	}

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Handle error duplikat email lebih spesifik
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(400).JSON(fiber.Map{"message": "Email already registered"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Failed to register user"})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login User
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parsing request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	var user models.User
	// Cek apakah email ada di database
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	// Cek apakah password benar
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	// Generate JWT Token dengan Permission
	token, err := middleware.GenerateToken(user.Email, user.Permission)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal generate token"})
	}

	// Return token
	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}

// Get Profile User (Hanya untuk user yang sudah login)
func GetProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	return c.JSON(user)
}
