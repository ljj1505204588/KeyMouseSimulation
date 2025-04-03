package svcComponent

import (
	"KeyMouseSimulation/common/windowsApi/windowsHook"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	rp_file "KeyMouseSimulation/pkg/file"
	"KeyMouseSimulation/share/topic"
	"sync"
	"time"
)

func GetRecordServer() *RecordServerT {
	R := RecordServerT{}
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
	l sync.Mutex

	notes     keyMouTool.MulNote // 记录
	saveNotes keyMouTool.MulNote // 待存储记录
	noteTime  int64

	mouseHs    []func(data interface{})
	keyBoardHs []func(data interface{})

	lastMoveEven *windowsHook.MouseEvent //最后移动事件，配合是否记录鼠标移动路径使用
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

	r.mouseHs = []func(data interface{}){}
	r.keyBoardHs = []func(data interface{}){}

	r.saveNotes = r.notes
	r.notes = keyMouTool.MulNote{}
}

// Save 存储
func (r *RecordServerT) Save(name string) {
	rp_file.FileControl.Save(name, r.saveNotes)
}

func (r *RecordServerT) registerHandler() {
	eventCenter.Event.Register(topic.WindowsMouseHook, func(data interface{}) (err error) {
		for _, per := range r.mouseHs {
			per(data)
		}
		return
	})
	eventCenter.Event.Register(topic.WindowsKeyBoardHook, func(data interface{}) (err error) {
		for _, per := range r.keyBoardHs {
			per(data)
		}
		return
	})

}

// -------------------------------------------- record --------------------------------------------

// mouseHandler 鼠标记录
func (r *RecordServerT) mouseHandler(data interface{}) {
	defer r.lockSelf()()

	var info = data.(*topic.WindowsMouseHookData)

	if !conf.RecordMouseTrackConf.GetValue() {
		if info.Date.Message == windowsHook.WM_MOUSEMOVE {
			r.lastMoveEven = info.Date
			return
		} else if r.lastMoveEven != nil {
			// 先把鼠标移动过去
			r.notes.AppendMouseNote(info.Date.RecordTime, r.lastMoveEven)
		}
	}

	r.notes.AppendMouseNote(r.noteTime, info.Date)
	r.noteTime = info.Date.RecordTime

	// 设置长度
	conf.RecordLen.SetValue(len(r.notes))
}

// keyBoardHandler 键盘记录
func (r *RecordServerT) keyBoardHandler(data interface{}) {
	defer r.lockSelf()()

	var info = data.(*topic.WindowsKeyBoardHookData)

	r.notes.AppendKeyBoardNote(r.noteTime, info.Date)
	r.noteTime = info.Date.RecordTime

	// 设置长度
	conf.RecordLen.SetValue(len(r.notes))
}

func (r *RecordServerT) lockSelf() func() {
	r.l.Lock()
	return r.l.Unlock
}
