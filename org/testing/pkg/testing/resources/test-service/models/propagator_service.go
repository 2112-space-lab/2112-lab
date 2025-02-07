package models

type SatellitePropagationRequest struct {
	NoradId         string `json:"noradId"`
	TleLine1        string `json:"tleLine1"`
	TleLine2        string `json:"tleLine2"`
	StartTime       string `json:"startTime"`
	DurationMinutes int32  `json:"durationMinutes"`
	IntervalSeconds int32  `json:"intervalSeconds"`
}
