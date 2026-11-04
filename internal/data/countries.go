package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed countries.json
var countriesJSON []byte

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var (
	countriesOnce sync.Once
	countries     []Country
	countriesErr  error
)

// Countries returns the shared list of countries.
// The JSON is parsed only once.
func Countries() ([]Country, error) {
	countriesOnce.Do(func() {
		countriesErr = json.Unmarshal(countriesJSON, &countries)
		if countriesErr != nil {
			countriesErr = fmt.Errorf(
				"unmarshal countries.json: %w",
				countriesErr,
			)
		}
	})

	if countriesErr != nil {
		return nil, countriesErr
	}

	// Return a copy so callers cannot modify the shared slice.
	result := make([]Country, len(countries))
	copy(result, countries)

	return result, nil
}
