package recordAndPlayBack

import (
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

	R.mouseDwMap = map[windowsHook.Message]keyMouTool.MouseInputDW{
		windowsHook.WM_MOUSEMOVE:         keyMouTool.DW_MOUSEEVENTF_MOVE,
		windowsHook.WM_MOUSELEFTDOWN:     keyMouTool.DW_MOUSEEVENTF_LEFTDOWN,
		windowsHook.WM_MOUSELEFTUP:       keyMouTool.DW_MOUSEEVENTF_LEFTUP,
		windowsHook.WM_MOUSERIGHTDOWN:    keyMouTool.DW_MOUSEEVENTF_RIGHTDOWN,
		windowsHook.WM_MOUSERIGHTUP:      keyMouTool.DW_MOUSEEVENTF_RIGHTUP,
		windowsHook.WM_MOUSEMIDDLEDOWN:   keyMouTool.DW_MOUSEEVENTF_MIDDLEDOWN,
		windowsHook.WM_MOUSEMIDDLEUP:     keyMouTool.DW_MOUSEEVENTF_MIDDLEUP,
		windowsHook.WM_MOUSELEFTSICEDOWN: keyMouTool.DW_MOUSEEVENTF_XDOWN,
		windowsHook.WM_MOUSELEFTSICEUP:   keyMouTool.DW_MOUSEEVENTF_XUP,
		windowsHook.WM_MOUSEWHEEL:        keyMouTool.DW_MOUSEEVENTF_WHEEL,
		windowsHook.WM_MOUSEHWHEEL:       keyMouTool.DW_MOUSEEVENTF_HWHEEL,
	}
	R.keyDwMap = map[windowsHook.Message]keyMouTool.KeyBoardInputDW{
		windowsHook.WM_KEYUP: keyMouTool.DW_KEYEVENTF_KEYUP,
	}

	go R.loop()
	return &R
}

type RecordServerT struct {
	status     ServerStatus //状态
	noteName   string       //记录名称
	notes      []noteT      //记录
	mouseDwMap map[windowsHook.Message]keyMouTool.MouseInputDW
	keyDwMap   map[windowsHook.Message]keyMouTool.KeyBoardInputDW

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
	case SERVER_TYPE_RECORD_PAUSE:
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
	go func() {
		logTool.DebugAJ(" record 发送变动消息：", event.String())
		// 发送变动消息
		if R.messageChan == nil {
			R.messageChan = make(chan RecordMessageT, 100)
		}

		R.messageChan <- RecordMessageT{
			Event: event,
			Value: value,
		}
	}()

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

	exit := make(chan struct{}, 0)
	go R.free(exit, kHook, mHook)
	nowStatus := SERVER_TYPE_FREE

	//记录主循环
	for {
		if R.status != nowStatus {
			logTool.DebugAJ("record 状态变动:" + nowStatus.String() + "->" + R.status.String())

			nowStatus = R.status
			exit <- struct{}{}
			switch R.status {
			case SERVER_TYPE_FREE:
				go R.free(exit, kHook, mHook)
			case SERVER_TYPE_RECORD:
				go R.record(exit, kHook, mHook)
			case SERVER_TYPE_RECORD_PAUSE:
				go R.stop(exit, kHook, mHook)
			}
		}
		time.Sleep(10 * time.Millisecond) //TODO 想通过通道去阻塞
	}

	return
}
func (R *RecordServerT) free(exit chan struct{}, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	if len(R.notes) != 0 {
		notes := make([]noteT, len(R.notes))
		_ = copy(notes, R.notes)
		R.notes = make([]noteT, 0)
		go R.recordNote(R.noteName, notes)
	}

	for {
		select {
		case <-exit:
			logTool.DebugAJ("record 退出记录空闲状态")
			return
		case kEvent := <-kHook:
			R.hotKeyDeal(kEvent)
		case _ = <-mHook:
			time.Sleep(10 * time.Nanosecond)
		}
	}
}
func (R *RecordServerT) record(exit chan struct{}, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	st := time.Now().UnixNano()
	timeGap := int64(0)
	for {
		select {
		case <-exit:
			logTool.DebugAJ("record 退出记录状态")
			return
		case kEvent := <-kHook:
			if R.hotKeyDeal(kEvent) {
				continue
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			R.notes = append(R.notes, noteT{
				NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
				KeyNote:  keyMouTool.KeyInputChanT{VK: keyMouTool.VKCode(kEvent.VkCode), DwFlags: R.transKeyDwFlags(kEvent.Message)},
				TimeGap:  timeGap,
			})
		case mEvent := <-mHook:
			if !R.recordMouseTrack {
				if mEvent.Message == windowsHook.WM_MOUSEMOVE {
					R.lastMoveEven = mEvent
					continue
				} else if R.lastMoveEven.Message == windowsHook.WM_MOUSEMOVE {
					R.notes = append(R.notes, noteT{
						NoteType:  keyMouTool.TYPE_INPUT_MOUSE,
						MouseNote: keyMouTool.MouseInputChanT{X: R.lastMoveEven.X, Y: R.lastMoveEven.Y, DWFlags: R.transMouseDwFlags(R.lastMoveEven.Message)},
						TimeGap:   0,
					})
				}
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			R.notes = append(R.notes, noteT{
				NoteType:  keyMouTool.TYPE_INPUT_MOUSE,
				MouseNote: keyMouTool.MouseInputChanT{X: mEvent.X, Y: mEvent.Y, DWFlags: R.transMouseDwFlags(mEvent.Message)},
				TimeGap:   timeGap,
			})
		}
	}

}
func (R *RecordServerT) stop(exit chan struct{}, kHook chan windowsHook.KeyboardEvent, mHook chan windowsHook.MouseEvent) {
	for {
		select {
		case <-exit:
			logTool.DebugAJ("record 退出记录停止状态")
			return
		case kEvent := <-kHook:
			R.hotKeyDeal(kEvent)
		case _ = <-mHook:
			time.Sleep(10 * time.Nanosecond)
		}
	}
}
func (R *RecordServerT) hotKeyDeal(event windowsHook.KeyboardEvent) (isHotKey bool) {
	if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(event.VkCode)]; exist {
		R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
		isHotKey = true
	}
	return
}
func (R *RecordServerT) recordNote(name string, notes []noteT) {
	logTool.DebugAJ("record 开始记录文件：" + "名称:" + name + " 长度：" + strconv.Itoa(len(notes)))

	if name == "" {
		R.sendMessage(RECORD_SAVE_FILE_ERROR, language.ErrorSaveFileNameNilStr)
		return
	}

	if len(notes) == 0 {
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
func (R *RecordServerT) transMouseDwFlags(message windowsHook.Message) (dw keyMouTool.MouseInputDW) {
	return R.mouseDwMap[message] //| keyMouTool.DW_MOUSEEVENTF_ABSOLUTE
}
func (R *RecordServerT) transKeyDwFlags(message windowsHook.Message) keyMouTool.KeyBoardInputDW {
	return R.keyDwMap[message]
}
