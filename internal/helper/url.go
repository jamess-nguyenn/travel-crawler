package helper

import "net/url"

func EscapePath(value string) string {
	return url.PathEscape(value)
}
