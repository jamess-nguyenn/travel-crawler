package user

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
)

var (
	upsertConflictColumns = []string{
		"auth0_id",
	}

	upsertUpdateColumns = []string{
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

	selectColumns = []string{
		"id",
		"email",
		"auth0_id",
		"email_verified_at",
		"first_name",
		"last_name",
		"full_name",
		"phone_number",
		"status",
		"avatar_url",
		"last_login_at",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
)
