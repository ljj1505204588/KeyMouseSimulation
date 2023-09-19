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
