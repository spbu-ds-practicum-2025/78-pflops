package main

import (
	"78-pflops/services/user_service/gen/proto"
	"78-pflops/services/user_service/internal/db"
	"78-pflops/services/user_service/internal/repository"
	"78-pflops/services/user_service/internal/service"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	proto.UnimplementedUserServiceServer
	service *service.UserService
}

func newServer() *userServiceServer {
	conn := db.Connect()
	repo := repository.NewUserRepository(conn)
	svc := service.NewUserService(repo)
	return &userServiceServer{service: svc}
}

func (s *userServiceServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	userID, token, err := s.service.Register(ctx, req.Email, req.Password, "New User")
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{UserId: userID, Token: token}, nil
}

func (s *userServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, newServer())

	fmt.Println("✅ UserService running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
