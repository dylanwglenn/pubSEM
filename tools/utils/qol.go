package utils

func MapValsToSlice[K comparable, V any](m map[K]V) []V {
	res := make([]V, len(m))
	i := 0
	for _, v := range m {
		res[i] = v
		i++
	}
	return res
}
