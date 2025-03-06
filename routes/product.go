package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-api/database"
	"go-fiber-api/models"
	"go-fiber-api/middleware"
)

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)
	return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Product not found"})
	}
	return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}
	database.DB.Create(&product)
	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Product not found"})
	}
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}
	database.DB.Save(&product)
	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Product not found"})
	}
	database.DB.Delete(&product)
	return c.JSON(fiber.Map{"message": "Product deleted"})
}

func SetupProductRoutes(app *fiber.App) {
    productGroup := app.Group("/products")

    productGroup.Get("/", GetProducts)
    productGroup.Get("/:id", GetProduct)

    // Hanya admin yang bisa melakukan operasi berikut
    productGroup.Post("/", middleware.AdminOnly, CreateProduct)
    productGroup.Put("/:id", middleware.AdminOnly, UpdateProduct)
    productGroup.Delete("/:id", middleware.AdminOnly, DeleteProduct)
}

