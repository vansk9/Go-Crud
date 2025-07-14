package dto

// ProductRequest digunakan saat create atau update product
type ProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Color       string  `json:"color" validate:"required"`
	Size        string  `json:"size" validate:"required"`
}

// ProductResponse adalah format response ke client
type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Color       string  `json:"color"`
	Size        string  `json:"size"`
}
