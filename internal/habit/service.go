package habit

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, name string, schedule uint8, description *string) error {
	hbt := &Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Schedule:    schedule,
		Description: description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return s.repo.Create(ctx, hbt)
}
