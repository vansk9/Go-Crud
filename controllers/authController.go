package controllers

import (
	"go-fiber-api/database"
	"go-fiber-api/middleware"
	"go-fiber-api/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register User
func Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parsing request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengenkripsi password"})
	}

	// Set default permission sebagai "user"
	user := models.User{
		Email:      data["email"],
		Password:   string(hashedPassword),
		Pin:        data["pin"],
		Permission: "user", // Default permission
	}

	// Simpan user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mendaftarkan user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil didaftarkan"})
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
