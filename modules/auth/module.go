package auth

import (
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vitor-franz/book-log/entities"
	"github.com/vitor-franz/book-log/infrastructure/router"
	"github.com/vitor-franz/book-log/utils/logger"
)

// CookieName is the name of the HTTP-only cookie that carries the session token.
const CookieName = "auth_token"

func NewModule(authUseCase UseCase, insecureCookie bool) router.Module {
	return module{authUseCase: authUseCase, insecureCookie: insecureCookie}
}

type module struct {
	authUseCase    UseCase
	insecureCookie bool
}

func (m module) Path() string {
	return "/auth"
}

func (m module) Setup(_ *mux.Router, router *mux.Router) {
	subRouter := router.PathPrefix(m.Path()).Subrouter()

	subRouter.HandleFunc("/login", m.handleLogin).Methods(http.MethodPost)
	subRouter.HandleFunc("/register", m.handleRegister).Methods(http.MethodPost)
	subRouter.HandleFunc("/logout", m.handleLogout).Methods(http.MethodPost)
}

func (m module) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", logger.Error(err))
		router.HandleError(w, entities.NewBadRequestError("failed to read request body"))
		return
	}

	var credentials entities.LoginCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request JSON", logger.Error(err))
		router.HandleError(w, entities.NewBadRequestError("invalid request JSON"))
		return
	}

	token, err := m.authUseCase.AttemptLogin(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to authenticate user", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	m.setTokenCookie(w, token, tokenDuration)

	body, err = json.Marshal(entities.NewDefaultResponse())
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode success response", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write success response", logger.Error(err))
	}
}

func (m module) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", logger.Error(err))
		router.HandleError(w, entities.NewBadRequestError("failed to read request body"))
		return
	}

	var credentials entities.RegisterCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode request JSON", logger.Error(err))
		router.HandleError(w, entities.NewBadRequestError("invalid request JSON"))
		return
	}

	err = m.authUseCase.AttemptRegister(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to register user", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	token, err := m.authUseCase.AttemptLogin(ctx, entities.LoginCredentials{
		Login:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to authenticate registered user", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	m.setTokenCookie(w, token, tokenDuration)

	body, err = json.Marshal(entities.NewDefaultResponse())
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode success response", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write success response", logger.Error(err))
	}
}

func (m module) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := json.Marshal(entities.NewDefaultResponse())
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode success response", logger.Error(err))
		router.HandleError(w, err)
		return
	}

	m.setTokenCookie(w, "", -time.Second)

	_, err = w.Write(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write success response", logger.Error(err))
	}
}

// setTokenCookie writes the session cookie with the given token and lifetime.
// A negative duration makes the browser delete the cookie immediately.
func (m module) setTokenCookie(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(duration),
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   !m.insecureCookie,
		SameSite: http.SameSiteStrictMode,
	})
}
