package auth

import "errors"

var (
	ErrVerificationEmailFailed = errors.New("Verification email was not accepted")
)
