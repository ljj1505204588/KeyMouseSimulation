package svcComponent

import (
	"KeyMouseSimulation/common/windowsApi/windowsHook"
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	conf "KeyMouseSimulation/pkg/config"
	eventCenter "KeyMouseSimulation/pkg/event"
	rp_file "KeyMouseSimulation/pkg/file"
	hk "KeyMouseSimulation/pkg/hotkey"
	"KeyMouseSimulation/share/topic"
	"fmt"
	"sync/atomic"
	"time"
)

func GetPlaybackServer() PlayBackServerI {
	p := playBackServerT{}

	p.input = map[keyMouTool.InputType]func(note *keyMouTool.NoteT){
		keyMouTool.TYPE_INPUT_MOUSE:    p.mouseInput,
		keyMouTool.TYPE_INPUT_KEYBOARD: p.keyBoardInput,
	}

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
	run bool //回放标识

	name       string             //读取文件名称
	notes      []keyMouTool.NoteT //回放数据
	notesIndex int64              //当前回放在文件内容位置

	input map[keyMouTool.InputType]func(note *keyMouTool.NoteT)

	playBackTime     int64 // 回放理论时间
	playBackRealTime int64 // 回放实际时间
}

// Start 开始
func (p *playBackServerT) Start(fileName string) {
	p.tryLoadFile(fileName)

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
			atomic.CompareAndSwapInt64(&p.notesIndex, index, index+1)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
func (p *playBackServerT) checkPlayBackFinish(index int64) bool {
	if index >= int64(len(p.notes)) {
		// 时间
		{
			fmt.Println("回放时间：", p.playBackRealTime/int64(time.Millisecond))
			p.playBackTime, p.playBackRealTime = 0, 0
		}

		atomic.SwapInt64(&p.notesIndex, 0)

		var remainTime = conf.PlaybackRemainTimesConf.GetValue()
		if remainTime < 0 {
			return false
		}

		// 剩余次数则减1
		if remainTime >= 1 {
			remainTime--
			conf.PlaybackRemainTimesConf.SetValue(remainTime)
		}

		// 0次则回放结束
		if remainTime == 0 {
			_ = eventCenter.Event.Publish(topic.PlaybackFinish, &topic.PlayBackFinishData{})
			p.run = false
			return true
		}

	}
	return false
}

func (p *playBackServerT) sleep(gap int64, mouseMove bool) {
	var slTime = int64(float64(gap) / conf.PlaybackSpeedConf.GetValue())

	// 如果真实比实际时间大于睡眠时间，则不睡眠
	if mouseMove && (p.playBackRealTime-p.playBackTime) > slTime {
		p.playBackTime += slTime
		return
	}

	var now = time.Now()
	time.Sleep(time.Duration(slTime))

	p.playBackTime += slTime
	p.playBackRealTime += time.Since(now).Nanoseconds()
}

func (p *playBackServerT) mouseInput(note *keyMouTool.NoteT) {
	var isMove = note.GetMouseDW() == windowsHook.WM_MOUSEMOVE

	p.sleep(note.TimeGap, isMove)

	if err := eventCenter.Event.Publish(topic.WindowsMouseInput, note.MouseNote); err != nil {
		publishServerError(err)
	}
}

func (p *playBackServerT) keyBoardInput(note *keyMouTool.NoteT) {
	p.sleep(note.TimeGap, false)

	if hk.Center.IsHotKey(note.KeyNote.VK) {
		return
	}

	if err := eventCenter.Event.Publish(topic.WindowsKeyBoardInput, note.KeyNote); err != nil {
		publishServerError(err)
	}
}

// 加载文件
func (p *playBackServerT) tryLoadFile(fileName string) {
	if p.name != fileName {
		p.run = false
		p.name = fileName
		p.notes = rp_file.FileControl.ReadFile(fileName)
	}
}
