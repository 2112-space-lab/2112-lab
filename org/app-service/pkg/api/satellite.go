package api_mappers

type ContextSatelliteRequest struct {
	Name           string   `json:"name"`
	SatelliteNames []string `json:"satelliteNames"`
}
