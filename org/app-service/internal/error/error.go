package errors

import "context"

type Cause interface {
	~uint
}

type Error[T Cause] struct {
	Cause     T
	rootCause any
	Message   string
}

func NewError[T Cause](ctx context.Context, cause T, message string, values ...any) Error[T] {
	return Error[T]{Cause: cause, Message: message}
}

func FromError[T Cause, R Cause](cause T, origin Error[R]) Error[T] {
	return Error[T]{Cause: cause, rootCause: origin.Cause, Message: origin.Message}
}

func Ok[T Cause]() Error[T] {
	return Error[T]{}
}

// Error message.
//
// Satisfies the error interface.
func (self Error[T]) Error() string {
	return self.String()
}

func (self Error[T]) IsOk() bool {
	return self.Cause == 0
}

func (self Error[T]) IsErr() bool {
	return self.Cause != 0
}

func (self Error[T]) String() string {
	return self.Message
}
