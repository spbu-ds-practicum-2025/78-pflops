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
	// TEMP: author id should come from JWT (metadata); using placeholder
	authorID := "00000000-0000-0000-0000-000000000001"
	ad, err := s.svc.CreateAd(ctx, authorID, req.Title, req.Description, req.Price, req.CategoryId, req.Condition)
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

func (s *adServer) SearchAds(ctx context.Context, req *pb.SearchAdsRequest) (*pb.SearchAdsResponse, error) {
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
	ads, total, err := s.svc.SearchAds(ctx, req.Text, categoryPtr, priceMinPtr, priceMaxPtr, conditionPtr, limit, offset)
	if err != nil {
		return nil, err
	}
	respAds := make([]*pb.Ad, 0, len(ads))
	for i := range ads {
		respAds = append(respAds, toPb(&ads[i]))
	}
	return &pb.SearchAdsResponse{Ads: respAds, Total: int32(total), Page: int32(page), PageSize: int32(limit)}, nil
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
