package model

type Product struct {
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Quantity    int     `db:"quantity" json:"quantity"`
	Price       float64 `db:"price" json:"price"`
	Color       string  `db:"color" json:"color"`
	Size        string  `db:"size" json:"size"`
}
