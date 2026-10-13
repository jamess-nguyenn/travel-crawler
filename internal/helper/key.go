package helper

import (
	"fmt"
	"strings"
)

func JoinKey(parts ...any) string {
	values := make([]string, len(parts))

	for i, part := range parts {
		values[i] = fmt.Sprint(part)
	}

	return strings.Join(values, ":")
}
