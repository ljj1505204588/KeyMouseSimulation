package component

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/commonTool"
	"KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
)

func init() {
	windowsService = windowsHookServerT{}
	windowsService.start()
}

var windowsService windowsHookServerT

type windowsHookServerT struct {
	mouseChan    chan *windowsHook.MouseEvent    //鼠标监听通道
	keyboardChan chan *windowsHook.KeyboardEvent //键盘监听通道

	keySend   chan *keyMouTool.KeyInputT   //键盘发送通道
	mouseSend chan *keyMouTool.MouseInputT //鼠标发送通道
}

func (s *windowsHookServerT) start() {
	//获取key、mouse发送通道2
	var err error
	s.keySend, err = keyMouTool.GetKeySendInputChan(3000)
	commonTool.MustNil(err)
	s.mouseSend, err = keyMouTool.GetMouseSendInputChan(3000)
	commonTool.MustNil(err)

	eventCenter.Event.Register(events.WindowsMouseInput, s.MouseInput)
	eventCenter.Event.Register(events.WindowsKeyBoardInput, s.KeyBoardInput)

	go s.MouseHook()
	go s.KeyBoardHook()
}

// -------------------------------------- 勾子 --------------------------------------

// todo 加个优雅退出

func (s *windowsHookServerT) MouseHook() {
	defer func() { go s.MouseHook() }()

	var err error
	if s.mouseChan == nil {
		s.mouseChan, err = windowsHook.MouseHook(nil)
		commonTool.MustNil(err)
	}

	for {
		var event = <-s.mouseChan
		_ = eventCenter.Event.Publish(events.WindowsMouseHook, events.WindowsMouseHookData{
			Date: event,
		})
	}
}
func (s *windowsHookServerT) KeyBoardHook() {
	defer func() { go s.KeyBoardHook() }()

	var err error
	if s.keyboardChan == nil {
		s.keyboardChan, err = windowsHook.KeyBoardHook(nil)
		commonTool.MustNil(err)
	}

	for {
		var event = <-s.keyboardChan
		_ = eventCenter.Event.Publish(events.WindowsKeyBoardHook, events.WindowsKeyBoardHookData{
			Date: event,
		})
	}
}

// -------------------------------------- 输入 --------------------------------------

func (s *windowsHookServerT) MouseInput(data interface{}) (err error) {
	var input = data.(events.WindowsMouseInputData)

	s.mouseSend <- input.Data

	return
}
func (s *windowsHookServerT) KeyBoardInput(data interface{}) (err error) {
	var input = data.(events.WindowsKeyBoardInputData)

	s.keySend <- input.Data

	return
}
