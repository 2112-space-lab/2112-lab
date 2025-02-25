package errors

type EventError uint

const (
	EventErrorInternalFailure = EventError(iota + 1)
	EventErrorNotFound
)
