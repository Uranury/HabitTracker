package checkin

import "errors"

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrCheckinNotFound = errors.New("checkin not found")
	ErrFutureDate      = errors.New("check-in date cannot be in the future")
	ErrInvalidDate     = errors.New("invalid date format, expected YYYY-MM-DD")
)
