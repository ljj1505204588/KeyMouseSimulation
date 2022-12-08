package keyMouTool_test

import (
	. "KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"testing"
)

func TestGetKeySendInputChan(t *testing.T) {
	c, err := GetKeySendInputChan(100)
	if err != nil {
		return
	}
	c <- KeyInputT{
		VK:      VK_6,
		DwFlags: 0,
	}
	c <- KeyInputT{
		VK:      VK_6,
		DwFlags: DW_KEYEVENTF_KEYUP,
	}
}

/*
	cpu: Intel(R) Core(TM) i5-6300HQ CPU @ 2.30GHz
	BenchmarkGetKeySendInputChan
	BenchmarkGetKeySendInputChan-4   	    6666	    307189 ns/op
*/
func BenchmarkGetKeySendInputChan(b *testing.B) {
	c, err := GetKeySendInputChan(100)
	if err != nil {
		return
	}

	event := KeyInputT{
		VK: VK_6,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		event.DwFlags = 0
		c <- event

		event.DwFlags = DW_KEYEVENTF_KEYUP
		c <- event
	}
}

func TestGeyMouseSendInputChan(t *testing.T) {
	c, err := GetMouseSendInputChan(100)
	if err != nil {
		return
	}
	c <- MouseInputT{
		X:       0,
		Y:       0,
		DWFlags: DW_MOUSEEVENTF_LEFTDOWN | DW_MOUSEEVENTF_LEFTUP,
	}
}

/*
	cpu: Intel(R) Core(TM) i5-6300HQ CPU @ 2.30GHz
	BenchmarkGeyMouseSendInputChan
	BenchmarkGeyMouseSendInputChan-4   	   10000	    411393 ns/op
*/
func BenchmarkGeyMouseSendInputChan(b *testing.B) {
	c, err := GetMouseSendInputChan(100)
	if err != nil {
		return
	}

	event := MouseInputT{
		DWFlags: DW_MOUSEEVENTF_LEFTDOWN | DW_MOUSEEVENTF_LEFTUP,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- event
	}
}
