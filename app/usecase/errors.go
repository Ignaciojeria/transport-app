package usecase

import "github.com/joomcode/errorx"

// Namespace y clase de errores específicos para casos de uso relacionados con organizaciones
var (
	organizationErrorNamespace   = errorx.NewNamespace("organization")
	ErrOrganizationAlreadyExists = organizationErrorNamespace.NewType("already_exists")
)
