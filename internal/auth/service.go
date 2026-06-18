package auth

import (
	"context"
	"errors"
	"github.com/Uranury/HabitTracker/internal/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStore interface {
}

type Service struct {
	userRepo user.Repository
	tokenSvc *TokenService
}

func NewService(userRepo user.Repository, tokenSvc *TokenService) *Service {
	return &Service{
		userRepo: userRepo,
		tokenSvc: tokenSvc,
	}
}

func (s *Service) Signup(ctx context.Context, username, password, timezone string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	u := &user.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(bcryptHash),
		TimeZone:  timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	token, err := s.tokenSvc.Generate(u.ID, timezone)
	if err != nil {
		return "", err
	}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (string, error) {
	usr, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if usr == nil {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := s.tokenSvc.Generate(usr.ID, usr.TimeZone)
	if err != nil {
		return "", err
	}
	return token, nil
}
