package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Pin        int    `gorm:"default:0"`              
	Permission string `gorm:"type:varchar(10);default:user"`
}
