package entities

// LoginCredentials contains the login and password supplied to login.
type LoginCredentials struct {
	// Login is the user's email address.
	Login string `json:"login"`

	// Password is the user's password.
	Password string `json:"password"`
}

// RegisterCredentials contains the details supplied to create a user.
type RegisterCredentials struct {
	// Name is the user's name.
	Name string `json:"name"`

	// Email is the user's email address, which is used to log in.
	Email string `json:"email"`

	// Password is the plaintext password to hash before persistence.
	Password string `json:"password"`
}

// User contains the persisted identity of an application user.
type User struct {
	// ID is the unique database identifier used as the JWT subject.
	ID int64 `json:"id"`

	// Name is the user's display name.
	Name string `json:"name"`

	// Email is the user's email address.
	Email string `json:"email"`
}
