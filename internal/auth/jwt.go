package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var AccessTokenTTL = time.Hour * 24

type TokenService struct {
	jwtKey []byte
}

func NewTokenService(jwtKey []byte) *TokenService {
	return &TokenService{jwtKey: jwtKey}
}

type Claims struct {
	UserID   uuid.UUID `json:"id"`
	TimeZone string    `json:"time_zone"`
	jwt.RegisteredClaims
}

func (t *TokenService) Generate(userID uuid.UUID, timeZone string) (string, error) {
	claims := Claims{
		UserID:   userID,
		TimeZone: timeZone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (t *TokenService) Validate(tokenString string) (uuid.UUID, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.jwtKey, nil
	})
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to parse token: %w", err)
	}
	if token == nil || !token.Valid {
		return uuid.Nil, "", fmt.Errorf("invalid token")
	}
	return claims.UserID, claims.TimeZone, nil
}
