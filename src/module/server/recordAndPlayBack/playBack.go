package recordAndPlayBack

import (
	"KeyMouseSimulation/common/logTool"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type PlayBackServerI interface {
	Start(name string) error
	Pause() error
	Stop() error

	GetPlayBackMessageChan() chan PlaybackMessageT

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
		status:        SERVER_TYPE_FREE,
		messageChan:   make(chan PlaybackMessageT, 100),
	}

	p.GetWindowRect()
	go p.loop()
	return &p
}

type PlayBackServerT struct {
	name        string //读取文件名称
	playbackPod int    //当前回放在文件内容位置

	speed         float64 //回放速度
	playbackTimes int     //回放次数
	currentTimes  int     //当前回放次数

	windowsX int          //电脑屏幕宽度
	windowsY int          //电脑屏幕长度
	status   ServerStatus //状态

	messageChan chan PlaybackMessageT //消息反馈通道
}

//Start 开始
func (p *PlayBackServerT) Start(name string) error {
	p.name = name

	return p.changeStatus(SERVER_TYPE_PLAYBACK)
}

//Pause 暂停
func (p *PlayBackServerT) Pause() error {
	return p.changeStatus(SERVER_TYPE_PLAYBACK_PAUSE)
}

//Stop 停止
func (p *PlayBackServerT) Stop() error {
	return p.changeStatus(SERVER_TYPE_FREE)
}
func (p *PlayBackServerT) changeStatus(status ServerStatus) error {
	switch p.status {
	case SERVER_TYPE_FREE:
		if status == SERVER_TYPE_PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		}
	case SERVER_TYPE_PLAYBACK:
	case SERVER_TYPE_PLAYBACK_PAUSE:
	}

	p.status = status
	p.sendMessage(PLAYBACK_EVENT_STATUS_CHANGE, status)

	return nil
}

//GetPlayBackMessageChan 获取消息反馈通道
func (p *PlayBackServerT) GetPlayBackMessageChan() chan PlaybackMessageT {
	if p.messageChan == nil {
		p.messageChan = make(chan PlaybackMessageT)
	}
	return p.messageChan
}
func (p *PlayBackServerT) sendMessage(event PlaybackEvent, value interface{}) {
	logTool.DebugAJ(" playback 发送变动消息：", event.String())

	if p.messageChan == nil {
		p.messageChan = make(chan PlaybackMessageT, 100)
	}

	p.messageChan <- PlaybackMessageT{
		Event: event,
		Value: value,
	}
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
		p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
	}
}

// ----------------------- playback 模块主体循环 -----------------------

func (p *PlayBackServerT) loop() {
	defer func() {
		if info := recover(); info != nil {
			go p.loop()
		} else {
			panic("playback 错误退出")
		}
	}()

	//获取key、mouse发送通道
	kSend, err := keyMouTool.GetKeySendInputChan(100)
	if err != nil {
		os.Exit(1)
	}
	mSend, err := keyMouTool.GetMouseSendInputChan(100)
	if err != nil {
		os.Exit(1)
	}

	//主循环所需参数
	exit := make(chan struct{}, 0)
	nowStatus := p.status
	go p.free(exit)

	p.currentTimes = p.playbackTimes
	for {
		if nowStatus != p.status {
			logTool.DebugAJ("playback 状态变动:" + nowStatus.String() + "->" + p.status.String())

			nowStatus = p.status
			exit <- struct{}{}
			switch p.status {
			case SERVER_TYPE_FREE:
				go p.free(exit)
			case SERVER_TYPE_PLAYBACK:
				go p.playback(exit, kSend, mSend)
			case SERVER_TYPE_PLAYBACK_PAUSE:
				go p.pause(exit)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}
func (p *PlayBackServerT) free(exit chan struct{}) {
	//回放节点设为初始、回放次数恢复设置次数
	p.playbackPod = 0
	p.currentTimes = p.playbackTimes
	p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)

	<-exit
	logTool.DebugAJ("playback 退出回放空闲状态")
}
func (p *PlayBackServerT) pause(exit chan struct{}) {
	<-exit // 哈哈进来就等着走！！！
	logTool.DebugAJ("playback 退出回放暂停状态")
}
func (p *PlayBackServerT) playback(exit chan struct{}, kSend chan keyMouTool.KeyInputChanT, mSend chan keyMouTool.MouseInputChanT) {
	var err error
	var notes = make([]noteT, 0)

	if notes, err = p.loadPlaybackNotes(p.name); err != nil || len(notes) == 0 {
		_ = p.changeStatus(SERVER_TYPE_FREE)
		if err != nil {
			p.sendMessage(PLAYBACK_EVENT_ERROR, err.Error())
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
				kSend <- notes[pos].KeyNote
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(notes[pos].timeGap / p.speed)))
				mSend <- notes[pos].MouseNote
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
	p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
	if p.currentTimes <= 0 {
		_ = p.changeStatus(SERVER_TYPE_FREE)
		isReturn = true
	}
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
		p.sendMessage(PLAYBACK_EVENT_ERROR, err.Error())
		return
	}
	y, _, err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics, windowsApi.SM_CYSCREEN)
	if err != nil {
		p.sendMessage(PLAYBACK_EVENT_ERROR, err.Error())
		return
	}
	p.windowsX, p.windowsY = int(x), int(y)
}
