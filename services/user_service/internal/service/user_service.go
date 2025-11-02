package service

import (
	"context"
	"errors"

	"78-pflops/services/user_service/internal/model"
	"78-pflops/services/user_service/internal/repository"
	"78-pflops/services/user_service/internal/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (string, string, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", "", err
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
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID)
}
