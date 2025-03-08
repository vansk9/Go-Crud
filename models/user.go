package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Pin int `gorm:"default:user"` 
	Permission string `gorm:"default:user"` // Default jadi user biasa

}
