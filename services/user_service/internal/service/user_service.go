package service

import (
	"errors"
	"user_service/internal/models"
	"user_service/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(email, password, name string) error {
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	}

	hash, _ := utils.HashPassword(password)
	user := models.User{Email: email, PasswordHash: hash, Name: name}
	return s.db.Create(&user).Error
}

func (s *UserService) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email)
	return token, nil
}

func (s *UserService) ValidateToken(token string) (*models.User, error) {
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
