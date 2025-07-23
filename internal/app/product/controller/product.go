package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"go-fiber-api/internal/app/product/service"
	"go-fiber-api/internal/shared/dto"
	"go-fiber-api/internal/shared/types"
	"go-fiber-api/middleware"
	"go-fiber-api/utils/web"

	"github.com/gorilla/schema"
)

type product struct {
	productService service.Product
	decoder        *schema.Decoder
}

func NewProductController(mux *http.ServeMux, productService service.Product) {
	p := &product{
		productService: productService,
		decoder:        schema.NewDecoder(),
	}
	p.decoder.IgnoreUnknownKeys(true)

	mux.HandleFunc("POST /v1/products", middleware.ValidateRole(types.RoleAdmin)(p.Create))
	mux.HandleFunc("GET /v1/products", middleware.ValidateRole(types.RoleAdmin)(p.GetAllProducts))
	mux.HandleFunc("GET /v1/products/{id}", middleware.ValidateRole(types.RoleAdmin)(p.GetProductsByID))
	mux.HandleFunc("PUT /v1/products/{id}", middleware.ValidateRole(types.RoleAdmin)(p.Update))
	mux.HandleFunc("DELETE /v1/products/{id}", middleware.ValidateRole(types.RoleAdmin)(p.Delete))
}

func (p *product) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetAllProducts called")
	products, err := p.productService.GetAllProducts(r.Context())
	if err != nil {
		slog.Error("GetAllProducts failed", "error", err)
		web.Err(w, err)
		return
	}
	slog.Info("GetAllProducts success", "count", len(products))
	web.OK(w, http.StatusOK, products)
}

func (p *product) GetProductsByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Invalid product ID", "id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid product ID", web.ErrValidation))
		return
	}
	slog.Info("GetProductsByID called", "id", id)

	product, err := p.productService.GetProductsByID(r.Context(), uint(id))
	if err != nil {
		slog.Error("GetProductsByID failed", "id", id, "error", err)
		web.Err(w, err)
		return
	}
	slog.Info("GetProductsByID success", "id", id)
	web.OK(w, http.StatusOK, product)
}

func (p *product) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	slog.Info("Create Product Body", "body", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Create failed - invalid JSON", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", web.ErrValidation))
		return
	}

	if err := web.Validator().Struct(&req); err != nil {
		slog.Error("Create failed - validation error", "error", err)
		web.Err(w, err)
		return
	}

	product, err := p.productService.Create(r.Context(), &req)
	if err != nil {
		slog.Error("Create failed", "error", err)
		web.Err(w, err)
		return
	}
	slog.Info("Create success", "product_id", product.ID)
	web.OK(w, http.StatusCreated, product)
}

func (p *product) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Invalid product ID for update", "id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid product ID", web.ErrValidation))
		return
	}
	slog.Info("Update called", "id", id)

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Update failed - invalid JSON", "error", err)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", web.ErrValidation))
		return
	}

	if err := web.Validator().Struct(&req); err != nil {
		slog.Error("Update failed - validation error", "error", err)
		web.Err(w, err)
		return
	}

	updated, err := p.productService.Update(r.Context(), uint(id), &req)
	if err != nil {
		slog.Error("Update failed", "id", id, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Update success", "id", id)
	web.OK(w, http.StatusOK, updated)
}

func (p *product) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Invalid product ID for delete", "id", idStr)
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid product ID", web.ErrValidation))
		return
	}
	slog.Info("Delete called", "id", id)

	if err := p.productService.Delete(r.Context(), uint(id)); err != nil {
		slog.Error("Delete failed", "id", id, "error", err)
		web.Err(w, err)
		return
	}

	slog.Info("Delete success", "id", id)
	web.OKNoContent(w, http.StatusOK)
}
