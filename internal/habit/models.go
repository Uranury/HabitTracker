package habit

import (
	"github.com/google/uuid"
	"time"
)

type Habit struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Schedule    uint8     `json:"schedule" db:"schedule"`
	Description *string   `json:"description" db:"description"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
