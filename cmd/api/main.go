package main

import (
	"context"

	"github.com/Uranury/HabitTracker/internal/app"
	"github.com/Uranury/HabitTracker/internal/auth"
	"github.com/Uranury/HabitTracker/internal/checkin"
	"github.com/Uranury/HabitTracker/internal/habit"
	"github.com/Uranury/HabitTracker/internal/habitgroup"
	"github.com/Uranury/HabitTracker/internal/server"
	"github.com/Uranury/HabitTracker/internal/user"
)

func main() {
	ctx := context.Background()
	infra, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	authHandler := auth.NewHandler(infra.AuthSvc)
	habitHandler := habit.NewHandler(infra.HabitSvc)
	habitGroupHandler := habitgroup.NewHandler(infra.HabitGroupSvc)
	userHandler := user.NewHandler(infra.UserSvc)
	checkinHandler := checkin.NewHandler(infra.CheckinSvc)

	serv := server.NewServer(infra.Middlw, authHandler, habitHandler, habitGroupHandler, userHandler, checkinHandler)
	if err := serv.Run(); err != nil {
		panic(err)
	}
}
