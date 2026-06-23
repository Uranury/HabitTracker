package habit

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, name string, schedule uint8, description *string, groupID *uuid.UUID, icon *string) error {
	now := time.Now().UTC()
	hbt := &Habit{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Schedule:    schedule,
		Description: description,
		GroupID:     groupID,
		Icon:        icon,
		CreatedAt:   now,
		UpdatedAt:   now,
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
			return nil, ErrHabitNotFound
		}
		return nil, err
	}
	return habit, nil
}

func (s *Service) UpdateHabit(ctx context.Context, userID, habitID uuid.UUID, name *string, schedule *uint8, description *string, groupID *uuid.UUID, icon *string) error {
	existing, err := s.repo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return err
	}
	if name != nil {
		existing.Name = *name
	}
	if schedule != nil {
		existing.Schedule = *schedule
	}
	if description != nil {
		existing.Description = description
	}
	if groupID != nil {
		existing.GroupID = groupID
	}
	if icon != nil {
		existing.Icon = icon
	}
	existing.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}

func (s *Service) DeleteHabit(ctx context.Context, userID, habitID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, habitID)
}
