package middleware

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"travel-crawler/internal/user"

	"github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/gorilla/mux"
	"travel-crawler/internal/auth"
	"travel-crawler/internal/requestcontext"
	"travel-crawler/internal/response"
)

func RequireUser(
	authService auth.Service,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			ctx := r.Context()

			claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)
			if err != nil {
				response.ErrorUnauthorized(
					w,
					"Invalid token claims",
				)
				return
			}

			auth0ID := claims.RegisteredClaims.Subject
			if auth0ID == "" {
				response.ErrorUnauthorized(
					w,
					"Missing token subject",
				)
				return
			}

			currentUser, err := authService.Get(
				ctx,
				auth0ID,
			)
			if err != nil {
				switch {
				case errors.Is(err, gorm.ErrRecordNotFound):
					response.ErrorUnauthorized(w, "User not found")

				case currentUser.Status != string(user.StatusActive):
					response.ErrorForbidden(w, "User is not active")

				default:
					response.ErrorInternalServer(w, nil)
				}
				return
			}

			ctx = requestcontext.WithUser(
				ctx,
				currentUser,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
