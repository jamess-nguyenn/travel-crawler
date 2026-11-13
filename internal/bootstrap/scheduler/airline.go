package scheduler

import (
	"fmt"
	"travel-crawler/internal/airline"
	airlinejob "travel-crawler/internal/job/airline"
	airlineprovider "travel-crawler/internal/provider/airline"
)

type AirlineModule struct{}

func (AirlineModule) Init(app *App, ctr *Container) error {
	providers, err := airlineprovider.New(ctr.http, *ctr.cfg, ctr.cache, ctr.logger)
	if err != nil {
		return fmt.Errorf("initialize airline providers: %w", err)
	}

	airlineService := airline.NewSyncService(
		ctr.txManager,
		ctr.airlineRepo,
		providers,
		ctr.logger,
	)

	app.AirlineSyncJob = airlinejob.NewSyncJob(
		airlineService,
		ctr.cfg.AirlineProviders,
		ctr.logger,
	)

	return nil
}
