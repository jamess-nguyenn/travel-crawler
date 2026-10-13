package auth0

import (
	"context"
	"go.uber.org/zap"
	"travel-crawler/internal/config"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/httpclient"
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
	logger *zap.SugaredLogger
}

func New(
	cfg config.Auth0,
	logger *zap.SugaredLogger,
) Client {
	return &client{
		http:   httpclient.New(logger),
		config: cfg,
		logger: logger,
	}
}

func (c *client) getAccessToken(ctx context.Context) (string, error) {
	var result TokenResponse

	req := &httpclient.HTTPRequest{
		Headers: httpclient.FormURLEncodedHeaders(),
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

	return result.AccessToken, nil
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

	url := c.config.Domain + URIVerificationEmail

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
