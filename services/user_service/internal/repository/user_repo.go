package repository

import (
	"context"

	"78-pflops/services/user_service/internal/model"

	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name FROM users WHERE id = $1`, id)
	var u model.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateName(ctx context.Context, id, name string) error {
	result, err := r.db.Exec(ctx,
		`UPDATE users SET name = $1 WHERE id = $2`, name, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
