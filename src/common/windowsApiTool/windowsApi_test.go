package windowsApi

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func Test_test(T *testing.T) {
	type RectT struct {
		Left   uint32
		Top    uint32
		Right  uint32
		Bottom uint32
	}
	for {
		time.Sleep(2 * time.Second)
		HWnd,_,err := DllUser.Call("GetForegroundWindow")
		if err != nil && err.Error() != "The operation completed successfully."{
			fmt.Println(err.Error())
		} else {
			fmt.Println("当前窗口HWnd：", HWnd)
		}
		var rect RectT
		_,_,err = DllUser.Call("GetWindowRect",HWnd,uintptr(unsafe.Pointer(&rect)))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(fmt.Sprintf("%v", rect))
		}

		//按键打出hello world
		var wParam, lParam uintptr
		var msg uint32
		//wParam = uintptr(uint32(0x48))
		wParam = uintptr(0x02)
		times := 0
		click := 0
		for {
			times ++ ; click ++
			fmt.Println("新一轮按键:" + fmt.Sprintf("%d", wParam))
			msg = 0x100 //keyDown
			_,_,err = DllUser.Call(FuncSendMessageW,HWnd,uintptr(msg),wParam,lParam)
			if err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep(50 * time.Millisecond)
			//msg = 0x0102 //keyUp
			//_, err = SendMessage(dll, HWnd, msg, wParam, lParam)
			//if err != nil {
			//	fmt.Println(err.Error())
			//}

			if wParam == 0x27 {
				wParam = 0x28
			}else {
				wParam = 0x27
			}
			if times == 100 {
				times = 0
				wParam = 0x25
			}
			if click == 10000 {
				wParam = 0x01
				click = 0
			}

			//wParam += 1
			//if wParam > 0x5A {
			//	wParam -= 24
			//}

		}

	}
}

func Test_test1(T *testing.T) {
	for {
		//DllUser.Call(FuncPostMessageA,0xffff,0,0x0200,100<<16+500)
		DllUser.Call(FuncSendInput)
		time.Sleep(3*time.Second)
		fmt.Println(time.Now())
	}
}

func abort(funcname string, err string) {
	panic(funcname + " failed: " + err)
}

func print_version(v uint32) {
	major := byte(v)
	minor := uint8(v >> 8)
	build := uint16(v >> 16)
	print("windows version ", major, ".", minor, " (Build ", build, ")\n")
}


