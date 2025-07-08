// package middleware

// import (
//     "github.com/gofiber/fiber/v2"
//     "github.com/golang-jwt/jwt/v4"
//     "os"
//     "strings"
// )

// // CustomClaims untuk JWT
// type CustomClaims struct {
//     Email      string `json:"email"`
//     Permission string `json:"permission"`
//     jwt.StandardClaims
// }

// func AdminOnly(c *fiber.Ctx) error {
//     authHeader := c.Get("Authorization")
//     if authHeader == "" {
//         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
//     }

//     tokenString := strings.TrimPrefix(authHeader, "Bearer ")

//     token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//         return []byte(os.Getenv("JWT_SECRET")), nil
//     })

//     if err != nil || !token.Valid {
//         return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
//     }

//     claims, ok := token.Claims.(*CustomClaims)
//     if !ok || claims.Permission != "admin" {
//         return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Anda bukan admin"})
//     }

//     return c.Next()
// }

package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
	jwt.StandardClaims
}
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid claims"})
	}
	c.Locals("user_id", claims.UserID)

	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims.Permission != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Anda bukan admin"})
	}

	return c.Next()
}
