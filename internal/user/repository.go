package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, username, password, time_zone, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.Password, u.TimeZone, u.Avatar, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET username=?, password=?, time_zone=?, avatar=?, updated_at=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, u.Username, u.Password, u.TimeZone, u.Avatar, u.UpdatedAt, u.ID)
	return err
}
