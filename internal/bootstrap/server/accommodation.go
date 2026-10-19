package server

import (
	"travel-crawler/internal/accommodation"
)

type AccommodationModule struct{}

func (AccommodationModule) Init(
	app *App,
	ctr *Container,
) {
	app.AccommodationService = accommodation.NewService(
		ctr.txManager,
		ctr.accommodationRepo,
		ctr.logger,
	)
}
