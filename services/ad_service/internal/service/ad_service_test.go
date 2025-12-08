package service

import (
	"context"
	"testing"
	"time"

	"78-pflops/services/ad_service/internal/model"
)

type stubRepo struct {
	createErr error
	deleteErr error
	getAd     *model.Ad
	searchAds []model.Ad
	searchCnt int
}

func (s *stubRepo) Create(ctx context.Context, ad *model.Ad) error {
	if s.createErr != nil {
		return s.createErr
	}
	// simulate persistence side effects
	if ad.ID == "" {
		ad.ID = "stub-id"
	}
	ad.CreatedAt = time.Unix(1000, 0)
	ad.UpdatedAt = ad.CreatedAt
	return nil
}
func (s *stubRepo) Get(ctx context.Context, id string) (*model.Ad, error) { return s.getAd, nil }
func (s *stubRepo) Search(ctx context.Context, text string, categoryID *string, priceMin, priceMax *int64, condition *string, limit, offset int) ([]model.Ad, int, error) {
	return s.searchAds, s.searchCnt, nil
}

func (s *stubRepo) Update(ctx context.Context, id string, authorID string, title, description *string, price *int64, categoryID, condition, status *string) error {
	return nil
}
func (s *stubRepo) Delete(ctx context.Context, id string, authorID string) error { return s.deleteErr }
func (s *stubRepo) AttachMedia(ctx context.Context, adID, mediaID string) error  { return nil }

func (s *stubRepo) ListImages(ctx context.Context, adID string) ([]model.AdImage, error) {
	return nil, nil
}

func (s *stubRepo) DetachMedia(ctx context.Context, adID, mediaID string) error { return nil }

func (s *stubRepo) ReplaceImages(ctx context.Context, adID string, mediaIDs []string) error {
	return nil
}

func TestCreateAd(t *testing.T) {
	repo := &stubRepo{}
	svc := &AdService{repo: repo}
	ad, err := svc.CreateAd(context.Background(), "author-1", "Title", "Desc", 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ad == nil {
		t.Fatalf("ad is nil")
	}
	if ad.Status != "ACTIVE" {
		t.Errorf("expected status ACTIVE got %s", ad.Status)
	}
	if ad.ID != "stub-id" {
		t.Errorf("expected id stub-id got %s", ad.ID)
	}
	if ad.Price != 123 {
		t.Errorf("price not set")
	}
}

func TestGetAd(t *testing.T) {
	expected := &model.Ad{ID: "x", Title: "T"}
	repo := &stubRepo{getAd: expected}
	svc := &AdService{repo: repo}
	ad, err := svc.GetAd(context.Background(), "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ad != expected {
		t.Errorf("expected same pointer")
	}
}

func TestListAds(t *testing.T) {
	list := []model.Ad{{ID: "1"}, {ID: "2"}}
	repo := &stubRepo{searchAds: list, searchCnt: 2}
	svc := &AdService{repo: repo}
	ads, cnt, err := svc.ListAds(context.Background(), Filters{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cnt != 2 {
		t.Errorf("expected count 2 got %d", cnt)
	}
	if len(ads) != 2 {
		t.Errorf("expected 2 ads")
	}
}

func TestDeleteAd(t *testing.T) {
	repo := &stubRepo{}
	svc := &AdService{repo: repo}
	if err := svc.DeleteAd(context.Background(), "id", "author"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateAd(t *testing.T) {
	repo := &stubRepo{}
	svc := &AdService{repo: repo}
	title := "New"
	if err := svc.UpdateAd(context.Background(), "ad1", "author-1", &title, nil, nil, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAttachMedia(t *testing.T) {
	repo := &stubRepo{}
	svc := &AdService{repo: repo}
	if err := svc.AttachMedia(context.Background(), "ad1", "media-123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
