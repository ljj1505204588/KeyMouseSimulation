package svcComponent

import (
	eventCenter "KeyMouseSimulation/common/Event"
	gene "KeyMouseSimulation/common/GenTool"
	events2 "KeyMouseSimulation/common/share/events"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	component "KeyMouseSimulation/module/baseComponent"
	"sync/atomic"
	"time"
)

func GetPlaybackServer() PlayBackServerI {
	p := playBackServerT{
		fileControl: component.FileControl,
	}

	p.input = map[keyMouTool.InputType]func(note *keyMouTool.NoteT){
		keyMouTool.TYPE_INPUT_MOUSE:    p.mouseInput,
		keyMouTool.TYPE_INPUT_KEYBOARD: p.keyBoardInput,
	}

	p.registerHandler()

	go p.playBack()
	return &p
}

type PlayBackServerI interface {
	Start(fileName string)
	Pause()
	Stop()
}

/*
*	---------------------------------------------------- PlayBackServerI实现接口 ----------------------------------------------------
 */

type playBackServerT struct {
	run         bool //回放标识
	fileControl component.FileControlI

	name       string             //读取文件名称
	notes      []keyMouTool.NoteT //回放数据
	notesIndex int64              //当前回放在文件内容位置

	speed       float64 // 回放速度
	times       int64   // 回放次数
	remainTimes int64   // 剩余次数

	input map[keyMouTool.InputType]func(note *keyMouTool.NoteT)
}

// Start 开始
func (p *playBackServerT) Start(fileName string) {
	p.tryLoadFile(fileName)
	component.PlaybackConfig.SetPlaybackRemainTimes(p.times)

	p.run = true
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
	defer func() {
		recover()
		go p.playBack()
	}()

	for {
		for p.run {
			var index = atomic.LoadInt64(&p.notesIndex)
			if p.checkPlayBackFinish(index) {
				break
			}

			var n = &p.notes[index]
			p.input[n.NoteType](n)
			p.sleep(n.TimeGap)
			atomic.CompareAndSwapInt64(&p.notesIndex, index, index+1)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
func (p *playBackServerT) checkPlayBackFinish(index int64) bool {
	if index >= int64(len(p.notes)) {
		atomic.SwapInt64(&p.notesIndex, 0)

		if p.remainTimes >= 1 {
			defer component.PlaybackConfig.SetPlaybackRemainTimes(p.remainTimes - 1)

			if p.remainTimes == 1 {
				_ = eventCenter.Event.Publish(events2.PlaybackFinish, events2.PlayBackFinishData{})
				p.run = false
				return true
			}
		}

	}
	return false
}

func (p *playBackServerT) sleep(gap int64) {
	// 性能优化时候看看能不能按位来操作
	var slTime = float64(gap) / p.speed
	time.Sleep(time.Duration(slTime))
}

func (p *playBackServerT) mouseInput(note *keyMouTool.NoteT) {
	if err := eventCenter.Event.Publish(events2.WindowsMouseInput, events2.WindowsMouseInputData{
		Data: &keyMouTool.MouseInputT{
			X:         note.MouseNote.X,
			Y:         note.MouseNote.Y,
			DWFlags:   note.MouseNote.DWFlags,
			MouseData: note.MouseNote.MouseData,
			Time:      note.MouseNote.Time,
		},
	}); err != nil {
		publishServerError(err)
	}
}

func (p *playBackServerT) keyBoardInput(note *keyMouTool.NoteT) {
	if err := eventCenter.Event.Publish(events2.WindowsKeyBoardInput, events2.WindowsKeyBoardInputData{
		Data: &keyMouTool.KeyInputT{
			VK:      note.KeyNote.VK,
			DwFlags: note.KeyNote.DwFlags,
		},
	}); err != nil {
		publishServerError(err)
	}
}

// 注册回调
func (p *playBackServerT) registerHandler() {
	component.PlaybackConfig.SetSpeedChange(true, func(speed float64) {
		p.speed = gene.Choose(speed > 0, speed, 1)
	})
	component.PlaybackConfig.SetPlaybackTimesChange(true, func(times int64) {
		p.times = gene.Choose(times > 0, times, 1)
	})
	component.PlaybackConfig.SetPlaybackRemainTimesChange(false, func(times int64) {
		p.remainTimes = gene.Choose(times >= 0, times, 0)
	})
}

// 加载文件
func (p *playBackServerT) tryLoadFile(fileName string) {
	if p.name != fileName {
		p.run = false
		p.name = fileName
		p.notes = p.fileControl.ReadFile(fileName)
	}
}
