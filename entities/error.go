package entities

import "fmt"

// ApplicationError contains the code and description of an application error
// and implements the Go error interface.
type ApplicationError struct {
	// Code identifies the error category, such as "BAD_REQUEST".
	Code string `json:"code"`

	// Message describes the error.
	Message string `json:"message"`
}

func (e ApplicationError) Error() string {
	return e.Message
}

func NewInternalServerError() ApplicationError {
	return ApplicationError{Code: "INTERNAL_SERVER_ERROR", Message: "unexpected error"}
}

func NewBadRequestError(message string) error {
	return ApplicationError{Code: "BAD_REQUEST", Message: message}
}

func NewUnauthorizedError(message string) error {
	return ApplicationError{Code: "UNAUTHORIZED", Message: message}
}

func NewNotFoundError(message string) error {
	return ApplicationError{Code: "NOT_FOUND", Message: message}
}

func NewInvalidCredentials() error {
	return NewUnauthorizedError("invalid login or password")
}

// JoinUnexpectedError wraps err with message, returning nil when err is nil.
func JoinUnexpectedError(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", message, err)
}
