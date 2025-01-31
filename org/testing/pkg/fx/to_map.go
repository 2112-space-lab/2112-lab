package fx

// ToMapByKeySelector turns an array into a map using result of key selector as index.
// use with caution, if items are clashing with same key, only the last item with this key will remain in map
func ToMapByKeySelector[T any, K comparable](inputArray []T, keySelector func(T) K) map[K]T {
	res := make(map[K]T)
	for _, v := range inputArray {
		k := keySelector(v)
		res[k] = v
	}
	return res
}
