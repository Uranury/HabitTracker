package habitgroup

import (
	"time"

	"github.com/google/uuid"
)

type HabitGroup struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Icon      *string   `json:"icon" db:"icon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (g *HabitGroup) SetTimes(createdAt, updatedAt time.Time) {
	g.CreatedAt = createdAt
	g.UpdatedAt = updatedAt
}
