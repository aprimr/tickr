// Package er defines shared application-wide errors used across all package
package er

import "errors"

var (
	ErrNotFound       = errors.New("requested resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden action")
	ErrInternalError  = errors.New("internal server error")
	ErrInvalidReqBody = errors.New("invalid request body")
)
