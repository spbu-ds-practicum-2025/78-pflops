package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	adpb "78-pflops/services/ad_service/pb/ad_service/pb"
	mediapb "78-pflops/services/http_gateway/mediapb"
	userpb "78-pflops/services/user_service/pb/user_service/pb"
)

type gateway struct {
	userSvcAddr  string
	adSvcAddr    string
	userHTTPBase string
	mediaSvcAddr string
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

type updateAdRequest struct {
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
	userHTTPBase := getenv("USER_HTTP_BASE", "http://user_service_app:8081")
	mediaSvcAddr := getenv("MEDIA_SERVICE_ADDR", "media_service_app:50053")

	g := &gateway{userSvcAddr: userSvcAddr, adSvcAddr: adSvcAddr, userHTTPBase: userHTTPBase, mediaSvcAddr: mediaSvcAddr}

	http.HandleFunc("/api/auth/register", g.handleRegister)
	http.HandleFunc("/api/auth/login", g.handleLogin)
	http.HandleFunc("/api/ads", g.handleAds)
	http.HandleFunc("/api/ads/", g.handleAdByID)

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

// handleAdByID обрабатывает запросы /api/ads/{id} для получения, обновления и удаления объявления.
func (g *gateway) handleAdByID(w http.ResponseWriter, r *http.Request) {
	// Ожидаем путь формата /api/ads/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/ads"), "/")
	if len(parts) < 2 || parts[1] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[1]

	switch r.Method {
	case http.MethodGet:
		g.getAd(w, r, id)
	case http.MethodPut, http.MethodPatch:
		g.updateAd(w, r, id)
	case http.MethodDelete:
		g.deleteAd(w, r, id)
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

	// сначала валидируем токен через HTTP /api/users/me и получаем user_id
	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, g.userHTTPBase+"/api/users/me", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meReq.Header.Set("Authorization", "Bearer "+req.Token)

	resp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var me struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil || me.UserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Загружаем изображения в MediaService и получаем media_ids/urls
	mediaConn, err := grpc.DialContext(ctx, g.mediaSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer mediaConn.Close()

	mediaClient := mediapb.NewMediaServiceClient(mediaConn)
	mediaIDs := make([]string, 0, len(req.Images))
	for idx, imgB64 := range req.Images {
		if imgB64 == "" {
			continue
		}
		// На фронте мы отправляем только base64 без префикса data:
		data, err := base64.StdEncoding.DecodeString(imgB64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fileName := filepath.Base("image")
		_, _ = idx, fileName
		upResp, err := mediaClient.UploadMedia(ctx, &mediapb.UploadMediaRequest{
			UserId:    me.UserID,
			FileBytes: data,
			MimeType:  "image/jpeg",
			FileName:  "image.jpg",
		})
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		// Используем URL, который вернул MediaService; если по какой-то причине
		// он пустой, сохраняем data:-URL как запасной вариант, чтобы картинки
		// продолжали отображаться.
		url := upResp.Url
		if url == "" {
			url = "data:image/jpeg;base64," + imgB64
		}
		mediaIDs = append(mediaIDs, url)
	}

	adConn, err := grpc.DialContext(ctx, g.adSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer adConn.Close()

	adClient := adpb.NewAdServiceClient(adConn)
	// Создаём объявление с привязанными изображениями
	createResp, err := adClient.CreateAdWithImages(ctx, &adpb.CreateAdWithImagesRequest{
		UserId:      me.UserID,
		Title:       req.Title,
		Description: req.Description,
		Price:       int64(req.Price),
		MediaIds:    mediaIDs,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(createResp)
}

// getAd возвращает одно объявление по ID.
func (g *gateway) getAd(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, g.adSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := adpb.NewAdServiceClient(conn)
	resp, err := client.GetAd(ctx, &adpb.GetAdRequest{Id: id})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// updateAd обновляет объявление (заголовок, описание, цену, категорию).
// Требуется заголовок Authorization: Bearer <token>.
func (g *gateway) updateAd(w http.ResponseWriter, r *http.Request, id string) {
	// Получаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Валидируем токен через /api/users/me и получаем user_id
	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, g.userHTTPBase+"/api/users/me", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meReq.Header.Set("Authorization", "Bearer "+token)

	meResp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var me struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil || me.UserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reqBody updateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Подготовка списка новых изображений (если переданы)
	newMediaIDs := make([]string, 0, len(reqBody.Images))
	if len(reqBody.Images) > 0 {
		mediaConn, err := grpc.DialContext(ctx, g.mediaSvcAddr, grpc.WithInsecure())
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer mediaConn.Close()

		mediaClient := mediapb.NewMediaServiceClient(mediaConn)
		for _, imgB64 := range reqBody.Images {
			if imgB64 == "" {
				continue
			}
			data, err := base64.StdEncoding.DecodeString(imgB64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			upResp, err := mediaClient.UploadMedia(ctx, &mediapb.UploadMediaRequest{
				UserId:    me.UserID,
				FileBytes: data,
				MimeType:  "image/jpeg",
				FileName:  "image.jpg",
			})
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			url := upResp.Url
			if url == "" {
				url = "data:image/jpeg;base64," + imgB64
			}
			newMediaIDs = append(newMediaIDs, url)
		}
	}

	if reqBody.Title == "" && reqBody.Description == "" && reqBody.Price == 0 && reqBody.Category == "" && len(newMediaIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	adConn, err := grpc.DialContext(ctx, g.adSvcAddr, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer adConn.Close()

	adClient := adpb.NewAdServiceClient(adConn)
	// Обновляем объявление: для простоты считаем, что все поля передаются целиком.
	updateReq := &adpb.UpdateAdRequest{
		AdId:   id,
		UserId: me.UserID,
	}
	if reqBody.Title != "" {
		updateReq.Title = wrapperspb.String(reqBody.Title)
	}
	if reqBody.Description != "" {
		updateReq.Description = wrapperspb.String(reqBody.Description)
	}
	if reqBody.Price > 0 {
		updateReq.Price = wrapperspb.Int64(int64(reqBody.Price))
	}
	if reqBody.Category != "" {
		updateReq.CategoryId = wrapperspb.String(reqBody.Category)
	}

	// Если есть изменения текста/цены/категории — отправляем UpdateAd
	if updateReq.Title != nil || updateReq.Description != nil || updateReq.Price != nil || updateReq.CategoryId != nil {
		if _, err := adClient.UpdateAd(ctx, updateReq); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	// Если переданы новые изображения — полностью заменяем их в объявлении
	if len(newMediaIDs) > 0 {
		if _, err := adClient.ReplaceImages(ctx, &adpb.ReplaceImagesRequest{
			AdId:     id,
			UserId:   me.UserID,
			MediaIds: newMediaIDs,
		}); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteAd удаляет объявление. Требуется заголовок Authorization: Bearer <token>.
func (g *gateway) deleteAd(w http.ResponseWriter, r *http.Request, id string) {
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, g.userHTTPBase+"/api/users/me", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meReq.Header.Set("Authorization", "Bearer "+token)

	meResp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var me struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil || me.UserID == "" {
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
	if _, err := adClient.DeleteAd(ctx, &adpb.DeleteAdRequest{AdId: id, UserId: me.UserID}); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
