package service

import (
	"context"
	"errors"

	"78-pflops/services/ad_service/internal/model"
	"78-pflops/services/ad_service/internal/repository"
)

// repoInterface abstracts persistence for testability.
type repoInterface interface {
	Create(ctx context.Context, ad *model.Ad) error
	Get(ctx context.Context, id string) (*model.Ad, error)
	Search(ctx context.Context, text string, categoryID *string, priceMin, priceMax *int64, condition *string, limit, offset int) ([]model.Ad, int, error)
	Update(ctx context.Context, id string, authorID string, title, description *string, price *int64, categoryID, condition, status *string) error
	Delete(ctx context.Context, id string, authorID string) error
	AttachMedia(ctx context.Context, adID, mediaID string) error
	ListImages(ctx context.Context, adID string) ([]model.AdImage, error)
	DetachMedia(ctx context.Context, adID, mediaID string) error
	ReplaceImages(ctx context.Context, adID string, mediaIDs []string) error
}

type AdService struct {
	repo repoInterface
}

// NewAdService keeps backward compatibility with concrete repository.
func NewAdService(repo *repository.AdRepository) *AdService { return &AdService{repo: repo} }

// Filters for listing ads. Keep minimal yet flexible.
type Filters struct {
	Text       string
	CategoryID *string
	PriceMin   *int64
	PriceMax   *int64
	Condition  *string
	Limit      int
	Offset     int
}

// CreateAd(user_id, title, description, price)
func (s *AdService) CreateAd(ctx context.Context, userID, title, description string, price int64) (*model.Ad, error) {
	// Minimal defaults to satisfy schema
	defaultCategory := "00000000-0000-0000-0000-000000000000"
	defaultCondition := "NEW"
	ad := &model.Ad{
		AuthorID:    userID,
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  defaultCategory,
		Condition:   defaultCondition,
		Status:      "ACTIVE",
	}
	if err := s.repo.Create(ctx, ad); err != nil {
		return nil, err
	}
	return ad, nil
}

// GetAd(ad_id)
func (s *AdService) GetAd(ctx context.Context, adID string) (*model.Ad, error) {
	ad, err := s.repo.Get(ctx, adID)
	if err != nil {
		return nil, err
	}
	images, err := s.repo.ListImages(ctx, adID)
	if err != nil {
		return nil, err
	}
	// attach images to ad model for use in transport layer
	ad.Images = images
	return ad, nil
}

// ListAds(filters)
func (s *AdService) ListAds(ctx context.Context, f Filters) ([]model.Ad, int, error) {
	return s.repo.Search(ctx, f.Text, f.CategoryID, f.PriceMin, f.PriceMax, f.Condition, f.Limit, f.Offset)
}

// UpdateAd(ad_id, user_id, title?, description?, price?, category_id?, condition?, status?)
func (s *AdService) UpdateAd(ctx context.Context, adID, userID string, title, description *string, price *int64, categoryID, condition, status *string) error {
	return s.repo.Update(ctx, adID, userID, title, description, price, categoryID, condition, status)
}

// DeleteAd(ad_id, user_id)
func (s *AdService) DeleteAd(ctx context.Context, adID, userID string) error {
	return s.repo.Delete(ctx, adID, userID)
}

// AttachMedia(ad_id, media_id)
func (s *AdService) AttachMedia(ctx context.Context, adID, mediaID string) error {
	return s.repo.AttachMedia(ctx, adID, mediaID)
}

// DetachMedia(ad_id, media_id)
func (s *AdService) DetachMedia(ctx context.Context, adID, mediaID string) error {
	return s.repo.DetachMedia(ctx, adID, mediaID)
}

// ReplaceImages(ad_id, media_ids)
func (s *AdService) ReplaceImages(ctx context.Context, adID, userID string, mediaIDs []string) error {
	// Авторизация по владельцу объявления (возможность добавить админа в будущем).
	ad, err := s.repo.Get(ctx, adID)
	if err != nil {
		return err
	}
	if ad.AuthorID != userID {
		return errors.New("not found or no permission")
	}
	return s.repo.ReplaceImages(ctx, adID, mediaIDs)
}

func (s *AdService) CreateAdWithImages(ctx context.Context, userID, title, description string, price int64, mediaIDs []string) (*model.Ad, error) {
	ad, err := s.CreateAd(ctx, userID, title, description, price)
	if err != nil {
		return nil, err
	}
	for _, mid := range mediaIDs {
		if mid == "" {
			continue
		}
		if err := s.AttachMedia(ctx, ad.ID, mid); err != nil {
			return nil, err
		}
	}
	return ad, nil
}
