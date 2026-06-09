package habit

import (
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

type CreateHabitRequest struct {
	Name        string  `json:"name" binding:"required"`
	Schedule    uint8   `json:"schedule" binding:"required"`
	Description *string `json:"description" binding:"optional"`
}

func (h *Handler) CreateHabit(c *gin.Context) {
	var req CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.Create(c.Request.Context(), userID, req.Name, req.Schedule, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
