package service

import (
	"context"
	"errors"
	"go-fiber-api/internal/app/user/model"
	"go-fiber-api/internal/app/user/repository"
	"go-fiber-api/internal/shared/dto"
	"go-fiber-api/utils"
	"strings"
)

type User interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type userService struct {
	repo repository.User
}

func NewUserService(repo repository.User) User {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal hash password")
	}

	user := &model.User{
		Username:    req.Username,
		Email:       strings.ToLower(req.Email),
		Password:    hashed,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: req.DateOfBirth,
		Role:        req.Role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: user.DateOfBirth,
		Role:        user.Role,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("email atau password salah")
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("email atau password salah")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			DateOfBirth: user.DateOfBirth,
			Role:        user.Role,
		},
	}, nil
}
