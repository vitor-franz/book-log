package logger

import (
	"context"
	"log/slog"
)

// contextKey identifies logger attributes attached to a request context.
type contextKey string

// logFields is the private context key for structured request attributes.
const logFields contextKey = "slog_fields"

// AddLogValueToContext returns a child context with additional logging attributes.
// The parent's attribute slice is copied to keep sibling contexts independent.
func AddLogValueToContext(parent context.Context, attributes ...slog.Attr) context.Context {
	existing := GetLogValuesFromContext(parent)
	values := make([]slog.Attr, 0, len(existing)+len(attributes))
	values = append(values, existing...)
	values = append(values, attributes...)

	return context.WithValue(parent, logFields, values)
}

// GetLogValuesFromContext returns logging attributes attached to ctx, or nil
// when no attributes are present. The returned slice must not be modified.
func GetLogValuesFromContext(ctx context.Context) []slog.Attr {
	attributes, _ := ctx.Value(logFields).([]slog.Attr)
	return attributes
}
