package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TokenValidator is satisfied by auth.TokenService. Defined here so middleware
// does not import the auth package, avoiding a dependency cycle with user.
type TokenValidator interface {
	Validate(tokenString string) (userID uuid.UUID, timeZone string, err error)
}

type Auth struct {
	validator TokenValidator
}

func NewAuth(v TokenValidator) *Auth {
	return &Auth{validator: v}
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
		userID, timeZone, err := m.validator.Validate(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(string(userIDKey), userID)
		c.Set(string(userTimeZoneKey), timeZone)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get(string(userIDKey))
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	uid, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID has invalid type")
	}
	return uid, nil
}

func GetUserTimeZone(c *gin.Context) (string, error) {
	val, exists := c.Get(string(userTimeZoneKey))
	if !exists {
		return "", errors.New("user time zone not found in context")
	}
	tz, ok := val.(string)
	if !ok {
		return "", errors.New("user time zone has invalid type")
	}
	return tz, nil
}
