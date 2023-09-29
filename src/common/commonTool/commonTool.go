package commonTool

import "runtime"

var sysPath string

func GetSysPthSep() string {
	if sysPath == "" {
		if runtime.GOOS == "linux" {
			sysPath = "/"
		} else {
			sysPath = "\\"
		}
	}
	return sysPath
}

func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}

type Number interface {
	int32 | int | int64 | float32 | float64
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}
