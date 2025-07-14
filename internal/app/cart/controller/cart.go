package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"go-fiber-api/internal/app/cart/service"
	"go-fiber-api/internal/shared/dto"
	"go-fiber-api/middleware"
	"go-fiber-api/utils/web"
)

type cart struct {
	service service.Cart
}

func NewCartController(mux *http.ServeMux, cartService service.Cart) {
	c := &cart{service: cartService}

	mux.Handle("GET /v1/cart", middleware.AuthMiddleware(http.HandlerFunc(c.GetAll)))
	mux.Handle("GET /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.GetByID)))
	mux.Handle("POST /v1/cart", middleware.AuthMiddleware(http.HandlerFunc(c.Create)))
	mux.Handle("POST /v1/cart/bulk", middleware.AuthMiddleware(http.HandlerFunc(c.CreateMany)))
	mux.Handle("PUT /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.Update)))
	mux.Handle("DELETE /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.Delete)))
	mux.Handle("DELETE /v1/cart", middleware.AuthMiddleware(http.HandlerFunc(c.DeleteMany)))

}

func (c *cart) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := web.GetUserID(r)
	if userID == 0 {
		slog.Warn("Unauthorized GetAll request: user ID not found in context")
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
		return
	}

	slog.Info("Fetching cart items", "user_id", userID)
	items, err := c.service.GetByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to fetch cart items", "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusOK, items)
}

func (c *cart) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("Invalid cart ID format", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID", web.ErrValidation))
		return
	}

	slog.Info("Fetching cart item by ID", "cart_id", id)
	item, err := c.service.GetByID(r.Context(), uint(id))
	if err != nil {
		slog.Error("Failed to fetch cart item", "cart_id", id, "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusOK, item)
}

func (c *cart) Create(w http.ResponseWriter, r *http.Request) {
	userID := web.GetUserID(r)
	if userID == 0 {
		slog.Warn("Unauthorized Create request: user ID not found in context")
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
		return
	}

	var input dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Failed to decode Create request body", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request", web.ErrValidation))
		return
	}

	slog.Info("Creating new cart item", "user_id", userID, "product_id", input.ProductID)
	item, err := c.service.Create(r.Context(), userID, &input)
	if err != nil {
		slog.Error("Failed to create cart item", "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusCreated, item)
}

func (c *cart) CreateMany(w http.ResponseWriter, r *http.Request) {
	userID := web.GetUserID(r)
	if userID == 0 {
		slog.Warn("Unauthorized CreateMany request: user ID not found")
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
		return
	}

	var inputs []dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		slog.Warn("Failed to decode CreateMany request body", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request", web.ErrValidation))
		return
	}

	slog.Info("Creating multiple cart items", "user_id", userID, "count", len(inputs))
	items, err := c.service.CreateMany(r.Context(), userID, inputs)
	if err != nil {
		slog.Error("Failed to create multiple cart items", "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusCreated, items)
}

func (c *cart) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("Invalid cart ID for update", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID", web.ErrValidation))
		return
	}

	var input dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Failed to decode update request", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request", web.ErrValidation))
		return
	}

	slog.Info("Updating cart item", "cart_id", id)
	item, err := c.service.Update(r.Context(), uint(id), &input)
	if err != nil {
		slog.Error("Failed to update cart item", "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusOK, item)
}

func (c *cart) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("Invalid cart ID for deletion", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID", web.ErrValidation))
		return
	}

	slog.Info("Deleting cart item", "cart_id", id)
	if err := c.service.Delete(r.Context(), uint(id)); err != nil {
		slog.Error("Failed to delete cart item", "error", err)
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusNoContent, nil)
}

func (c *cart) DeleteMany(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDs []uint `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("Failed to decode DeleteMany request", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request", web.ErrValidation))
		return
	}

	slog.Info("Deleting multiple cart items", "ids", body.IDs)
	if err := c.service.DeleteMany(r.Context(), body.IDs); err != nil {
		slog.Error("Failed to delete multiple cart items", "error", err)
		web.Err(w, err)
		return
	}
	web.OKNoContent(w, http.StatusOK)
}
