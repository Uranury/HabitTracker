package util

import (
	"time"
)

type Timestamped interface {
	SetTimes(createdAt, updatedAt time.Time)
}

func ParseTime(obj Timestamped, createdAtStr, updatedAtStr string) error {
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return err
	}
	obj.SetTimes(createdAt, updatedAt)
	return nil
}
