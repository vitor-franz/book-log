package session

import (
	"context"

	"github.com/vitor-franz/book-log/entities"
)

// contextKey isolates session context values from keys defined by other packages.
type contextKey string

// currentUserKey identifies the authenticated user in the request context.
const currentUserKey contextKey = "current_user"

// WithCurrentUser returns a context containing user and retaining the parent's
// cancellation, deadlines, and values.
func WithCurrentUser(ctx context.Context, user entities.User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

// CurrentUser returns the authenticated user and true when present, or a zero
// User and false when the context contains no current user.
func CurrentUser(ctx context.Context) (entities.User, bool) {
	user, ok := ctx.Value(currentUserKey).(entities.User)
	return user, ok
}
