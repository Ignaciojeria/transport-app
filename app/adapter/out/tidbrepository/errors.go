package tidbrepository

import "github.com/joomcode/errorx"

// Namespace y clases de error específicos para organizaciones
var (
	organizationErrorNamespace = errorx.NewNamespace("organization")
	ErrOrganizationNotFound    = organizationErrorNamespace.NewType("not_found")
	ErrOrganizationDatabase    = organizationErrorNamespace.NewType("database_error")
)
