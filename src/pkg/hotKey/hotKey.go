package hk

import (
	"KeyMouseSimulation/common/windowsApi"
	"KeyMouseSimulation/common/windowsApi/windowsHook"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/topic"
	"syscall"
	"time"
	"unsafe"
)

type hotKeyI interface {
	set(key enum.HotKey, mod, code keyMouTool.VKCode) // 设置热键
	show() string                                     // 显示热键标识
	vkCode() keyMouTool.VKCode                        // 获取vkcode
}

type hotKeyT struct {
	key  enum.HotKey    // 热键
	mod  uint32         // 修饰键
	vk   uint32         // 主键
	hook syscall.Handle // 钩子句柄

	hotKeyEffectTime int64 // 热键生效时间
}

// Set 设置热键
func (h *hotKeyT) set(key enum.HotKey, mod, code keyMouTool.VKCode) {
	if h.hook != 0 {
		h.cleanup()
	}

	h.key = key
	h.mod = uint32(mod)
	h.vk = uint32(code)

	// 创建钩子回调函数
	hookProc := func(nCode int32, wparam, lparam uintptr) uintptr {
		nowTime := time.Now().UnixMilli()
		if nCode >= 0 && nowTime-h.hotKeyEffectTime >= 100 {
			// 获取键盘事件结构
			kbdStruct := (*windowsHook.STRUCT_KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
			// 只处理按键按下事件，避免重复触发
			if wparam == uintptr(windowsHook.WM_KEYDOWN) || wparam == uintptr(windowsHook.WM_SYSKEYDOWN) {
				// 检查是否是目标热键
				if kbdStruct.VkCode == h.vk {
					// 检查修饰键状态
					if h.mod == 0 || (kbdStruct.Flags&h.mod) == h.mod {
						// 触发热键事件
						_ = eventCenter.Event.Publish(topic.HotKeyEffect, &topic.HotKeyEffectData{
							HotKey: h.key,
						})

						h.hotKeyEffectTime = nowTime
					}
				}
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

// vkCode
func (h *hotKeyT) vkCode() keyMouTool.VKCode {
	return keyMouTool.VKCode(h.vk)
}
