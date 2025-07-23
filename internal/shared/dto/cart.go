package dto

import "time"

// CartItemRequest represents the request payload for creating/updating cart items
type CartItemRequest struct {
	ProductID uint   `json:"product_id" validate:"required,min=1"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=999"`
	Color     string `json:"color,omitempty" validate:"omitempty,max=50"`
	Size      string `json:"size,omitempty" validate:"omitempty,max=50"`
}

// CartItemTotal represents individual item totals in cart calculation
type CartItemTotal struct {
	CartItemID  uint    `json:"cart_item_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

// CartTotalResponse represents the complete cart total calculation
type CartTotalResponse struct {
	UserID       uint            `json:"user_id"`
	TotalItems   int             `json:"total_items"`
	TotalAmount  float64         `json:"total_amount"`
	Currency     string          `json:"currency"`
	Items        []CartItemTotal `json:"items"`
	CalculatedAt time.Time       `json:"calculated_at"`
}

// CartSummary provides a quick overview of the cart
type CartSummary struct {
	UserID      uint    `json:"user_id"`
	ItemCount   int     `json:"item_count"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
	IsEmpty     bool    `json:"is_empty"`
}

// CartItemResponse represents the response for cart item operations
type CartItemResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Price       float64   `json:"price,omitempty"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BulkCartRequest represents bulk operations on cart items
type BulkCartRequest struct {
	Items []CartItemRequest `json:"items" validate:"required,min=1,max=50"`
}

// BulkDeleteRequest represents bulk delete operations
type BulkDeleteRequest struct {
	IDs []uint `json:"ids" validate:"required,min=1,max=50"`
}

// CartActionResponse represents responses for cart actions
type CartActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// Pagination for cart items if needed
type CartPaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// CartListResponse represents paginated cart items
type CartListResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCount int                `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
