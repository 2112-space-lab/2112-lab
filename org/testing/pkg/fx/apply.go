package fx

import (
	"context"
)

// ApplyIfNoError executes Mapper function with in parameter if there is no inErr, otherwise short-circuit return directly inErr
func ApplyIfNoError[T any, U any](in T, inErr error, f Mapper[T, U]) (res U, err error) {
	if inErr != nil {
		return res, inErr
	}
	res, err = f(in)
	return
}

// ApplyIfNoErrorContext executes Mapper function with in parameter if there is no inErr, otherwise short-circuit return directly inErr
func ApplyIfNoErrorContext[T any, U any](ctx context.Context, in T, inErr error, f MapperCtx[T, U]) (res U, err error) {
	if inErr != nil {
		return res, inErr
	}
	res, err = f(ctx, in)
	return
}
