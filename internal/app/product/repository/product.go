package repository

import (
	"context"
	"go-fiber-api/internal/app/product/model"

	"github.com/jmoiron/sqlx"
)

type Product interface {
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	Create(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, id uint, p *model.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) Product {
	return &productRepo{db}
}

func (r *productRepo) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	query := `SELECT * FROM products ORDER BY id DESC`
	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepo) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM products WHERE id = $1`
	if err := r.db.GetContext(ctx, &product, query, id); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) Create(ctx context.Context, p *model.Product) error {
	query := `
		INSERT INTO products (name, description, quantity, price, color, size)
		VALUES (:name, :description, :quantity, :price, :color, :size)
		RETURNING id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&p.ID)
	}
	return nil
}

func (r *productRepo) Update(ctx context.Context, id uint, p *model.Product) error {
	query := `
		UPDATE products
		SET name = :name, description = :description, quantity = :quantity,
			price = :price, color = :color, size = :size
		WHERE id = :id
	`
	p.ID = id
	_, err := r.db.NamedExecContext(ctx, query, p)
	return err
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
