package habitgroup

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, name string, icon *string) error {
	now := time.Now().UTC()
	g := &HabitGroup{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Icon:      icon,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, g)
}

func (s *Service) GetByID(ctx context.Context, userID, groupID uuid.UUID) (*HabitGroup, error) {
	return s.repo.GetByID(ctx, userID, groupID)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]*HabitGroup, error) {
	groups, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return []*HabitGroup{}, nil
	}
	return groups, nil
}

func (s *Service) Update(ctx context.Context, userID, groupID uuid.UUID, name *string, icon *string) error {
	existing, err := s.repo.GetByID(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if name != nil {
		existing.Name = *name
	}
	if icon != nil {
		existing.Icon = icon
	}
	existing.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}

func (s *Service) Delete(ctx context.Context, userID, groupID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, groupID)
}
