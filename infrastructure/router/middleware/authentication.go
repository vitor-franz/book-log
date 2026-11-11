package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitor-franz/book-log/entities"
	"github.com/vitor-franz/book-log/infrastructure/router"
	"github.com/vitor-franz/book-log/modules/auth"
	"github.com/vitor-franz/book-log/utils/logger"
	"github.com/vitor-franz/book-log/utils/session"
)

// NewAuthenticationMiddleware returns middleware that requires a valid
// auth_token cookie. It loads the persisted user and stores it in the request
// context, readable through session.CurrentUser. Missing, invalid, or expired
// tokens and deleted users return UNAUTHORIZED.
func NewAuthenticationMiddleware(authUseCase auth.UseCase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie(auth.CookieName)
			if err != nil {
				router.HandleError(w, entities.NewUnauthorizedError("authentication required"))
				return
			}

			userID, err := authUseCase.Authenticate(cookie.Value)
			if err != nil {
				slog.ErrorContext(ctx, "failed to authenticate request", logger.Error(err))
				router.HandleError(w, err)
				return
			}

			user, err := authUseCase.GetUserByID(ctx, userID)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get authenticated user", logger.Error(err))
				router.HandleError(w, err)
				return
			}

			ctx = session.WithCurrentUser(ctx, *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
