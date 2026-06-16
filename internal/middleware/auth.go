package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Uranury/HabitTracker/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Auth struct {
	authSvc *auth.TokenService
}

func NewAuth(authSvc *auth.TokenService) *Auth {
	return &Auth{
		authSvc: authSvc,
	}
}

type contextKey string

const (
	userIDKey       contextKey = "user_id"
	userTimeZoneKey contextKey = "user_time_zone"
)

func (m *Auth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.authSvc.Validate(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Set(userTimeZoneKey, claims.TimeZone)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get(string(userIDKey))
	if !exists {
		return uuid.Nil, errors.New("user ID not found")
	}
	uid, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID has invalid type")
	}
	return uid, nil
}

func GetUserTimeZone(c *gin.Context) (timezone string, err error) {
	val, exists := c.Get(string(userTimeZoneKey))
	if !exists {
		return "", errors.New("user time zone not found")
	}
	userTimeZone, ok := val.(string)
	if !ok {
		return "", errors.New("user time zone has invalid type")
	}
	return userTimeZone, nil
}
