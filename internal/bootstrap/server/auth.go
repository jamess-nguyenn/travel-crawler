package server

import (
	"travel-crawler/internal/auth"
)

type AuthModule struct{}

func (AuthModule) Init(
	app *App,
	ctr *Container,
) {
	app.AuthService = auth.NewService(
		ctr.txManager,
		ctr.userRepo,
		ctr.auth0Client,
		ctr.logger,
	)
}
