package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"78-pflops/services/user_service/internal/model"
	"78-pflops/services/user_service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// fakeRepo is an in-memory implementation of UserRepository for unit tests.
type fakeRepo struct {
	byEmail map[string]*model.User
	byID    map[string]*model.User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		byEmail: make(map[string]*model.User),
		byID:    make(map[string]*model.User),
	}
}

func (r *fakeRepo) Create(ctx context.Context, user *model.User) error {
	if _, exists := r.byEmail[user.Email]; exists {
		return errors.New("duplicate email")
	}
	// simulate DB insert
	r.byEmail[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *fakeRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	return u, nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	return u, nil
}

func (r *fakeRepo) UpdateName(ctx context.Context, id, name string) error {
	u, ok := r.byID[id]
	if !ok {
		return pgx.ErrNoRows
	}
	u.Name = name
	return nil
}

func (r *fakeRepo) Delete(ctx context.Context, id string) error {
	u, ok := r.byID[id]
	if !ok {
		return pgx.ErrNoRows
	}
	delete(r.byID, id)
	delete(r.byEmail, u.Email)
	return nil
}

func setupService() *UserService {
	repo := newFakeRepo()
	return NewUserService(repo)
}

func TestRegisterSuccess(t *testing.T) {
	svc := setupService()
	email := "user@example.com"
	password := "ValidPass!" // meets validation (>=8 + special char)
	name := "Test User"

	id, token, err := svc.Register(context.Background(), email, password, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non-empty id")
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
	// verify user stored
	stored, err := svc.repo.GetByEmail(context.Background(), email)
	if err != nil || stored == nil {
		t.Fatalf("user not stored: %v", err)
	}
	if stored.PasswordHash == password {
		t.Fatalf("password should be hashed")
	}
	// decode JWT
	parsed, perr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte("secret"), nil })
	if perr != nil || !parsed.Valid {
		t.Fatalf("token invalid: %v", perr)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	if claims["user_id"].(string) != id {
		t.Fatalf("token user_id mismatch: %v vs %v", claims["user_id"], id)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	svc := setupService()
	email := "dup@example.com"
	password := "ValidPass!"
	name := "A"
	if _, _, err := svc.Register(context.Background(), email, password, name); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if _, _, err := svc.Register(context.Background(), email, password, name); err == nil {
		t.Fatalf("expected duplicate email error")
	}
}

func TestRegisterInvalidEmail(t *testing.T) {
	svc := setupService()
	_, _, err := svc.Register(context.Background(), "invalid-email", "ValidPass!", "Name")
	if err == nil || err.Error() != "invalid email format" {
		t.Fatalf("expected invalid email format error, got: %v", err)
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	svc := setupService()
	// missing special char
	_, _, err := svc.Register(context.Background(), "a@b.com", "weakpass", "Name")
	if err == nil || err.Error() != "password must be at least 8 characters long and contain at least one special character" {
		t.Fatalf("expected password validation error, got: %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	svc := setupService()
	email := "login@example.com"
	pass := "ValidPass!"
	name := "Login User"
	if _, _, err := svc.Register(context.Background(), email, pass, name); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	token, err := svc.Login(context.Background(), email, pass)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Fatalf("expected token on login")
	}
}

func TestLoginEmailNotFound(t *testing.T) {
	svc := setupService()
	token, err := svc.Login(context.Background(), "nope@example.com", "AnyPass!")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got err=%v token=%s", err, token)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc := setupService()
	email := "wrongpass@example.com"
	pass := "ValidPass!"
	name := "User"
	if _, _, err := svc.Register(context.Background(), email, pass, name); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	token, err := svc.Login(context.Background(), email, "OtherPass!")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got err=%v token=%s", err, token)
	}
}

// Simple timing test to ensure hashing not absurdly fast (basic regression guard)
func TestPasswordHashingCost(t *testing.T) {
	start := time.Now()
	hash, err := utils.HashPassword("ValidPass!")
	if err != nil || hash == "" {
		t.Fatalf("hash failed: %v", err)
	}
	if dur := time.Since(start); dur < 5*time.Millisecond {
		t.Fatalf("hashing unexpectedly fast, duration=%v (cost might have been reduced)", dur)
	}
}

func TestValidateToken_Success(t *testing.T) {
	svc := setupService()
	email := "val@example.com"
	pass := "ValidPass!"
	name := "Val User"
	id, token, err := svc.Register(context.Background(), email, pass, name)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	uid, valid, err := svc.Validate(context.Background(), token)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if !valid {
		t.Fatalf("expected token to be valid")
	}
	if uid != id {
		t.Fatalf("expected userID %s, got %s", id, uid)
	}
}

func TestGetProfile_UpdateProfile_DeleteUser(t *testing.T) {
	svc := setupService()
	email := "prof@example.com"
	pass := "ValidPass!"
	name := "Profile User"
	id, _, err := svc.Register(context.Background(), email, pass, name)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// GetProfile
	user, err := svc.GetProfile(context.Background(), id)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if user.Name != name {
		t.Fatalf("expected name %s, got %s", name, user.Name)
	}

	// UpdateProfile
	newName := "New Name"
	if err := svc.UpdateProfile(context.Background(), id, newName); err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	user2, err := svc.GetProfile(context.Background(), id)
	if err != nil {
		t.Fatalf("GetProfile after update failed: %v", err)
	}
	if user2.Name != newName {
		t.Fatalf("expected updated name %s, got %s", newName, user2.Name)
	}

	// DeleteUser
	if err := svc.DeleteUser(context.Background(), id); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if _, err := svc.GetProfile(context.Background(), id); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}
