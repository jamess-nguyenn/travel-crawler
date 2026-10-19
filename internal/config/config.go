package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"net"
	"strings"
)

type App struct {
	Name string `env:"NAME"`
	Key  string
	Env  string `env:"ENV"`
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (a *App) Init() {
	a.Key = strings.ReplaceAll(a.Name, "-", "_")
}

func (a *App) Address() string {
	return net.JoinHostPort(
		a.Host,
		a.Port,
	)
}

type Log struct {
	Level string `env:"LEVEL" env-default:"info"`
}

type Auth0 struct {
	Domain              string `env:"DOMAIN"`
	Audience            string `env:"AUDIENCE"`
	ClientID            string `env:"CLIENT_ID"`
	ClientSecret        string `env:"CLIENT_SECRET"`
	ManagerAudience     string `env:"MANAGER_AUDIENCE"`
	ManagerClientID     string `env:"MANAGER_CLIENT_ID"`
	ManagerClientSecret string `env:"MANAGER_CLIENT_SECRET"`
}

type MySQL struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Timezone string `env:"TIMEZONE" env-default:"UTC"`
	LogLevel string `env:"LOG_LEVEL" env-default:"error"`
}

type Redis struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

func (r *Redis) Address() string {
	return net.JoinHostPort(
		r.Host,
		r.Port,
	)
}

type Kafka struct {
	Brokers string `env:"BROKERS"`
}

type Nuitee struct {
	BaseURL string `env:"BASE_URL"`
	APIKey  string `env:"API_KEY"`
}

type Config struct {
	App                       App    `env-prefix:"APP_"`
	Log                       Log    `env-prefix:"LOG_"`
	Auth0                     Auth0  `env-prefix:"AUTH0_"`
	MySQL                     MySQL  `env-prefix:"DB_MYSQL_"`
	Redis                     Redis  `env-prefix:"REDIS_"`
	Kafka                     Kafka  `env-prefix:"KAFKA_"`
	Nuitee                    Nuitee `env-prefix:"NUITEE_"`
	AccommodationProvidersKey string `env:"ACCOMMODATION_PROVIDERS"`
	AccommodationProviders    []string
}

func (c *Config) Init() {
	c.AccommodationProviders = strings.Split(strings.Trim(c.AccommodationProvidersKey, ","), ",")
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.App.Init()
	cfg.Init()

	return &cfg, nil
}
