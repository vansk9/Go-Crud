package controllers

import (
	"go-fiber-api/database" // Keep middleware import if needed in controller logic
	"go-fiber-api/models"

	"github.com/gofiber/fiber/v2"
)

// Mendapatkan semua item di keranjang berdasarkan user_id
func GetCartItems(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", userID).Find(&cartItems)
	return c.JSON(cartItems)
}

// Mendapatkan satu item di keranjang berdasarkan ID
func GetCartItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var cartItem models.CartItem
	if err := database.DB.First(&cartItem, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Cart item not found"})
	}
	return c.JSON(cartItem)
}

// Menambahkan satu item ke keranjang
func AddCartItem(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var input models.CartItem
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Cek apakah produk sudah ada di keranjang user
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existingItem)

	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Product already in cart", "product_id": input.ProductID})
	}

	// Jika produk belum ada, tambahkan ke keranjang
	input.UserID = userID
	database.DB.Create(&input)
	return c.JSON(fiber.Map{"message": "Item added to cart", "cart_item": input})
}

// Menambahkan beberapa item ke keranjang sekaligus
func AddMultipleCartItems(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var items []models.CartItem
	if err := c.BodyParser(&items); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Cek apakah ada product_id yang sudah ada di keranjang user
	for _, item := range items {
		var existingItem models.CartItem
		result := database.DB.Where("user_id = ? AND product_id = ?", userID, item.ProductID).First(&existingItem)

		if result.RowsAffected > 0 {
			return c.Status(400).JSON(fiber.Map{"message": "Duplicate product in cart", "product_id": item.ProductID})
		}
	}

	// Jika tidak ada duplikasi, tambahkan semua item
	for _, item := range items {
		item.UserID = userID
		database.DB.Create(&item)
	}

	return c.JSON(fiber.Map{"message": "Items added to cart", "items": items})
}

// Mengupdate item di keranjang berdasarkan ID
func UpdateCartItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var cartItem models.CartItem
	if err := database.DB.First(&cartItem, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Cart item not found"})
	}

	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}
	database.DB.Save(&cartItem)
	return c.JSON(cartItem)
}

// Menghapus satu item dari keranjang berdasarkan ID
func DeleteCartItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var cartItem models.CartItem
	if err := database.DB.First(&cartItem, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Cart item not found"})
	}
	database.DB.Delete(&cartItem)
	return c.JSON(fiber.Map{"message": "Cart item deleted"})
}

// Menghapus beberapa item dari keranjang berdasarkan daftar ID
func DeleteMultipleCartItems(c *fiber.Ctx) error {
	var input struct {
		IDs []uint `json:"ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	database.DB.Delete(&models.CartItem{}, input.IDs)
	return c.JSON(fiber.Map{"message": "Selected cart items deleted"})
}