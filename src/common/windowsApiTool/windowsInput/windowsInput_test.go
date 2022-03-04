package windowsInput_test

import (
	"fmt"
	. "KeyMouseSimulation/common/windowsApiTool/windowsInput"
	"syscall"
	"testing"
	"unsafe"
)

func TestKeyboardInput(t *testing.T) {
	var input = make([]KeyboardInputT,4)
	input[0].Type = 1
	input[1].Type = 1
	input[2].Type = 1
	input[3].Type = 1

	input[0].Ki.WVk = 0x5B
	input[1].Ki.WVk = 0x44
	input[2].Ki.WVk = 0x5B
	input[2].Ki.DwFlags = 0x0002
	input[3].Ki.WVk = 0x44
	input[3].Ki.DwFlags = 0x0002

	fmt.Println(KeyboardInput(input...))
}
func BenchmarkKeyboardInput(b *testing.B) {
	var input = make([]KeyboardInputT,2)
	input[0].Type = 1
	input[1].Type = 1

	input[0].Ki.WVk = 0x44
	input[1].Ki.WVk = 0x44
	input[1].Ki.DwFlags = 0x0002

	for i := 0; i < b.N; i++ {
		KeyboardInput(input...)
		//fmt.Println(KeyboardInput(input...))
	}
}

func TestMouseInput(t *testing.T) {
	input := make([]MouseInputT,2)

	input[0].Type = 0
	input[1].Type = 0

	input[0].Mi.X = 3042 * 65555 / 1920
	input[0].Mi.Y = 475 * 65555 / 1080
	input[0].Mi.DwFlags = 0x8007

	input[1].Mi.X = 3105 * 65555 / 1920
	input[1].Mi.Y = 357 * 65555 / 1080
	input[1].Mi.DwFlags = 0x8007


	fmt.Println(MouseInput(input...))
	return
}

func BenchmarkMouseInput(b *testing.B) {
	input := MouseInputT{}

	input.Type = 0
	input.Mi.DwFlags = 0x0007

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MouseInput(input)
	}
}
func BenchmarkMouseInput2(b *testing.B) {
	input := MouseInputT{}

	input.Type = 0
	input.Mi.DwFlags = 0x0007

	dll,_ := syscall.LoadDLL("user32.dll")
	proc,_ := dll.FindProc("SendInput")


	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc.Call(uintptr(1),uintptr(unsafe.Pointer(&input)),unsafe.Sizeof(MouseInputT{}))

	}
}

func TestHardWareInput(t *testing.T) {
	return
}