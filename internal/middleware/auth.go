package middleware

import (
	"errors"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Auth struct {
	authSvc *auth.TokenService
}

type contextKey string

const (
	userIDKey contextKey = "user_id"
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
