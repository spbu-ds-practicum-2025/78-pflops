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
	pb "78-pflops/services/ad_service/pb/ad_service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: import generated proto after we create ad.proto
// import "78-pflops/services/ad_service/pb/ad_service/pb"

type adServer struct {
	pb.UnimplementedAdServiceServer
	svc *service.AdService
}

func newServer() *adServer {
	pool := db.Connect()
	repo := repository.NewAdRepository(pool)
	svc := service.NewAdService(repo)
	return &adServer{svc: svc}
}

// helper: convert domain model to protobuf
func toPb(ad *model.Ad) *pb.Ad {
	if ad == nil {
		return nil
	}
	var rating float64
	if ad.SellerRatingCached != nil {
		rating = *ad.SellerRatingCached
	}
	return &pb.Ad{
		Id:           ad.ID,
		AuthorId:     ad.AuthorID,
		Title:        ad.Title,
		Description:  ad.Description,
		Price:        ad.Price,
		CategoryId:   ad.CategoryID,
		Condition:    ad.Condition,
		ImageUrls:    []string{}, // will fill when image repository added
		SellerRating: rating,
		CreatedAt:    ad.CreatedAt.Unix(),
		UpdatedAt:    ad.UpdatedAt.Unix(),
	}
}

// CreateAd implements gRPC CreateAd
func (s *adServer) CreateAd(ctx context.Context, req *pb.CreateAdRequest) (*pb.CreateAdResponse, error) {
	ad, err := s.svc.CreateAd(ctx, req.UserId, req.Title, req.Description, req.Price)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAdResponse{Ad: toPb(ad)}, nil
}

func (s *adServer) GetAd(ctx context.Context, req *pb.GetAdRequest) (*pb.GetAdResponse, error) {
	ad, err := s.svc.GetAd(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAdResponse{Ad: toPb(ad)}, nil
}

func (s *adServer) ListAds(ctx context.Context, req *pb.ListAdsRequest) (*pb.ListAdsResponse, error) {
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
	respAds := make([]*pb.Ad, 0, len(ads))
	for i := range ads {
		respAds = append(respAds, toPb(&ads[i]))
	}
	return &pb.ListAdsResponse{Ads: respAds, Total: int32(total), Page: int32(page), PageSize: int32(limit)}, nil
}

func (s *adServer) UpdateAd(ctx context.Context, req *pb.UpdateAdRequest) (*pb.UpdateAdResponse, error) {
	var titlePtr, descPtr *string
	var pricePtr *int64
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
	if err := s.svc.UpdateAd(ctx, req.AdId, req.UserId, titlePtr, descPtr, pricePtr); err != nil {
		return nil, err
	}
	return &pb.UpdateAdResponse{}, nil
}

func (s *adServer) DeleteAd(ctx context.Context, req *pb.DeleteAdRequest) (*pb.DeleteAdResponse, error) {
	if err := s.svc.DeleteAd(ctx, req.AdId, req.UserId); err != nil {
		return nil, err
	}
	return &pb.DeleteAdResponse{}, nil
}

func (s *adServer) AttachMedia(ctx context.Context, req *pb.AttachMediaRequest) (*pb.AttachMediaResponse, error) {
	if err := s.svc.AttachMedia(ctx, req.AdId, req.MediaId); err != nil {
		return nil, err
	}
	return &pb.AttachMediaResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdServiceServer(grpcServer, newServer())

	reflection.Register(grpcServer)

	fmt.Println("🚀 AdService running on port 50052 (gRPC)")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Placeholder to avoid unused imports warning until implementation
var _ = time.Now
