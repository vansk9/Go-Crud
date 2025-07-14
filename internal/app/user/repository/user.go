package repository

import (
	"context"
	"go-fiber-api/internal/app/user/model"

	"github.com/jmoiron/sqlx"
)

type User interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) User {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, phone_number, date_of_birth, role)
		VALUES (:username, :email, :password, :phone_number, :date_of_birth, :role)
		RETURNING id
	`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return err
		}
	}

	return nil
}
