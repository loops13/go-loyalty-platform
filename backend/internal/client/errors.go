package client

import "fmt"

// Error types for client domain.
type ClientError struct {
	Code    string
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrNotFound         = &ClientError{Code: "CLIENT_NOT_FOUND", Message: "client does not exist"}
	ErrEmptyName        = &ClientError{Code: "EMPTY_NAME", Message: "name is required"}
	ErrEmptyEmail       = &ClientError{Code: "EMPTY_EMAIL", Message: "email is required"}
	ErrInvalidAwardType = &ClientError{Code: "INVALID_AWARD_TYPE", Message: "unknown award type"}
	ErrInsufficientPts  = &ClientError{Code: "INSUFFICIENT_POINTS", Message: "insufficient point balance"}
)
