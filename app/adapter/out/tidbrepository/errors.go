package tidbrepository

import "github.com/cockroachdb/errors"

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationDatabase = errors.New("organization database error")
)
