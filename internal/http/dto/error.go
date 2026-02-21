package dto

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/apperr"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/pkg/zerror"
)

var _ huma.StatusError = (*ErrorResponse)(nil)

var (
	InternalServerErrResponse = &ErrorResponse{
		Code:       apperr.InternalServerErr.MsgID(),
		Message:    apperr.InternalServerErr.Msg(),
		statusCode: http.StatusInternalServerError,
	}
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`

	statusCode int `json:"-"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return errorToErrorResponse(err)
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code=%s, message=%s", e.Code, e.Message)
}

func (e *ErrorResponse) GetStatus() int {
	return e.statusCode
}

func errorToErrorResponse(err error) *ErrorResponse {
	zErr, ok := errors.AsType[*zerror.ZError](err)
	if ok {
		return &ErrorResponse{
			Code:       zErr.MsgID(),
			Message:    zErr.Msg(),
			statusCode: zErrorStatusToHTTPStatus(zErr.Status()),
		}
	}

	return InternalServerErrResponse
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
