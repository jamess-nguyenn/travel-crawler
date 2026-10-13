package helper

import (
	"strconv"
)

func IntValue(value *int) int {
	if value == nil || *value < 1 {
		return 0
	}

	return *value
}

func ParseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func ParseUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}
