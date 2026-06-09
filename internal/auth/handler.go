package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	authSvc *Service
}

func NewHandler(authSvc *Service) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TimeZone string `json:"time_zone"`
}

func (h *Handler) Signup(c *gin.Context) {
	var signupRequest SignupRequest
	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authSvc.Signup(c.Request.Context(), signupRequest.Username, signupRequest.Password, signupRequest.TimeZone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
