package main

import (
	"78-pflops/services/user_service/gen/proto"
	"78-pflops/services/user_service/internal/db"
	"78-pflops/services/user_service/internal/repository"
	"78-pflops/services/user_service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

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

// HTTP handlers

type httpServer struct {
	svc *service.UserService
}

func newHTTPServer(svc *service.UserService) *httpServer {
	return &httpServer{svc: svc}
}

func (h *httpServer) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	ctx := r.Context()
	id, token, err := h.svc.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":    id,
		"token": token,
		"email": req.Email,
		"name":  req.Name,
	})
}

func (h *httpServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	ctx := r.Context()
	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	})
}

func (h *httpServer) meHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing bearer token"})
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	ctx := r.Context()
	userID, valid, err := h.svc.Validate(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid token"})
		return
	}
	user, err := h.svc.GetProfile(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	})
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
	// shared service layer
	conn := db.Connect()
	repo := repository.NewUserRepository(conn)
	svc := service.NewUserService(repo)

	// gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		proto.RegisterUserServiceServer(grpcServer, &userServiceServer{service: svc})
		reflection.Register(grpcServer)
		fmt.Println("✅ UserService gRPC running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// HTTP server
	httpSrv := newHTTPServer(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", httpSrv.registerHandler)
	mux.HandleFunc("/api/auth/login", httpSrv.loginHandler)
	mux.HandleFunc("/api/users/me", httpSrv.meHandler)

	fmt.Println("✅ UserService HTTP running on port 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
