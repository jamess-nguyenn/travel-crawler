package auth

type CreateRequest struct {
	Email         string  `json:"email" validate:"required,email,max=255"`
	Auth0ID       string  `json:"auth0_id" validate:"required,max=255"`
	EmailVerified bool    `json:"email_verified"`
	FirstName     string  `json:"first_name" validate:"required,max=100"`
	LastName      string  `json:"last_name" validate:"required,max=100"`
	PhoneNumber   *string `json:"phone_number" validate:"omitempty,max=100"`
	AvatarURL     *string `json:"avatar_url" validate:"omitempty,max=500"`
}
