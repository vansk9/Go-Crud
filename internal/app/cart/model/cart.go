package model

type CartItem struct {
	ID        uint    `db:"id" json:"id"`
	UserID    uint    `db:"user_id" json:"user_id"`
	ProductID uint    `db:"product_id" json:"product_id"`
	Name      string  `db:"name" json:"name"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
	Color     string  `db:"color" json:"color"`
	Size      string  `db:"size" json:"size"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt string  `db:"updated_at" json:"updated_at"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at,omitempty"`
}
