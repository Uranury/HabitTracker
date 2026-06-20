package main

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/checkin"
	"github.com/Uranury/HabitTracker/internal/user"

	"github.com/Uranury/HabitTracker/internal/app"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/server"
)

func main() {
	ctx := context.Background()
	infra, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	authHandler := auth.NewHandler(infra.AuthSvc)
	habitHandler := habit.NewHandler(infra.HabitSvc)
	userHandler := user.NewHandler(infra.UserSvc)
	checkinHandler := checkin.NewHandler(infra.CheckinSvc)
	defer cleanup()
	serv := server.NewServer(infra.Middlw, authHandler, habitHandler, userHandler, checkinHandler)
	if err := serv.Run(); err != nil {
		panic(err)
	}
}
