package repository

import (
	"context"
	"go-fiber-api/internal/app/product/model"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Product interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetProductsByID(ctx context.Context, id uint) (*model.Product, error)
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

func (r *productRepo) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	query := `SELECT * FROM products ORDER BY id DESC`

	slog.Info("Executing query GetAll", "query", query)
	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		slog.Error("Failed to get all products", "error", err)
		return nil, err
	}
	return products, nil
}

func (r *productRepo) GetProductsByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM products WHERE id = $1`

	slog.Info("Executing query GetByID", "query", query, "id", id)
	if err := r.db.GetContext(ctx, &product, query, id); err != nil {
		slog.Error("Failed to get product by ID", "id", id, "error", err)
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

	slog.Info("Executing query Create", "query", query, "product", p)
	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		slog.Error("Failed to create product", "error", err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.ID); err != nil {
			slog.Error("Failed to scan created product ID", "error", err)
			return err
		}
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

	slog.Info("Executing query Update", "query", query, "id", id, "product", p)
	_, err := r.db.NamedExecContext(ctx, query, p)
	if err != nil {
		slog.Error("Failed to update product", "id", id, "error", err)
	}
	return err
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM products WHERE id = $1`

	slog.Info("Executing query Delete", "query", query, "id", id)
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete product", "id", id, "error", err)
	}
	return err
}
