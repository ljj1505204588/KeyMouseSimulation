package component

import (
	"KeyMouseSimulation/common/common"
	"KeyMouseSimulation/common/elegantExit"
	"KeyMouseSimulation/common/windowsApi/windowsHook"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	eventCenter "KeyMouseSimulation/pkg/event"
	"KeyMouseSimulation/share/topic"
	"fmt"
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
	common.MustNil(err)
	s.mouseSend, err = keyMouTool.GetMouseSendInputChan(3000)
	common.MustNil(err)

	eventCenter.Event.Register(topic.WindowsMouseInput, s.MouseInput)
	eventCenter.Event.Register(topic.WindowsKeyBoardInput, s.KeyBoardInput)

	go s.MouseHook()
	go s.KeyBoardHook()

	elegantExit.AddElegantExit(func() {
		_ = windowsHook.MouseUnhook()
		_ = windowsHook.KeyBoardUnhook()

		fmt.Println("windows回调已删除.")
	})
}

// -------------------------------------- 勾子 --------------------------------------

// todo 加个优雅退出

func (s *windowsHookServerT) MouseHook() {
	defer func() { go s.MouseHook() }()

	var err error
	if s.mouseChan == nil {
		s.mouseChan, err = windowsHook.MouseHook(nil)
		common.MustNil(err)
	}

	for {
		var event = <-s.mouseChan
		_ = eventCenter.Event.Publish(topic.WindowsMouseHook, event)
	}
}
func (s *windowsHookServerT) KeyBoardHook() {
	defer func() { go s.KeyBoardHook() }()

	var err error
	if s.keyboardChan == nil {
		s.keyboardChan, err = windowsHook.KeyBoardHook(nil)
		common.MustNil(err)
	}

	for {
		var event = <-s.keyboardChan
		_ = eventCenter.Event.Publish(topic.WindowsKeyBoardHook, event)
	}
}

// -------------------------------------- 输入 --------------------------------------

func (s *windowsHookServerT) MouseInput(data interface{}) (err error) {

	s.mouseSend <- data.(*keyMouTool.MouseInputT)

	return
}
func (s *windowsHookServerT) KeyBoardInput(data interface{}) (err error) {
	s.keySend <- data.(*keyMouTool.KeyInputT)

	return
}
