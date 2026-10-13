package user

import "errors"

var (
	ErrEmailExists      = errors.New("Email already exists")
	ErrAuth0IDExists    = errors.New("Auth0 ID already exists")
	ErrEmailNotVerified = errors.New("Email not verified")
)
