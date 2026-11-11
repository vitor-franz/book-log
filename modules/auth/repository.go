package auth

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/vitor-franz/book-log/entities"
	"golang.org/x/crypto/bcrypt"
)

//go:embed sql/insert_user.sql
var insertUserQuery string

//go:embed sql/select_password_by_login.sql
var selectPasswordByLoginQuery string

//go:embed sql/select_user_by_id.sql
var selectUserByIDQuery string

//go:embed sql/select_user_by_login.sql
var selectUserByLoginQuery string

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}

type repository struct {
	db *sql.DB
}

func (r repository) GetUserByLogin(ctx context.Context, login string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRowContext(ctx, selectUserByLoginQuery, login).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, entities.JoinUnexpectedError(err, "failed to query user by login")
	}

	return &user, nil
}

func (r repository) CheckCredentials(ctx context.Context, credentials entities.LoginCredentials) error {
	var passwordHash string
	err := r.db.QueryRowContext(ctx, selectPasswordByLoginQuery, credentials.Login).Scan(&passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewInvalidCredentials()
		}

		return entities.JoinUnexpectedError(err, "failed to query user password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(credentials.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return entities.NewInvalidCredentials()
		}

		return entities.JoinUnexpectedError(err, "failed to compare user password")
	}

	return nil
}

func (r repository) AttemptRegister(ctx context.Context, credentials entities.RegisterCredentials) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.JoinUnexpectedError(err, "failed to hash user password")
	}

	_, err = r.db.ExecContext(ctx, insertUserQuery, credentials.Name, credentials.Email, string(passwordHash))
	if err != nil {
		return entities.JoinUnexpectedError(err, "failed to insert user")
	}

	return nil
}

func (r repository) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRowContext(ctx, selectUserByIDQuery, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, entities.JoinUnexpectedError(err, "failed to query user by ID")
	}

	return &user, nil
}
