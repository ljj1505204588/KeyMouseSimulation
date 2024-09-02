package svcComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/share/events"
	"sync"
	"time"
)

func GetRecordServer() *RecordServerT {
	R := RecordServerT{
		fileControl: component.FileControl,
	}
	R.registerHandler()
	return &R
}

type RecordServerI interface {
	Start()           // 开始
	Pause()           // 暂停
	Stop()            // 停止
	Save(name string) // 存储
}

/*
*	---------------------------------------------------- RecordServerI ----------------------------------------------------
 */

type RecordServerT struct {
	l           sync.Mutex
	fileControl component.FileControlI

	notes     component.MulNote // 记录
	saveNotes component.MulNote // 待存储记录
	noteTime  int64

	mouseHs    []func(data interface{})
	keyBoardHs []func(data interface{})

	recordMouseTrack bool                    //是否记录鼠标移动路径使用
	lastMoveEven     *windowsHook.MouseEvent //最后移动事件，配合是否记录鼠标移动路径使用
}

// Start 开始
func (r *RecordServerT) Start() {
	defer r.lockSelf()()

	r.mouseHs = append([]func(data interface{}){}, r.mouseHandler)
	r.keyBoardHs = append([]func(data interface{}){}, r.keyBoardHandler)

	r.noteTime = time.Now().UnixNano()
}

// Pause 暂停
func (r *RecordServerT) Pause() {
	defer r.lockSelf()()

	r.mouseHs = []func(data interface{}){}
	r.keyBoardHs = []func(data interface{}){}
}

// Stop 停止
func (r *RecordServerT) Stop() {
	defer r.lockSelf()()

	r.saveNotes = r.notes
	r.notes = component.MulNote{}
	_ = eventCenter.Event.Publish(events.RecordFinish, events.RecordFinishData{})
}

// Save 存储
func (r *RecordServerT) Save(name string) {
	r.fileControl.Save(name, r.saveNotes)
}

func (r *RecordServerT) registerHandler() {
	eventCenter.Event.Register(events.WindowsMouseHook, func(data interface{}) (err error) {
		for _, per := range r.mouseHs {
			per(data)
		}
		return
	})
	eventCenter.Event.Register(events.WindowsKeyBoardHook, func(data interface{}) (err error) {
		for _, per := range r.keyBoardHs {
			per(data)
		}
		return
	})

	component.RecordConfig.SetMouseTrackChange(true, func(record bool) {
		r.recordMouseTrack = record
	})
}

// mouseHandler 鼠标记录
func (r *RecordServerT) mouseHandler(data interface{}) {
	defer r.lockSelf()()

	var info = data.(events.WindowsMouseHookData)
	r.notes.AppendMouseNote(r.noteTime, info.Date)

	r.noteTime = info.Date.RecordTime
}

// keyBoardHandler 键盘记录
func (r *RecordServerT) keyBoardHandler(data interface{}) {
	defer r.lockSelf()()

	var info = data.(events.WindowsKeyBoardHookData)

	r.notes.AppendKeyBoardNote(r.noteTime, info.Date)
	r.noteTime = info.Date.RecordTime
}

func (r *RecordServerT) lockSelf() func() {
	r.l.Lock()
	return r.l.Unlock
}
