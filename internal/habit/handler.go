package habit

import (
	"errors"
	"net/http"

	"github.com/Uranury/HabitTracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type CreateHabitRequest struct {
	Name        string  `json:"name" binding:"required"`
	Schedule    uint8   `json:"schedule" binding:"-"`
	Description *string `json:"description" binding:"-"`
	Type        *string `json:"type" binding:"-"`
	Icon        *string `json:"icon" binding:"-"`
}

func (h *Handler) CreateHabit(c *gin.Context) {
	var req CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Schedule == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schedule must have at least one day"})
		return
	}
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.Create(c.Request.Context(), userID, req.Name, req.Schedule, req.Description, req.Type, req.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) ListHabits(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	habits, err := h.svc.ListByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"habits": habits})
}

func (h *Handler) GetHabit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	param := c.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hbt, err := h.svc.GetByID(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, ErrHabitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"habit": hbt})
}

type UpdateHabitRequest struct {
	Name        *string `json:"name" binding:"-"`
	Schedule    *uint8  `json:"schedule" binding:"-"`
	Description *string `json:"description" binding:"-"`
	Type        *string `json:"type" binding:"-"`
	Icon        *string `json:"icon" binding:"-"`
}

func (h *Handler) UpdateHabit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req UpdateHabitRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Schedule != nil && *req.Schedule == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schedule must have at least one day"})
		return
	}
	if err = h.svc.UpdateHabit(c.Request.Context(), userID, habitID, req.Name, req.Schedule, req.Description, req.Type, req.Icon); err != nil {
		if errors.Is(err, ErrHabitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) DeleteHabit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.svc.DeleteHabit(c.Request.Context(), userID, habitID)
	if err != nil {
		if errors.Is(err, ErrHabitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
