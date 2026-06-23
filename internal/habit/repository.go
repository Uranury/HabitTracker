package habit

import (
	"context"
	"database/sql"

	"github.com/Uranury/HabitTracker/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Repository interface {
	Create(ctx context.Context, h *Habit) error
	Update(ctx context.Context, h *Habit) error
	Delete(ctx context.Context, userID, habitID uuid.UUID) error
	GetHabitByID(ctx context.Context, userID, habitID uuid.UUID) (*Habit, error)
	GetHabitsByUserID(ctx context.Context, userID uuid.UUID) ([]*Habit, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func groupIDVal(id *uuid.UUID) interface{} {
	if id == nil {
		return nil
	}
	return id.String()
}

func scanHabit(scan func(...any) error) (*Habit, error) {
	var h Habit
	var createdAtStr, updatedAtStr string
	var groupIDStr sql.NullString
	if err := scan(&h.ID, &h.UserID, &h.Name, &h.Schedule, &h.Description, &groupIDStr, &h.Icon, &createdAtStr, &updatedAtStr); err != nil {
		return nil, err
	}
	if groupIDStr.Valid {
		id, err := uuid.Parse(groupIDStr.String)
		if err != nil {
			return nil, err
		}
		h.GroupID = &id
	}
	if err := util.ParseTime(&h, createdAtStr, updatedAtStr); err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *repository) Create(ctx context.Context, h *Habit) error {
	query := `INSERT INTO habits (id, user_id, name, schedule, description, group_id, icon, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.UserID, h.Name, h.Schedule, h.Description, groupIDVal(h.GroupID), h.Icon, h.CreatedAt.Format(time.RFC3339), h.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *repository) Update(ctx context.Context, h *Habit) error {
	query := `UPDATE habits SET name = ?, schedule = ?, description = ?, group_id = ?, icon = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, query, h.Name, h.Schedule, h.Description, groupIDVal(h.GroupID), h.Icon, h.UpdatedAt.Format(time.RFC3339), h.ID, h.UserID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrHabitNotFound
	}
	return nil
}

func (r *repository) GetHabitByID(ctx context.Context, userID, habitID uuid.UUID) (*Habit, error) {
	query := `SELECT id, user_id, name, schedule, description, group_id, icon, created_at, updated_at FROM habits WHERE user_id = ? AND id = ?`
	row := r.db.QueryRowxContext(ctx, query, userID, habitID)
	h, err := scanHabit(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHabitNotFound
		}
		return nil, err
	}
	return h, nil
}

func (r *repository) GetHabitsByUserID(ctx context.Context, userID uuid.UUID) (_ []*Habit, err error) {
	query := `SELECT id, user_id, name, schedule, description, group_id, icon, created_at, updated_at FROM habits WHERE user_id = ?`
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = errors.Wrap(cerr, "rows.Close()")
		}
	}()

	var habits []*Habit
	for rows.Next() {
		h, err := scanHabit(rows.Scan)
		if err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *repository) Delete(ctx context.Context, userID, habitID uuid.UUID) error {
	query := `DELETE FROM habits WHERE user_id = ? AND id = ?`
	res, err := r.db.ExecContext(ctx, query, userID, habitID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrHabitNotFound
	}
	return nil
}
