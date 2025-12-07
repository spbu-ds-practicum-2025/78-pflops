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
	"google.golang.org/grpc/reflection"
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
	userID, token, err := s.service.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{Id: userID, Token: token, Email: req.Email, Name: req.Name}, nil
}

func (s *userServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}

func (s *userServiceServer) ValidateToken(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	userID, valid, err := s.service.Validate(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &proto.ValidateResponse{UserId: userID, Valid: valid}, nil
}

func (s *userServiceServer) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
	user, err := s.service.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.GetProfileResponse{UserId: user.ID, Name: user.Name}, nil
}

func (s *userServiceServer) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.UpdateProfileResponse, error) {
	if err := s.service.UpdateProfile(ctx, req.UserId, req.Name); err != nil {
		return &proto.UpdateProfileResponse{Success: false, Message: err.Error()}, nil
	}
	return &proto.UpdateProfileResponse{Success: true, Message: "profile updated"}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// В упрощённом варианте пароль не проверяем, только наличие токена/ID
	if err := s.service.DeleteUser(ctx, req.UserId); err != nil {
		return &proto.DeleteUserResponse{Success: false, Message: err.Error()}, nil
	}
	return &proto.DeleteUserResponse{Success: true, Message: "user deleted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, newServer())

	reflection.Register(grpcServer)

	fmt.Println("✅ UserService running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
