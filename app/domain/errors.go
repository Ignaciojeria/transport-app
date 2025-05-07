package domain

import "github.com/cockroachdb/errors"

var (
	ErrInvalidDateFormat    = errors.New("invalid date format")
	ErrInvalidTimeFormat    = errors.New("invalid time format")
	ErrInvalidPackageFormat = errors.New("invalid package format")
)
