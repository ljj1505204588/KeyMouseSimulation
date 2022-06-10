package recordAndPlayBack

import (
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RecordServerI interface {
	Start() error           //开始
	Pause() error           //暂停
	Stop(name string) error //停止

	GetRecordMessageChan() chan RecordMessageT       //获取事件反馈通道
	SetHotKey(k HotKey, vks keyMouTool.VKCode) error //设置热键
	SetIfTrackMouseMove(sign bool)                   //设置是否记录鼠标移动路径
}

/*
*	RecordServerI 实现接口
 */

func GetRecordServer() *RecordServerT {
	R := RecordServerT{
		hotKeyM:          make(map[keyMouTool.VKCode]HotKey),
		recordMouseTrack: true,
		status:           SERVER_TYPE_FREE,
		messageChan:      make(chan RecordMessageT, 100),
	}

	go R.loop()
	return &R
}

type RecordServerT struct {
	status   ServerStatus //状态
	noteName string       //记录名称

	recordMouseTrack bool                   //是否记录鼠标移动路径使用
	lastMoveEven     windowsHook.MouseEvent //最后移动事件，配合是否记录鼠标移动路径使用

	hotKeyM map[keyMouTool.VKCode]HotKey //热键信息

	messageChan chan RecordMessageT //消息通道
}

//Start 开始
func (R *RecordServerT) Start() error {
	return R.changeStatus(SERVER_TYPE_RECORD)
}

//Pause 暂停
func (R *RecordServerT) Pause() error {
	return R.changeStatus(SERVER_TYPE_RECORD_PAUSE)
}

//Stop 停止
func (R *RecordServerT) Stop(name string) error {
	R.noteName = name

	return R.changeStatus(SERVER_TYPE_FREE)
}
func (R *RecordServerT) changeStatus(status ServerStatus) error {
	//判断状态更改是否合理
	switch R.status {
	case SERVER_TYPE_FREE:
		if status == SERVER_TYPE_RECORD_PAUSE {

			return fmt.Errorf(language.ErrorFreeToRecordPause)
		}
	case SERVER_TYPE_RECORD:
		return nil
	case SERVER_TYPE_RECORD_PAUSE:
		return nil
	}

	//修改状态
	R.status = status
	R.sendMessage(RECORD_EVENT_STATUS_CHANGE, status)

	return nil
}

//GetRecordMessageChan 获取消息监听通道
func (R *RecordServerT) GetRecordMessageChan() chan RecordMessageT {
	if R.messageChan == nil {
		R.messageChan = make(chan RecordMessageT, 100)
	}

	return R.messageChan
}
func (R *RecordServerT) sendMessage(event RecordEvent, value interface{}) {
	// 发送变动消息
	if R.messageChan == nil {
		R.messageChan = make(chan RecordMessageT, 100)
	}

	R.messageChan <- RecordMessageT{
		Event: event,
		Value: value,
	}
}

//SetHotKey 设置热键
func (R *RecordServerT) SetHotKey(k HotKey, vks keyMouTool.VKCode) error {
	M := make(map[keyMouTool.VKCode]HotKey)

	if R.hotKeyM != nil {
		for k, v := range R.hotKeyM {
			M[k] = v
		}
	}

	M[vks] = k
	R.hotKeyM = M

	return nil
}

//SetIfTrackMouseMove 设置是否追踪鼠标记录
func (R *RecordServerT) SetIfTrackMouseMove(sign bool) {
	R.recordMouseTrack = sign
}

// ----------------------- record 模块主体循环 -----------------------

func (R *RecordServerT) loop() {
	defer func() {
		if info := recover(); info != nil {
			go R.loop()
		} else {
			panic("record 错误退出")
		}
	}()

	if R.hotKeyM == nil {
		R.hotKeyM = make(map[keyMouTool.VKCode]HotKey)
	}

	//挂上钩子
	kHook, err1 := windowsHook.KeyBoardHook(nil)
	mHook, err2 := windowsHook.MouseHook(nil)
	if err1 != nil || err2 != nil {
		os.Exit(1)
	}
	defer func() { _ = windowsHook.KeyBoardUnhook() }()
	defer func() { _ = windowsHook.MouseUnhook() }()

	ctx, down := context.WithCancel(context.Background())
	go R.free(ctx, kHook, mHook)
	nowStatus := SERVER_TYPE_FREE

	//记录主循环
	for {
		if R.status != nowStatus {
			nowStatus = R.status
			down()
			ctx, down = context.WithCancel(context.Background())
			switch R.status {
			case SERVER_TYPE_FREE:
				go R.free(ctx, kHook, mHook)
			case SERVER_TYPE_RECORD:
				go R.record(ctx, kHook, mHook)
			case SERVER_TYPE_RECORD_PAUSE:
				go R.stop(ctx, kHook, mHook)
			}
		}
		time.Sleep(10 * time.Millisecond) //TODO 想通过通道去阻塞
	}

	return
}
func (R *RecordServerT) free(ctx context.Context, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	for {
		select {
		case kEvent := <-kHook:
			if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(kEvent.VkCode)]; exist {
				R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
			}
		case _ = <-mHook:
			time.Sleep(10 * time.Nanosecond)
		case <-ctx.Done():
			return
		}
	}

}
func (R *RecordServerT) record(ctx context.Context, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	var notes []noteT
	st := time.Now().UnixNano()
	timeGap := int64(0)
	defer func() {
		if len(notes) != 0 && R.noteName != "" {
			go R.recordNote(R.noteName, notes)
		}
	}()
	for {
		select {
		case kEvent := <-kHook:
			if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(kEvent.VkCode)]; exist {
				R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
				continue
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			notes = append(notes, noteT{
				NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
				KeyNote:  keyMouTool.KeyInputChanT{VK: keyMouTool.VKCode(kEvent.VkCode), DwFlags: transKeyDwFlags(kEvent.Message)},
				TimeGap:  timeGap,
			})
		case mEvent := <-mHook:
			if !R.recordMouseTrack {
				if mEvent.Message == windowsHook.WM_MOUSEMOVE {
					R.lastMoveEven = mEvent
					continue
				} else if R.lastMoveEven.Message == windowsHook.WM_MOUSEMOVE {
					notes = append(notes, noteT{
						NoteType:  keyMouTool.TYPE_INPUT_MOUSE,
						MouseNote: keyMouTool.MouseInputChanT{X: R.lastMoveEven.X, Y: R.lastMoveEven.Y, DWFlags: transMouseDwFlags(R.lastMoveEven.Message)},
						TimeGap:   0,
					})
				}
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			notes = append(notes, noteT{
				NoteType:  keyMouTool.TYPE_INPUT_MOUSE,
				MouseNote: keyMouTool.MouseInputChanT{X: mEvent.X, Y: mEvent.Y, DWFlags: transMouseDwFlags(mEvent.Message)},
				TimeGap:   timeGap})
		case <-ctx.Done():
			return
		}
	}

}
func (R *RecordServerT) stop(ctx context.Context, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	for {
		select {
		case kEvent := <-kHook:
			if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(kEvent.VkCode)]; exist {
				R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
				continue
			}
		case _ = <-mHook:
			time.Sleep(10 * time.Nanosecond)
		case <-ctx.Done():
			return
		}
	}
}
func (R *RecordServerT) recordNote(name string, notes []noteT) {
	if name == "" {
		R.sendMessage(RECORD_SAVE_FILE_ERROR, language.ErrorSaveFileNameNilStr)
		return
	}

	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
	if err != nil {
		R.sendMessage(RECORD_SAVE_FILE_ERROR, err.Error())
		return
	}
	defer func() { _ = file.Close() }()

	js, err := json.Marshal(notes)
	if err != nil {
		R.sendMessage(RECORD_SAVE_FILE_ERROR, err.Error())
		return
	}

	_, err = file.Write(js)
	if err != nil {
		R.sendMessage(RECORD_SAVE_FILE_ERROR, err.Error())
		return
	}
}

// ----------------------- 被调用模块 -----------------------
func transMouseDwFlags(message windowsHook.Message) (dw keyMouTool.MouseInputDW) {
	//TODO 这个变更考虑做到加载时候
	switch message {
	case windowsHook.WM_MOUSEMOVE:
		dw = keyMouTool.DW_MOUSEEVENTF_MOVE
	case windowsHook.WM_MOUSELEFTDOWN:
		dw = keyMouTool.DW_MOUSEEVENTF_LEFTDOWN
	case windowsHook.WM_MOUSELEFTUP:
		dw = keyMouTool.DW_MOUSEEVENTF_LEFTUP
	case windowsHook.WM_MOUSERIGHTDOWN:
		dw = keyMouTool.DW_MOUSEEVENTF_RIGHTDOWN
	case windowsHook.WM_MOUSERIGHTUP:
		dw = keyMouTool.DW_MOUSEEVENTF_RIGHTUP
	case windowsHook.WM_MOUSEMIDDLEDOWN:
		dw = keyMouTool.DW_MOUSEEVENTF_MIDDLEDOWN
	case windowsHook.WM_MOUSEMIDDLEUP:
		dw = keyMouTool.DW_MOUSEEVENTF_MIDDLEUP
	}
	return dw | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE
}
func transKeyDwFlags(message windowsHook.Message) keyMouTool.KeyBoardInputDW {
	if message == windowsHook.WM_KEYUP {
		return keyMouTool.DW_KEYEVENTF_KEYUP
	}
	return 0
}
