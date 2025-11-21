package service

import (
	"context"
	"errors"

	"78-pflops/services/user_service/internal/model"
	"78-pflops/services/user_service/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (string, string, error) {
	// Check if email already exists to return a Go error before hitting DB constraints
	if existing, err := s.repo.GetByEmail(ctx, email); err == nil && existing != nil {
		return "", "", errors.New("email already exists")
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", "", err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", "", err
	}

	if !utils.IsValidEmail(email) {
		return "", "", errors.New("invalid email format")
	}

	if !utils.IsValidPassword(password) {
		return "", "", errors.New("password must be at least 8 characters long and contain at least one special character")
	}

	user := &model.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Name:         name,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return "", "", err
	}

	token, err := utils.GenerateToken(user.ID)
	return user.ID, token, err
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("invalid credentials")
	} else if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID)
}
