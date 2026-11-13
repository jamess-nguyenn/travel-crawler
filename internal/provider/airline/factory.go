package airline

import (
	"fmt"
	"travel-crawler/internal/airline"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/integration/nuitee"
	airlinenuitee "travel-crawler/internal/integration/nuitee/airline"
	"travel-crawler/internal/redis"

	"go.uber.org/zap"
)

type Factory struct {
	http   *httpclient.Client
	config config.Config
	cache  *redis.Cache
	logger *zap.SugaredLogger
}

func NewFactory(
	http *httpclient.Client,
	cfg config.Config,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) *Factory {
	return &Factory{
		http:   http,
		config: cfg,
		cache:  cache,
		logger: logger,
	}
}

func (f *Factory) Create(code string) (airline.Provider, error) {
	switch code {
	case constant.ProviderCodeNuitee:
		api := nuitee.NewClient(
			f.http,
			f.config.Nuitee,
			f.cache,
			f.logger,
		)

		client := airlinenuitee.NewClient(api)

		limiter, err := nuitee.NewRedisRequestLimiter(
			f.cache,
			f.config.Nuitee.RequestPerSecond,
		)
		if err != nil {
			return nil, fmt.Errorf("initialize %q rate limiter: %w", code, err)
		}

		return airlinenuitee.New(client, nuitee.NewRequestGuard(limiter)), nil

	default:
		return nil, fmt.Errorf("unsupported airline provider %q", code)
	}
}
