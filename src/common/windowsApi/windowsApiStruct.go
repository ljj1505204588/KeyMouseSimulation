package windowsApi

import (
	"sync"
	"syscall"
)

type dllLazyT struct {
	sync.Mutex
	dllName    string
	dll        *syscall.DLL
	dllFuncMap map[string]*syscall.Proc
}

// 包含Point信息
type POINT struct {
	X int32
	Y int32
}

// 包含GetMessage中MSG
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uint32
	LParam  uint32
	Time    uint32
	POINT
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}
