package fx

// Result generic functional type to explicitly handle potential errors
type Result[T any] struct {
	Success T
	Err     error
}

// IsError returns true if result is created from NewFailResult constructor
func (r Result[T]) IsError() bool {
	return r.Err != nil
}

// IsSuccess returns true if result is created from NewSuccessResult constructor
func (r Result[T]) IsSuccess() bool {
	return r.Err == nil
}

// NewSuccessResult initializes an Result with success value
func NewSuccessResult[T any](value T) Result[T] {
	return Result[T]{
		Success: value,
		Err:     nil,
	}
}

// NewFailResult initializes an Result with an error
func NewFailResult[T any](err error) Result[T] {
	return Result[T]{
		Err: err,
	}
}
