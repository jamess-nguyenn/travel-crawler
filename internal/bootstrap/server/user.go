package server

import (
	"travel-crawler/internal/user"
)

type UserModule struct{}

func (UserModule) Init(
	app *App,
	ctr *Container,
) {
	app.UserService = user.NewService(
		ctr.txManager,
		ctr.userRepo,
		ctr.logger,
	)
}
