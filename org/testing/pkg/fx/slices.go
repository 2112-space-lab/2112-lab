package fx

// ConcatSlices combines any slices of type []T into a single slice with all elements
func ConcatSlices[T any](slices ...[]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
}

func NegateFilter[T any](f func(item T) bool) (n func(T) bool) {
	return func(item T) bool {
		x := f(item)
		return !x
	}
}

// FilterSlice returns filtered slice with items matching shouldKeep result
func FilterSlice[T any](slice []T, shouldKeep func(item T) bool) []T {
	filtered := make([]T, 0, len(slice))
	for _, item := range slice {
		if shouldKeep(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// ToSliceWithMapper applies Mapper function to all elements in a map and returns as a slice
// by ignoring all keys and keeping just the values of the input map
func ToSliceWithMapper[K comparable, V any, U any](m map[K]V, mapper Mapper[V, U]) ([]U, error) {
	errs := []error{}
	values := make([]U, 0, len(m))
	for _, v := range m {
		mapped, err := mapper(v)
		if err != nil {
			errs = append(errs, err)
		}
		values = append(values, mapped)
	}

	err := FlattenErrorsIfAny(errs...)
	return values, err
}

// SliceContains returns true if the item is present in the slice otherwise returns false
func SliceContains[T comparable](slice []T, item T) bool {
	for _, r := range slice {
		if r == item {
			return true
		}
	}
	return false
}
