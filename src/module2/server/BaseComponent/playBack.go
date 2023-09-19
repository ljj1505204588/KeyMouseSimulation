package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
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

func GetPlaybackServer() *PlayBackServerT {
	p := PlayBackServerT{}

	return &p
}

type PlayBackServerT struct {
	playbackSign bool //回放标识

	name        string  //读取文件名称
	notes       []noteT //回放数据
	playbackPod int     //当前回放在文件内容位置

	//speed       float64 //回放速度

}

// Start 开始
func (p *PlayBackServerT) Start() {

}

// Pause 暂停
func (p *PlayBackServerT) Pause() {

}

// Stop 停止
func (p *PlayBackServerT) Stop() {
}

// ----------------------- playBack 模块主体功能函数 -----------------------

func (p *PlayBackServerT) playBack() {
	var (
		pod                int
		keyInput, mouInput = keyMouTool.KeyInputT{}, keyMouTool.MouseInputT{}
	)

	for {
		switch {
		case !p.playbackSign:
			logTool.DebugAJ("playback 退出回放状态")
			return
		default:
			if p.playbackPod >= len(p.notes) {
				p.playbackPod = 0
				p.playbackSign = false
				p.publishPlaybackFinish()
				return
			}
			pod = p.playbackPod
			switch p.notes[pod].NoteType {
			case keyMouTool.TYPE_INPUT_KEYBOARD:
				time.Sleep(time.Duration(int(p.notes[pod].timeGap / p.speed)))
				keyInput.VK = p.notes[pod].KeyNote.VK
				keyInput.DwFlags = p.notes[pod].KeyNote.DwFlags

				p.keySend <- &keyInput
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(p.notes[pod].timeGap / p.speed)))
				mouInput.X = p.notes[pod].MouseNote.X
				mouInput.Y = p.notes[pod].MouseNote.Y
				mouInput.DWFlags = p.notes[pod].MouseNote.DWFlags

				p.mouseSend <- &mouInput
			}
			p.playbackPod += 1
		}
	}
}

// ----------------------- Util -----------------------

// 发布回放结束事件
func (p *PlayBackServerT) publishPlaybackFinish() {
	_ = eventCenter.Event.Publish(events.PlayBackFinish, events.PlayBackFinishData{})
}

// 发布服务错误事件
func (p *PlayBackServerT) tryPublishServerError(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{
			ErrInfo: err.Error(),
		})
	}
}
