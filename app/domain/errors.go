package domain

import "github.com/joomcode/errorx"

// Namespace y clases de errores específicos para validación de órdenes
var (
	orderValidationErrorNamespace = errorx.NewNamespace("order_validation")
	ErrInvalidDateFormat          = orderValidationErrorNamespace.NewType("invalid_date_format")
	ErrInvalidTimeFormat          = orderValidationErrorNamespace.NewType("invalid_time_format")
)
