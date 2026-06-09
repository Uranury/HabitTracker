package server

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	midlw      *middleware.Auth
}

func NewServer(middlw *middleware.Auth) *Server {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger())

	server := &Server{
		router: router,
		httpServer: &http.Server{
			Addr:         ":8080",
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		midlw: middlw,
	}
	return server
}

func (s *Server) setupRoutes() {
	protected := s.router.Group("/protected", s.midlw.JWTAuth())
	protected.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}

func (s *Server) Run() error {
	s.setupRoutes()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
