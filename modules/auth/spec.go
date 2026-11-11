package auth

import (
	"context"

	"github.com/vitor-franz/book-log/entities"
)

// Repository defines the persistence operations for login and registration.
type Repository interface {
	// GetUserByLogin finds a user by email and returns
	// only the user's ID.
	// A missing user returns nil, nil; database failures return an error.
	// The query respects ctx cancellation and deadlines.
	GetUserByLogin(ctx context.Context, login string) (*entities.User, error)

	// CheckCredentials loads the password hash by login and compares it with bcrypt.
	// It returns nil on success, an UNAUTHORIZED ApplicationError for invalid
	// credentials, or a wrapped database or password-hash error.
	// The query respects ctx cancellation and deadlines; bcrypt has no context support.
	CheckCredentials(ctx context.Context, credentials entities.LoginCredentials) error

	// AttemptRegister inserts the user's name, email, and bcrypt password hash.
	// It returns nil on success or a hashing or database error, including duplicate identifiers.
	// The insert respects ctx cancellation and deadlines; bcrypt has no context support.
	AttemptRegister(ctx context.Context, credentials entities.RegisterCredentials) error

	// GetUserByID loads the full user identified by id.
	// A missing user returns nil, nil; database failures return an error.
	// The query respects ctx cancellation and deadlines.
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
}

// UseCase defines the application's login and registration operations.
type UseCase interface {
	// Authenticate validates a signed session token and returns its user ID.
	// Invalid or expired tokens return zero and an UNAUTHORIZED ApplicationError.
	// This operation performs no I/O and does not load user data.
	Authenticate(token string) (int64, error)

	// AttemptLogin authenticates a user by email and password.
	// It returns a signed JWT on success or an UNAUTHORIZED ApplicationError
	// for invalid credentials. Unexpected failures retain their wrapped causes.
	// ctx carries the request's cancellation and deadline to the repository.
	AttemptLogin(ctx context.Context, credentials entities.LoginCredentials) (string, error)

	// AttemptRegister creates a user without returning user data or a token.
	// It returns nil on success or a wrapped hashing or database error on failure.
	// ctx carries the request's cancellation and deadline to the repository.
	AttemptRegister(ctx context.Context, credentials entities.RegisterCredentials) error

	// GetUserByID loads the user identified by an authenticated token subject.
	// It returns the user on success, or nil and an UNAUTHORIZED error when the
	// user no longer exists. Database failures return nil and a wrapped error.
	// ctx carries cancellation and deadlines to the repository.
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
}
