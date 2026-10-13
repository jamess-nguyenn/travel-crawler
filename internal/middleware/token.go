package middleware

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v3"
	"github.com/gorilla/mux"
)

func RequireToken(
	jwtMiddleware *jwtmiddleware.JWTMiddleware,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return jwtMiddleware.CheckJWT(next)
	}
}
