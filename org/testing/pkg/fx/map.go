package fx

import "context"

// Mapper generic function signature for mapping functions
type Mapper[T any, U any] func(T) (U, error)

// MapperCtx generic function signature for mapping functions with context
type MapperCtx[T any, U any] func(context.Context, T) (U, error)

// MapAll applies Mapper function to each element and returns slice of new elements
// or a flatten error containing all encountered mapping errors
func MapAll[T any, U any](inputs []T, mapper Mapper[T, U]) ([]U, error) {
	res := make([]U, 0, len(inputs))
	errs := make([]error, 0)
	for _, v := range inputs {
		mapped, err := mapper(v)
		if err != nil {
			errs = append(errs, err)
		}
		res = append(res, mapped)
	}
	if len(errs) > 0 {
		errSum := FlattenErrorsIfAny(errs...)
		return res, errSum
	}
	return res, nil
}

func MapStrings[T, U ~string](inputs []T) []U {
	res := make([]U, 0, len(inputs))
	for _, v := range inputs {
		mapped := U(v)
		res = append(res, mapped)
	}
	return res
}

func MapNumbers[T, U ~int32 | ~int64 | ~float32 | ~float64](inputs []T) []U {
	res := make([]U, 0, len(inputs))
	for _, v := range inputs {
		mapped := U(v)
		res = append(res, mapped)
	}
	return res
}

// PrepareMapper returns a function with embedded input mapper
func PrepareMapper[T any, U any](mapper Mapper[T, U]) func([]T) ([]U, error) {
	return func(inputs []T) ([]U, error) {
		mapped, err := MapAll(inputs, mapper)
		return mapped, err
	}
}

// MapAllMap applies Mapper function to each element in a map and returns a map of the new elements
// or a flatten error containing all encountered mapping errors
func MapAllMap[K comparable, U any, V any](m map[K]U, mapper Mapper[U, V]) (map[K]V, error) {
	errs := []error{}

	res := make(map[K]V, len(m))
	for k, v := range m {
		mapped, err := mapper(v)
		if err != nil {
			errs = append(errs, err)
		}

		res[k] = mapped
	}

	err := FlattenErrorsIfAny(errs...)
	return res, err
}

// Mappable allows for a type to convert into type M
type Mappable[M any] interface {
	MapTo() (M, error)
}

type mapFunc[M any] func() (M, error)

// MapTo apply mapping to type M
func (f mapFunc[M]) MapTo() (M, error) {
	return f()
}

// AsMappable turns a value of type S and a mapping function into a Mappable[M]
func AsMappable[S, M any](src S, mapper func(S) (M, error)) Mappable[M] {
	return mapFunc[M](func() (M, error) {
		return mapper(src)
	})
}

type mapFuncOption[M any] func() (Option[M], error)

// MapTo apply mapping to type M
func (f mapFuncOption[M]) MapTo() (Option[M], error) {
	return f()
}

// MappableOption allows for a type to convert into type M
type MappableOption[M any] interface {
	MapTo() (Option[M], error)
}

// AsMappableOption turns a value of type Option[S] and a mapping function into a MappableOption[M]
func AsMappableOption[S, M any](src Option[S], mapper func(S) (M, error)) MappableOption[M] {
	optMapper := ToOptionalMapper(mapper)
	return mapFuncOption[M](func() (Option[M], error) {
		return optMapper(src)
	})
}

// ToOptionalMapper convert a mapper func to handler Optional value
func ToOptionalMapper[S, M any](mapper func(S) (M, error)) func(Option[S]) (Option[M], error) {
	return func(s Option[S]) (Option[M], error) {
		if s.HasValue {
			mapped, err := mapper(s.Value)
			return NewValueOption(mapped), err
		}
		return NewEmptyOption[M](), nil
	}
}
