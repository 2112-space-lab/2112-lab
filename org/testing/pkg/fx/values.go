package fx

// Values takes a map as an input and returns a slice of all values within this map
func Values[K comparable, V any](input map[K]V) []V {
	res := make([]V, 0, len(input))
	for _, v := range input {
		res = append(res, v)
	}
	return res
}

// FilterValues returns filtered slice with items matching shouldKeep result
func FilterValues[K comparable, V any](input map[K]V, shouldKeep func(item V) bool) []V {
	filtered := make([]V, 0, len(input))
	for _, item := range input {
		if shouldKeep(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// FilterAndMapValues returns filtered slice with items matching shouldKeep result and apply mapper to them
func FilterAndMapValues[K comparable, V any, W any](input map[K]V, shouldKeep func(item V) bool, mapper Mapper[V, W]) ([]W, error) {
	errs := []error{}
	filteredAndMapped := make([]W, 0, len(input))
	for _, item := range input {
		if shouldKeep(item) {
			mapped, err := mapper(item)
			if err != nil {
				errs = append(errs, err)
			}
			filteredAndMapped = append(filteredAndMapped, mapped)
		}
	}
	err := FlattenErrorsIfAny(errs...)
	if err != nil {
		return filteredAndMapped, err
	}
	return filteredAndMapped, nil
}

// MapValues returns items and apply mapper to them
func MapValues[K comparable, V any, W any](input map[K]V, mapper Mapper[V, W]) ([]W, error) {
	errs := []error{}
	filteredAndMapped := make([]W, 0, len(input))
	for _, item := range input {
		mapped, err := mapper(item)
		if err != nil {
			errs = append(errs, err)
		}
		filteredAndMapped = append(filteredAndMapped, mapped)
	}
	err := FlattenErrorsIfAny(errs...)
	if err != nil {
		return filteredAndMapped, err
	}
	return filteredAndMapped, nil
}

// MergeMaps merges input maps into a single one
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	res := make(map[K]V, 0)
	for _, toAdd := range maps {
		for k, v := range toAdd {
			res[k] = v
		}
	}
	return res
}
