package auth

import (
	"net/url"

	"github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"travel-crawler/internal/config"
)

func New(
	cfg config.Auth0,
) (*jwtmiddleware.JWTMiddleware, error) {
	issuerURL, err := url.Parse(
		cfg.Domain + "/",
	)
	if err != nil {
		return nil, err
	}

	provider, err := jwks.NewCachingProvider(
		jwks.WithIssuerURL(issuerURL),
	)
	if err != nil {
		return nil, err
	}

	jwtValidator, err := validator.New(
		validator.WithKeyFunc(provider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(issuerURL.String()),
		validator.WithAudience(cfg.Audience),
		validator.WithCustomClaims(func() *CustomClaims {
			return &CustomClaims{}
		}),
	)
	if err != nil {
		return nil, err
	}

	return jwtmiddleware.New(
		jwtmiddleware.WithValidator(jwtValidator),
	)
}
