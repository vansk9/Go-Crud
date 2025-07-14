package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string, permission string) (string, error) {
	claims := &CustomClaims{
		UserID:     userID,
		Email:      email,
		Permission: permission,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// // Middleware untuk validasi token dan mendapatkan user role
// func AuthMiddleware(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})

// 	if err != nil || !token.Valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
// 	}

// 	claims, ok := token.Claims.(*Claims)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Claims"})
// 	}

// 	// Simpan data user ke context
// 	c.Locals("user_email", claims.Email)
// 	c.Locals("user_permission", claims.Permission)

// 	return c.Next()
// }

// // Middleware untuk mengecek apakah user adalah admin
// func AdminMiddleware(c *fiber.Ctx) error {
// 	permission := c.Locals("user_permission")
// 	if permission != "admin" {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses ditolak, hanya admin yang bisa"})
// 	}
// 	return c.Next()
// }
