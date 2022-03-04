package recordAndPlayBack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	windowsApi "KeyMouseSimulation/common/windowsApiTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
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
	go p.playBack()
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

func (p *PlayBackServerT) changeStatus(s ServerStatus) {
	p.status = s
	p.sendMessage(PLAYBACK_EVENT_STATUS_CHANGE, s)
}
func (p *PlayBackServerT) judgeStatus(s ServerStatus) error {
	switch p.status {
	case SERVER_TYPE_FREE:
		if s == SERVER_TYPE_PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		}
	case SERVER_TYPE_PLAYBACK:
		return nil
	case SERVER_TYPE_PLAYBACK_PAUSE:
		return nil
	}
	return nil
}
func (p *PlayBackServerT) Start(name string) error {
	if err := p.judgeStatus(SERVER_TYPE_PLAYBACK); err != nil {
		return err
	}

	p.changeStatus(SERVER_TYPE_PLAYBACK)
	p.name = name
	return nil
}
func (p *PlayBackServerT) Pause() error {
	if err := p.judgeStatus(SERVER_TYPE_PLAYBACK_PAUSE); err != nil {
		return err
	}

	p.changeStatus(SERVER_TYPE_PLAYBACK_PAUSE)
	return nil
}
func (p *PlayBackServerT) Stop() error {
	if err := p.judgeStatus(SERVER_TYPE_FREE); err != nil {
		return err
	}

	p.changeStatus(SERVER_TYPE_FREE)
	return nil
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
func (p *PlayBackServerT) GetPlayBackMessageChan() chan PlaybackMessageT {
	if p.messageChan == nil {
		p.messageChan = make(chan PlaybackMessageT)
	}
	return p.messageChan
}
func (p *PlayBackServerT) SetSpeed(speed float64) {
	p.speed = speed
}
func (p *PlayBackServerT) SetPlaybackTimes(playbackTimes int) {
	if playbackTimes > 0 || playbackTimes == -1 {
		p.playbackTimes = playbackTimes
		p.currentTimes = playbackTimes
		p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
	}
}

// ----------------------- playback 模块主体循环 -----------------------

func (p *PlayBackServerT) playBack() {
	defer func() {
		if info := recover(); info != nil {
			go p.playBack()
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

	//窗口大小
	x,y := int32(p.windowsX),int32(p.windowsY)

	notes := make([]noteT, 100)
	pos := 0
	p.currentTimes = p.playbackTimes
	keyInput, mouInput := keyMouTool.KeyInputChanT{}, keyMouTool.MouseInputChanT{}
	for {
		switch p.status {
		case SERVER_TYPE_FREE:
			//warning GC 协作式抢占，如果没有time.sleep的话，在函数前后的检查抢占信号将不会被发现，GC时候会发生假死现象
			if len(notes) != 0 {
				notes = []noteT{}
				pos = 0
				p.currentTimes = p.playbackTimes
				p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
			}
			time.Sleep(10 * time.Nanosecond)
		case SERVER_TYPE_PLAYBACK:
			if len(notes) == 0 {
				if notes,err = p.playbackNotes(p.name); err != nil {
					p.changeStatus(SERVER_TYPE_FREE)
					p.sendMessage(PLAYBACK_EVENT_ERROR,err.Error())
					continue
				}
				p.currentTimes = p.playbackTimes
			}
			if pos >= len(notes) {
				if p.currentTimes != -1 {
					p.currentTimes--
					p.sendMessage(PLAYBACK_EVENT_CURRENT_TIMES_CHANGE, p.currentTimes)
					if p.currentTimes <= 0 {
						p.changeStatus(SERVER_TYPE_FREE)
					}
				}
				pos = 0
				continue
			}
			switch notes[pos].NoteType {
			case keyMouTool.TYPE_INPUT_KEYBOARD:
				//TODO 考虑这样转化会不会耗时很久
				time.Sleep(time.Duration(int64((float64(notes[pos].TimeGap)) / p.speed)))
				keyInput.VK = notes[pos].KeyNote.Vk
				keyInput.DwFlags = notes[pos].KeyNote.DWFlags

				kSend <- keyInput
			case keyMouTool.TYPE_INPUT_MOUSE:
				time.Sleep(time.Duration(int(float64(notes[pos].TimeGap) / p.speed)))
				mouInput.X = notes[pos].MouseNote.X * 65535 / x
				mouInput.Y = notes[pos].MouseNote.Y * 65535 / y
				mouInput.DWFlags = notes[pos].MouseNote.DWFlags

				mSend <- mouInput
			}
			pos += 1
		case SERVER_TYPE_PLAYBACK_PAUSE:
			time.Sleep(50 * time.Nanosecond)
		default:
			time.Sleep(50 * time.Nanosecond)
		}
	}

}
func (p *PlayBackServerT) playbackNotes(name string) ([]noteT,error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0772)
	if err != nil {
		return nil,err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil,err
	}

	n := make([]noteT, 100)
	_ = json.Unmarshal(b, &n)
	return n,err
}

/*
	获取信息
*/

func (p *PlayBackServerT)GetWindowRect(){
	p.windowsX,p.windowsY =  1920,1080
	x,_,err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics,windowsApi.SM_CXSCREEN)
	if err != nil {
		p.sendMessage(PLAYBACK_EVENT_ERROR,err.Error())
		return
	}
	y,_,err := windowsApi.DllUser.Call(windowsApi.FuncGetSystemMetrics,windowsApi.SM_CYSCREEN)
	if err != nil {
		p.sendMessage(PLAYBACK_EVENT_ERROR,err.Error())
		return
	}
	p.windowsX,p.windowsY =  int(x),int(y)
}