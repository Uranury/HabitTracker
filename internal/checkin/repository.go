package checkin

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Repository interface {
	Record(ctx context.Context, c *CheckIn) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*CheckIn, error)
	UpdateStatus(ctx context.Context, checkinID uuid.UUID, status Status) error
	GetByUserAndHabitID(ctx context.Context, userID, habitID uuid.UUID) ([]*CheckIn, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Record(ctx context.Context, c *CheckIn) error {
	query := `INSERT INTO checkins (id, user_id, habit_id, status, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.UserID, c.HabitID, c.Status, c.Date, c.CreatedAt.Format(time.RFC3339), c.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) (_ []*CheckIn, err error) {
	query := `SELECT id, user_id, habit_id, status, date, created_at, updated_at FROM checkins WHERE user_id = ?`
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = errors.Wrapf(cerr, "rows.Close()")
		}
	}()
	checkins := []*CheckIn{}
	for rows.Next() {
		var c CheckIn
		var createdAtStr, updatedAtStr string
		if err = rows.Scan(&c.ID, &c.UserID, &c.HabitID, &c.Status, &c.Date, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}
		c.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		c.UpdatedAt, err = time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			return nil, err
		}

		checkins = append(checkins, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return checkins, nil
}

func (r *repository) UpdateStatus(ctx context.Context, checkinID uuid.UUID, status Status) error {
	rows, err := r.db.ExecContext(ctx, "UPDATE checkins SET status = ? WHERE id = ?", status, checkinID)
	if err != nil {
		return err
	}
	n, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *repository) GetByUserAndHabitID(ctx context.Context, userID, habitID uuid.UUID) (_ []*CheckIn, err error) {
	query := `SELECT id, user_id, habit_id, status, date, created_at, updated_at FROM checkins WHERE user_id = ? AND habit_id = ? ORDER BY date DESC`
	rows, err := r.db.QueryxContext(ctx, query, userID, habitID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = errors.Wrapf(cerr, "rows.Close()")
		}
	}()

	checkins := []*CheckIn{}
	for rows.Next() {
		var c CheckIn
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&c.ID, &c.UserID, &c.HabitID, &c.Status, &c.Date, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}
		c.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		c.UpdatedAt, err = time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			return nil, err
		}
		checkins = append(checkins, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return checkins, nil
}
