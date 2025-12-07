package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"78-pflops/services/ad_service/internal/db"
	"78-pflops/services/ad_service/internal/model"
	"78-pflops/services/ad_service/internal/repository"
	"78-pflops/services/ad_service/internal/service"
	adpb "78-pflops/services/ad_service/pb/ad_service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO: import generated proto after we create ad.proto
// import "78-pflops/services/ad_service/pb/ad_service/pb"

type adServer struct {
	adpb.UnimplementedAdServiceServer
	svc *service.AdService
}

func newServer() *adServer {
	pool := db.Connect()
	repo := repository.NewAdRepository(pool)
	svc := service.NewAdService(repo)
	return &adServer{svc: svc}
}

// helper: convert domain model to protobuf
func toPb(ad *model.Ad) *adpb.Ad {
	if ad == nil {
		return nil
	}
	var rating float64
	if ad.SellerRatingCached != nil {
		rating = *ad.SellerRatingCached
	}
	imageURLs := make([]string, 0, len(ad.Images))
	for _, img := range ad.Images {
		if img.URL != "" {
			imageURLs = append(imageURLs, img.URL)
		}
	}
	return &adpb.Ad{
		Id:           ad.ID,
		AuthorId:     ad.AuthorID,
		Title:        ad.Title,
		Description:  ad.Description,
		Price:        ad.Price,
		CategoryId:   ad.CategoryID,
		Condition:    ad.Condition,
		ImageUrls:    imageURLs,
		SellerRating: rating,
		CreatedAt:    ad.CreatedAt.Unix(),
		UpdatedAt:    ad.UpdatedAt.Unix(),
	}
}

// CreateAd implements gRPC CreateAd
func (s *adServer) CreateAd(ctx context.Context, req *adpb.CreateAdRequest) (*adpb.CreateAdResponse, error) {
	ad, err := s.svc.CreateAd(ctx, req.UserId, req.Title, req.Description, req.Price)
	if err != nil {
		return nil, err
	}
	return &adpb.CreateAdResponse{Ad: toPb(ad)}, nil
}

func (s *adServer) GetAd(ctx context.Context, req *adpb.GetAdRequest) (*adpb.GetAdResponse, error) {
	ad, err := s.svc.GetAd(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &adpb.GetAdResponse{Ad: toPb(ad)}, nil
}

func (s *adServer) ListAds(ctx context.Context, req *adpb.ListAdsRequest) (*adpb.ListAdsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	var categoryPtr *string
	if req.CategoryId != "" {
		categoryPtr = &req.CategoryId
	}
	var priceMinPtr *int64
	if req.PriceMin > 0 {
		v := req.PriceMin
		priceMinPtr = &v
	}
	var priceMaxPtr *int64
	if req.PriceMax > 0 {
		v := req.PriceMax
		priceMaxPtr = &v
	}
	var conditionPtr *string
	if req.Condition != "" {
		conditionPtr = &req.Condition
	}
	ads, total, err := s.svc.ListAds(ctx, service.Filters{Text: req.Text, CategoryID: categoryPtr, PriceMin: priceMinPtr, PriceMax: priceMaxPtr, Condition: conditionPtr, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	respAds := make([]*adpb.Ad, 0, len(ads))
	for i := range ads {
		respAds = append(respAds, toPb(&ads[i]))
	}
	return &adpb.ListAdsResponse{Ads: respAds, Total: int32(total), Page: int32(page), PageSize: int32(limit)}, nil
}

func (s *adServer) UpdateAd(ctx context.Context, req *adpb.UpdateAdRequest) (*adpb.UpdateAdResponse, error) {
	var titlePtr, descPtr *string
	var pricePtr *int64
	var categoryPtr, conditionPtr, statusPtr *string
	if req.Title != nil {
		v := req.Title.Value
		titlePtr = &v
	}
	if req.Description != nil {
		v := req.Description.Value
		descPtr = &v
	}
	if req.Price != nil {
		v := req.Price.Value
		pricePtr = &v
	}
	if req.CategoryId != nil {
		v := req.CategoryId.Value
		categoryPtr = &v
	}
	if req.Condition != nil {
		v := req.Condition.Value
		conditionPtr = &v
	}
	if req.Status != nil {
		v := req.Status.Value
		statusPtr = &v
	}
	if err := s.svc.UpdateAd(ctx, req.AdId, req.UserId, titlePtr, descPtr, pricePtr, categoryPtr, conditionPtr, statusPtr); err != nil {
		return nil, err
	}
	return &adpb.UpdateAdResponse{}, nil
}

func (s *adServer) DeleteAd(ctx context.Context, req *adpb.DeleteAdRequest) (*adpb.DeleteAdResponse, error) {
	if err := s.svc.DeleteAd(ctx, req.AdId, req.UserId); err != nil {
		return nil, err
	}
	return &adpb.DeleteAdResponse{}, nil
}

func (s *adServer) AttachMedia(ctx context.Context, req *adpb.AttachMediaRequest) (*adpb.AttachMediaResponse, error) {
	if req.AdId == "" {
		return nil, status.Error(codes.InvalidArgument, "ad_id is required")
	}
	if req.MediaId == "" {
		return nil, status.Error(codes.InvalidArgument, "media_id is required")
	}
	if err := s.svc.AttachMedia(ctx, req.AdId, req.MediaId); err != nil {
		return nil, err
	}
	return &adpb.AttachMediaResponse{}, nil
}

func (s *adServer) CreateAdWithImages(ctx context.Context, req *adpb.CreateAdWithImagesRequest) (*adpb.CreateAdWithImagesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	ad, err := s.svc.CreateAdWithImages(ctx, req.UserId, req.Title, req.Description, req.Price, req.MediaIds)
	if err != nil {
		return nil, err
	}

	// пока игнорируем фактическую загрузку изображений, mediaIDs = nil
	return &adpb.CreateAdWithImagesResponse{Ad: toPb(ad)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	adpb.RegisterAdServiceServer(grpcServer, newServer())

	reflection.Register(grpcServer)

	fmt.Println("🚀 AdService running on port 50052 (gRPC)")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Placeholder to avoid unused imports warning until implementation
var _ = time.Now
