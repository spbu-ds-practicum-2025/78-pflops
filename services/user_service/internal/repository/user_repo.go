package repository

import (
	"context"

	"78-pflops/services/user_service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, user.PasswordHash, user.Name)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name FROM users WHERE email = $1`, email)
	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}
