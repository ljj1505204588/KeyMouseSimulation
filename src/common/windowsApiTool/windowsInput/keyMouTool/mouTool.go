package keyMouTool

import (
	KeyboardMouseInput "KeyMouseSimulation/pkg/common/windowsApiTool/windowsInput"
	"sync"
)

var mouseChan chan *MouseInputT
var MouseInputPool = sync.Pool{
	New: func() interface{} {
		return new(KeyboardMouseInput.MouseInputT)
	},
}

// GetMouseSendInputChan 创建一个向计算机底层发送信息的通道
func GetMouseSendInputChan(size int) (chan *MouseInputT, error) {

	var createRoutine = false

	if mouseChan == nil {
		mouseChan = make(chan *MouseInputT, size)
		createRoutine = true
	}
	if len(mouseChan) != size {
		close(mouseChan)
		mouseChan = make(chan *MouseInputT, size)
		createRoutine = true
	}

	if createRoutine {
		go mouseRoutine(mouseChan)
	}

	return mouseChan, nil
}

func mouseRoutine(c chan *MouseInputT) {

	defer func() {
		if p := recover(); p != nil {
			go mouseRoutine(c)
		}
	}()

	//TODO 考虑批量
	for even := range c {
		MouseInput := MouseInputPool.Get().(*KeyboardMouseInput.MouseInputT)

		MouseInput.Type = uint32(TYPE_INPUT_MOUSE)
		MouseInput.Mi.X, MouseInput.Mi.Y = even.X, even.Y
		MouseInput.Mi.DwFlags = uint32(even.DWFlags)
		MouseInput.Mi.MouseData = even.MouseData
		MouseInput.Mi.Time = even.Time

		_, err := KeyboardMouseInput.MouseInput(*MouseInput)
		if err != nil {
			//No Importance
		}

		MouseInputPool.Put(MouseInput)
	}
}
