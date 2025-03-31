package gene

func Keys[K comparable, V any](m map[K]V) (ks []K) {
	for k := range m {
		ks = append(ks, k)
	}
	return
}

func Values[K comparable, V any](m map[K]V) (vs []V) {
	for k := range m {
		vs = append(vs, m[k])
	}
	return
}

func Reverse[K comparable, V comparable](m map[K]V) map[V]K {
	r := make(map[V]K)
	for k, v := range m {
		r[v] = k
	}
	return r
}
