package airline

var (
	upsertConflictColumns = []string{
		"provider_code",
		"iata_code",
	}

	upsertUpdateColumns = []string{
		"name",
		"alliance",
		"active",
		"logo",
		"metadata",
		"updated_at",
		"updated_by",
	}

	selectColumns = []string{
		"id",
		"provider_code",
		"iata_code",
		"name",
		"alliance",
		"active",
		"logo",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
)
