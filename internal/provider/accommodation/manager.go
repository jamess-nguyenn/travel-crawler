package accommodation

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/config"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/redis"
)

type Manager struct {
	providers map[string]accommodation.Provider
	mu        sync.RWMutex
}

func New(
	http *httpclient.Client,
	cfg config.Config,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) (*Manager, error) {
	if len(cfg.AccommodationProviders) == 0 {
		return nil, fmt.Errorf(
			"no accommodation providers configured",
		)
	}

	factory := NewFactory(
		http,
		cfg,
		cache,
		logger,
	)

	manager := NewManager()

	for _, code := range cfg.AccommodationProviders {
		provider, err := factory.Create(code)
		if err != nil {
			return nil, fmt.Errorf(
				"create accommodation provider %q: %w",
				code,
				err,
			)
		}

		if err := manager.Register(provider); err != nil {
			return nil, fmt.Errorf(
				"register accommodation provider %q: %w",
				code,
				err,
			)
		}
	}

	return manager, nil
}

func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]accommodation.Provider),
	}
}

func (m *Manager) Register(p accommodation.Provider) error {
	if p == nil {
		return fmt.Errorf("cannot register nil provider")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	code := p.Code()
	if code == "" {
		return fmt.Errorf("provider code cannot be empty")
	}

	if _, exists := m.providers[code]; exists {
		return fmt.Errorf(
			"provider %q is already registered",
			code,
		)
	}

	m.providers[code] = p

	return nil
}

func (m *Manager) Get(code string) (accommodation.Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.providers[code]
	if !exists {
		return nil, fmt.Errorf(
			"provider %q is not registered",
			code,
		)
	}

	return p, nil
}

func (m *Manager) List() []accommodation.Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providers := make([]accommodation.Provider, 0, len(m.providers))

	for _, p := range m.providers {
		providers = append(providers, p)
	}

	return providers
}
