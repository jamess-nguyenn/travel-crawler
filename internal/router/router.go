package router

import (
	"github.com/auth0/go-jwt-middleware/v3"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"travel-crawler/internal/auth"
	"travel-crawler/internal/bootstrap/server"
	"travel-crawler/internal/middleware"
	"travel-crawler/internal/user"
)

func New(
	app *server.App,
	jwtMiddleware *jwtmiddleware.JWTMiddleware,
	logger *zap.SugaredLogger,
) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.RequestContext())
	router.Use(middleware.AccessLog(logger))

	api := router.PathPrefix("/api/v1").Subrouter()

	//JWT Required

	tokenAPI := api.PathPrefix("").Subrouter()

	tokenAPI.Use(
		middleware.RequireToken(
			jwtMiddleware,
		),
	)

	tokenAPI.
		HandleFunc(
			"/auth/send-verify-email",
			auth.SendVerifyEmail(
				app.AuthService,
			),
		).
		Methods(http.MethodPost).
		Name("SendVerifyEmail")

	tokenAPI.
		HandleFunc(
			"/auth/verify",
			auth.Verify(
				app.AuthService,
			),
		).
		Methods(http.MethodPost).
		Name("VerifyAuth")

	//User Required

	protectedAPI := tokenAPI.PathPrefix("").Subrouter()

	protectedAPI.Use(
		middleware.RequireUser(
			app.AuthService,
		),
	)

	protectedAPI.
		HandleFunc(
			"/me",
			user.GetMe(),
		).
		Methods(http.MethodGet).
		Name("GetMe")

	//Users

	apiUser := protectedAPI.PathPrefix("/users").Subrouter()

	apiUser.
		HandleFunc(
			"/{id:[0-9]+}",
			user.Get(
				app.UserService,
			),
		).
		Methods(http.MethodGet).
		Name("GetUser")

	apiUser.
		HandleFunc(
			"/{id:[0-9]+}",
			user.Update(
				app.UserService,
			),
		).
		Methods(http.MethodPut).
		Name("UpdateUser")

	apiUser.
		HandleFunc(
			"",
			user.List(
				app.UserService,
			),
		).
		Methods(http.MethodGet).
		Name("ListUser")

	return router
}
