package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	grpc "google.golang.org/grpc"

	adpb "78-pflops/services/ad_service/pb/ad_service/pb"
	userpb "78-pflops/services/user_service/pb/user_service/pb"
)

type gateway struct {
	userSvcAddr string
	adSvcAddr   string
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createAdRequest struct {
	Token       string   `json:"token"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	userSvcAddr := getenv("USER_SERVICE_ADDR", "user_service_app:50051")
	adSvcAddr := getenv("AD_SERVICE_ADDR", "ad_service_app:50052")
	port := getenv("HTTP_GATEWAY_PORT", "8081")

	g := &gateway{userSvcAddr: userSvcAddr, adSvcAddr: adSvcAddr}

	http.HandleFunc("/api/auth/register", g.handleRegister)
	http.HandleFunc("/api/auth/login", g.handleLogin)
	http.HandleFunc("/api/ads", g.handleAds)

	log.Printf("HTTP gateway listening on %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

func (g *gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, g.userSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	resp, err := client.Register(ctx, &userpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g *gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, g.userSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	resp, err := client.Login(ctx, &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g *gateway) handleAds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		g.listAds(w, r)
	case http.MethodPost:
		g.createAd(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (g *gateway) listAds(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var minPrice, maxPrice int64
	if v := q.Get("min_price"); v != "" {
		p, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		minPrice = p
	}
	if v := q.Get("max_price"); v != "" {
		p, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		maxPrice = p
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, g.adSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := adpb.NewAdServiceClient(conn)
	resp, err := client.ListAds(ctx, &adpb.ListAdsRequest{
		Text:       q.Get("query"),
		CategoryId: q.Get("category"),
		PriceMin:   minPrice,
		PriceMax:   maxPrice,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g *gateway) createAd(w http.ResponseWriter, r *http.Request) {
	var req createAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// сначала валидируем токен и получаем user_id
	userConn, err := grpc.DialContext(ctx, g.userSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	valResp, err := userClient.Validate(ctx, &userpb.ValidateRequest{Token: req.Token})
	if err != nil || !valResp.GetValid() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	adConn, err := grpc.DialContext(ctx, g.adSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer adConn.Close()

	adClient := adpb.NewAdServiceClient(adConn)
	// Пока создаём объявление без картинок через CreateAd
	createResp, err := adClient.CreateAd(ctx, &adpb.CreateAdRequest{
		UserId:      fmt.Sprintf("%d", valResp.GetUserId()),
		Title:       req.Title,
		Description: req.Description,
		Price:       int64(req.Price),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(createResp)
}
