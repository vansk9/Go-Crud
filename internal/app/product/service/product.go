package service

import (
	"context"
	"errors"
	"log/slog"

	"go-fiber-api/internal/app/product/model"
	"go-fiber-api/internal/app/product/repository"
	"go-fiber-api/internal/shared/dto"
)

type Product interface {
	Create(ctx context.Context, req *dto.ProductRequest) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error)
	GetProductsByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
	Update(ctx context.Context, id uint, req *dto.ProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) Product {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, req *dto.ProductRequest) (*dto.ProductResponse, error) {
	slog.Info("Creating product", "name", req.Name)

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Color:       req.Color,
		Size:        req.Size,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		slog.Error("Failed to create product", "error", err)
		return nil, err
	}

	slog.Info("Product created successfully", "product_id", product.ID)
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Color:       product.Color,
		Size:        product.Size,
	}, nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	slog.Info("Fetching all products")

	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		slog.Error("Failed to fetch products", "error", err)
		return nil, err
	}

	var result []*dto.ProductResponse
	for _, product := range products {
		result = append(result, &dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Quantity:    product.Quantity,
			Price:       product.Price,
			Color:       product.Color,
			Size:        product.Size,
		})
	}

	slog.Info("Fetched products successfully", "count", len(result))
	return result, nil
}

func (s *productService) GetProductsByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	slog.Info("Fetching product by ID", "product_id", id)

	product, err := s.repo.GetProductsByID(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch product by ID", "product_id", id, "error", err)
		return nil, err
	}

	slog.Info("Product found", "product_id", id)
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Color:       product.Color,
		Size:        product.Size,
	}, nil
}

func (s *productService) Update(ctx context.Context, id uint, req *dto.ProductRequest) (*dto.ProductResponse, error) {
	slog.Info("Updating product", "product_id", id)

	product, err := s.repo.GetProductsByID(ctx, id)
	if err != nil {
		slog.Warn("Product not found", "product_id", id)
		return nil, errors.New("product not found")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Quantity = req.Quantity
	product.Price = req.Price
	product.Color = req.Color
	product.Size = req.Size

	if err := s.repo.Update(ctx, id, product); err != nil {
		slog.Error("Failed to update product", "product_id", id, "error", err)
		return nil, err
	}

	slog.Info("Product updated successfully", "product_id", id)
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Color:       product.Color,
		Size:        product.Size,
	}, nil
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	slog.Info("Deleting product", "product_id", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete product", "product_id", id, "error", err)
		return err
	}

	slog.Info("Product deleted successfully", "product_id", id)
	return nil
}
