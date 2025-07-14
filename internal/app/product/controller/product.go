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

	mux.HandleFunc("GET /v1/products", p.GetAll)
	mux.HandleFunc("GET /v1/products/{id}", p.GetByID)
	mux.HandleFunc("POST /v1/products", p.Create)
	mux.HandleFunc("PUT /v1/products/{id}", p.Update)
	mux.HandleFunc("DELETE /v1/products/{id}", p.Delete)
}

func (p *product) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := p.productService.GetAll(r.Context())
	if err != nil {
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusOK, products)
}

func (p *product) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	product, err := p.productService.GetByID(r.Context(), uint(id))
	if err != nil {
		web.Err(w, err)
		return
	}
	web.OK(w, http.StatusOK, product)
}

func (p *product) Create(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	slog.Info("Create Product Body", "body", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", web.ErrValidation))
		return
	}

	if err := web.Validator().Struct(&req); err != nil {
		web.Err(w, err)
		return
	}

	product, err := p.productService.Create(r.Context(), &req)
	if err != nil {
		web.Err(w, err)
		return
	}

	web.OK(w, http.StatusCreated, product)
}

func (p *product) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.Err(w, web.NewHTTPError(http.StatusBadRequest, "Invalid request body", web.ErrValidation))
		return
	}

	if err := web.Validator().Struct(&req); err != nil {
		web.Err(w, err)
		return
	}

	updated, err := p.productService.Update(r.Context(), uint(id), &req)
	if err != nil {
		web.Err(w, err)
		return
	}

	web.OK(w, http.StatusOK, updated)
}

func (p *product) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	if err := p.productService.Delete(r.Context(), uint(id)); err != nil {
		web.Err(w, err)
		return
	}

	web.OKNoContent(w, http.StatusOK)
}
