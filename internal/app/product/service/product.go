package service

import (
	"context"
	"errors"

	"go-fiber-api/internal/app/product/model"
	"go-fiber-api/internal/app/product/repository"
	"go-fiber-api/internal/shared/dto"
)

type Product interface {
	Create(ctx context.Context, req *dto.ProductRequest) (*dto.ProductResponse, error)
	GetAll(ctx context.Context) ([]*dto.ProductResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
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

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Color:       req.Color,
		Size:        req.Size,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

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

func (s *productService) GetAll(ctx context.Context) ([]*dto.ProductResponse, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
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

	return result, nil
}

func (s *productService) GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

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
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Quantity = req.Quantity
	product.Price = req.Price
	product.Color = req.Color
	product.Size = req.Size

	if err := s.repo.Update(ctx, id, product); err != nil {

		return nil, err
	}

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
	return s.repo.Delete(ctx, id)
}
