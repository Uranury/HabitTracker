package habitgroup

import (
	"context"
	"database/sql"
	"time"

	"github.com/Uranury/HabitTracker/pkg/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	Create(ctx context.Context, g *HabitGroup) error
	GetByID(ctx context.Context, userID, groupID uuid.UUID) (*HabitGroup, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*HabitGroup, error)
	Update(ctx context.Context, g *HabitGroup) error
	Delete(ctx context.Context, userID, groupID uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, g *HabitGroup) error {
	query := `INSERT INTO habit_groups (id, user_id, name, icon, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, g.ID, g.UserID, g.Name, g.Icon, g.CreatedAt.Format(time.RFC3339), g.UpdatedAt.Format(time.RFC3339))
	return err
}

func (r *repository) GetByID(ctx context.Context, userID, groupID uuid.UUID) (*HabitGroup, error) {
	var g HabitGroup
	var createdAtStr, updatedAtStr string
	query := `SELECT id, user_id, name, icon, created_at, updated_at FROM habit_groups WHERE id = ? AND user_id = ?`
	err := r.db.QueryRowxContext(ctx, query, groupID, userID).Scan(&g.ID, &g.UserID, &g.Name, &g.Icon, &createdAtStr, &updatedAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	if err = util.ParseTime(&g, createdAtStr, updatedAtStr); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) (_ []*HabitGroup, err error) {
	query := `SELECT id, user_id, name, icon, created_at, updated_at FROM habit_groups WHERE user_id = ?`
	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = errors.Wrap(cerr, "rows.Close()")
		}
	}()

	var groups []*HabitGroup
	for rows.Next() {
		var g HabitGroup
		var createdAtStr, updatedAtStr string
		if err = rows.Scan(&g.ID, &g.UserID, &g.Name, &g.Icon, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}
		if err = util.ParseTime(&g, createdAtStr, updatedAtStr); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *repository) Update(ctx context.Context, g *HabitGroup) error {
	query := `UPDATE habit_groups SET name = ?, icon = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, query, g.Name, g.Icon, g.UpdatedAt.Format(time.RFC3339), g.ID, g.UserID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, userID, groupID uuid.UUID) error {
	query := `DELETE FROM habit_groups WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, query, groupID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}
