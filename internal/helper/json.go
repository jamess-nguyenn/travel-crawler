package helper

import "encoding/json"

func FormatJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return ""
	}

	return string(data)
}
