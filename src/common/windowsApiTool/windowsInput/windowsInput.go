package windowsInput

import (
	"fmt"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"unsafe"
)

func KeyboardInput(k ...KeyboardInputT)(int,error){
	if len(k) == 0 {
		return 0,fmt.Errorf("k len must bigger than 0 ")
	}

	r,_,err := windowsApi.DllUser.Call(windowsApi.FuncSendInput,uintptr(len(k)),uintptr(unsafe.Pointer(&k[0])),unsafe.Sizeof(KeyboardInputT{}))

	return int(r),err
}

func MouseInput(m ...MouseInputT)(int,error){
	if len(m) == 0 {
		return 0,fmt.Errorf("m len must bigger than 0 ")
	}

	r,_,err := windowsApi.DllUser.Call(windowsApi.FuncSendInput,uintptr(len(m)),uintptr(unsafe.Pointer(&m[0])),unsafe.Sizeof(MouseInputT{}))

	return int(r),err
}

func HardWareInput()(int,error){
	return 0,fmt.Errorf("Unrealized.")
}