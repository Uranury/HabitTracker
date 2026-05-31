package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Uranury/HabitTracker/pkg/config"
	"github.com/Uranury/HabitTracker/pkg/database"

	"github.com/jmoiron/sqlx"
)

type Infra struct {
	DBConn *sqlx.DB
	Config *config.Config
	Logger *slog.Logger
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

	infra := &Infra{
		DBConn: dbConn,
		Config: cfg,
		Logger: logger,
	}

	cleanup := func() {
		if err := dbConn.Close(); err != nil {
			logger.Warn("failed to close database connection", "error", err)
		}
		logger.Info("infra closed up")
	}
	return infra, cleanup, nil
}
