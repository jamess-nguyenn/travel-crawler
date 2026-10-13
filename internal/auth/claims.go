package auth

import "context"

type CustomClaims struct {
	Email         string       `json:"email"`
	EmailVerified bool         `json:"email_verified"`
	UserMetadata  UserMetadata `json:"user_metadata"`
	Picture       string       `json:"picture"`
}

type UserMetadata struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}
