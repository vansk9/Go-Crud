package repository

import (
	"context"
	"go-fiber-api/internal/app/cart/model"

	"github.com/jmoiron/sqlx"
)

type Cart interface {
	FindByUserID(ctx context.Context, userID uint) ([]model.CartItem, error)
	FindByID(ctx context.Context, id uint) (*model.CartItem, error)
	FindByUserAndProductID(ctx context.Context, userID, productID uint) (*model.CartItem, error)
	Create(ctx context.Context, item *model.CartItem) error
	Update(ctx context.Context, item *model.CartItem) error
	Delete(ctx context.Context, id uint) error
	DeleteMany(ctx context.Context, ids []uint) error
	FindByUserProductColorSize(ctx context.Context, userID, productID uint, color, size string) (*model.CartItem, error)
}

type cartRepo struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) Cart {
	return &cartRepo{db: db}
}

func (r *cartRepo) FindByUserID(ctx context.Context, userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	query := `SELECT * FROM cart_items WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &items, query, userID)
	return items, err
}

func (r *cartRepo) FindByID(ctx context.Context, id uint) (*model.CartItem, error) {
	var item model.CartItem
	query := `SELECT * FROM cart_items WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *cartRepo) FindByUserAndProductID(ctx context.Context, userID, productID uint) (*model.CartItem, error) {
	var item model.CartItem
	query := `SELECT * FROM cart_items WHERE user_id = $1 AND product_id = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &item, query, userID, productID)
	return &item, err
}

func (r *cartRepo) Create(ctx context.Context, item *model.CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, product_id, name, quantity, price, color, size)
		VALUES (:user_id, :product_id, :name, :quantity, :price, :color, :size)
		RETURNING id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&item.ID)
	}
	return nil
}

func (r *cartRepo) Update(ctx context.Context, item *model.CartItem) error {
	query := `
		UPDATE cart_items
		SET product_id = :product_id, name = :name, quantity = :quantity,
		    price = :price, color = :color, size = :size
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *cartRepo) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE id = $1`, id)
	return err
}

func (r *cartRepo) DeleteMany(ctx context.Context, ids []uint) error {
	query := `DELETE FROM cart_items WHERE id = ANY($1)`
	_, err := r.db.ExecContext(ctx, query, ids)
	return err
}

func (r *cartRepo) FindByUserProductColorSize(ctx context.Context, userID, productID uint, color, size string) (*model.CartItem, error) {
	var item model.CartItem
	query := `
		SELECT * FROM cart_items
		WHERE user_id = $1 AND product_id = $2 AND color = $3 AND size = $4
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &item, query, userID, productID, color, size)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
