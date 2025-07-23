package model

import (
	"time"
)

type User struct {
	ID          int64     `db:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	PhoneNumber string    `db:"phone_number"`
	DateOfBirth time.Time `db:"date_of_birth"` // Bisa juga time.Time
	Role        int       `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
