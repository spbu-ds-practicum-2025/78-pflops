package main

import (
	"ad_service/gen"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type memoryStorage struct {
	ads      map[string]*gen.Ad //хранилище объявлений
	mutex    sync.RWMutex       //защита от гонки данных
	nextAdID int32              //счетчик для ID объявлений (артикулы)
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		ads:      make(map[string]*gen.Ad),
		nextAdID: 1,
	}
}

func (s *memoryStorage) createAd(ad *gen.Ad) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ad.Id = strconv.Itoa(int(s.nextAdID))
	ad.CreatedAt = time.Now().Format(time.RFC3339)
	ad.UpdatedAt = ad.CreatedAt
	s.ads[ad.Id] = ad
	s.nextAdID++

	return nil
}

func (s *memoryStorage) getAd(id string) (*gen.Ad, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	ad, exists := s.ads[id]
	if !exists {
		return nil, nil
	}
	return ad, nil
}

func (s *memoryStorage) updateAd(ad *gen.Ad) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.ads[ad.Id]; !exists {
		return fmt.Errorf("ad not found")
	}
	ad.UpdatedAt = time.Now().Format(time.RFC3339)
	s.ads[ad.Id] = ad
	return nil
}

func (s *memoryStorage) deleteAd(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.ads[id]; !exists {
		return fmt.Errorf("ad not found")
	}
	delete(s.ads, id)
	return nil
}

func (s *memoryStorage) listAds() ([]*gen.Ad, int32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	ads := make([]*gen.Ad, 0, len(s.ads))
	for _, ad := range s.ads {
		ads = append(ads, ad)
	}

	return ads, int32(len(ads)), nil
}

// вспомогательная функция для поиска
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// gRPC сервер
type adServiceServer struct {
	gen.UnimplementedAdServiceServer
	storage *memoryStorage
}

func newAdServiceServer() *adServiceServer {
	return &adServiceServer{
		storage: newMemoryStorage(),
	}
}

func (s *adServiceServer) ListAds(ctx context.Context, req *gen.ListAdsRequest) (*gen.ListAdsResponse, error) {
	ads, total, err := s.storage.listAds()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get ads")
	}

	return &gen.ListAdsResponse{
		Ads:   ads,
		Total: total,
	}, nil
}

func (s *adServiceServer) GetAd(ctx context.Context, req *gen.GetAdRequest) (*gen.GetAdResponse, error) {

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	//получаем объявление
	ad, err := s.storage.getAd(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get ad")
	}

	if ad == nil {
		return nil, status.Error(codes.NotFound, "ad not found")
	}

	if ad.Status == gen.AdStatus_DRAFT {
		//TODO: добавить проверку авторизации - только автор может видеть черновики
		//пока возвращаем ошибку для всех
		return nil, status.Error(codes.PermissionDenied, "draft ads are only visible to owners")
	}

	return &gen.GetAdResponse{Ad: ad}, nil
}

func (s *adServiceServer) CreateAd(ctx context.Context, req *gen.CreateAdRequest) (*gen.CreateAdResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetPrice() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be positive")
	}
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// TODO: проверка user_id через UserService

	// создаем объявление
	ad := &gen.Ad{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		UserId:      req.GetUserId(),
		Status:      gen.AdStatus_DRAFT, // По умолчанию черновик
	}

	//сохраняем в хранилище
	if err := s.storage.createAd(ad); err != nil {
		return nil, status.Error(codes.Internal, "failed to create ad")
	}

	return &gen.CreateAdResponse{Ad: ad}, nil
}

func (s *adServiceServer) UpdateAd(ctx context.Context, req *gen.UpdateAdRequest) (*gen.UpdateAdResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetPrice() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be positive")
	}

	existingAd, err := s.storage.getAd(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get ad")
	}
	if existingAd == nil {
		return nil, status.Error(codes.NotFound, "ad not found")
	}

	//TODO: проверка владельца

	//обновляем поля
	existingAd.Title = req.GetTitle()
	existingAd.Description = req.GetDescription()
	existingAd.Price = req.GetPrice()

	if err := s.storage.updateAd(existingAd); err != nil {
		return nil, status.Error(codes.Internal, "failed to update ad")
	}

	return &gen.UpdateAdResponse{Ad: existingAd}, nil
}

func (s *adServiceServer) DeleteAd(ctx context.Context, req *gen.DeleteAdRequest) (*gen.DeleteAdResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.storage.deleteAd(req.GetId()); err != nil {
		return nil, status.Error(codes.NotFound, "ad not found")
	}

	return &gen.DeleteAdResponse{Success: true}, nil
}

func (s *adServiceServer) HealthCheck(ctx context.Context, req *gen.HealthCheckRequest) (*gen.HealthCheckResponse, error) {
	return &gen.HealthCheckResponse{
		Status: "OK",
	}, nil
}
