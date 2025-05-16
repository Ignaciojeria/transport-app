package tidbrepository

import "github.com/cockroachdb/errors"

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantDatabase = errors.New("tenant database error")
)
