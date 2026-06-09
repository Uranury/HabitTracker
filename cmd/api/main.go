package main

import (
	"context"
	"github.com/Uranury/HabitTracker/internal/server"

	"github.com/Uranury/HabitTracker/internal/app"
)

func main() {
	ctx := context.Background()
	infra, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	serv := server.NewServer(infra.Middlw)
	if err := serv.Run(); err != nil {
		panic(err)
	}
}
