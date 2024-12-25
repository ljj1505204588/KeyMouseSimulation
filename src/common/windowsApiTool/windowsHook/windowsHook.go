package windowsHook

import (
	windowsApi "KeyMouseSimulation/pkg/common/windowsApiTool"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

// todo 加一个优雅退出

type HOOKPROC func(code int32, wParam, lParam uintptr) uintptr

type hookT struct {
	sync.RWMutex
	M        map[IdHook]uintptr
	HandlerM map[IdHook]interface{}
}

var hook hookT

func install(id IdHook, fn interface{}) error {
	hook.Lock()
	defer hook.Unlock()

	if hook.M == nil {
		hook.M = make(map[IdHook]uintptr)
	}

	if hook.HandlerM == nil {
		hook.HandlerM = make(map[IdHook]interface{})
	}

	if hook.M[id] != 0 {
		return fmt.Errorf("hook function is already installed")
	}

	go func() {
		hhk, _, _ := windowsApi.DllUser.Call(windowsApi.FuncSetWindowsHookExW, uintptr(id), syscall.NewCallback(fn), 0, 0)

		if hhk == 0 {
			panic(" failed to install hook function")
		}

		hook.Lock()
		hook.M[id] = hhk
		hook.Unlock()

		var msg *windowsApi.MSG

		for {
			hook.RLock()
			if hook.M[id] == 0 {
				break
			}
			hook.RUnlock()

			r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncGetMessageW, uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if r != 0 {
				_, _, _ = windowsApi.DllUser.Call(windowsApi.FuncTranslateMessage, uintptr(unsafe.Pointer(&msg)))
				_, _, _ = windowsApi.DllUser.Call(windowsApi.FuncDispatchMessageW, uintptr(unsafe.Pointer(&msg)))
			}
		}
	}()

	return nil
}
func uninstall(id IdHook) error {
	hook.Lock()
	defer hook.Unlock()

	if hook.M[id] == 0 {
		return fmt.Errorf("unhook")
	}

	if r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncUnhookWindowsHookEx, hook.M[id]); r == 0 {
		return fmt.Errorf(" failed to uninstall hook function")
	}

	hook.M[id] = 0
	hook.HandlerM[id] = nil

	return nil
}
