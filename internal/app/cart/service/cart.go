package service

import (
	"context"
	"errors"
	"fmt"
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
	GetCartTotal(ctx context.Context, userID uint) (float64, error)
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
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}

	slog.Info("Cart items fetched successfully", "user_id", userID, "count", len(items))
	return items, nil
}

func (s *cartService) GetByID(ctx context.Context, id uint) (*model.CartItem, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch cart item", "cart_id", id, "error", err)
		return nil, fmt.Errorf("failed to fetch cart item: %w", err)
	}
	return item, nil
}

// calculateTotalPrice menghitung total harga berdasarkan harga satuan dan quantity
func (s *cartService) calculateTotalPrice(unitPrice float64, quantity int) float64 {
	return unitPrice * float64(quantity)
}

// validateCartInput memvalidasi input cart item
func (s *cartService) validateCartInput(input *dto.CartItemRequest) error {
	if input.Quantity <= 0 {
		return errors.New("quantity harus lebih dari 0")
	}
	if input.ProductID == 0 {
		return errors.New("product ID tidak boleh kosong")
	}
	return nil
}

func (s *cartService) Create(ctx context.Context, userID uint, input *dto.CartItemRequest) (*model.CartItem, error) {
	// Validasi input
	if err := s.validateCartInput(input); err != nil {
		slog.Error("Invalid cart input", "user_id", userID, "error", err)
		return nil, err
	}

	// Cek apakah produk ada
	product, err := s.productRepo.GetProductsByID(ctx, input.ProductID)
	if err != nil {
		slog.Error("Product not found", "product_id", input.ProductID, "error", err)
		return nil, errors.New("produk tidak ditemukan")
	}

	// Cek apakah produk sudah ada di cart dengan size dan color yang sama
	existingItems, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Failed to check existing cart items", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to check existing cart items: %w", err)
	}

	// Jika item sudah ada, update quantity-nya
	for _, existingItem := range existingItems {
		if existingItem.ProductID == input.ProductID &&
			existingItem.Color == input.Color &&
			existingItem.Size == input.Size {

			newQuantity := existingItem.Quantity + input.Quantity
			updateInput := &dto.CartItemRequest{
				ProductID: input.ProductID,
				Quantity:  newQuantity,
				Color:     input.Color,
				Size:      input.Size,
			}

			slog.Info("Updating existing cart item", "cart_id", existingItem.ID, "new_quantity", newQuantity)
			return s.Update(ctx, existingItem.ID, updateInput)
		}
	}

	// Buat item baru jika belum ada
	totalPrice := s.calculateTotalPrice(product.Price, input.Quantity)

	item := &model.CartItem{
		UserID:    userID,
		ProductID: product.ID,
		Name:      product.Name,
		Quantity:  input.Quantity,
		Price:     totalPrice, // Total harga = harga satuan * quantity
		Color:     input.Color,
		Size:      input.Size,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		slog.Error("Failed to create cart item", "user_id", userID, "product_id", input.ProductID, "error", err)
		return nil, fmt.Errorf("failed to create cart item: %w", err)
	}

	slog.Info("Cart item created successfully", "user_id", userID, "product_id", input.ProductID, "quantity", input.Quantity, "total_price", totalPrice)
	return item, nil
}

func (s *cartService) CreateMany(ctx context.Context, userID uint, inputs []dto.CartItemRequest) ([]model.CartItem, error) {
	if len(inputs) == 0 {
		return []model.CartItem{}, nil
	}

	var result []model.CartItem

	// Bungkus dalam transaction jika repository mendukung
	for i, input := range inputs {
		// Validasi setiap input
		if err := s.validateCartInput(&input); err != nil {
			slog.Error("Invalid cart input in bulk create", "index", i, "error", err)
			return nil, fmt.Errorf("invalid input at index %d: %w", i, err)
		}

		product, err := s.productRepo.GetProductsByID(ctx, input.ProductID)
		if err != nil {
			slog.Error("Product not found for CreateMany", "product_id", input.ProductID, "error", err)
			return nil, fmt.Errorf("produk tidak ditemukan untuk item ke-%d", i+1)
		}

		totalPrice := s.calculateTotalPrice(product.Price, input.Quantity)

		item := model.CartItem{
			UserID:    userID,
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  input.Quantity,
			Price:     totalPrice,
			Color:     input.Color,
			Size:      input.Size,
		}

		if err := s.repo.Create(ctx, &item); err != nil {
			slog.Error("Failed to create cart item (bulk)", "product_id", input.ProductID, "index", i, "error", err)
			return nil, fmt.Errorf("failed to create cart item at index %d: %w", i, err)
		}

		result = append(result, item)
	}

	slog.Info("Bulk cart items created successfully", "user_id", userID, "count", len(result))
	return result, nil
}

func (s *cartService) Update(ctx context.Context, id uint, input *dto.CartItemRequest) (*model.CartItem, error) {
	// Validasi input
	if err := s.validateCartInput(input); err != nil {
		slog.Error("Invalid cart update input", "cart_id", id, "error", err)
		return nil, err
	}

	// Cek apakah cart item ada
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("Cart item not found", "cart_id", id, "error", err)
		return nil, fmt.Errorf("cart item tidak ditemukan: %w", err)
	}

	// Cek apakah produk ada
	product, err := s.productRepo.GetProductsByID(ctx, input.ProductID)
	if err != nil {
		slog.Error("Product not found for update", "product_id", input.ProductID, "error", err)
		return nil, errors.New("produk tidak ditemukan")
	}

	// Update item dengan harga yang dihitung ulang
	totalPrice := s.calculateTotalPrice(product.Price, input.Quantity)

	item.ProductID = product.ID
	item.Name = product.Name
	item.Quantity = input.Quantity
	item.Price = totalPrice
	item.Color = input.Color
	item.Size = input.Size

	if err := s.repo.Update(ctx, item); err != nil {
		slog.Error("Failed to update cart item", "cart_id", id, "error", err)
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	slog.Info("Cart item updated successfully", "cart_id", id, "quantity", input.Quantity, "total_price", totalPrice)
	return item, nil
}

func (s *cartService) Delete(ctx context.Context, id uint) error {
	// Cek apakah item ada sebelum dihapus
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("Cart item not found for deletion", "cart_id", id, "error", err)
		return fmt.Errorf("cart item tidak ditemukan: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete cart item", "cart_id", id, "error", err)
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	slog.Info("Cart item deleted successfully", "cart_id", id)
	return nil
}

func (s *cartService) DeleteMany(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	if err := s.repo.DeleteMany(ctx, ids); err != nil {
		slog.Error("Failed to delete cart items", "count", len(ids), "error", err)
		return fmt.Errorf("failed to delete cart items: %w", err)
	}

	slog.Info("Cart items deleted successfully", "count", len(ids))
	return nil
}

// GetCartTotal menghitung total harga semua item di cart user
func (s *cartService) GetCartTotal(ctx context.Context, userID uint) (float64, error) {
	items, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Price
	}

	return total, nil
}
