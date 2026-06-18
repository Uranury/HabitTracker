package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Uranury/HabitTracker/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, u *User) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error
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
	var createdAtStr, updatedAtStr string
	query := "SELECT id, username, password, time_zone, avatar, created_at, updated_at FROM users WHERE id = ?"
	err := r.db.QueryRowxContext(ctx, query, id).Scan(&u.ID, &u.Username, &u.Password, &u.TimeZone, &u.Avatar, &createdAtStr, &updatedAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	err = util.ParseTime(&u, createdAtStr, updatedAtStr)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	var createdAtStr, updatedAtStr string
	query := "SELECT id, username, password, time_zone, avatar, created_at, updated_at FROM users WHERE username = ?"
	err := r.db.QueryRowxContext(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password, &u.TimeZone, &u.Avatar, &createdAtStr, &updatedAtStr)
	if err != nil {
		return nil, err
	}
	err = util.ParseTime(&u, createdAtStr, updatedAtStr)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, username, password, time_zone, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Username, u.Password, u.TimeZone, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *repository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatar string) error {
	query := `UPDATE users SET avatar=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, userID, avatar)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error {
	query := `UPDATE users SET time_zone=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, timezone, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	query := `UPDATE users SET username=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, username, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	query := `UPDATE users SET password=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, password, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
