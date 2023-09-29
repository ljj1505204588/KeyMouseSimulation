package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"sync/atomic"
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
	p := PlayBackServerT{
		fileControl: GetFileControl(),
	}

	p.input = map[keyMouTool.InputType]func(note *noteT){
		keyMouTool.TYPE_INPUT_MOUSE:    p.mouseInput,
		keyMouTool.TYPE_INPUT_KEYBOARD: p.keyBoardInput,
	}

	return &p
}

type PlayBackServerT struct {
	run         bool //回放标识
	fileControl FileControlI

	name       string  //读取文件名称
	notes      []noteT //回放数据
	notesIndex int64   //当前回放在文件内容位置

	input map[keyMouTool.InputType]func(note *noteT)
}

// Start 开始
func (p *PlayBackServerT) Start(fileName string) {
	p.tryLoadFile(fileName)

	go p.playBack()
	p.run = true
}

// Pause 暂停
func (p *PlayBackServerT) Pause() {
	p.run = false
}

// Stop 停止
func (p *PlayBackServerT) Stop() {
	p.run = false
	atomic.SwapInt64(&p.notesIndex, 0)
}

// ----------------------- playBack 模块主体功能函数 -----------------------

func (p *PlayBackServerT) playBack() {

	for p.run {
		var index = atomic.LoadInt64(&p.notesIndex)
		if p.checkPlayBackFinish(index) {
			return
		}

		var n = &p.notes[index]
		p.input[n.NoteType](n)
		atomic.CompareAndSwapInt64(&p.notesIndex, index, index+1)
	}
}
func (p *PlayBackServerT) checkPlayBackFinish(index int64) (finish bool) {

	if index >= int64(len(p.notes)) {
		atomic.SwapInt64(&p.notesIndex, 0)
		_ = eventCenter.Event.Publish(events.PlayBackFinish, events.PlayBackFinishData{})
		return true
	}

	return false
}
func (p *PlayBackServerT) mouseInput(note *noteT) {
	if err := eventCenter.Event.Publish(events.WindowsMouseInput, events.WindowsMouseInputData{
		Data: &keyMouTool.MouseInputT{},
	}); err != nil {
		p.tryPublishServerError(err)
	}
}
func (p *PlayBackServerT) keyBoardInput(note *noteT) {
	if err := eventCenter.Event.Publish(events.WindowsKeyBoardInput, events.WindowsKeyBoardInputData{
		Data: &keyMouTool.KeyInputT{},
	}); err != nil {
		p.tryPublishServerError(err)
	}
}

// 加载文件
func (p *PlayBackServerT) tryLoadFile(fileName string) {
	if p.name != fileName {
		p.run = false
		p.name = fileName
		p.notes = p.fileControl.ReadFile(fileName)
	}
}

// ----------------------- Util -----------------------

// 发布服务错误事件
func (p *PlayBackServerT) tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}
