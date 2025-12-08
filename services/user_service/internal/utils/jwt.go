package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTKey() []byte {
	if env := os.Getenv("JWT_SECRET"); env != "" {
		return []byte(env)
	}
	panic("JWT_SECRET environment variable is not set. Application cannot start without a secure JWT secret.")
}

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTKey())
}

// ValidateToken parses and validates JWT and returns (userID, valid, error).
// valid=false with nil error means token корректно разобран, но не валиден (например, истёк).
func ValidateToken(tokenStr string) (string, bool, error) {
	if tokenStr == "" {
		return "", false, errors.New("empty token")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTKey(), nil
	})
	if err != nil {
		// например, истёк срок действия
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", false, nil
		}
		return "", false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", false, nil
	}

	uid, _ := claims["user_id"].(string)
	if uid == "" {
		return "", false, errors.New("user_id missing in token")
	}

	return uid, true, nil
}

func GetJWTKey() []byte {
	return getJWTKey()
}
