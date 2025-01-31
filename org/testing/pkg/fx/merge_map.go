package fx

// MergeMapOverrides creates a copy of original map and apply overrides for each key from subsequent maps
func MergeMapOverrides[TKey comparable, TVal any](original map[TKey]TVal, overrides ...map[TKey]TVal) map[TKey]TVal {
	res := make(map[TKey]TVal, len(original))
	for k, v := range original {
		res[k] = v
	}
	for _, o := range overrides {
		if o == nil {
			continue
		}
		for k, v := range o {
			res[k] = v
		}
	}
	return res
}
