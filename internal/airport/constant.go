package airport

var (
	upsertConflictColumns = []string{
		"iata_code",
	}

	upsertUpdateColumns = []string{
		"city",
		"country",
		"metadata",
		"updated_at",
		"updated_by",
	}

	selectColumns = []string{
		"id",
		"iata_code",
		"city",
		"country",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
)
