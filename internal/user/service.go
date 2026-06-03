package user

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Avatar = avatar
	user.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, user)
}
