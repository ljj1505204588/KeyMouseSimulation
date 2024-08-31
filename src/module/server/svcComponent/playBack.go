package svcComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	component "KeyMouseSimulation/module/baseComponent"
	"KeyMouseSimulation/share/events"
	"sync/atomic"
	"time"
)

type PlayBackServerI interface {
	Start(fileName string)
	Pause()
	Stop()
}

/*
*	PlayBackServerI实现接口
 */

func GetPlaybackServer() PlayBackServerI {
	p := playBackServerT{
		fileControl: component.FileControl,
	}

	p.input = map[keyMouTool.InputType]func(note *component.NoteT){
		keyMouTool.TYPE_INPUT_MOUSE:    p.mouseInput,
		keyMouTool.TYPE_INPUT_KEYBOARD: p.keyBoardInput,
	}

	return &p
}

type playBackServerT struct {
	run         bool //回放标识
	fileControl component.FileControlI

	name       string            //读取文件名称
	notes      []component.NoteT //回放数据
	notesIndex int64             //当前回放在文件内容位置

	input map[keyMouTool.InputType]func(note *component.NoteT)
}

// Start 开始
func (p *playBackServerT) Start(fileName string) {
	p.tryLoadFile(fileName)
	p.run = true

	go p.playBack()
}

// Pause 暂停
func (p *playBackServerT) Pause() {
	p.run = false
}

// Stop 停止
func (p *playBackServerT) Stop() {
	p.run = false
	atomic.SwapInt64(&p.notesIndex, 0)
}

// ----------------------- playBack 模块主体功能函数 -----------------------

func (p *playBackServerT) playBack() {
	for p.run {
		var index = atomic.LoadInt64(&p.notesIndex)
		if p.checkPlayBackFinish(index) {
			return
		}

		var n = &p.notes[index]
		p.input[n.NoteType](n)
		time.Sleep(time.Duration(n.TimeGap))
		atomic.CompareAndSwapInt64(&p.notesIndex, index, index+1)
	}
}
func (p *playBackServerT) checkPlayBackFinish(index int64) (finish bool) {

	if index >= int64(len(p.notes)) {
		atomic.SwapInt64(&p.notesIndex, 0)
		_ = eventCenter.Event.Publish(events.PlayBackFinish, events.PlayBackFinishData{})
		return true
	}

	return false
}
func (p *playBackServerT) mouseInput(note *component.NoteT) {
	if err := eventCenter.Event.Publish(events.WindowsMouseInput, events.WindowsMouseInputData{
		Data: &keyMouTool.MouseInputT{
			X:         note.MouseNote.X,
			Y:         note.MouseNote.Y,
			DWFlags:   note.MouseNote.DWFlags,
			MouseData: note.MouseNote.MouseData,
			Time:      note.MouseNote.Time,
		},
	}); err != nil {
		p.tryPublishServerError(err)
	}
}
func (p *playBackServerT) keyBoardInput(note *component.NoteT) {
	if err := eventCenter.Event.Publish(events.WindowsKeyBoardInput, events.WindowsKeyBoardInputData{
		Data: &keyMouTool.KeyInputT{
			VK:      note.KeyNote.VK,
			DwFlags: note.KeyNote.DwFlags,
		},
	}); err != nil {
		p.tryPublishServerError(err)
	}
}

// 加载文件
func (p *playBackServerT) tryLoadFile(fileName string) {
	if p.name != fileName {
		p.run = false
		p.name = fileName
		p.notes = p.fileControl.ReadFile(fileName)
	}
}

// ----------------------- Util -----------------------

// 发布服务错误事件
func (p *playBackServerT) tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}
