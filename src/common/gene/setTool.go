package gene

// Intersection 取交集
func Intersection[T comparable](ls ...[]T) (res []T) {
	var m = make(map[T]int)
	for lIndex, l := range ls {
		for _, per := range l {
			// 集合"互异性"
			if m[per] == lIndex {
				m[per] += 1
			}
		}
	}

	var lsLen = len(ls)
	for per, num := range m {
		if num == lsLen {
			res = append(res, per)
		}
	}
	return
}

// UnionSet 取并集
func UnionSet[T comparable](ls ...[]T) (res []T) {
	var m = make(map[T]struct{})
	for _, l := range ls {
		for _, per := range l {
			if _, ok := m[per]; !ok {
				m[per] = struct{}{}
				res = append(res, per)
			}
		}
	}
	return res
}

// Exclude 剔除
func Exclude[T comparable](sour []T, ex ...[]T) (res []T) {
	var m = make(map[T]int32)
	for _, perEx := range ex {
		for _, per := range perEx {
			m[per] -= 1
		}
	}

	for _, per := range sour {
		if m[per] == 0 {
			res = append(res, per)
		}
	}

	return res
}

// RemoveDuplicate 去重
func RemoveDuplicate[T comparable](ls ...[]T) (res []T) {
	var m = make(map[T]struct{})
	for _, l := range ls {
		for _, per := range l {
			m[per] = struct{}{}
		}
	}

	for per := range m {
		res = append(res, per)
	}
	return
}

// Contain 包含
func Contain[T comparable](l []T, n T) bool {
	for _, per := range l {
		if per == n {
			return true
		}
	}
	return false
}

// Equal 判断多个数组是否相等
func Equal[T comparable](ls ...[]T) bool {
	var m = make(map[int]T)
	for lPos, l := range ls {
		if len(l) != len(ls[0]) {
			return false
		}

		for index, per := range l {
			if lPos == 0 {
				m[index] = per
			} else if m[index] != per {
				return false
			}
		}
	}

	return true
}

// Choose 选择
func Choose[T any](j bool, x, y T) T {
	if j {
		return x
	}
	return y
}
