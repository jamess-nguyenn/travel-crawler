package accommodation

var (
	UpsertConflictColumns = []string{
		"provider_code",
		"external_id",
	}

	UpsertUpdateColumns = []string{
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
)
