package fx

// Keys return a slice of all keys inside a map
func Keys[T comparable, U any](input map[T]U) []T {
	keys := make([]T, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}
