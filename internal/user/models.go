package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	Avatar   *string   `json:"avatar" db:"avatar"`
	TimeZone string    `json:"time_zone" db:"time_zone"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) SetTimes(createdAt, updatedAt time.Time) {
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt
}
