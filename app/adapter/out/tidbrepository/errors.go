package tidbrepository

import "github.com/cockroachdb/errors"

var (
	ErrTenantNotFound            = errors.New("tenant not found")
	ErrTenantDatabase            = errors.New("tenant database error")
	ErrClientCredentialsNotFound = errors.New("client credentials not found")
	ErrClientCredentialsDatabase = errors.New("client credentials database error")
)
