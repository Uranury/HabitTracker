package habit

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, h *Habit) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, h *Habit) error {
	query := `INSERT INTO habits (id, user_id, name, schedule, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.UserID, h.Name, h.Schedule, h.CreatedAt, h.UpdatedAt)
	return err
}
