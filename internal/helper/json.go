package helper

import (
	"encoding/json"
	"gorm.io/datatypes"
)

const (
	emptyJSONObject = `{}`
	emptyJSONArray  = `[]`
)

func toJSON(v any, fallback string) datatypes.JSON {
	data, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON(fallback)
	}

	if !json.Valid(data) {
		return datatypes.JSON(fallback)
	}

	return datatypes.JSON(data)
}

func ToJSONObject(v any) datatypes.JSON {
	data := toJSON(v, emptyJSONObject)

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil || result == nil {
		return datatypes.JSON(emptyJSONObject)
	}

	return data
}

func ToJSONArray(v any) datatypes.JSON {
	data := toJSON(v, emptyJSONArray)

	var result []any
	if err := json.Unmarshal(data, &result); err != nil || result == nil {
		return datatypes.JSON(emptyJSONArray)
	}

	return data
}

func FormatJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return ""
	}

	return string(data)
}
