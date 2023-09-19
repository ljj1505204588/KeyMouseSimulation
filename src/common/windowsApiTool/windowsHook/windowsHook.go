package windowsHook

import (
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"fmt"
	"sync"
	"syscall"
	"time"
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
				if r < 0 {
					// We don't care what's went wrong, ignore the result value.
					continue
				} else {
					_, _, _ = windowsApi.DllUser.Call(windowsApi.FuncTranslateMessage, uintptr(unsafe.Pointer(&msg)))
					_, _, _ = windowsApi.DllUser.Call(windowsApi.FuncDispatchMessageW, uintptr(unsafe.Pointer(&msg)))
				}
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

/*
 * 具体实现  ---  键盘
 */

type keyBoardHookHandler func() (HOOKPROC, chan *KeyboardEvent)

func keyBoardDefaultHookHandler() (HOOKPROC, chan *KeyboardEvent) {
	c := make(chan *KeyboardEvent, 3000)
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			keyboardEvent := KeyboardEvent{
				Message:                Message(wParam),
				STRUCT_KBDLLHOOKSTRUCT: *(*STRUCT_KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
				RecordTime:             time.Now().UnixNano(),
			}

			select {
			case c <- &keyboardEvent:
				r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
				return r
			default:
				r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
				return r
			}
		}
		r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
		return r
	}, c
}

func KeyBoardHook(h keyBoardHookHandler) (chan *KeyboardEvent, error) {
	hook.Lock()
	defer hook.Unlock()

	if h == nil {
		h = keyBoardDefaultHookHandler
	}
	if hook.HandlerM == nil {
		hook.HandlerM = make(map[IdHook]interface{})
	}

	fn, ch := h()
	hook.Unlock()
	err := install(WH_KEYBOARD_LL, fn)
	hook.Lock()
	if err != nil {
		return nil, err
	}

	hook.HandlerM[WH_KEYBOARD_LL] = h
	return ch, nil
}
func KeyBoardUnhook() error {
	return uninstall(WH_KEYBOARD_LL)
}

/*
 * 具体实现  ---   鼠标
 */
type mouseHookHandler func() (HOOKPROC, chan *MouseEvent)

func mouseDefaultHookHandler() (HOOKPROC, chan *MouseEvent) {
	c := make(chan *MouseEvent, 3000)
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			mouseEvent := MouseEvent{
				Message:               Message(wParam),
				STRUCT_MSLLHOOKSTRUCT: *(*STRUCT_MSLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
				RecordTime:            time.Now().UnixNano(),
			}

			select {
			case c <- &mouseEvent:
				r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
				return r
			default:
				r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
				return r
			}
		}

		r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(code), wParam, lParam)
		return r
	}, c
}
func MouseHook(h mouseHookHandler) (chan *MouseEvent, error) {
	hook.Lock()
	defer hook.Unlock()

	if h == nil {
		h = mouseDefaultHookHandler
	}
	if hook.HandlerM == nil {
		hook.HandlerM = make(map[IdHook]interface{})
	}

	fn, ch := h()
	hook.Unlock()
	err := install(WH_MOUSE_LL, fn)
	hook.Lock()
	if err != nil {
		return nil, err
	}

	hook.HandlerM[WH_MOUSE_LL] = h
	return ch, nil
}

func MouseUnhook() error {
	return uninstall(WH_MOUSE_LL)
}
