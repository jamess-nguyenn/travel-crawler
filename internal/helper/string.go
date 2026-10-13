package helper

import (
	"strconv"
	"strings"
)

func StringValue(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func ParseString(value int) string {
	return strconv.Itoa(value)
}
