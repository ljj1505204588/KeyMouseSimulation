package hk

import (
	"KeyMouseSimulation/common/windowsApi"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/event_topic"
	"syscall"
)

type hotKeyT struct {
	key  enum.HotKey    // 热键
	mod  uint32         // 修饰键
	vk   uint32         // 主键
	hook syscall.Handle // 钩子句柄
}

// Set 设置热键
func (h *hotKeyT) set(key enum.HotKey, mod, code keyMouTool.VKCode) {
	h.key = key
	h.mod = uint32(mod)
	h.vk = uint32(code)

	// 创建钩子回调函数
	hookProc := func(nCode int32, wparam, lparam uintptr) uintptr {
		if nCode >= 0 {
			// 检查是否是目标热键
			if wparam == uintptr(h.vk) {
				// 触发热键事件
				eventCenter.Event.ASyncPublish(event_topic.HotKeyEffect, event_topic.HotKeyEffectData{
					HotKey: h.key,
				})
			}
		}
		// 调用下一个钩子
		r, _, _ := windowsApi.DllUser.Call(windowsApi.FuncCallNextHookEx, 0, uintptr(nCode), wparam, lparam)
		return r
	}

	// 注册钩子
	hookHandle, _, err := windowsApi.DllUser.Call(
		windowsApi.FuncSetWindowsHookExW,
		uintptr(13), // WH_KEYBOARD_LL
		syscall.NewCallback(hookProc),
		0,
		0,
	)
	if hookHandle == 0 {
		panic(err)
	}
	h.hook = syscall.Handle(hookHandle)
}

// 清理钩子
func (h *hotKeyT) cleanup() {
	if h.hook != 0 {
		_, _, _ = windowsApi.DllUser.Call(windowsApi.FuncUnhookWindowsHookEx, uintptr(h.hook))
		h.hook = 0
	}
}

// 显热键对应按键
func (h *hotKeyT) show() (res string) {
	if h.mod != 0 {
		res = keyMouTool.VKCodeStringMapReverse[keyMouTool.VKCode(h.mod)] + "+"
	}

	res += keyMouTool.VKCodeStringMapReverse[keyMouTool.VKCode(h.vk)]

	return
}
