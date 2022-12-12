package recordAndPlayBack

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/events"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type PlayBackServerI interface {
	Start(name string)
	Pause()
	Stop()

	SetSpeed(speed float64)
}

/*
*	PlayBackServerI实现接口
 */

func GetPlaybackServer() *PlayBackServerT {
	p := PlayBackServerT{
		speed: 30,
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

	//获取窗口信息
	p.getWindowRect()

	//初始化

	return &p
}

type PlayBackServerT struct {
	keySend      chan keyMouTool.KeyInputT   //键盘发送通道
	mouseSend    chan keyMouTool.MouseInputT //鼠标发送通道
	playbackSign *bool                       //回放标识

	name        string  //读取文件名称
	notes       []noteT //回放数据
	playbackPod int     //当前回放在文件内容位置
	speed       float64 //回放速度

	windowsX int //电脑屏幕宽度
	windowsY int //电脑屏幕长度
}

// Start 开始
func (p *PlayBackServerT) Start(name string) {
	var err error
	if p.name != name {
		p.playbackPod = 0
		p.name = name
		if p.notes, err = p.loadPlaybackNotes(name); err != nil {
			p.publishPlaybackFinish()
			p.tryPublishServerError(err)
			return
		}
	}

	var sign = true
	p.playbackSign = &sign
	go p.playback()

	return
}

// Pause 暂停
func (p *PlayBackServerT) Pause() {
	*p.playbackSign = false

	logTool.DebugAJ("playback 回放暂停状态")
}

// Stop 停止
func (p *PlayBackServerT) Stop() {
	p.playbackPod = 0
	*p.playbackSign = false

	logTool.DebugAJ("playback 退出回放状态")
}

// SetSpeed 设置回放速度
func (p *PlayBackServerT) SetSpeed(speed float64) {
	p.speed = speed
}

// ----------------------- playback 模块主体功能函数 -----------------------

func (p *PlayBackServerT) playback() {
	sign := p.playbackSign
	for {
		switch {
		case !*sign:
			logTool.DebugAJ("playback 退出回放状态")
			return
		default:
			if p.playbackPod >= len(p.notes) {
				p.playbackPod = 0
				p.publishPlaybackFinish()
				return
			}
			note := p.notes[p.playbackPod]
			switch note.NoteType {
			case keyMouTool.TYPE_INPUT_KEYBOARD:
				time.Sleep(time.Duration(int(note.timeGap / p.speed)))
				p.keySend <- *note.KeyNote
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(note.timeGap / p.speed)))
				p.mouseSend <- *note.MouseNote
			}
			p.playbackPod += 1
		}
	}
}

// loadPlaybackNotes 加载回放记录
func (p *PlayBackServerT) loadPlaybackNotes(name string) ([]noteT, error) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0772)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	nodes := make([]noteT, 100)
	err = json.Unmarshal(b, &nodes)

	for nodePos := range nodes {
		nodes[nodePos].timeGap = float64(nodes[nodePos].TimeGap)
		if nodes[nodePos].NoteType == keyMouTool.TYPE_INPUT_MOUSE {
			nodes[nodePos].MouseNote.X = nodes[nodePos].MouseNote.X * 65535 / int32(p.windowsX)
			nodes[nodePos].MouseNote.Y = nodes[nodePos].MouseNote.Y * 65535 / int32(p.windowsY)
		}
	}

	if err == nil {
		logTool.DebugAJ("playback 加载文件成功：" + "名称:" + name + " 长度：" + strconv.Itoa(len(nodes)))
	}

	return nodes, err
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

// 获取windows窗口大小
func (p *PlayBackServerT) getWindowRect() {
	p.windowsX, p.windowsY = 1920, 1080
	x, _, err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics, windowsApi.SM_CXSCREEN)
	if err != nil {
		return
	}
	y, _, err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics, windowsApi.SM_CYSCREEN)
	if err != nil {
		return
	}
	p.windowsX, p.windowsY = int(x), int(y)
}
