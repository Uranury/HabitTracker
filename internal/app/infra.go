package app

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/Uranury/HabitTracker/internal/checkin"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/middleware"
	"github.com/Uranury/HabitTracker/internal/user"
	"log/slog"
	"os"

	"github.com/Uranury/HabitTracker/pkg/config"
	"github.com/Uranury/HabitTracker/pkg/database"

	"github.com/jmoiron/sqlx"
)

type Infra struct {
	DBConn     *sqlx.DB
	Config     *config.Config
	Logger     *slog.Logger
	UserSvc    *user.Service
	AuthSvc    *auth.Service
	HabitSvc   *habit.Service
	CheckinSvc *checkin.Service
	Middlw     *middleware.Auth
}

func New(ctx context.Context) (*Infra, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler).With("app", "habitTracker")

	if err := database.RunMigrations(cfg.MigrationsPath, cfg.Database.Driver, cfg.Database.DSN()); err != nil {
		return nil, nil, err
	}

	dbConn, err := database.InitDB(ctx, cfg.Database.Driver, cfg.Database.DSN())
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewRepository(dbConn)
	userSvc := user.NewService(userRepo)

	tokenSvc := auth.NewTokenService([]byte(cfg.JWTSecret))
	authSvc := auth.NewService(userRepo, tokenSvc)
	habitRepo := habit.NewRepository(dbConn)
	habitSvc := habit.NewService(habitRepo)
	checkinRepo := checkin.NewRepository(dbConn)
	checkinSvc := checkin.NewService(checkinRepo, habitRepo)

	middlw := middleware.NewAuth(tokenSvc)

	infra := &Infra{
		DBConn:     dbConn,
		Config:     cfg,
		Logger:     logger,
		UserSvc:    userSvc,
		AuthSvc:    authSvc,
		HabitSvc:   habitSvc,
		CheckinSvc: checkinSvc,
		Middlw:     middlw,
	}

	cleanup := func() {
		if err := dbConn.Close(); err != nil {
			logger.Warn("failed to close database connection", "error", err)
		}
		logger.Info("infra closed up")
	}
	return infra, cleanup, nil
}
