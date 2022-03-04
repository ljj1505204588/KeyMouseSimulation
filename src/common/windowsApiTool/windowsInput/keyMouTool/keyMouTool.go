package keyMouTool

import (
	KeyboardMouseInput "KeyMouseSimulation/common/windowsApiTool/windowsInput"
	"sync"
)

/*
*  --------------------------------------- KEYBOARD INPUT ---------------------------------------
 */

var keyChan chan KeyInputChanT
var KeyboardInputPool = sync.Pool{
	New: func() interface{} {
		return new(KeyboardMouseInput.KeyboardInputT)
	},
}

//创建一个向计算机底层发送信息的通道
func GetKeySendInputChan(size int) (chan KeyInputChanT, error) {

	var createRoutine = false

	if keyChan == nil {
		keyChan = make(chan KeyInputChanT, size)
		createRoutine = true
	}
	if len(keyChan) != size {
		close(keyChan)
		keyChan = make(chan KeyInputChanT, size)
		createRoutine = true
	}

	if createRoutine {
		go keyRoutine(keyChan)
	}

	return keyChan, nil
}

func keyRoutine(c chan KeyInputChanT) {

	defer func() {
		if p := recover(); p != nil {
			go keyRoutine(c)
		}
	}()

	//TODO 考虑批量
	for even := range c {
		KeyboardInput := KeyboardInputPool.Get().(*KeyboardMouseInput.KeyboardInputT)

		KeyboardInput.Type = uint32(TYPE_INPUT_KEYBOARD)
		KeyboardInput.Ki.WVk = int16(even.VK)
		KeyboardInput.Ki.DwFlags = uint32(even.DwFlags)

		_, err := KeyboardMouseInput.KeyboardInput(*KeyboardInput)
		if err != nil {
			//No Importance
		}

		KeyboardInputPool.Put(KeyboardInput)
	}
}

/*
*  --------------------------------------- MOUSE INPUT ---------------------------------------
 */

var mouseChan chan MouseInputChanT
var MouseInputPool = sync.Pool{
	New: func() interface{} {
		return new(KeyboardMouseInput.MouseInputT)
	},
}

//创建一个向计算机底层发送信息的通道
func GetMouseSendInputChan(size int) (chan MouseInputChanT, error) {

	var createRoutine = false

	if mouseChan == nil {
		mouseChan = make(chan MouseInputChanT, size)
		createRoutine = true
	}
	if len(mouseChan) != size {
		close(mouseChan)
		mouseChan = make(chan MouseInputChanT, size)
		createRoutine = true
	}

	if createRoutine {
		go mouseRoutine(mouseChan)
	}

	return mouseChan, nil
}

func mouseRoutine(c chan MouseInputChanT) {

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

		_, err := KeyboardMouseInput.MouseInput(*MouseInput)
		if err != nil {
			//No Importance
		}

		MouseInputPool.Put(MouseInput)
	}
}