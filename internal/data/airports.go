package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed airports.json
var airportsJSON []byte

type Airport struct {
	Country string `json:"country"`
	City    string `json:"city"`
	IATA    string `json:"iata"`
}

var (
	airportsOnce sync.Once
	airports     []Airport
	airportsErr  error
)

// Airports returns the shared list of airports.
// The JSON is parsed only once.
func Airports() ([]Airport, error) {
	airportsOnce.Do(func() {
		airportsErr = json.Unmarshal(airportsJSON, &airports)
		if airportsErr != nil {
			airportsErr = fmt.Errorf(
				"unmarshal airports.json: %w",
				airportsErr,
			)
		}
	})

	if airportsErr != nil {
		return nil, airportsErr
	}

	// Return a copy so callers cannot modify the shared slice.
	result := make([]Airport, len(airports))
	copy(result, airports)

	return result, nil
}
