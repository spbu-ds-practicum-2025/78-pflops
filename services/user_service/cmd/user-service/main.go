package main

import (
	"78-pflops/services/user_service/gen/proto"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	proto.UnimplementedUserServiceServer
}

func (s *userServiceServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return &proto.RegisterResponse{
		UserId: "123",
		Token:  "fake-jwt-token",
	}, nil
}

func (s *userServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return &proto.LoginResponse{Token: "fake-jwt-token"}, nil
}

func (s *userServiceServer) ValidateToken(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	return &proto.ValidateResponse{UserId: "123", Valid: true}, nil
}

func (s *userServiceServer) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
	return &proto.GetProfileResponse{
		UserId: req.UserId,
		Name:   "Test User",
	}, nil
}

func (s *userServiceServer) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.UpdateProfileResponse, error) {
	return &proto.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated",
	}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return &proto.DeleteUserResponse{
		Success: true,
		Message: "User deleted",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	fmt.Println("UserService gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
