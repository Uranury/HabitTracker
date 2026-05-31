package main

import (
	"context"

	"github.com/Uranury/HabitTracker/internal/app"
)

func main() {
	ctx := context.Background()
	_, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()
}
