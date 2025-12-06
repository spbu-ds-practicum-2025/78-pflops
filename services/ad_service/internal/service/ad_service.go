package service

import (
	"context"

	"78-pflops/services/ad_service/internal/model"
	"78-pflops/services/ad_service/internal/repository"
)

type AdService struct {
	repo *repository.AdRepository
}

func NewAdService(repo *repository.AdRepository) *AdService { return &AdService{repo: repo} }

func (s *AdService) CreateAd(ctx context.Context, authorID, title, description string, price int64, categoryID, condition string) (*model.Ad, error) {
	ad := &model.Ad{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		Condition:   condition,
		Status:      "ACTIVE",
	}
	if err := s.repo.Create(ctx, ad); err != nil {
		return nil, err
	}
	return ad, nil
}

func (s *AdService) GetAd(ctx context.Context, id string) (*model.Ad, error) {
	return s.repo.Get(ctx, id)
}

func (s *AdService) SearchAds(ctx context.Context, text string, categoryID *string, priceMin, priceMax *int64, condition *string, limit, offset int) ([]model.Ad, int, error) {
	return s.repo.Search(ctx, text, categoryID, priceMin, priceMax, condition, limit, offset)
}

func (s *AdService) DeleteAd(ctx context.Context, id string, authorID string) error {
	return s.repo.Delete(ctx, id, authorID)
}
