package helper

func IsNonEmptyString(value *string) bool {
	return value != nil && *value != ""
}
