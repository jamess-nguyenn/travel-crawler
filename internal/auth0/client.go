package auth0

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"travel-crawler/internal/config"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/redis"
)

const (
	GrantTypeClientCredentials = "client_credentials"
)

const (
	URIOauthToken        = "/oauth/token"
	URIVerificationEmail = "/jobs/verification-email"
	URIUsers             = "/users"
)

type Client interface {
	SendVerificationEmail(ctx context.Context, auth0ID string) (*VerificationEmailResponse, error)
	GetUserByAuth0ID(ctx context.Context, auth0ID string) (*UserResponse, error)
}

type client struct {
	http   *httpclient.Client
	config config.Auth0
	cache  *redis.Cache
	logger *zap.SugaredLogger

	mu sync.Mutex
}

func New(
	cfg config.Auth0,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) Client {
	return &client{
		http:   httpclient.New(logger),
		config: cfg,
		cache:  cache,
		logger: logger,
	}
}

func (c *client) cacheKey() string {
	return helper.JoinKey(
		"auth0",
		"token",
		c.config.ClientID,
	)
}

func (c *client) getCachedAccessToken(ctx context.Context, cacheKey string) (string, bool) {
	var token string

	found, err := c.cache.Get(ctx, cacheKey, &token)
	if err != nil {
		c.logger.Warnw(
			"Fail to get access token from cache",
			"key", cacheKey,
			"error", err,
		)
	}

	return token, found
}

func (c *client) cacheAccessToken(ctx context.Context, cacheKey string, accessToken string) {
	if err := c.cache.Set(ctx, cacheKey, accessToken, helper.DayTTL(1)); err != nil {
		c.logger.Warnw(
			"Fail to cache access token",
			"key", cacheKey,
			"error", err,
		)
	}
}

func (c *client) getAccessToken(ctx context.Context) (string, error) {
	cacheKey := c.cacheKey()

	// First cache check: avoid locking when the token exists
	if token, found := c.getCachedAccessToken(ctx, cacheKey); found {
		return token, nil
	}

	// Only one goroutine can fetch a token at a time
	c.mu.Lock()
	defer c.mu.Unlock()

	// Second cache check: another goroutine may have fetched
	// and cached the token while this goroutine was waiting
	if token, found := c.getCachedAccessToken(ctx, cacheKey); found {
		return token, nil
	}

	// No cached token: request one from Auth0
	var result TokenResponse

	req := &httpclient.HTTPRequest{
		Headers: httpclient.JSONHeaders(),
		Body: TokenRequest{
			Audience:     c.config.ManagerAudience + "/",
			ClientID:     c.config.ManagerClientID,
			ClientSecret: c.config.ManagerClientSecret,
			GrantType:    GrantTypeClientCredentials,
		},
	}

	url := c.config.Domain + URIOauthToken

	if err := c.http.Post(ctx, url, req, &result); err != nil {
		return "", err
	}

	accessToken := result.AccessToken
	if accessToken == "" {
		return "", fmt.Errorf("Auth0 returned an empty access token")
	}

	c.cacheAccessToken(ctx, cacheKey, accessToken)

	return accessToken, nil
}

func (c *client) buildAuthHeaders(
	ctx context.Context,
) (map[string]string, error) {
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	headers := httpclient.JSONHeaders()
	headers[httpclient.HeaderAuthorization] = "Bearer " + accessToken

	return headers, nil
}

func (c *client) SendVerificationEmail(ctx context.Context, auth0ID string) (*VerificationEmailResponse, error) {
	headers, err := c.buildAuthHeaders(ctx)
	if err != nil {
		return nil, err
	}

	var result VerificationEmailResponse

	req := &httpclient.HTTPRequest{
		Headers: headers,
		Body: VerificationEmailRequest{
			UserID:   auth0ID,
			ClientID: c.config.ClientID,
		},
	}

	url := c.config.ManagerAudience + URIVerificationEmail

	if err := c.http.Post(ctx, url, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *client) GetUserByAuth0ID(ctx context.Context, auth0ID string) (*UserResponse, error) {
	headers, err := c.buildAuthHeaders(ctx)
	if err != nil {
		return nil, err
	}

	var result UserResponse

	req := &httpclient.HTTPRequest{
		Headers: headers,
	}

	url := c.config.ManagerAudience + URIUsers + "/" + helper.EscapePath(auth0ID)

	if err := c.http.Get(ctx, url, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
