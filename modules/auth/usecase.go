package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vitor-franz/book-log/entities"
)

// tokenDuration is how long a signed session token and its cookie stay valid.
const tokenDuration = time.Hour

func NewUseCase(authRepository Repository, signingKey string) UseCase {
	return useCase{authRepository: authRepository, signingKey: signingKey}
}

type useCase struct {
	authRepository Repository
	signingKey     string
}

func (u useCase) Authenticate(signedToken string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(signedToken, claims, func(_ *jwt.Token) (any, error) {
		return []byte(u.signingKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return 0, entities.NewUnauthorizedError("invalid authentication token")
	}

	if !token.Valid {
		return 0, entities.NewUnauthorizedError("invalid authentication token")
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, entities.NewUnauthorizedError("invalid authentication token")
	}

	return id, nil
}

func (u useCase) AttemptLogin(ctx context.Context, credentials entities.LoginCredentials) (string, error) {
	user, err := u.authRepository.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return "", entities.JoinUnexpectedError(err, "failed to get user by login")
	}

	if user == nil {
		return "", entities.NewInvalidCredentials()
	}

	err = u.authRepository.CheckCredentials(ctx, credentials)
	if err != nil {
		return "", entities.JoinUnexpectedError(err, "failed to check credentials")
	}

	return u.signToken(user.ID)
}

func (u useCase) AttemptRegister(ctx context.Context, credentials entities.RegisterCredentials) error {
	err := u.authRepository.AttemptRegister(ctx, credentials)
	if err != nil {
		return entities.JoinUnexpectedError(err, "failed to register user")
	}

	return nil
}

func (u useCase) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	user, err := u.authRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, entities.JoinUnexpectedError(err, "failed to get user")
	}

	if user == nil {
		return nil, entities.NewUnauthorizedError("authenticated user no longer exists")
	}

	return user, nil
}

func (u useCase) signToken(userID int64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
	})

	signedToken, err := token.SignedString([]byte(u.signingKey))
	if err != nil {
		return "", entities.JoinUnexpectedError(err, "failed to sign authentication token")
	}

	return signedToken, nil
}
