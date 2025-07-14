package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"go-fiber-api/internal/app/user/service"
	"go-fiber-api/internal/shared/dto"
	"go-fiber-api/internal/shared/types"
	utils "go-fiber-api/utils/helper"
	"go-fiber-api/utils/web"

	"github.com/gorilla/schema"
)

type user struct {
	userService service.User
	decoder     *schema.Decoder
}

func NewUserController(mux *http.ServeMux, userService service.User) {
	u := &user{
		userService: userService,
		decoder:     schema.NewDecoder(),
	}
	u.decoder.IgnoreUnknownKeys(true)

	mux.HandleFunc("POST /v1/users/register", u.register)
	mux.HandleFunc("POST /v1/users/login", u.login)

}

func (u *user) register(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	slog.Info("Raw Body", "body", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "err", err.Error())
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", web.ErrValidation))
		return
	}

	// Validasi input dengan validator
	if err := web.Validator().Struct(&req); err != nil {
		web.Err(w, err)
		return
	}

	// Validasi manual (jaga-jaga)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "username, email, and password are required", web.ErrValidation))
		return
	}

	if req.PhoneNumber == "" || !utils.IsNumeric(req.PhoneNumber) {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "phone number is required and must be numeric", web.ErrValidation))
		return
	}

	if req.DateOfBirth == "" {
		slog.Info("User registration failed - Date of birth not provided")
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "date of birth is required", web.ErrValidation))
		return
	}

	// Proteksi jika user ingin membuat akun dengan role ADMIN
	if req.Role == int(types.RoleAdmin) {
		claims, err := utils.GetClaims(r.Header.Get("Authorization"))

		if err != nil {
			slog.Error("Unauthorized: invalid token format", "error", err.Error())
			web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized: invalid token format", web.ErrAuthentication))
			return
		}

		if claims.Role != int(types.RoleAdmin) {
			slog.Error("Unauthorized: insufficient permissions", "role", claims.Role)
			web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized: insufficient permissions, cannot register as admin with user role", web.ErrPermission))
			return
		}
	}

	// Default role = user
	if req.Role == 0 {
		req.Role = int(types.RoleUser)
	}

	// Call to service layer
	res, err := u.userService.Register(r.Context(), &req)
	if err != nil {
		web.Err(w, err)
		return
	}

	web.OK(w, http.StatusCreated, res)
}

func (u *user) login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", 1001))
		return
	}

	if err := web.Validator().Struct(&req); err != nil {
		web.Err(w, err)
		return
	}

	res, err := u.userService.Login(r.Context(), &req)
	if err != nil {
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, err.Error(), 1002))
		return
	}

	web.OK(w, http.StatusOK, res)
}
