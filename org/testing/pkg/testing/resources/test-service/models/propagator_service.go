package models

type PropagatorSettings struct {
	NoradId         string
	TleLine1        string
	TleLine2        string
	StartTime       string
	DurationMinutes int32
	IntervalSeconds int32
}
