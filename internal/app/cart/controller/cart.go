package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"go-fiber-api/internal/app/cart/model"
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

	// Cart routes
	mux.Handle("GET /v1/cart", middleware.AuthMiddleware(http.HandlerFunc(c.GetAll)))
	mux.Handle("GET /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.GetByID)))
	mux.Handle("GET /v1/cart/total", middleware.AuthMiddleware(http.HandlerFunc(c.GetCartTotal)))
	mux.Handle("POST /v1/cart", middleware.AuthMiddleware(http.HandlerFunc(c.Create)))
	mux.Handle("POST /v1/cart/bulk", middleware.AuthMiddleware(http.HandlerFunc(c.CreateMany)))
	mux.Handle("PUT /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.Update)))
	mux.Handle("DELETE /v1/cart/", middleware.AuthMiddleware(http.HandlerFunc(c.Delete)))
	mux.Handle("DELETE /v1/cart/bulk", middleware.AuthMiddleware(http.HandlerFunc(c.DeleteMany)))
}

// validateUserAccess memvalidasi apakah user memiliki akses ke cart item tertentu
func (c *cart) validateUserAccess(r *http.Request, cartID uint) error {
	userID := web.GetUserID(r)
	if userID == 0 {
		return web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication)
	}

	// Cek apakah cart item milik user yang sedang login
	item, err := c.service.GetByID(r.Context(), cartID)
	if err != nil {
		return web.NewHTTPError(http.StatusNotFound, "Cart item not found", web.ErrNotFound)
	}

	if item.UserID != userID {
		return web.NewHTTPError(http.StatusForbidden, "Access denied to this cart item", web.ErrForbidden)
	}

	return nil
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
		slog.Error("Failed to fetch cart items", "user_id", userID, "error", err)
		web.Err(w, err)
		return
	}

	// Return empty array instead of null for better frontend handling
	if items == nil {
		items = []model.CartItem{}
	}

	web.OK(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"count": len(items),
	})
}

func (c *cart) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	if idStr == "" || idStr == "total" { // Hindari konflik dengan endpoint lain
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Cart ID is required", web.ErrValidation))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("Invalid cart ID format", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID format", web.ErrValidation))
		return
	}

	cartID := uint(id)

	// Validasi akses user
	if err := c.validateUserAccess(r, cartID); err != nil {
		web.Err(w, err)
		return
	}

	slog.Info("Fetching cart item by ID", "cart_id", cartID)
	item, err := c.service.GetByID(r.Context(), cartID)
	if err != nil {
		slog.Error("Failed to fetch cart item", "cart_id", cartID, "error", err)
		web.Err(w, err)
		return
	}

	web.OK(w, http.StatusOK, item)
}

func (c *cart) GetCartTotal(w http.ResponseWriter, r *http.Request) {
	userID := web.GetUserID(r)
	if userID == 0 {
		slog.Warn("Unauthorized GetCartTotal request: user ID not found")
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
		return
	}

	slog.Info("Calculating cart total", "user_id", userID)
	total, err := c.service.GetCartTotal(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to calculate cart total", "user_id", userID, "error", err)
		web.Err(w, err)
		return
	}

	web.OK(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"total":   total,
	})
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
		slog.Warn("Failed to decode Create request body", "user_id", userID, "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid JSON format", web.ErrValidation))
		return
	}

	// Basic validation
	if input.ProductID == 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Product ID is required", web.ErrValidation))
		return
	}

	if input.Quantity <= 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Quantity must be greater than 0", web.ErrValidation))
		return
	}

	slog.Info("Creating new cart item", "user_id", userID, "product_id", input.ProductID, "quantity", input.Quantity)
	item, err := c.service.Create(r.Context(), userID, &input)
	if err != nil {
		slog.Error("Failed to create cart item", "user_id", userID, "product_id", input.ProductID, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Cart item created successfully", "user_id", userID, "cart_id", item.ID, "product_id", item.ProductID)
	web.OK(w, http.StatusCreated, map[string]interface{}{
		"message": "Item added to cart successfully",
		"item":    item,
	})
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
		slog.Warn("Failed to decode CreateMany request body", "user_id", userID, "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid JSON format", web.ErrValidation))
		return
	}

	if len(inputs) == 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "At least one item is required", web.ErrValidation))
		return
	}

	if len(inputs) > 50 { // Batasi maksimal 50 item sekaligus
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Maximum 50 items allowed per request", web.ErrValidation))
		return
	}

	// Validasi semua input
	for i, input := range inputs {
		if input.ProductID == 0 {
			web.Err(w, web.NewHTTPError(http.StatusBadRequest,
				fmt.Sprintf("Product ID is required for item %d", i+1), web.ErrValidation))
			return
		}
		if input.Quantity <= 0 {
			web.Err(w, web.NewHTTPError(http.StatusBadRequest,
				fmt.Sprintf("Quantity must be greater than 0 for item %d", i+1), web.ErrValidation))
			return
		}
	}

	slog.Info("Creating multiple cart items", "user_id", userID, "count", len(inputs))
	items, err := c.service.CreateMany(r.Context(), userID, inputs)
	if err != nil {
		slog.Error("Failed to create multiple cart items", "user_id", userID, "count", len(inputs), "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Multiple cart items created successfully", "user_id", userID, "created_count", len(items))
	web.OK(w, http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprintf("%d items added to cart successfully", len(items)),
		"items":   items,
		"count":   len(items),
	})
}

func (c *cart) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	if idStr == "" {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Cart ID is required", web.ErrValidation))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("Invalid cart ID for update", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID format", web.ErrValidation))
		return
	}

	cartID := uint(id)

	// Validasi akses user
	if err := c.validateUserAccess(r, cartID); err != nil {
		web.Err(w, err)
		return
	}

	var input dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Failed to decode update request", "cart_id", cartID, "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid JSON format", web.ErrValidation))
		return
	}

	// Basic validation
	if input.ProductID == 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Product ID is required", web.ErrValidation))
		return
	}

	if input.Quantity <= 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Quantity must be greater than 0", web.ErrValidation))
		return
	}

	slog.Info("Updating cart item", "cart_id", cartID, "product_id", input.ProductID, "quantity", input.Quantity)
	item, err := c.service.Update(r.Context(), cartID, &input)
	if err != nil {
		slog.Error("Failed to update cart item", "cart_id", cartID, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Cart item updated successfully", "cart_id", cartID, "product_id", item.ProductID)
	web.OK(w, http.StatusOK, map[string]interface{}{
		"message": "Cart item updated successfully",
		"item":    item,
	})
}

func (c *cart) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/cart/")
	if idStr == "" {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Cart ID is required", web.ErrValidation))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.Warn("Invalid cart ID for deletion", "cart_id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid cart ID format", web.ErrValidation))
		return
	}

	cartID := uint(id)

	// Validasi akses user
	if err := c.validateUserAccess(r, cartID); err != nil {
		web.Err(w, err)
		return
	}

	slog.Info("Deleting cart item", "cart_id", cartID)
	if err := c.service.Delete(r.Context(), cartID); err != nil {
		slog.Error("Failed to delete cart item", "cart_id", cartID, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Cart item deleted successfully", "cart_id", cartID)
	web.OK(w, http.StatusOK, map[string]interface{}{
		"message": "Cart item deleted successfully",
	})
}

func (c *cart) DeleteMany(w http.ResponseWriter, r *http.Request) {
	userID := web.GetUserID(r)
	if userID == 0 {
		slog.Warn("Unauthorized DeleteMany request: user ID not found")
		web.Err(w, web.NewHTTPError(http.StatusUnauthorized, "Unauthorized", web.ErrAuthentication))
		return
	}

	var body struct {
		IDs []uint `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("Failed to decode DeleteMany request", "user_id", userID, "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid JSON format", web.ErrValidation))
		return
	}

	if len(body.IDs) == 0 {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "At least one cart ID is required", web.ErrValidation))
		return
	}

	if len(body.IDs) > 50 { // Batasi maksimal 50 item sekaligus
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Maximum 50 items allowed per request", web.ErrValidation))
		return
	}

	// Validasi bahwa semua cart items milik user yang sedang login
	for _, cartID := range body.IDs {
		item, err := c.service.GetByID(r.Context(), cartID)
		if err != nil {
			slog.Error("Cart item not found for bulk delete", "cart_id", cartID, "user_id", userID)
			web.Err(w, web.NewHTTPError(http.StatusNotFound,
				fmt.Sprintf("Cart item %d not found", cartID), web.ErrNotFound))
			return
		}

		if item.UserID != userID {
			slog.Warn("Access denied for bulk delete", "cart_id", cartID, "owner_id", item.UserID, "requester_id", userID)
			web.Err(w, web.NewHTTPError(http.StatusForbidden,
				fmt.Sprintf("Access denied to cart item %d", cartID), web.ErrForbidden))
			return
		}
	}

	slog.Info("Deleting multiple cart items", "user_id", userID, "ids", body.IDs, "count", len(body.IDs))
	if err := c.service.DeleteMany(r.Context(), body.IDs); err != nil {
		slog.Error("Failed to delete multiple cart items", "user_id", userID, "ids", body.IDs, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Multiple cart items deleted successfully", "user_id", userID, "count", len(body.IDs))
	web.OK(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%d cart items deleted successfully", len(body.IDs)),
		"count":   len(body.IDs),
	})
}
