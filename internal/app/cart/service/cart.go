package service

import (
	"context"
	"errors"
	"go-fiber-api/internal/app/cart/model"
	cartRepo "go-fiber-api/internal/app/cart/repository"
	productRepo "go-fiber-api/internal/app/product/repository"
	"go-fiber-api/internal/shared/dto"
	"log/slog"
)

type Cart interface {
	GetByUserID(ctx context.Context, userID uint) ([]model.CartItem, error)
	GetByID(ctx context.Context, id uint) (*model.CartItem, error)
	Create(ctx context.Context, userID uint, input *dto.CartItemRequest) (*model.CartItem, error)
	CreateMany(ctx context.Context, userID uint, inputs []dto.CartItemRequest) ([]model.CartItem, error)
	Update(ctx context.Context, id uint, input *dto.CartItemRequest) (*model.CartItem, error)
	Delete(ctx context.Context, id uint) error
	DeleteMany(ctx context.Context, ids []uint) error
}

type cartService struct {
	repo        cartRepo.Cart
	productRepo productRepo.Product
}

func NewCartService(cartRepo cartRepo.Cart, productRepo productRepo.Product) Cart {
	return &cartService{
		repo:        cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) GetByUserID(ctx context.Context, userID uint) ([]model.CartItem, error) {
	slog.Info("Fetching cart items", "user_id", userID)

	items, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Failed to fetch cart items", "user_id", userID, "error", err)
		return nil, err
	}

	slog.Info("Cart items fetched", "count", len(items))
	return items, nil
}

func (s *cartService) GetByID(ctx context.Context, id uint) (*model.CartItem, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *cartService) Create(ctx context.Context, userID uint, input *dto.CartItemRequest) (*model.CartItem, error) {
	product, err := s.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		slog.Error("Product not found", "product_id", input.ProductID, "error", err)
		return nil, errors.New("produk tidak ditemukan")
	}

	item := &model.CartItem{
		UserID:    userID,
		ProductID: product.ID,
		Name:      product.Name,
		Quantity:  input.Quantity,
		Price:     product.Price * float64(input.Quantity), // Total harga = harga satuan * quantity
		Color:     input.Color,
		Size:      input.Size,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		slog.Error("Failed to create cart item", "user_id", userID, "error", err)
		return nil, err
	}

	return item, nil
}

func (s *cartService) CreateMany(ctx context.Context, userID uint, inputs []dto.CartItemRequest) ([]model.CartItem, error) {
	var result []model.CartItem

	for _, input := range inputs {
		product, err := s.productRepo.GetByID(ctx, input.ProductID)
		if err != nil {
			slog.Error("Product not found for CreateMany", "product_id", input.ProductID, "error", err)
			return nil, errors.New("produk tidak ditemukan")
		}

		item := model.CartItem{
			UserID:    userID,
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  input.Quantity,
			Price:     product.Price * float64(input.Quantity),
			Color:     input.Color,
			Size:      input.Size,
		}

		if err := s.repo.Create(ctx, &item); err != nil {
			slog.Error("Failed to create cart item (bulk)", "product_id", input.ProductID, "error", err)
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *cartService) Update(ctx context.Context, id uint, input *dto.CartItemRequest) (*model.CartItem, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		slog.Error("Product not found for update", "product_id", input.ProductID, "error", err)
		return nil, errors.New("produk tidak ditemukan")
	}

	item.ProductID = product.ID
	item.Name = product.Name
	item.Quantity = input.Quantity
	item.Price = product.Price * float64(input.Quantity)
	item.Color = input.Color
	item.Size = input.Size

	if err := s.repo.Update(ctx, item); err != nil {
		slog.Error("Failed to update cart item", "cart_id", id, "error", err)
		return nil, err
	}

	return item, nil
}

func (s *cartService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *cartService) DeleteMany(ctx context.Context, ids []uint) error {
	return s.repo.DeleteMany(ctx, ids)
}
