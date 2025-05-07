package apperrors

import "github.com/cockroachdb/errors"

type alertableError struct {
	error
}

type retryableError struct {
	error
}

func MarkAsAlertable(err error) error {
	return &alertableError{error: err}
}

func MarkAsRetryable(err error) error {
	return &retryableError{error: err}
}

func MarkAsAlertableAndRetryable(err error) error {
	return MarkAsAlertable(MarkAsRetryable(err))
}

func IsAlertable(err error) bool {
	var ae *alertableError
	return errors.As(err, &ae)
}

func IsRetryable(err error) bool {
	var re *retryableError
	return errors.As(err, &re)
}
