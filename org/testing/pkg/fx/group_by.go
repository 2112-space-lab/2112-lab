package fx

// GroupBy turns an array into a map using result of key selector as index.
// Items with the same key are appended to an array of values
func GroupBy[T any, K comparable](inputArray []T, keySelector func(T) K) map[K][]T {
	res := make(map[K][]T)
	for _, v := range inputArray {
		k := keySelector(v)
		if x, ok := res[k]; ok {
			res[k] = append(x, v)
		} else {
			res[k] = []T{v}
		}
	}
	return (res)
}
