package recordAndPlayBack

import (
	"encoding/json"
	"fmt"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"os"
	"time"
)

type RecordServerI interface {
	Start() error
	Pause() error
	Stop(name string) error

	GetRecordMessageChan() chan RecordMessageT
	SetHotKey(k HotKey, vks keyMouTool.VKCode) error
	SetIfTrackMouseMove(sign bool)
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

	go R.record()
	return &R
}

type RecordServerT struct {
	status   ServerStatus
	noteName string

	lastMoveEven     windowsHook.MouseEvent
	recordMouseTrack bool

	hotKeyM map[keyMouTool.VKCode]HotKey

	messageChan chan RecordMessageT
}

func (R *RecordServerT) changeStatus(s ServerStatus) {
	R.status = s
	R.sendMessage(RECORD_EVENT_STATUS_CHANGE, s)
}
func (R *RecordServerT) judgeStatus(s ServerStatus) error {
	switch R.status {
	case SERVER_TYPE_FREE:
		if s == SERVER_TYPE_RECORD_PAUSE {

			return fmt.Errorf(language.ErrorFreeToRecordPause)
		}
	case SERVER_TYPE_RECORD:
		return nil
	case SERVER_TYPE_RECORD_PAUSE:
		return nil
	}
	return nil
}
func (R *RecordServerT) Start() error {

	if err := R.judgeStatus(SERVER_TYPE_RECORD); err != nil {
		return err
	}

	R.changeStatus(SERVER_TYPE_RECORD)
	return nil
}
func (R *RecordServerT) Pause() error {

	if err := R.judgeStatus(SERVER_TYPE_RECORD_PAUSE); err != nil {
		return err
	}

	R.changeStatus(SERVER_TYPE_RECORD_PAUSE)
	return nil
}
func (R *RecordServerT) Stop(name string) error {

	R.noteName = name
	if err := R.judgeStatus(SERVER_TYPE_FREE); err != nil {
		return err
	}

	R.changeStatus(SERVER_TYPE_FREE)
	return nil
}

func (R *RecordServerT) sendMessage(event RecordEvent, value interface{}) {
	if R.messageChan == nil {
		R.messageChan = make(chan RecordMessageT, 100)
	}

	R.messageChan <- RecordMessageT{
		Event: event,
		Value: value,
	}
}
func (R *RecordServerT) GetRecordMessageChan() chan RecordMessageT {
	if R.messageChan == nil {
		R.messageChan = make(chan RecordMessageT, 100)
	}

	return R.messageChan
}
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
func (R *RecordServerT) SetIfTrackMouseMove(sign bool) {
	R.recordMouseTrack = sign
}

// ----------------------- record 模块主体循环 -----------------------

func (R *RecordServerT) record() {
	defer func() {
		if info := recover(); info != nil {
			go R.record()
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
	defer windowsHook.KeyBoardUnhook()
	defer windowsHook.MouseUnhook()

	var notes []noteT
	timeGap := int64(0)
	st := time.Now().UnixNano()
	//记录主循环
	for {
		switch R.status {
		case SERVER_TYPE_FREE:
			if len(notes) != 0 {
				go R.recordNote(R.noteName, notes)
				notes = []noteT{}
			}
			select {
			case kEvent := <-kHook:
				if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(kEvent.VkCode)]; exist {
					R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
					continue
				}
			case _ = <-mHook:
				time.Sleep(10 * time.Nanosecond)
			default:
				time.Sleep(10 * time.Nanosecond)
			}
		case SERVER_TYPE_RECORD:
			if len(notes) == 0 {
				st = time.Now().UnixNano()
			}
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
					KeyNote: keyBoardNoteT{Vk:      keyMouTool.VKCode(kEvent.VkCode), DWFlags: transKeyDwFlags(kEvent.Message),},
					TimeGap: timeGap,
				})
			case mEvent := <-mHook:
				if !R.recordMouseTrack {
					if mEvent.Message == windowsHook.WM_MOUSEMOVE {
						R.lastMoveEven = mEvent
						continue
					} else if R.lastMoveEven.Message == windowsHook.WM_MOUSEMOVE{
						notes = append(notes, noteT{
							NoteType: keyMouTool.TYPE_INPUT_MOUSE,
							MouseNote: mouseNoteT{X: R.lastMoveEven.X, Y: R.lastMoveEven.Y, DWFlags: transMouseDwFlags(R.lastMoveEven.Message)},
							TimeGap: 0,
						})
					}
				}
				timeGap = time.Now().UnixNano()
				st, timeGap = timeGap, timeGap-st
				notes = append(notes, noteT{
					NoteType: keyMouTool.TYPE_INPUT_MOUSE,
					MouseNote: mouseNoteT{X: mEvent.X, Y: mEvent.Y, DWFlags: transMouseDwFlags(mEvent.Message)},
					TimeGap: timeGap})
			default:
				time.Sleep(10 * time.Nanosecond)
			}
		case SERVER_TYPE_RECORD_PAUSE:
			st = time.Now().UnixNano()
			select {
			case kEvent := <-kHook:
				if hotKey, exist := R.hotKeyM[keyMouTool.VKCode(kEvent.VkCode)]; exist {
					R.sendMessage(RECORD_EVENT_HOTKEY_DOWN, hotKey)
					continue
				}
			case _ = <-mHook:
				time.Sleep(10 * time.Nanosecond)
			default:
				time.Sleep(10 * time.Nanosecond)
			}
		default:
			time.Sleep(10 * time.Nanosecond)
		}
	}

	return
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
	defer file.Close()

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
