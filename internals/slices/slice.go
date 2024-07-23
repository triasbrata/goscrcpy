package slices

func Entries[T any, V any](s []T, cb func(it T) (string, V)) map[string]V {
	res := make(map[string]V)
	for _, it := range s {
		key, val := cb(it)
		res[key] = val
	}
	return res
}

func Map[T any, V any](s []T, cb func(it T) V) []V {
	res := make([]V, 0)
	for _, it := range s {
		val := cb(it)
		res = append(res, val)
	}
	return res
}
