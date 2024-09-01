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
