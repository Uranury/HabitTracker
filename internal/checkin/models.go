package checkin

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	Checked Status = iota
	Skipped
	Missed
)

type CheckIn struct {
	ID      uuid.UUID `json:"id" db:"id"`
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	HabitID uuid.UUID `json:"habit_id" db:"habit_id"`
	Status  Status    `json:"status" db:"status"`
	Date    time.Time `json:"date" db:"date"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
