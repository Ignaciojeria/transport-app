package usecase

import "github.com/cockroachdb/errors"

var (
	ErrOrganizationAlreadyExists = errors.New("organization already exists")
)
