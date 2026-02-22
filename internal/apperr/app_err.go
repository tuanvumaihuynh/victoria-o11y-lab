package apperr

import "github.com/tuanvumaihuynh/victoria-o11y-lab/pkg/zerror"

var (
	InternalServerErr = zerror.NewInternalServerError("internal_server_error", "Internal server error")
	ValidationError   = zerror.NewValidationFailed("validation_failed", "Validation failed")
)
