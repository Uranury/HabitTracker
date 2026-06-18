package user

import (
	"errors"
	"github.com/Uranury/HabitTracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(c *gin.Context) {
	id, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": u})
}

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	id, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req UpdateAvatarRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.svc.UploadAvatar(c.Request.Context(), id, req.Avatar)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

type UpdateTimezoneRequest struct {
	Timezone string `json:"timezone" binding:"required"`
}

func (h *Handler) UpdateTimezone(c *gin.Context) {
	id, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req UpdateTimezoneRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.svc.UpdateTimezone(c.Request.Context(), id, req.Timezone)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

func (h *Handler) UpdateUsername(c *gin.Context) {
	id, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req UpdateUsernameRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.svc.UpdateUsername(c.Request.Context(), id, req.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	id, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req UpdatePasswordRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old password and new password must be provided"})
		return
	}
	err = h.svc.UpdatePassword(c.Request.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
