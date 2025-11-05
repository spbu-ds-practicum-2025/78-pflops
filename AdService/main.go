package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"ad_service/gen"

	"google.golang.org/grpc"
)

func main() {
	//запускаем gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	adService := newAdServiceServer()
	gen.RegisterAdServiceServer(grpcServer, adService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("gRPC Server started on :50051")
	fmt.Println("Available methods:")
	fmt.Println("   - ListAds")
	fmt.Println("   - GetAd")
	fmt.Println("   - CreateAd")
	fmt.Println("   - UpdateAd")
	fmt.Println("   - DeleteAd")
	fmt.Println("   - SearchAds")
	fmt.Println("   - GetUserAds")
	fmt.Println("   - ChangeAdStatus")
	fmt.Println("   - ListCategories")
	fmt.Println("   - HealthCheck")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// интерцептор для логирования запросов
func loggingInterceptor(
	ctx context.Context, //контекст вызова
	req interface{}, //входящий запрос
	info *grpc.UnaryServerInfo, //информация о методе
	handler grpc.UnaryHandler, //следующий обработчик
) (interface{}, error) { //возвращает результат
	log.Printf("gRPC method %s called", info.FullMethod)
	return handler(ctx, req)
}
