package dto

type CartItemRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	Color     string `json:"color" validate:"required"`
	Size      string `json:"size" validate:"required"`
}

type CartItemResponse struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
}
