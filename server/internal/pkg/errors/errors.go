// Package er defines shared application-wide errors used across all package
package er

import "errors"

var (
	ErrNotFound       = errors.New("requested resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden action")
	ErrInternalError  = errors.New("something went wrong")
	ErrInvalidReqBody = errors.New("invalid request body")
)
