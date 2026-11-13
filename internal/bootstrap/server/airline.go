package server

import (
	"travel-crawler/internal/airline"
)

type AirlineModule struct{}

func (AirlineModule) Init(
	app *App,
	ctr *Container,
) {
	app.AirlineService = airline.NewService(
		ctr.txManager,
		ctr.airlineRepo,
		ctr.logger,
	)
}
