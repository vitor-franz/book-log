package logger

import (
	"errors"
	"log/slog"

	"github.com/vitor-franz/book-log/entities"
)

// Error returns a structured error attribute with the application code, message,
// and full error details. Errors without an application error use the internal
// server error code and message. A nil error returns a null error attribute.
func Error(err error) slog.Attr {
	if err == nil {
		return slog.Any("error", nil)
	}

	applicationError := entities.NewInternalServerError()
	errors.As(err, &applicationError)

	return slog.Group("error",
		slog.String("code", applicationError.Code),
		slog.String("message", applicationError.Message),
		slog.String("details", err.Error()),
	)
}
