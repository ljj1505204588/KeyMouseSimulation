package windowsHook

import (
	windowsApi "KeyMouseSimulation/pkg/common/windowsApiTool"
	"time"
	"unsafe"
)

type keyBoardHookHandler func() (HOOKPROC, chan *KeyboardEvent)

func keyBoardDefaultHookHandler() (HOOKPROC, chan *KeyboardEvent) {
	c := make(chan *KeyboardEvent, 30000)
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
