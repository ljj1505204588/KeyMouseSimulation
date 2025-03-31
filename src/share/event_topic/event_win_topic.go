package event_topic

import (
	"KeyMouseSimulation/common/windowsApi/windowsHook"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	"KeyMouseSimulation/pkg/event"
)

const WindowsMouseHook eventCenter.Topic = "windows_mouse_hook"        // 鼠标钩子
const WindowsKeyBoardHook eventCenter.Topic = "windows_key_board_hook" // 键盘钩子

const WindowsMouseInput eventCenter.Topic = "windows_mouse_input"        // 鼠标输入
const WindowsKeyBoardInput eventCenter.Topic = "windows_key_board_input" // 键盘输入

type WindowsMouseHookData struct {
	Date *windowsHook.MouseEvent // 系统鼠标勾子
}
type WindowsKeyBoardHookData struct {
	Date *windowsHook.KeyboardEvent // 系统键盘勾子
}

type WindowsMouseInputData struct {
	Data *keyMouTool.MouseInputT // 鼠标输入数据
}
type WindowsKeyBoardInputData struct {
	Data *keyMouTool.KeyInputT // 键盘输入数据
}
