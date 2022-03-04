package windowsHook

type IdHook int

const (
	//Installs a hook procedure that monitors messages before the system sends them to the destination window procedure. For more information, see the CallWndProc hook procedure.
	WH_CALLWNDPROC = IdHook(4)
	//Installs a hook procedure that monitors messages after they have been processed by the destination window procedure. For more information, see the CallWndRetProc hook procedure.
	WH_CALLWNDPROCRET = IdHook(12)
	//Installs a hook procedure that receives notifications useful to a CBT application. For more information, see the CBTProc hook procedure.
	WH_CBT = IdHook(5)
	//Installs a hook procedure useful for debugging other hook procedures. For more information, see the DebugProc hook procedure.
	WH_DEBUG = IdHook(9)
	//Installs a hook procedure that will be called when the application's foreground thread is about to become idle. This hook is useful for performing low priority tasks during idle time. For more information, see the ForegroundIdleProc hook procedure.
	WH_FOREGROUNDIDLE = IdHook(11)
	//Installs a hook procedure that monitors messages posted to a message queue. For more information, see the GetMsgProc hook procedure.
	WH_GETMESSAGE = IdHook(3)
	//Installs a hook procedure that posts messages previously recorded by a WH_JOURNALRECORD hook procedure. For more information, see the JournalPlaybackProc hook procedure.
	WH_JOURNALPLAYBACK = IdHook(1)
	//Installs a hook procedure that records input messages posted to the system message queue. This hook is useful for recording macros. For more information, see the JournalRecordProc hook procedure.
	WH_JOURNALRECORD = IdHook(0)
	//Installs a hook procedure that monitors keystroke messages. For more information, see the KeyboardProc hook procedure.
	WH_KEYBOARD = IdHook(2)
	//Installs a hook procedure that monitors low-level keyboard input events. For more information, see the LowLevelKeyboardProc hook procedure.
	WH_KEYBOARD_LL = IdHook(13)
	//Installs a hook procedure that monitors mouse messages. For more information, see the MouseProc hook procedure.
	WH_MOUSE = IdHook(7)
	//Installs a hook procedure that monitors low-level mouse input events. For more information, see the LowLevelMouseProc hook procedure.
	WH_MOUSE_LL = IdHook(14)
	//Installs a hook procedure that monitors messages generated as a result of an input event in a dialog box, message box, menu, or scroll bar. For more information, see the MessageProc hook procedure.
	WH_MSGFILTER = IdHook(-1)
	//Installs a hook procedure that receives notifications useful to shell applications. For more information, see the ShellProc hook procedure.
	WH_SHELL = IdHook(10)
	//Installs a hook procedure that monitors messages generated as a result of an input event in a dialog box, message box, menu, or scroll bar. The hook procedure monitors these messages for all applications in the same desktop as the calling thread. For more information, see the SysMsgProc hook procedure.
	WH_SYSMSGFILTER = IdHook(6)
)

type Message uintptr

const (
	WM_MOUSEMOVE         Message = 0x0200
	WM_MOUSERIGHTDOWN    Message = 0x0204
	WM_MOUSERIGHTUP      Message = 0x0205
	WM_MOUSEMIDDLEDOWN   Message = 0x0207
	WM_MOUSEMIDDLEUP     Message = 0x0208
	WM_MOUSELEFTDOWN     Message = 0x0201
	WM_MOUSELEFTUP       Message = 0x0202
	WM_MOUSELEFTSICEDOWN Message = 0x020B
	WM_MOUSELEFTSICEUP   Message = 0x020C
	WM_MOUSEPULLEY       Message = 0x020A

	WM_KEYDOWN    Message = 0x0100
	WM_KEYUP      Message = 0x0101
	WM_SYSKEYDOWN Message = 0x0104
	WM_SYSKEYUP   Message = 0x0105
)
