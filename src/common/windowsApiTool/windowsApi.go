package windowsApi

import (
	"errors"
	"strconv"
	"syscall"
	"unsafe"
)

var DllUser dllLazyT

func init() {
	DllUser = dllLazyT{
		dllName:    "user32.dll",
		dll:        nil,
		dllFuncMap: make(map[string]*syscall.Proc),
	}
	err := DllUser.dllLazyInit(DllUser.dllName)
	if err != nil {
		panic(err.Error())
	}
}

//DllLazyInit 初始化
func (T *dllLazyT) dllLazyInit(dllName string) error {
	T.Lock()
	defer T.Unlock()

	var err error
	DllUser.dll, err = syscall.LoadDLL(dllName)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

//GetUserDllFunc get window which is Active
func (T *dllLazyT) getUserDllFunc(funcName string) (*syscall.Proc, error) {
	T.Lock()
	defer T.Unlock()

	if DllUser.dllFuncMap[funcName] == nil {
		var err error
		if T.dll == nil {
			err = T.dllLazyInit(T.dllName)
			if err != nil {
				return nil, err
			}
		}
		DllUser.dllFuncMap[funcName], err = T.dll.FindProc(funcName)
		if err != nil {
			return nil, err
		}
	}
	return DllUser.dllFuncMap[funcName], nil
}

//Call use func in by funcName
func (T *dllLazyT) Call(funcName string, a ...uintptr) (r1, r2 uintptr, lastErr error) {
	if len(a) > 18 {
		return 0, 0, errors.New("Call " + funcName + " with too many arguments " + strconv.Itoa(len(a)) + ".")
	}

	proc, err := T.getUserDllFunc(funcName)
	if err != nil {
		return 0, 0, err
	}

	r1, r2, lastErr = proc.Call(a...)
	if lastErr.Error() == "The operation completed successfully." {
		lastErr = nil
	}

	return
}
func (T *dllLazyT) CallUnsafePoint(funcName string, a ...unsafe.Pointer) (r1, r2 uintptr, lastErr error) {
	if len(a) > 18 {
		return 0, 0, errors.New("Call " + funcName + " with too many arguments " + strconv.Itoa(len(a)) + ".")
	}

	proc, err := T.getUserDllFunc(funcName)
	if err != nil {
		return 0, 0, err
	}

	aGroup := make([]uintptr, len(a))
	for p := range a {
		aGroup = append(aGroup, uintptr(a[p]))
	}

	r1, r2, lastErr = proc.Call(aGroup...)
	if lastErr.Error() == "The operation completed successfully." {
		lastErr = nil
	}

	return
}
