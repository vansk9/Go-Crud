package service

import (
	"context"
	"errors"
	"go-fiber-api/internal/app/user/model"
	"go-fiber-api/internal/app/user/repository"
	"go-fiber-api/internal/shared/dto"
	"net/http"

	// "go-fiber-api/utils"
	utils "go-fiber-api/utils/jwt"
	"go-fiber-api/utils/web"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
		DateOfBirth: req.DateOfBirth,
		Role:        user.Role,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	slog.Info("Login attempt", "email", req.Email)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.Warn("Gagal menemukan user saat login", "email", req.Email, "error", err)
		return nil, errors.New("email atau password salah")
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		slog.Info("Login failed - wrong password", "userID", user.ID)
		return nil, web.NewHTTPError(http.StatusUnauthorized, "Password Incorrect", web.ErrPasswordIncorrect)
	}

	userID := uint(user.ID)

	token, err := utils.GenerateJWT(userID, user.Role)
	if err != nil {
		slog.Error("Gagal generate JWT", "user_id", userID, "error", err)
		return nil, errors.New("gagal membuat token")
	}

	slog.Info("Login berhasil", "user_id", userID, "role", user.Role)

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       strings.ToLower(user.Email),
			PhoneNumber: user.PhoneNumber,
			DateOfBirth: user.DateOfBirth,
			Role:        user.Role,
		},
	}, nil
}
