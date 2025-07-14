package model

type Product struct {
	ID          uint    `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Quantity    int     `db:"quantity" json:"quantity"`
	Price       float64 `db:"price" json:"price"`
	Color       string  `db:"color" json:"color"`
	Size        string  `db:"size" json:"size"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	UpdatedAt   string  `db:"updated_at" json:"updated_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at,omitempty"` // Gunakan pointer untuk nullable field
}
