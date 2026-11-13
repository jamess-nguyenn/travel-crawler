package scheduler

import (
	"travel-crawler/internal/airport"
	airportjob "travel-crawler/internal/job/airport"
)

type AirportModule struct{}

func (AirportModule) Init(app *App, ctr *Container) error {

	airportService := airport.NewSyncService(
		ctr.txManager,
		ctr.airportRepo,
		ctr.logger,
	)

	app.AirportSyncJob = airportjob.NewSyncJob(
		airportService,
		ctr.logger,
	)

	return nil
}
