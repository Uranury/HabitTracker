package server

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/Uranury/HabitTracker/internal/checkin"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/middleware"
	"github.com/Uranury/HabitTracker/internal/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	midlw          *middleware.Auth
	authHandler    *auth.Handler
	habitHandler   *habit.Handler
	userHandler    *user.Handler
	checkinHandler *checkin.Handler
}

func NewServer(middlw *middleware.Auth, authHandler *auth.Handler, habitHandler *habit.Handler, userHandler *user.Handler, checkinHandler *checkin.Handler) *Server {
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
		midlw:          middlw,
		authHandler:    authHandler,
		habitHandler:   habitHandler,
		userHandler:    userHandler,
		checkinHandler: checkinHandler,
	}
	return server
}

func (s *Server) setupRoutes() {
	s.router.StaticFile("/", "./web/index.html")

	authGroup := s.router.Group("/auth")
	authGroup.POST("/signup", s.authHandler.Signup)
	authGroup.POST("/login", s.authHandler.Login)

	api := s.router.Group("/api")
	api.Use(s.midlw.JWTAuth())
	{
		habitsGroup := api.Group("/habits")
		{
			habitsGroup.POST("", s.habitHandler.CreateHabit)
			habitsGroup.PATCH("/:id", s.habitHandler.UpdateHabit)
			habitsGroup.GET("/:id", s.habitHandler.GetHabit)
			habitsGroup.GET("", s.habitHandler.ListHabits)
			habitsGroup.DELETE("/:id", s.habitHandler.DeleteHabit)

			habitsGroup.POST("/:id/checkin", s.checkinHandler.CheckIn)
			habitsGroup.GET("/:id/checkins", s.checkinHandler.GetCheckins)
			habitsGroup.GET("/:id/streak", s.checkinHandler.GetStreak)
		}

		usersGroup := api.Group("/users")
		{
			usersGroup.GET("/me", s.userHandler.GetProfile)
			usersGroup.PATCH("/me/avatar", s.userHandler.UpdateAvatar)
			usersGroup.POST("/me/avatar", s.userHandler.UpdateAvatar)
			usersGroup.PATCH("/me/password", s.userHandler.UpdatePassword)
			usersGroup.PATCH("/me/username", s.userHandler.UpdateUsername)
			usersGroup.PATCH("/me/timezone", s.userHandler.UpdateTimezone)
		}
	}
}

func (s *Server) Run() error {
	s.setupRoutes()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
