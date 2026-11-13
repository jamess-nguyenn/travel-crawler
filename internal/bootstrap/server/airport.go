package server

import (
	"travel-crawler/internal/airport"
)

type AirportModule struct{}

func (AirportModule) Init(
	app *App,
	ctr *Container,
) {
	app.AirportService = airport.NewService(
		ctr.txManager,
		ctr.airportRepo,
		ctr.logger,
	)
}
