package server

import (
	"context"
	"user_service/internal/service"
	"user_service/pb/services/user_service/pb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserServer(s *service.UserService) *UserServer {
	return &UserServer{userService: s}
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.userService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.userService.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *UserServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	user, err := s.userService.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}
	return &pb.ValidateResponse{
		Valid:  true,
		UserId: uint64(user.ID),
		Email:  user.Email,
	}, nil
}
