package airline

type Airline struct {
	IATA     string  `json:"iata"`
	Name     string  `json:"name"`
	Alliance *string `json:"alliance"`
	Active   bool    `json:"active"`
	Logo     *string `json:"logo"`
}

type AirlineSet struct {
	Airlines []Airline `json:"airlines"`
}

type AirlineResponse struct {
	Data []AirlineSet `json:"data"`
}
