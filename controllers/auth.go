package controllers

import (
	"go-fiber-api/database"
	"go-fiber-api/models"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	type Request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Semua field wajib diisi")
	}

	hashed, err := utils.HashPassword(body.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: hashed,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email sudah terdaftar")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Email tidak ditemukan")
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Password salah")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
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
