package user

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
)

var (
	UpsertConflictColumns = []string{
		"auth0_id",
	}

	UpsertUpdateColumns = []string{
		"email",
		"email_verified_at",
		"password_hash",
		"first_name",
		"last_name",
		"full_name",
		"phone_number",
		"status",
		"avatar_url",
		"metadata",
		"updated_at",
		"updated_by",
	}
)
