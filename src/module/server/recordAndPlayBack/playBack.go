package recordAndPlayBack

import (
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	name          string
	speed         float64
	playbackTimes int
	currentTimes  int
	windowsX      int
	windowsY      int
	status        ServerStatus

	messageChan chan PlaybackMessageT
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
		return nil
	case SERVER_TYPE_PLAYBACK_PAUSE:
		return nil
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
	ctx, done := context.WithCancel(context.Background())
	nowStatus := p.status

	p.currentTimes = p.playbackTimes
	for {
		if nowStatus != p.status {
			nowStatus = p.status
			done()
			ctx, done = context.WithCancel(context.Background())
			switch p.status {
			case SERVER_TYPE_FREE:
				go p.free(ctx)
			case SERVER_TYPE_PLAYBACK:
				go p.playback(ctx, kSend, mSend)
			case SERVER_TYPE_PLAYBACK_PAUSE:
				go p.pause(ctx)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}
func (p *PlayBackServerT) free(ctx context.Context) {
	p.currentTimes = p.playbackTimes
	p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
	<-ctx.Done()
}
func (p *PlayBackServerT) pause(ctx context.Context) {
	<-ctx.Done() // 哈哈进来就等着走！！！
}
func (p *PlayBackServerT) playback(ctx context.Context, kSend chan keyMouTool.KeyInputChanT, mSend chan keyMouTool.MouseInputChanT) {
	var pos = -1
	var err error
	var notes = make([]noteT, 100)

	if notes, err = p.loadPlaybackNotes(p.name); err != nil || len(notes) == 0 {
		_ = p.changeStatus(SERVER_TYPE_FREE)
		p.sendMessage(PLAYBACK_EVENT_ERROR, err.Error())
		return
	}
	p.currentTimes = p.playbackTimes

	for {
		select {
		//TODO 说实话觉得这样好傻，有啥子更好的办法么
		case <-ctx.Done():
			return
		default:
			pos += 1

			if pos >= len(notes) {
				pos = 0
				p.dealPlayBackTimes()
			}
			switch notes[pos].NoteType {
			case keyMouTool.TYPE_INPUT_KEYBOARD:
				time.Sleep(time.Duration(int(notes[pos].timeGap / p.speed)))
				kSend <- notes[pos].KeyNote
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(notes[pos].timeGap / p.speed)))
				mSend <- notes[pos].MouseNote
			}
		}
	}
}
func (p *PlayBackServerT) dealPlayBackTimes() {
	if p.currentTimes == -1 {
		return
	}
	p.currentTimes--
	p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
	if p.currentTimes <= 0 {
		_ = p.changeStatus(SERVER_TYPE_FREE)
	}
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
