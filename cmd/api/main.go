package main

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/auth"

	"github.com/Uranury/HabitTracker/internal/app"
	"github.com/Uranury/HabitTracker/internal/server"
)

func main() {
	ctx := context.Background()
	infra, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	authHandler := auth.NewHandler(infra.AuthSvc)
	defer cleanup()
	serv := server.NewServer(infra.Middlw, authHandler)
	if err := serv.Run(); err != nil {
		panic(err)
	}
}
