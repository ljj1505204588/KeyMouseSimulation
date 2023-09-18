package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"time"
)

type PlayBackServerI interface {
	Start()
	Pause()
	Stop()

	LoadFile(name string)
	SetSpeed(speed float64)
}

/*
*	PlayBackServerI实现接口
 */

func GetPlaybackServer() *PlayBackServerT {
	p := PlayBackServerT{
		speed: 1,
	}

	//获取key、mouse发送通道
	var err error
	p.keySend, err = keyMouTool.GetKeySendInputChan(3000)
	if err != nil {
		panic(err.Error())
	}
	p.mouseSend, err = keyMouTool.GetMouseSendInputChan(3000)
	if err != nil {
		panic(err.Error())
	}



	return &p
}

type PlayBackServerT struct {
	keySend      chan keyMouTool.KeyInputT   //键盘发送通道
	mouseSend    chan keyMouTool.MouseInputT //鼠标发送通道
	playbackSign bool                        //回放标识

	name        string  //读取文件名称
	notes       []noteT //回放数据
	playbackPod int     //当前回放在文件内容位置
	speed       float64 //回放速度


}

// Start 开始
func (p *PlayBackServerT) Start() {

	//var err error
	//if p.name != name {
	//	p.playbackPod = 0
	//	p.name = name
	//	if p.notes, err = p.loadPlaybackNotes(name); err != nil {
	//		p.publishPlaybackFinish()
	//		p.tryPublishServerError(err)
	//		return
	//	}
	//}

	p.playbackSign = true
	go p.playBack()
	return
}

// Pause 暂停
func (p *PlayBackServerT) Pause() {
	p.playbackSign = false

	logTool.DebugAJ("playback 回放暂停状态")
}

// Stop 停止
func (p *PlayBackServerT) Stop() {
	p.playbackPod = 0
	p.playbackSign = false

	logTool.DebugAJ("playback 退出回放状态")
}

// SetSpeed 设置回放速度
func (p *PlayBackServerT) SetSpeed(speed float64) {
	p.speed = speed
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

				p.keySend <- keyInput
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(p.notes[pod].timeGap / p.speed)))
				mouInput.X = p.notes[pod].MouseNote.X
				mouInput.Y = p.notes[pod].MouseNote.Y
				mouInput.DWFlags = p.notes[pod].MouseNote.DWFlags

				p.mouseSend <- mouInput
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


