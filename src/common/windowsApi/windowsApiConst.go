package windowsApi

/*Windows and Messages*/
const (
	FuncAdjustWindowRect             = "AdjustWindowRect"
	FuncAdjustWindowRectEx           = "AdjustWindowRectEx"
	FuncAllowSetForegroundWindow     = "AllowSetForegroundWindow"
	FuncAnimateWindow                = "AnimateWindow"
	FuncAnyPopup                     = "AnyPopup"
	FuncArrangeIconicWindows         = "ArrangeIconicWindows"
	FuncBeginDeferWindowPos          = "BeginDeferWindowPos"
	FuncBringWindowToTop             = "BringWindowToTop"
	FuncBroadcastSystemMessage       = "BroadcastSystemMessage"
	FuncBroadcastSystemMessageA      = "BroadcastSystemMessageA"
	FuncBroadcastSystemMessageExA    = "BroadcastSystemMessageExA"
	FuncBroadcastSystemMessageExW    = "BroadcastSystemMessageExW"
	FuncBroadcastSystemMessageW      = "BroadcastSystemMessageW"
	FuncCalculatePopupWindowPosition = "CalculatePopupWindowPosition"
	FuncCallMsgFilterA               = "CallMsgFilterA"
	FuncCallMsgFilterW               = "CallMsgFilterW"
	FuncCallNextHookEx               = "CallNextHookEx"
	FuncCallWindowProcA              = "CallWindowProcA"
	FuncCallWindowProcW              = "CallWindowProcW"
	FuncCascadeWindows               = "CascadeWindows"
	FuncChangeWindowMessageFilter    = "ChangeWindowMessageFilter"
	FuncChangeWindowMessageFilterEx  = "ChangeWindowMessageFilterEx"
	FuncChildWindowFromPoint         = "ChildWindowFromPoint"
	FuncChildWindowFromPointEx       = "ChildWindowFromPointEx"
	FuncCloseWindow                  = "CloseWindow"
	FuncCreateMDIWindowA             = "CreateMDIWindowA"
	FuncCreateMDIWindowW             = "CreateMDIWindowW"
	FuncCreateWindowExA              = "CreateWindowExA"
	FuncCreateWindowExW              = "CreateWindowExW"
	FuncDeferWindowPos               = "DeferWindowPos"
	FuncDefFrameProcA                = "DefFrameProcA"
	FuncDefFrameProcW                = "DefFrameProcW"
	FuncDefMDIChildProcA             = "DefMDIChildProcA"
	FuncDefMDIChildProcW             = "DefMDIChildProcW"
	FuncDefWindowProcA               = "DefWindowProcA"
	FuncDefWindowProcW               = "DefWindowProcW"
	FuncDeregisterShellHookWindow    = "DeregisterShellHookWindow"
	FuncDestroyWindow                = "DestroyWindow"
	FuncDispatchMessage              = "DispatchMessage"
	FuncDispatchMessageA             = "DispatchMessageA"
	FuncDispatchMessageW             = "DispatchMessageW"
	FuncEndDeferWindowPos            = "EndDeferWindowPos"
	FuncEndTask                      = "EndTask"
	FuncEnumChildWindows             = "EnumChildWindows"
	FuncEnumPropsA                   = "EnumPropsA"
	FuncEnumPropsExA                 = "EnumPropsExA"
	FuncEnumPropsExW                 = "EnumPropsExW"
	FuncEnumPropsW                   = "EnumPropsW"
	FuncEnumThreadWindows            = "EnumThreadWindows"
	FuncEnumWindows                  = "EnumWindows"
	FuncFindWindowA                  = "FindWindowA"
	FuncFindWindowExA                = "FindWindowExA"
	FuncFindWindowExW                = "FindWindowExW"
	FuncFindWindowW                  = "FindWindowW"
	FuncGetAltTabInfoA               = "GetAltTabInfoA"
	FuncGetAltTabInfoW               = "GetAltTabInfoW"
	FuncGetAncestor                  = "GetAncestor"
	FuncGetClassInfoA                = "GetClassInfoA"
	FuncGetClassInfoExA              = "GetClassInfoExA"
	FuncGetClassInfoExW              = "GetClassInfoExW"
	FuncGetClassInfoW                = "GetClassInfoW"
	FuncGetClassLongA                = "GetClassLongA"
	FuncGetClassLongPtrA             = "GetClassLongPtrA"
	FuncGetClassLongPtrW             = "GetClassLongPtrW"
	FuncGetClassLongW                = "GetClassLongW"
	FuncGetClassName                 = "GetClassName"
	FuncGetClassNameA                = "GetClassNameA"
	FuncGetClassNameW                = "GetClassNameW"
	FuncGetClassWord                 = "GetClassWord"
	FuncGetClientRect                = "GetClientRect"
	FuncGetDesktopWindow             = "GetDesktopWindow"
	FuncGetForegroundWindow          = "GetForegroundWindow"
	FuncGetGUIThreadInfo             = "GetGUIThreadInfo"
	FuncGetInputState                = "GetInputState"
	FuncGetLastActivePopup           = "GetLastActivePopup"
	FuncGetLayeredWindowAttributes   = "GetLayeredWindowAttributes"
	FuncGetMessage                   = "GetMessage"
	FuncGetMessageA                  = "GetMessageA"
	FuncGetMessageExtraInfo          = "GetMessageExtraInfo"
	FuncGetMessagePos                = "GetMessagePos"
	FuncGetMessageTime               = "GetMessageTime"
	FuncGetMessageW                  = "GetMessageW"
	FuncGetParent                    = "GetParent"
	FuncGetProcessDefaultLayout      = "GetProcessDefaultLayout"
	FuncGetPropA                     = "GetPropA"
	FuncGetPropW                     = "GetPropW"
	FuncGetQueueStatus               = "GetQueueStatus"
	FuncGetShellWindow               = "GetShellWindow"
	FuncGetSysColor                  = "GetSysColor"
	FuncGetSystemMetrics             = "GetSystemMetrics"
	FuncGetTitleBarInfo              = "GetTitleBarInfo"
	FuncGetTopWindow                 = "GetTopWindow"
	FuncGetWindow                    = "GetWindow"
	FuncGetWindowDisplayAffinity     = "GetWindowDisplayAffinity"
	FuncGetWindowInfo                = "GetWindowInfo"
	FuncGetWindowLongA               = "GetWindowLongA"
	FuncGetWindowLongPtrA            = "GetWindowLongPtrA"
	FuncGetWindowLongPtrW            = "GetWindowLongPtrW"
	FuncGetWindowLongW               = "GetWindowLongW"
	FuncGetWindowModuleFileNameA     = "GetWindowModuleFileNameA"
	FuncGetWindowModuleFileNameW     = "GetWindowModuleFileNameW"
	FuncGetWindowPlacement           = "GetWindowPlacement"
	FuncGetWindowRect                = "GetWindowRect"
	FuncGetWindowTextA               = "GetWindowTextA"
	FuncGetWindowTextLengthA         = "GetWindowTextLengthA"
	FuncGetWindowTextLengthW         = "GetWindowTextLengthW"
	FuncGetWindowTextW               = "GetWindowTextW"
	FuncGetWindowThreadProcessId     = "GetWindowThreadProcessId"
	FuncInSendMessage                = "InSendMessage"
	FuncInSendMessageEx              = "InSendMessageEx"
	FuncInternalGetWindowText        = "InternalGetWindowText"
	FuncIsChild                      = "IsChild"
	FuncIsGUIThread                  = "IsGUIThread"
	FuncIsHungAppWindow              = "IsHungAppWindow"
	FuncIsIconic                     = "IsIconic"
	FuncIsProcessDPIAware            = "IsProcessDPIAware"
	FuncIsWindow                     = "IsWindow"
	FuncIsWindowUnicode              = "IsWindowUnicode"
	FuncIsWindowVisible              = "IsWindowVisible"
	FuncIsZoomed                     = "IsZoomed"
	FuncKillTimer                    = "KillTimer"
	FuncLockSetForegroundWindow      = "LockSetForegroundWindow"
	FuncLogicalToPhysicalPoint       = "LogicalToPhysicalPoint"
	FuncMoveWindow                   = "MoveWindow"
	FuncOpenIcon                     = "OpenIcon"
	FuncPeekMessageA                 = "PeekMessageA"
	FuncPeekMessageW                 = "PeekMessageW"
	FuncPhysicalToLogicalPoint       = "PhysicalToLogicalPoint"
	FuncPostMessageA                 = "PostMessageA"
	FuncPostMessageW                 = "PostMessageW"
	FuncPostQuitMessage              = "PostQuitMessage"
	FuncPostThreadMessageA           = "PostThreadMessageA"
	FuncPostThreadMessageW           = "PostThreadMessageW"
	FuncRealChildWindowFromPoint     = "RealChildWindowFromPoint"
	FuncRealGetWindowClassA          = "RealGetWindowClassA"
	FuncRealGetWindowClassW          = "RealGetWindowClassW"
	FuncRegisterClassA               = "RegisterClassA"
	FuncRegisterClassExA             = "RegisterClassExA"
	FuncRegisterClassExW             = "RegisterClassExW"
	FuncRegisterClassW               = "RegisterClassW"
	FuncRegisterShellHookWindow      = "RegisterShellHookWindow"
	FuncRegisterWindowMessageA       = "RegisterWindowMessageA"
	FuncRegisterWindowMessageW       = "RegisterWindowMessageW"
	FuncRemovePropA                  = "RemovePropA"
	FuncRemovePropW                  = "RemovePropW"
	FuncReplyMessage                 = "ReplyMessage"
	FuncSendMessage                  = "SendMessage"
	FuncSendMessageA                 = "SendMessageA"
	FuncSendMessageCallbackA         = "SendMessageCallbackA"
	FuncSendMessageCallbackW         = "SendMessageCallbackW"
	FuncSendMessageTimeoutA          = "SendMessageTimeoutA"
	FuncSendMessageTimeoutW          = "SendMessageTimeoutW"
	FuncSendMessageW                 = "SendMessageW"
	FuncSendNotifyMessageA           = "SendNotifyMessageA"
	FuncSendNotifyMessageW           = "SendNotifyMessageW"
	FuncSetClassLongA                = "SetClassLongA"
	FuncSetClassLongPtrA             = "SetClassLongPtrA"
	FuncSetClassLongPtrW             = "SetClassLongPtrW"
	FuncSetClassLongW                = "SetClassLongW"
	FuncSetClassWord                 = "SetClassWord"
	FuncSetCoalescableTimer          = "SetCoalescableTimer"
	FuncSetForegroundWindow          = "SetForegroundWindow"
	FuncSetLayeredWindowAttributes   = "SetLayeredWindowAttributes"
	FuncSetMessageExtraInfo          = "SetMessageExtraInfo"
	FuncSetParent                    = "SetParent"
	FuncSetProcessDefaultLayout      = "SetProcessDefaultLayout"
	FuncSetProcessDPIAware           = "SetProcessDPIAware"
	FuncSetPropA                     = "SetPropA"
	FuncSetPropW                     = "SetPropW"
	FuncSetSysColors                 = "SetSysColors"
	FuncSetTimer                     = "SetTimer"
	FuncSetWindowDisplayAffinity     = "SetWindowDisplayAffinity"
	FuncSetWindowLongA               = "SetWindowLongA"
	FuncSetWindowLongPtrA            = "SetWindowLongPtrA"
	FuncSetWindowLongPtrW            = "SetWindowLongPtrW"
	FuncSetWindowLongW               = "SetWindowLongW"
	FuncSetWindowPlacement           = "SetWindowPlacement"
	FuncSetWindowPos                 = "SetWindowPos"
	FuncSetWindowsHookExA            = "SetWindowsHookExA"
	FuncSetWindowsHookExW            = "SetWindowsHookExW"
	FuncSetWindowTextA               = "SetWindowTextA"
	FuncSetWindowTextW               = "SetWindowTextW"
	FuncShowOwnedPopups              = "ShowOwnedPopups"
	FuncShowWindow                   = "ShowWindow"
	FuncShowWindowAsync              = "ShowWindowAsync"
	FuncSoundSentry                  = "SoundSentry"
	FuncSwitchToThisWindow           = "SwitchToThisWindow"
	FuncSystemParametersInfoA        = "SystemParametersInfoA"
	FuncSystemParametersInfoW        = "SystemParametersInfoW"
	FuncTileWindows                  = "TileWindows"
	FuncTranslateMDISysAccel         = "TranslateMDISysAccel"
	FuncTranslateMessage             = "TranslateMessage"
	FuncUnhookWindowsHookEx          = "UnhookWindowsHookEx"
	FuncUnregisterClassA             = "UnregisterClassA"
	FuncUnregisterClassW             = "UnregisterClassW"
	FuncUpdateLayeredWindow          = "UpdateLayeredWindow"
	FuncWaitMessage                  = "WaitMessage"
	FuncWindowFromPhysicalPoint      = "WindowFromPhysicalPoint"
	FuncWindowFromPoint              = "WindowFromPoint"

	FuncCallBackHOOKPROC        = "HOOKPROC"
	FuncCallBackPROPENUMPROCA   = "PROPENUMPROCA"
	FuncCallBackPROPENUMPROCEXA = "PROPENUMPROCEXA"
	FuncCallBackPROPENUMPROCEXW = "PROPENUMPROCEXW"
	FuncCallBackPROPENUMPROCW   = "PROPENUMPROCW"
	FuncCallBackSENDASYNCPROC   = "SENDASYNCPROC"
	FuncCallBackTIMERPROC       = "TIMERPROC"
	FuncCallBackWNDPROC         = "WNDPROC"

	/*Windows and Messages*/

)

/*			Winuser.h  结构体
详情看 https://docs.microsoft.com/zh-cn/windows/win32/api/winuser/ns-winuser-msllhookstruct
ALTTABINFO structure
ANIMATIONINFO structure
AUDIODESCRIPTION structure
BSMINFO structure
CBT_CREATEWNDA structure
CBT_CREATEWNDW structure
CBTACTIVATESTRUCT structure
CHANGEFILTERSTRUCT structure
CLIENTCREATESTRUCT structure
CREATESTRUCTA structure
CREATESTRUCTW structure
CWPRETSTRUCT structure
CWPSTRUCT structure
DEBUGHOOKINFO structure
EVENTMSG structure
GUITHREADINFO structure
STRUCT_KBDLLHOOKSTRUCT structure
MDICREATESTRUCTA structure
MDICREATESTRUCTW structure
MINIMIZEDMETRICS structure
MINMAXINFO structure
MOUSEHOOKSTRUCT structure
MOUSEHOOKSTRUCTEX structure
MSG structure
STRUCT_MSLLHOOKSTRUCT structure
NCCALCSIZE_PARAMS structure
NONCLIENTMETRICSA structure
NONCLIENTMETRICSW structure
STYLESTRUCT structure
TITLEBARINFO structure
TITLEBARINFOEX structure
UPDATELAYEREDWINDOWINFO structure
WINDOWINFO structure
WINDOWPLACEMENT structure
WINDOWPOS structure
WNDCLASSA structure
WNDCLASSEXA structure
WNDCLASSEXW structure
WNDCLASSW structure
*/
/*			Winuser.h 宏
详情看 https://docs.microsoft.com/zh-cn/windows/win32/api/winuser/ns-winuser-msllhookstruct
CreateWindowA macro
CreateWindowW macro
GetNextWindow macro
MAKELPARAM macro
MAKELRESULT macro
MAKEWPARAM macro
*/

/* Keyboard and Mouse Input*/
const (
	FuncActivateKeyboardLayout       = "ActivateKeyboardLayout"
	FuncBlockInput                   = "BlockInput"
	FuncDefRawInputProc              = "DefRawInputProc"
	FuncDragDetect                   = "DragDetect"
	FuncEnableWindow                 = "EnableWindow"
	FuncGetActiveWindow              = "GetActiveWindow"
	FuncGetAsyncKeyState             = "GetAsyncKeyState"
	FuncGetCapture                   = "GetCapture"
	FuncGetDoubleClickTime           = "GetDoubleClickTime"
	FuncGetFocus                     = "GetFocus"
	FuncGetKBCodePage                = "GetKBCodePage"
	FuncGetKeyboardLayout            = "GetKeyboardLayout"
	FuncGetKeyboardLayoutList        = "GetKeyboardLayoutList"
	FuncGetKeyboardLayoutNameA       = "GetKeyboardLayoutNameA"
	FuncGetKeyboardLayoutNameW       = "GetKeyboardLayoutNameW"
	FuncGetKeyboardState             = "GetKeyboardState"
	FuncGetKeyboardType              = "GetKeyboardType"
	FuncGetKeyNameTextA              = "GetKeyNameTextA"
	FuncGetKeyNameTextW              = "GetKeyNameTextW"
	FuncGetKeyState                  = "GetKeyState"
	FuncGetLastInputInfo             = "GetLastInputInfo"
	FuncGetMouseMovePointsEx         = "GetMouseMovePointsEx"
	FuncGetRawInputBuffer            = "GetRawInputBuffer"
	FuncGetRawInputData              = "GetRawInputData"
	FuncGetRawInputDeviceInfoA       = "GetRawInputDeviceInfoA"
	FuncGetRawInputDeviceInfoW       = "GetRawInputDeviceInfoW"
	FuncGetRawInputDeviceList        = "GetRawInputDeviceList"
	FuncGetRegisteredRawInputDevices = "GetRegisteredRawInputDevices"
	FuncIsWindowEnabled              = "IsWindowEnabled"
	Funckeybd_event                  = "keybd_event"
	FuncLoadKeyboardLayoutA          = "LoadKeyboardLayoutA"
	FuncLoadKeyboardLayoutW          = "LoadKeyboardLayoutW"
	FuncMapVirtualKeyA               = "MapVirtualKeyA"
	FuncMapVirtualKeyExA             = "MapVirtualKeyExA"
	FuncMapVirtualKeyExW             = "MapVirtualKeyExW"
	FuncMapVirtualKeyW               = "MapVirtualKeyW"
	Funcmouse_event                  = "mouse_event"
	FuncOemKeyScan                   = "OemKeyScan"
	FuncRegisterHotKey               = "RegisterHotKey"
	FuncRegisterRawInputDevices      = "RegisterRawInputDevices"
	FuncReleaseCapture               = "ReleaseCapture"
	FuncSendInput                    = "SendInput"
	FuncSetActiveWindow              = "SetActiveWindow"
	FuncSetCapture                   = "SetCapture"
	FuncSetDoubleClickTime           = "SetDoubleClickTime"
	FuncSetFocus                     = "SetFocus"
	FuncSetKeyboardState             = "SetKeyboardState"
	FuncSwapMouseButton              = "SwapMouseButton"
	FuncToAscii                      = "ToAscii"
	FuncToAsciiEx                    = "ToAsciiEx"
	FuncToUnicode                    = "ToUnicode"
	FuncToUnicodeEx                  = "ToUnicodeEx"
	FuncTrackMouseEvent              = "TrackMouseEvent"
	FuncUnloadKeyboardLayout         = "UnloadKeyboardLayout"
	FuncUnregisterHotKey             = "UnregisterHotKey"
	FuncVkKeyScanA                   = "VkKeyScanA"
	FuncVkKeyScanExA                 = "VkKeyScanExA"
	FuncVkKeyScanExW                 = "VkKeyScanExW"
	FuncVkKeyScanW                   = "VkKeyScanW"
)

/*
	结构体
	详情看 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
	HARDWAREINPUT structure
	INPUT structure
	KEYBDINPUT structure
	LASTINPUTINFO structure
	MOUSEINPUT structure
	MOUSEMOVEPOINT structure
	RAWHID structure
	RAWINPUT structure
	RAWINPUTDEVICE structure
	RAWINPUTDEVICELIST structure
	RAWINPUTHEADER structure
	RAWKEYBOARD structure
	RAWMOUSE structure
	RID_DEVICE_INFO structure
	RID_DEVICE_INFO_HID structure
	RID_DEVICE_INFO_KEYBOARD structure
	RID_DEVICE_INFO_MOUSE structure
	TRACKMOUSEEVENT structure
*/
/*
	宏
	详情看 https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
	GET_APPCOMMAND_LPARAM macro
	GET_DEVICE_LPARAM macro
	GET_FLAGS_LPARAM macro
	GET_KEYSTATE_LPARAM macro
	GET_KEYSTATE_WPARAM macro
	GET_NCHITTEST_WPARAM macro
	GET_RAWINPUT_CODE_WPARAM macro
	GET_WHEEL_DELTA_WPARAM macro
	GET_XBUTTON_WPARAM macro
	NEXTRAWINPUTBLOCK macro
*/

/*
GetSystemMetrics function  详情:https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getsystemmetrics
*/
const (
	SM_ARRANGE                     = 56
	SM_CLEANBOOT                   = 67
	SM_CMONITORS                   = 80
	SM_CMOUSEBUTTONS               = 43
	SM_CONVERTIBLESLATEMODE        = 0x2003
	SM_CXBORDER                    = 5
	SM_CXCURSOR                    = 13
	SM_CXDLGFRAME                  = 7
	SM_CXDOUBLECLK                 = 36
	SM_CXDRAG                      = 68
	SM_CXEDGE                      = 45
	SM_CXFIXEDFRAME                = 7
	SM_CXFOCUSBORDER               = 83
	SM_CXFRAME                     = 32
	SM_CXFULLSCREEN                = 16
	SM_CXHSCROLL                   = 21
	SM_CXHTHUMB                    = 10
	SM_CXICON                      = 11
	SM_CXICONSPACING               = 38
	SM_CXMAXIMIZED                 = 61
	SM_CXMAXTRACK                  = 59
	SM_CXMENUCHECK                 = 71
	SM_CXMENUSIZE                  = 54
	SM_CXMIN                       = 28
	SM_CXMINIMIZED                 = 57
	SM_CXMINSPACING                = 47
	SM_CXMINTRACK                  = 34
	SM_CXPADDEDBORDER              = 92
	SM_CXSCREEN                    = 0
	SM_CXSIZE                      = 30
	SM_CXSIZEFRAME                 = 32
	SM_CXSMICON                    = 49
	SM_CXSMSIZE                    = 52
	SM_CXVIRTUALSCREEN             = 78
	SM_CXVSCROLL                   = 2
	SM_CYBORDER                    = 6
	SM_CYCAPTION                   = 4
	SM_CYCURSOR                    = 14
	SM_CYDLGFRAME                  = 8
	SM_CYDOUBLECLK                 = 37
	SM_CYDRAG                      = 69
	SM_CYEDGE                      = 46
	SM_CYFIXEDFRAME                = 8
	SM_CYFOCUSBORDER               = 84
	SM_CYFRAME                     = 33
	SM_CYFULLSCREEN                = 17
	SM_CYHSCROLL                   = 3
	SM_CYICON                      = 12
	SM_CYICONSPACING               = 39
	SM_CYKANJIWINDOW               = 18
	SM_CYMAXIMIZED                 = 62
	SM_CYMAXTRACK                  = 60
	SM_CYMENU                      = 15
	SM_CYMENUCHECK                 = 72
	SM_CYMENUSIZE                  = 55
	SM_CYMIN                       = 29
	SM_CYMINIMIZED                 = 58
	SM_CYMINSPACING                = 48
	SM_CYMINTRACK                  = 35
	SM_CYSCREEN                    = 1
	SM_CYSIZE                      = 31
	SM_CYSIZEFRAME                 = 33
	SM_CYSMCAPTION                 = 51
	SM_CYSMICON                    = 50
	SM_CYSMSIZE                    = 53
	SM_CYVIRTUALSCREEN             = 79
	SM_CYVSCROLL                   = 20
	SM_CYVTHUMB                    = 9
	SM_DBCSENABLED                 = 42
	SM_DEBUG                       = 22
	SM_DIGITIZER                   = 94
	SM_IMMENABLED                  = 82
	SM_MAXIMUMTOUCHES              = 95
	SM_MEDIACENTER                 = 87
	SM_MENUDROPALIGNMENT           = 40
	SM_MIDEASTENABLED              = 74
	SM_MOUSEPRESENT                = 19
	SM_MOUSEHORIZONTALWHEELPRESENT = 91
	SM_MOUSEWHEELPRESENT           = 75
	SM_NETWORK                     = 63
	SM_PENWINDOWS                  = 41
	SM_REMOTECONTROL               = 0x2001
	SM_REMOTESESSION               = 0x1000
	SM_SAMEDISPLAYFORMAT           = 81
	SM_SECURE                      = 44
	SM_SERVERR2                    = 89
)
