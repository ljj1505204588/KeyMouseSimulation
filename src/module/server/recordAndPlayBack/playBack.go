package recordAndPlayBack

import (
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type PlayBackServerI interface {
	Start(name string) error
	Pause() error
	Stop() error

	SetSpeed(speed float64)
	SetPlaybackTimes(playbackTimes int)
}

/*
*	PlayBackServerI实现接口
 */

func GetPlaybackServer() *PlayBackServerT {
	p := PlayBackServerT{
		speed:         30,
		playbackTimes: 1,
	}

	//获取key、mouse发送通道
	var err error
	p.keySend, err = keyMouTool.GetKeySendInputChan(3000)
	if err != nil {
		os.Exit(1)
	}
	p.mouseSend, err = keyMouTool.GetMouseSendInputChan(3000)
	if err != nil {
		os.Exit(1)
	}

	//获取窗口信息
	p.GetWindowRect()

	//初始化
	p.exit = make(chan struct{})
	go p.free(p.exit)

	return &p
}

type PlayBackServerT struct {
	keySend   chan keyMouTool.KeyInputChanT   //键盘发送通道
	mouseSend chan keyMouTool.MouseInputChanT //鼠标发送通道
	exit      chan struct{}                   //退出通道

	name        string //读取文件名称
	playbackPod int    //当前回放在文件内容位置

	speed         float64 //回放速度
	playbackTimes int     //回放次数
	currentTimes  int     //当前回放次数

	windowsX int //电脑屏幕宽度
	windowsY int //电脑屏幕长度
}

//Start 开始
func (p *PlayBackServerT) Start(name string) (err error) {
	p.name = name

	exit := p.exitNowDeal()
	go p.playback(exit)

	return
}

//Pause 暂停
func (p *PlayBackServerT) Pause() (err error) {

	exit := p.exitNowDeal()
	go p.pause(exit)

	return
}

//Stop 停止
func (p *PlayBackServerT) Stop() (err error) {

	exit := p.exitNowDeal()
	go p.free(exit)

	return
}

//SetSpeed 设置回放速度
func (p *PlayBackServerT) SetSpeed(speed float64) {
	p.speed = speed
}

//SetPlaybackTimes 设置回放次数
func (p *PlayBackServerT) SetPlaybackTimes(playbackTimes int) {
	if playbackTimes > 0 || playbackTimes == -1 {
		p.playbackTimes = playbackTimes
		p.currentTimes = playbackTimes
	}
}

// ----------------------- playback 模块主体功能函数 -----------------------

func (p *PlayBackServerT) free(exit chan struct{}) {
	//回放节点设为初始、回放次数恢复设置次数
	p.playbackPod = 0
	p.currentTimes = p.playbackTimes

	<-exit
	logTool.DebugAJ("playback 退出回放空闲状态")
}
func (p *PlayBackServerT) pause(exit chan struct{}) {
	<-exit
	logTool.DebugAJ("playback 退出回放暂停状态")
}
func (p *PlayBackServerT) playback(exit chan struct{}) {
	var err error
	var notes = make([]noteT, 0)

	if notes, err = p.loadPlaybackNotes(p.name); err != nil || len(notes) == 0 {
		if err != nil {
		}
		<-exit
		return
	}

	for {
		select {
		case <-exit:
			logTool.DebugAJ("playback 退出回放状态")
			return
		default:
			if p.playbackPod >= len(notes) {
				p.playbackPod = 0
				if p.dealPlayBackTimes() {
					<-exit
					return
				}
			}
			pos := p.playbackPod
			switch notes[pos].NoteType {
			case keyMouTool.TYPE_INPUT_KEYBOARD:
				time.Sleep(time.Duration(int(notes[pos].timeGap / p.speed)))
				p.keySend <- *notes[pos].KeyNote
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(notes[pos].timeGap / p.speed)))
				p.mouseSend <- *notes[pos].MouseNote
			}
			p.playbackPod += 1
		}
	}
}
func (p *PlayBackServerT) dealPlayBackTimes() (isReturn bool) {
	if p.currentTimes == -1 {
		return
	}

	p.currentTimes--

	if p.currentTimes <= 0 {
		isReturn = true
	}
	//TODO 发布次数变动事件
	return
}
func (p *PlayBackServerT) exitNowDeal() (exit chan struct{}) {
	if p.exit != nil {
		p.exit <- struct{}{}
	}

	exit = make(chan struct{})
	p.exit = exit

	return
}

//loadPlaybackNotes 加载回放记录
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

/*
	获取信息
*/

func (p *PlayBackServerT) GetWindowRect() {
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
