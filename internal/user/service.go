package user

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	return s.repo.FindByID(ctx, id)
}

func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	return s.repo.UpdateAvatar(ctx, userID, avatar)
}

func (s *Service) UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error {
	return s.repo.UpdateTimezone(ctx, userID, timezone)
}

func (s *Service) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	return s.repo.UpdateUsername(ctx, userID, username)
}

func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error {
	u, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPass))
	if err != nil {
		return ErrInvalidPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, string(hashed))
}
