package utils

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
	jwt.RegisteredClaims
}

func GetClaims(authHeader string) (*JWTClaims, error) {
	if authHeader == "" {
		return nil, errors.New("Authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("Invalid Authorization header format")
	}

	tokenStr := parts[1]
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}

// IsNumeric memeriksa apakah string hanya terdiri dari angka
func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// IsValidEmail memeriksa format email valid dengan regex sederhana
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
