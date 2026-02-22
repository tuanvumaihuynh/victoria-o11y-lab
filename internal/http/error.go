package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/apperr"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/dto"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/pkg/zerror"
)

var internalServerErrResponse = &humaErrorResponse{
	ErrorResponse: dto.ErrorResponse{
		Code:    apperr.InternalServerErr.Code(),
		Message: apperr.InternalServerErr.Msg(),
	},
	statusCode: http.StatusInternalServerError,
}

func newHumaError(logger *slog.Logger) func(status int, message string, errs ...error) huma.StatusError {
	return func(status int, message string, errs ...error) huma.StatusError {
		if len(errs) == 0 {
			return internalServerErrResponse
		}

		// Huma returns multiple errors only for validation failures.
		// If huma behavior changes, revisit this logic.
		// https://github.com/danielgtaylor/huma/blob/887f7d43222686b060805a934ab33a417b44e2fc/huma.go#L1071-L1085
		if len(errs) > 1 {
			return validationErrorsToErrorResponse(errs)
		}

		// Only handle the first error
		err := errs[0]
		errResp := errorsToErrorResponse(err)
		if errResp.GetStatus() >= 500 {
			logger.Error("handler error", slog.Any("error", err))
		}

		return errResp
	}
}

func newHumaErrorWithContext(logger *slog.Logger) func(hctx huma.Context, status int, message string, errs ...error) huma.StatusError {
	return func(hctx huma.Context, status int, message string, errs ...error) huma.StatusError {
		ctx := hctx.Context()

		if len(errs) == 0 {
			if status != 0 {
				logger.ErrorContext(
					ctx,
					"huma error handler called with no errors",
					slog.Int("status", status),
					slog.String("message", message),
				)
			}
			return internalServerErrResponse
		}

		// Huma returns multiple errors only for validation failures.
		// If huma behavior changes, revisit this logic.
		// https://github.com/danielgtaylor/huma/blob/887f7d43222686b060805a934ab33a417b44e2fc/huma.go#L1071-L1085
		if len(errs) > 1 {
			return validationErrorsToErrorResponse(errs)
		}

		// Only handle the first error
		err := errs[0]
		errResp := errorsToErrorResponse(err)
		if errResp.GetStatus() >= 500 {
			logger.ErrorContext(ctx, "handler error", slog.Any("error", err))
		}

		return errResp
	}
}

func validationErrorsToErrorResponse(errs []error) *humaErrorResponse {
	convertFunc := func(err error) *dto.ErrorDetail {
		humaErr, ok := errors.AsType[*huma.ErrorDetail](err)
		if ok {
			return &dto.ErrorDetail{
				Field:   humaErr.Location,
				Message: humaErr.Message,
			}
		}
		return nil
	}

	details := make([]dto.ErrorDetail, 0, len(errs))

	for i := len(errs) - 1; i >= 0; i-- {
		detail := convertFunc(errs[i])
		if detail != nil {
			details = append(details, *detail)
		}
	}

	return &humaErrorResponse{
		ErrorResponse: dto.ErrorResponse{
			Code:         apperr.ValidationError.Code(),
			Message:      apperr.ValidationError.Msg(),
			ErrorDetails: details,
		},
		statusCode: http.StatusUnprocessableEntity,
	}
}

func errorsToErrorResponse(err error) *humaErrorResponse {
	zErr, ok := errors.AsType[*zerror.ZError](err)
	if ok {
		return &humaErrorResponse{
			ErrorResponse: dto.ErrorResponse{
				Code:    zErr.Code(),
				Message: zErr.Msg(),
			},
			statusCode: zErrorStatusToHTTPStatus(zErr.Status()),
		}
	}

	// Unhandled error always return internal server error
	return &humaErrorResponse{
		ErrorResponse: dto.ErrorResponse{
			Code:    apperr.InternalServerErr.Code(),
			Message: apperr.InternalServerErr.Msg(),
		},
		statusCode: http.StatusInternalServerError,
	}
}

var _ huma.StatusError = (*humaErrorResponse)(nil)

type humaErrorResponse struct {
	dto.ErrorResponse
	statusCode int `json:"-"`
}

func (e *humaErrorResponse) Error() string {
	return fmt.Sprintf("code=%s, message=%s", e.Code, e.Message)
}

func (e *humaErrorResponse) GetStatus() int {
	return e.statusCode
}

func zErrorStatusToHTTPStatus(status zerror.Status) int {
	switch status {
	case zerror.StatusUnauthorized:
		return http.StatusUnauthorized
	case zerror.StatusForbidden:
		return http.StatusForbidden
	case zerror.StatusNotFound:
		return http.StatusNotFound
	case zerror.StatusUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case zerror.StatusConflict:
		return http.StatusConflict
	case zerror.StatusTooManyRequests:
		return http.StatusTooManyRequests
	case zerror.StatusBadRequest:
		return http.StatusBadRequest
	case zerror.StatusValidationFailed:
		return http.StatusBadRequest
	case zerror.StatusUnknown, zerror.StatusInternalServerError:
		return http.StatusInternalServerError
	case zerror.StatusTimeout:
		return http.StatusGatewayTimeout
	case zerror.StatusNotImplemented:
		return http.StatusNotImplemented
	case zerror.StatusBadGateway:
		return http.StatusBadGateway
	case zerror.StatusServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
