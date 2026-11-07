package accommodation

var (
	upsertConflictColumns = []string{
		"provider_code",
		"external_id",
	}

	upsertUpdateColumns = []string{
		"external_parent_id",
		"name",
		"description",
		"property_type",
		"chain",
		"currency",
		"country",
		"city",
		"latitude",
		"longitude",
		"address",
		"zip",
		"main_photo",
		"thumbnail",
		"stars",
		"rating",
		"review_count",
		"metadata",
		"updated_at",
		"updated_by",
	}

	selectColumns = []string{
		"id",
		"provider_code",
		"external_id",
		"external_parent_id",
		"name",
		"property_type",
		"chain",
		"currency",
		"country",
		"city",
		"latitude",
		"longitude",
		"address",
		"zip",
		"main_photo",
		"thumbnail",
		"stars",
		"rating",
		"review_count",
		"created_at",
		"created_by",
		"updated_at",
		"updated_by",
		"deleted_at",
		"deleted_by",
	}
)
