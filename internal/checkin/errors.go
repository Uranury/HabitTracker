package checkin

import "errors"

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrCheckinNotFound = errors.New("checkin not found")
)
