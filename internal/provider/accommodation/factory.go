package accommodation

import (
	"fmt"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/integration/accommodation/nuitee"
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

func (f *Factory) Create(code string) (accommodation.Provider, error) {
	switch code {
	case constant.ProviderCodeNuitee:
		client := nuitee.NewClient(
			f.http,
			f.config.Nuitee,
			f.cache,
			f.logger,
		)

		limiter, err := nuitee.NewRedisRequestLimiter(
			f.cache,
			f.config.Nuitee.RequestInterval,
		)
		if err != nil {
			return nil, fmt.Errorf("initialize %q rate limiter: %w", code, err)
		}

		return nuitee.New(client, limiter), nil

	default:
		return nil, fmt.Errorf("unsupported accommodation provider %q", code)
	}
}
