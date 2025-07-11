package repository

import (
	"context"
	"go-fiber-api/internal/app/user/model"

	"gorm.io/gorm"
)

type User interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) User {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
