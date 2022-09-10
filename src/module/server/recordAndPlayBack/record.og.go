package recordAndPlayBack

import (
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/UI"
	"KeyMouseSimulation/share/enum"
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

	SetHotKey(k enum.HotKey, vks keyMouTool.VKCode) error //设置热键
	SetIfTrackMouseMove(sign bool)                        //设置是否记录鼠标移动路径
}

/*
*	RecordServerI 实现接口
 */

func GetRecordServer() *RecordServerT {
	R := RecordServerT{
		hotKeyM:          make(map[keyMouTool.VKCode]enum.HotKey),
		recordMouseTrack: true,
	}

	R.mouseDwMap = map[windowsHook.Message]keyMouTool.MouseInputDW{
		windowsHook.WM_MOUSEMOVE:         keyMouTool.DW_MOUSEEVENTF_MOVE | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTDOWN:     keyMouTool.DW_MOUSEEVENTF_LEFTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTUP:       keyMouTool.DW_MOUSEEVENTF_LEFTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSERIGHTDOWN:    keyMouTool.DW_MOUSEEVENTF_RIGHTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSERIGHTUP:      keyMouTool.DW_MOUSEEVENTF_RIGHTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEMIDDLEDOWN:   keyMouTool.DW_MOUSEEVENTF_MIDDLEDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEMIDDLEUP:     keyMouTool.DW_MOUSEEVENTF_MIDDLEUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTSICEDOWN: keyMouTool.DW_MOUSEEVENTF_XDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSELEFTSICEUP:   keyMouTool.DW_MOUSEEVENTF_XUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEWHEEL:        keyMouTool.DW_MOUSEEVENTF_WHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
		windowsHook.WM_MOUSEHWHEEL:       keyMouTool.DW_MOUSEEVENTF_HWHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE, //这个暂时不知道是啥，
	}
	R.keyDwMap = map[windowsHook.Message]keyMouTool.KeyBoardInputDW{
		windowsHook.WM_KEYUP: keyMouTool.DW_KEYEVENTF_KEYUP,
	}

	go R.loop()
	return &R
}

type RecordServerT struct {
	noteName   string  //记录名称
	notes      []noteT //记录
	mouseDwMap map[windowsHook.Message]keyMouTool.MouseInputDW
	keyDwMap   map[windowsHook.Message]keyMouTool.KeyBoardInputDW

	recordMouseTrack bool                   //是否记录鼠标移动路径使用
	lastMoveEven     windowsHook.MouseEvent //最后移动事件，配合是否记录鼠标移动路径使用

	hotKeyM map[keyMouTool.VKCode]enum.HotKey //热键信息

}

//Start 开始
func (R *RecordServerT) Start() error {
	return nil
}

//Pause 暂停
func (R *RecordServerT) Pause() error {
	return nil
}

//Stop 停止
func (R *RecordServerT) Stop(name string) error {
	R.noteName = name

	return nil
}

//SetHotKey 设置热键
func (R *RecordServerT) SetHotKey(k enum.HotKey, vks keyMouTool.VKCode) error {
	M := make(map[keyMouTool.VKCode]enum.HotKey)

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
		R.hotKeyM = make(map[keyMouTool.VKCode]enum.HotKey)
	}

	//挂上钩子 -- 鼠标钩子只在记录中使用
	kHook, err := windowsHook.KeyBoardHook(nil)
	if err != nil {
		os.Exit(1)
	}

	defer func() { _ = windowsHook.KeyBoardUnhook() }()

	exit := make(chan struct{}, 0)
	go R.free(exit, kHook)
	nowStatus := SERVER_TYPE_FREE

	//记录主循环
	for {
		if R.status != nowStatus {
			logTool.DebugAJ("record 状态变动:" + nowStatus.String() + "->" + R.status.String())

			nowStatus = R.status
			exit <- struct{}{}
			switch R.status {
			case SERVER_TYPE_FREE:
				go R.free(exit, kHook)
			case SERVER_TYPE_RECORD:
				go R.record(exit, kHook)
			case SERVER_TYPE_RECORD_PAUSE:
				go R.stop(exit, kHook)
			}
		}
		time.Sleep(10 * time.Millisecond) //TODO 想通过通道去阻塞
	}

	return
}
func (R *RecordServerT) free(exit chan struct{}, kHook chan windowsHook.KeyboardEvent) {
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
		}
	}
}
func (R *RecordServerT) record(exit chan struct{}, kHook chan windowsHook.KeyboardEvent) {
	mHook, err := windowsHook.MouseHook(nil)
	if err != nil {
		return
	}
	defer func() { _ = windowsHook.MouseUnhook() }()

	st := time.Now().UnixNano()
	timeGap := int64(0)
	for {
		select {
		case <-exit:
			logTool.DebugAJ("record 退出记录状态")
			return
		case kEvent := <-kHook:
			//fmt.Println(fmt.Sprintf("%v",kEvent))
			if R.hotKeyDeal(kEvent) {
				continue
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			R.notes = append(R.notes, noteT{
				NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
				KeyNote:  &keyMouTool.KeyInputChanT{VK: keyMouTool.VKCode(kEvent.VkCode), DwFlags: R.transKeyDwFlags(kEvent.Message)},
				TimeGap:  timeGap,
			})
		case mEvent := <-mHook:
			//fmt.Println(fmt.Sprintf("%v",mEvent))
			if !R.recordMouseTrack {
				if mEvent.Message == windowsHook.WM_MOUSEMOVE {
					R.lastMoveEven = mEvent
					continue
				} else if R.lastMoveEven.Message == windowsHook.WM_MOUSEMOVE {
					R.notes = append(R.notes, noteT{
						NoteType: keyMouTool.TYPE_INPUT_MOUSE,
						MouseNote: &keyMouTool.MouseInputChanT{X: R.lastMoveEven.X,
							Y:       R.lastMoveEven.Y,
							DWFlags: R.transMouseDwFlags(R.lastMoveEven.Message),
							Time:    mEvent.Time,
						},
						TimeGap: 0,
					})
				}
			}
			timeGap = time.Now().UnixNano()
			st, timeGap = timeGap, timeGap-st
			R.notes = append(R.notes, noteT{
				NoteType: keyMouTool.TYPE_INPUT_MOUSE,
				MouseNote: &keyMouTool.MouseInputChanT{X: mEvent.X,
					Y:         mEvent.Y,
					DWFlags:   R.transMouseDwFlags(mEvent.Message),
					MouseData: mEvent.MouseData,
				},
				TimeGap: timeGap,
			})
		}
	}

}
func (R *RecordServerT) stop(exit chan struct{}, kHook chan windowsHook.KeyboardEvent) {
	for {
		select {
		case <-exit:
			logTool.DebugAJ("record 退出记录停止状态")
			return
		case kEvent := <-kHook:
			R.hotKeyDeal(kEvent)
		}
	}
}
func (R *RecordServerT) hotKeyDeal(event windowsHook.KeyboardEvent) (isHotKey bool) {
	if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(event.VkCode)]; exist {
		isHotKey = true
	}
	return
}
func (R *RecordServerT) recordNote(name string, notes []noteT) {
	logTool.DebugAJ("record 开始记录文件：" + "名称:" + name + " 长度：" + strconv.Itoa(len(notes)))

	if name == "" {
		return
	}

	if len(notes) == 0 {
		return
	}

	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0772)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	js, err := json.Marshal(notes)
	if err != nil {
		return
	}
	_, err = file.Write(js)
	if err != nil {
		return
	}
}

// ----------------------- 被调用模块 -----------------------
func (R *RecordServerT) transMouseDwFlags(message windowsHook.Message) (dw keyMouTool.MouseInputDW) {

	return R.mouseDwMap[message]
}
func (R *RecordServerT) transKeyDwFlags(message windowsHook.Message) keyMouTool.KeyBoardInputDW {

	return R.keyDwMap[message]
}
