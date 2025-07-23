package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// type RegisterRequest struct {
// 	Username    string    `json:"username" validate:"required"`
// 	Email       string    `json:"email" validate:"required,email"`
// 	Password    string    `json:"password" validate:"required,min=6"`
// 	PhoneNumber string    `json:"phone_number" validate:"required"`
// 	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
// 	Role        int       `json:"role,omitempty"`
// }

type UserResponse struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Role        int       `json:"role"`
}

type JWTClaims struct {
	ID        string `json:"id"`
	Username  string `json:"username,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Role      int    `json:"role"`
	jwt.RegisteredClaims
}
