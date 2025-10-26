package main

import (
	"log"
	"net"

	"user_service/internal/config"
	"user_service/internal/repository"
	"user_service/internal/server"
	"user_service/internal/service"
	"user_service/pb/services/user_service/pb"

	"google.golang.org/grpc"
)

func main() {
	// 1. Загружаем конфигурацию (.env)
	cfg := config.Load()

	// 2. Подключаемся к БД
	db := repository.InitDB(cfg.DatabaseURL)

	// 3. Создаем слой бизнес-логики
	userService := service.NewUserService(db)

	// 4. Создаем gRPC сервер и регистрируем UserServiceServer
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, server.NewUserServer(userService))

	// 5. Слушаем порт
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("User Service gRPC running on port %s", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
