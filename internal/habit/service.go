package habit

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Service) ListByID(ctx context.Context, userID uuid.UUID) ([]*Habit, error) {
	habits, err := s.repo.GetHabitsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return []*Habit{}, nil
	}
	return habits, nil
}

func (s *Service) GetByID(ctx context.Context, userID, habitID uuid.UUID) (*Habit, error) {
	habit, err := s.repo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return habit, nil
}

func (s *Service) UpdateHabit(ctx context.Context, userID, habitID uuid.UUID, name *string, schedule *uint8, description *string) error {
	hbt := &Habit{
		ID:          habitID,
		UserID:      userID,
		Description: description,
	}
	if name != nil {
		hbt.Name = *name
	}
	if schedule != nil {
		hbt.Schedule = *schedule
	}
	return s.repo.Update(ctx, hbt)
}

func (s *Service) DeleteHabit(ctx context.Context, userID, habitID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, habitID)
}
