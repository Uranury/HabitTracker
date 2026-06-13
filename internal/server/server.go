package server

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/Uranury/HabitTracker/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	httpServer  *http.Server
	midlw       *middleware.Auth
	authHandler *auth.Handler
}

func NewServer(middlw *middleware.Auth, authHandler *auth.Handler) *Server {
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
		midlw:       middlw,
		authHandler: authHandler,
	}
	return server
}

func (s *Server) setupRoutes() {
	authGroup := s.router.Group("/auth")
	authGroup.POST("/signup", s.authHandler.Signup)
	authGroup.POST("/login", s.authHandler.Login)
}

func (s *Server) Run() error {
	s.setupRoutes()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
