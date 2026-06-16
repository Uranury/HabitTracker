package habit

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository interface {
	Create(ctx context.Context, h *Habit) error
	GetHabitByID(ctx context.Context, userID, habitID uuid.UUID) (*Habit, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, h *Habit) error {
	query := `INSERT INTO habits (id, user_id, name, schedule, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.UserID, h.Name, h.Schedule, h.Description, h.CreatedAt.Format(time.RFC3339), h.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *repository) GetHabitByID(ctx context.Context, userID, habitID uuid.UUID) (*Habit, error) {
	var h Habit
	var createdAtStr, updatedAtStr string
	query := `SELECT id, name, schedule, description, created_at, updated_at FROM habits WHERE user_id = ? AND id = ?`
	err := r.db.QueryRowxContext(ctx, query, userID, habitID).Scan(&h.ID, &h.Name, &h.Schedule, &h.Description, &createdAtStr, &updatedAtStr)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, err
	}
	h.CreatedAt = createdAt
	h.UpdatedAt = updatedAt
	return &h, nil
}
