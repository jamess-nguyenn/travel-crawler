package helper

import "net/url"

func EscapeQuery(value string) string {
	return url.QueryEscape(value)
}

func EscapePath(value string) string {
	return url.PathEscape(value)
}
