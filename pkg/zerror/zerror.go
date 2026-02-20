package zerror

import (
	"fmt"
)

// ZError represents the base error structure.
// Use [NewZError] to create a new instance of ZError.
type ZError struct {
	parent error
	status Status
	msgID  string
	msg    string
}

// NewZError initializes a ZError instance.
func NewZError(parent error, status Status, msgID, msg string) *ZError {
	return &ZError{
		parent: parent,
		status: status,
		msgID:  msgID,
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
		return fmt.Sprintf("ID=%s, Msg=%s, Parent=(%v)", e.msgID, e.msg, e.parent)
	}
	return fmt.Sprintf("ID=%s, Msg=%s", e.msgID, e.msg)
}

func (e *ZError) Unwrap() error {
	return e.parent
}

func (e ZError) Status() Status {
	return e.status
}

func (e ZError) MsgID() string {
	return e.msgID
}

func (e ZError) Msg() string {
	return e.msg
}

func (e ZError) Parent() error {
	return e.parent
}

func NewUnauthorized(msgID, msg string) *ZError {
	return NewZError(nil, StatusUnauthorized, msgID, msg)
}

func NewForbidden(msgID, msg string) *ZError {
	return NewZError(nil, StatusForbidden, msgID, msg)
}

func NewNotFound(msgID, msg string) *ZError {
	return NewZError(nil, StatusNotFound, msgID, msg)
}

func NewUnprocessableEntity(msgID, msg string) *ZError {
	return NewZError(nil, StatusUnprocessableEntity, msgID, msg)
}

func NewConflict(msgID, msg string) *ZError {
	return NewZError(nil, StatusConflict, msgID, msg)
}

func NewTooManyRequests(msgID, msg string) *ZError {
	return NewZError(nil, StatusTooManyRequests, msgID, msg)
}

func NewBadRequest(msgID, msg string) *ZError {
	return NewZError(nil, StatusBadRequest, msgID, msg)
}

func NewValidationFailed(msg string) *ZError {
	return NewZError(nil, StatusValidationFailed, "validationFailed", msg)
}

func NewInternalServerError(msgID, msg string) *ZError {
	return NewZError(nil, StatusInternalServerError, msgID, msg)
}

func NewTimeout(msgID, msg string) *ZError {
	return NewZError(nil, StatusTimeout, msgID, msg)
}

func NewNotImplemented(msgID, msg string) *ZError {
	return NewZError(nil, StatusNotImplemented, msgID, msg)
}

func NewBadGateway(msgID, msg string) *ZError {
	return NewZError(nil, StatusBadGateway, msgID, msg)
}

func NewServiceUnavailable(msgID, msg string) *ZError {
	return NewZError(nil, StatusServiceUnavailable, msgID, msg)
}
