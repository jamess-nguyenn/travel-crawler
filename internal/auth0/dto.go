package auth0

import "time"

type VerificationEmailJobStatus string

const (
	VerificationEmailJobStatusPending   VerificationEmailJobStatus = "pending"
	VerificationEmailJobStatusCompleted VerificationEmailJobStatus = "completed"
	VerificationEmailJobStatusFailed    VerificationEmailJobStatus = "failed"
)

type TokenRequest struct {
	Audience     string `json:"audience"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type VerificationEmailRequest struct {
	UserID   string `json:"user_id"`
	ClientID string `json:"client_id"`
}

type VerificationEmailResponse struct {
	Type      string                     `json:"type"`
	Status    VerificationEmailJobStatus `json:"status"`
	CreatedAt time.Time                  `json:"created_at"`
	ID        string                     `json:"id"`
}

type UserResponse struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture,omitempty"`
}
