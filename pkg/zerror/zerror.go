package zerror

import (
	"fmt"
)

// ZError represents the base error structure.
// Use [NewZError] to create a new instance of ZError.
type ZError struct {
	parent error
	status Status
	code   string
	msg    string
}

// NewZError initializes a ZError instance.
func NewZError(parent error, status Status, code, msg string) *ZError {
	return &ZError{
		parent: parent,
		status: status,
		code:   code,
		msg:    msg,
	}
}

// WithParent attaches an underlying error to an existing predefined XError.
func WithParent(e ZError, parent error) ZError {
	if parent == nil {
		return e
	}
	e.parent = parent
	return e
}

// WithMsg creates a new XError with a modified message while preserving status, code, and parent.
func WithMsg(e ZError, msg string) ZError {
	e.msg = msg
	return e
}

func (e ZError) Error() string {
	if e.parent != nil {
		return fmt.Sprintf("Code=%s, Msg=%s, Parent=(%v)", e.code, e.msg, e.parent)
	}
	return fmt.Sprintf("Code=%s, Msg=%s", e.code, e.msg)
}

func (e *ZError) Unwrap() error {
	return e.parent
}

func (e ZError) Status() Status {
	return e.status
}

func (e ZError) Code() string {
	return e.code
}

func (e ZError) Msg() string {
	return e.msg
}

func (e ZError) Parent() error {
	return e.parent
}

func NewUnauthorized(code, msg string) *ZError {
	return NewZError(nil, StatusUnauthorized, code, msg)
}

func NewForbidden(code, msg string) *ZError {
	return NewZError(nil, StatusForbidden, code, msg)
}

func NewNotFound(code, msg string) *ZError {
	return NewZError(nil, StatusNotFound, code, msg)
}

func NewUnprocessableEntity(code, msg string) *ZError {
	return NewZError(nil, StatusUnprocessableEntity, code, msg)
}

func NewConflict(code, msg string) *ZError {
	return NewZError(nil, StatusConflict, code, msg)
}

func NewTooManyRequests(code, msg string) *ZError {
	return NewZError(nil, StatusTooManyRequests, code, msg)
}

func NewBadRequest(code, msg string) *ZError {
	return NewZError(nil, StatusBadRequest, code, msg)
}

func NewValidationFailed(code, msg string) *ZError {
	return NewZError(nil, StatusValidationFailed, code, msg)
}

func NewInternalServerError(code, msg string) *ZError {
	return NewZError(nil, StatusInternalServerError, code, msg)
}

func NewTimeout(code, msg string) *ZError {
	return NewZError(nil, StatusTimeout, code, msg)
}

func NewNotImplemented(code, msg string) *ZError {
	return NewZError(nil, StatusNotImplemented, code, msg)
}

func NewBadGateway(code, msg string) *ZError {
	return NewZError(nil, StatusBadGateway, code, msg)
}

func NewServiceUnavailable(code, msg string) *ZError {
	return NewZError(nil, StatusServiceUnavailable, code, msg)
}
