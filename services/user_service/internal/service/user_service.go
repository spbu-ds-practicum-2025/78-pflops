package service

import (
	"context"
	"errors"

	"78-pflops/services/user_service/internal/model"
	"78-pflops/services/user_service/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	UpdateName(ctx context.Context, id, name string) error
	Delete(ctx context.Context, id string) error
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

// Validate проверяет JWT и возвращает userID и флаг валидности.
func (s *UserService) Validate(ctx context.Context, token string) (string, bool, error) {
	userID, valid, err := utils.ValidateToken(token)
	if err != nil {
		return "", false, err
	}
	if !valid {
		return "", false, nil
	}
	// опционально можно проверить, что пользователь ещё существует
	if _, err := s.repo.GetByID(ctx, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return userID, true, nil
}

// GetProfile возвращает профиль пользователя по ID.
func (s *UserService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// UpdateProfile обновляет только имя пользователя.
func (s *UserService) UpdateProfile(ctx context.Context, userID, name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return s.repo.UpdateName(ctx, userID, name)
}

// DeleteUser удаляет пользователя по ID.
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}
