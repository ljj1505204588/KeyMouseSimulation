package windowsHook

import windowsApi "KeyMouseSimulation/common/windowsApiTool"

//详情看这 https://docs.microsoft.com/zh-cn/windows/win32/winmsg/hooks

//在创建窗口之前，包含传递给 符合 _ CBT 挂钩过程 CBTProc的信息。
type STRUCT_CBT_CREATEWND struct {

}

//在激活窗口之前，包含传递给 符合 _ CBT 挂钩过程 CBTProc的信息。
type STRUCT_CBTACTIVATESTRUCT struct {

}

//定义传递给 符合 _ CALLWNDPROCRET 挂钩过程 CallWndRetProc的消息参数。
type STRUCT_CWPRETSTRUCT struct {

}

//定义传递给 符合 _ CALLWNDPROC 挂钩过程 CALLWNDPROC的消息参数。
type STRUCT_CWPSTRUCT struct {

}

//包含传递给 WH _ DEBUG 挂钩过程 DebugProc 的调试信息。
type STRUCT_DEBUGHOOKINFO struct {

}

//包含有关发送到系统消息队列的硬件消息的信息。 此结构用于存储 JournalPlaybackProc 回调函数的消息信息。
type STRUCT_EVENTMSG struct {

}

//包含有关低级别键盘输入事件的信息。
type STRUCT_KBDLLHOOKSTRUCT struct {
	VkCode uint32
	ScanCode uint32
	Flags uint32
	Time uint32
	DwExtraInfo uint32
}

//包含有关传递给 WH _ MOUSE 挂钩过程 MouseProc 的鼠标事件的信息。
type STRUCT_MOUSEHOOKSTRUCT struct {

}

//包含有关传递给 WH _ MOUSE 挂钩过程 MouseProc 的鼠标事件的信息。
type STRUCT_MOUSEHOOKSTRUCTEX struct {

}

//包含有关低级别鼠标输入事件的信息。
type STRUCT_MSLLHOOKSTRUCT struct {
	windowsApi.POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DWExtraInfo uint32
}

/*
	------------------------------------------ Hook EVENT ------------------------------------------
*/

// KeyboardEvent contains information about keyboard input event.
type KeyboardEvent struct {
	Message Message
	STRUCT_KBDLLHOOKSTRUCT
}

// MouseEvent contains information about mouse input event.
type MouseEvent struct {
	Message Message
	STRUCT_MSLLHOOKSTRUCT
}