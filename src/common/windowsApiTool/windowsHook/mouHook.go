package windowsHook

import (
	windowsApi "KeyMouseSimulation/pkg/common/windowsApiTool"
	"time"
	"unsafe"
)

type mouseHookHandler func() (HOOKPROC, chan *MouseEvent)

func mouseDefaultHookHandler() (HOOKPROC, chan *MouseEvent) {
	c := make(chan *MouseEvent, 30000)
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
