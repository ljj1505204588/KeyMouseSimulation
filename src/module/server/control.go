package server

import (
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/module/server/recordAndPlayBack"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"
)

/*
* 对接UI层
 */

type ControlI interface {
	StartRecord() error
	StartPlayback() error
	Pause() error
	Stop() error

	GetMessageChan() chan MessageT
	GetKeyList() []string
	SetHotKey(k HotKey, key string)
	SetFileName(fileName string)
	SetSpeed(speed float64)
	SetPlaybackTimes(times int)
	SetIfTrackMouseMove(sign bool)
	ScanFile() []string
}

/*
	windows 实现接口
*/

func GetWinControl() *WinControlT {
	c := WinControlT{
		playBack:    recordAndPlayBack.GetPlaybackServer(),
		record:      recordAndPlayBack.GetRecordServer(),
		messageChan: make(chan MessageT, 100),
		keyM:        getKeyM(),
		status:      CONTROL_TYPE_FREE,
		fileName:    "",
	}

	c.SetHotKey(HOT_KEY_RECORD_START, "F7")
	c.SetHotKey(HOT_KEY_PLAYBACK_START, "F8")
	c.SetHotKey(HOT_KEY_PUASE, "F9")
	c.SetHotKey(HOT_KEY_STOP, "F10")

	c.changeStatus(c.status)

	go c.Monitor()

	return &c
}

type WinControlT struct {
	//绑定 playback record 实例
	playBack recordAndPlayBack.PlayBackServerI
	record   recordAndPlayBack.RecordServerI

	//事件 监听通道
	messageChan         chan MessageT
	recordMonitorChan   chan recordAndPlayBack.RecordMessageT
	playbackMonitorChan chan recordAndPlayBack.PlaybackMessageT

	keyM map[string]keyMouTool.VKCode

	status   ControlStatus
	fileName string
}

func (c *WinControlT) judgeStatus(s ControlStatus) error {

	ns := c.status

	switch ns {
	case CONTROL_TYPE_FREE:
		if s == CONTROL_TYPE_RECORD_PAUSE {
			return fmt.Errorf(language.ErrorFreeToRecordPause)
		} else if s == CONTROL_TYPE_PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		} else if s == CONTROL_TYPE_FREE {
			return fmt.Errorf(language.ErrorFreeToFree)
		}
	case CONTROL_TYPE_PLAYBACK:
		if s == CONTROL_TYPE_RECORDING || s == CONTROL_TYPE_RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackToRecordOrRecordPause)
		}
	case CONTROL_TYPE_PLAYBACK_PAUSE:
		if s == CONTROL_TYPE_RECORDING || s == CONTROL_TYPE_RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackPauseToRecordOrRecordPause)
		}
	case CONTROL_TYPE_RECORDING:
		if s == CONTROL_TYPE_PLAYBACK || s == CONTROL_TYPE_PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordToPlaybackOrPlaybackPause)
		}
	case CONTROL_TYPE_RECORD_PAUSE:
		if s == CONTROL_TYPE_PLAYBACK || s == CONTROL_TYPE_PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordPauseToPlaybackOrPlaybackPause)
		}
	}
	fmt.Println("control", s.String())
	return nil
}
func (c *WinControlT) changeStatus(s ControlStatus) {
	c.status = s
	c.sendMessage(CONTROL_EVENT_STATUS_CHANGE, int(c.status), s.String())
}

func (c *WinControlT) StartRecord() error {
	if err := c.judgeStatus(CONTROL_TYPE_RECORDING); err != nil {
		return err
	}

	err := c.record.Start()
	if err != nil {
		return err
	}

	c.changeStatus(CONTROL_TYPE_RECORDING)
	return nil
}
func (c *WinControlT) StartPlayback() error {
	if err := c.judgeStatus(CONTROL_TYPE_PLAYBACK); err != nil {
		return err
	}

	err := c.playBack.Start(c.fileName)
	if err != nil {
		return err
	}

	c.changeStatus(CONTROL_TYPE_PLAYBACK)
	return nil
}
func (c *WinControlT) Pause() error {

	var status = CONTROL_TYPE_FREE
	if c.status == CONTROL_TYPE_RECORDING {
		status = CONTROL_TYPE_RECORD_PAUSE
	} else if c.status == CONTROL_TYPE_PLAYBACK {
		status = CONTROL_TYPE_PLAYBACK_PAUSE
	} else {
		return nil
	}

	var err error
	if err = c.judgeStatus(status); err != nil {
		return err
	}

	switch status {
	case CONTROL_TYPE_PLAYBACK_PAUSE:
		err = c.playBack.Pause()
	case CONTROL_TYPE_RECORD_PAUSE:
		err = c.record.Pause()
	}

	if err != nil {
		return err
	}

	c.changeStatus(status)
	return nil
}
func (c *WinControlT) Stop() error {

	if err := c.judgeStatus(CONTROL_TYPE_FREE); err != nil {
		return err
	}

	switch {
	case c.status == CONTROL_TYPE_RECORDING || c.status == CONTROL_TYPE_RECORD_PAUSE:
		if err := c.record.Stop(c.fileName); err != nil {
			return err
		}
	case c.status == CONTROL_TYPE_PLAYBACK || c.status == CONTROL_TYPE_PLAYBACK_PAUSE:
		if err := c.playBack.Stop(); err != nil {
			return err
		}
	}

	c.changeStatus(CONTROL_TYPE_FREE)
	return nil
}

func (c *WinControlT) Monitor() {
	defer func() {
		if info := recover(); info != nil {
			go func() { c.Monitor() }()
		}
	}()

	if c.playbackMonitorChan == nil {
		c.playbackMonitorChan = c.playBack.GetPlayBackMessageChan()
	}
	if c.recordMonitorChan == nil {
		c.recordMonitorChan = c.record.GetRecordMessageChan()
	}

	lastHotkeyEvenTime := time.Now().UnixNano()
	for {
		select {
		case playbackMsg, ok := <-c.playbackMonitorChan:
			if !ok {
				c.playbackMonitorChan = c.playBack.GetPlayBackMessageChan()
			}
			switch playbackMsg.Event {
			case recordAndPlayBack.PLAYBACK_EVENT_STATUS_CHANGE:
				if value, ok := playbackMsg.Value.(recordAndPlayBack.ServerStatus); ok {
					if c.status == CONTROL_TYPE_PLAYBACK && value == recordAndPlayBack.SERVER_TYPE_FREE {
						c.changeStatus(CONTROL_TYPE_FREE)
					}
				}
			case recordAndPlayBack.PLAYBACK_EVENT_CURRENT_TIMES_CHANGE:
				if value, ok := playbackMsg.Value.(int); ok {
					c.sendMessage(CONTROL_EVENT_PLAYBACK_TIMES_CHANGE, value)
				}
			}
		case recordMsg, ok := <-c.recordMonitorChan:
			if !ok {
				c.recordMonitorChan = c.record.GetRecordMessageChan()
			}
			switch recordMsg.Event {
			case recordAndPlayBack.RECORD_EVENT_STATUS_CHANGE:
				continue
			case recordAndPlayBack.RECORD_SAVE_FILE_ERROR:
				if value, ok := recordMsg.Value.(string); ok {
					c.sendMessage(CONTROL_EVENT_SAVE_FILE_ERROR, value)
				}
			case recordAndPlayBack.RECORD_EVENT_HOTKEY_DOWN:
				if time.Now().UnixNano()-lastHotkeyEvenTime < 200*int64(time.Millisecond) {
					continue
				}
				if value, ok := recordMsg.Value.(recordAndPlayBack.HotKey); ok {
					switch value {
					case recordAndPlayBack.HOT_KEY_PLAYBACK_START:
						c.sendMessage(CONTROL_EVENT_HOTKEY_DOWN, HOT_KEY_PLAYBACK_START)
					case recordAndPlayBack.HOT_KEY_RECORD_START:
						c.sendMessage(CONTROL_EVENT_HOTKEY_DOWN, HOT_KEY_RECORD_START)
					case recordAndPlayBack.HOT_KEY_PAUSE:
						c.sendMessage(CONTROL_EVENT_HOTKEY_DOWN, HOT_KEY_PUASE)
					case recordAndPlayBack.HOT_KEY_STOP:
						c.sendMessage(CONTROL_EVENT_HOTKEY_DOWN, HOT_KEY_STOP)
					}
				}
				lastHotkeyEvenTime = time.Now().UnixNano()
			}
		}
	}
}
func (c *WinControlT) sendMessage(event Event, values ...interface{}) {
	if c.messageChan == nil {
		c.messageChan = make(chan MessageT, 100)
	}
	if len(values) == 1 {
		c.messageChan <- MessageT{
			Event: event,
			Value: values[0],
		}
	} else if len(values) == 2 {
		c.messageChan <- MessageT{
			Event:  event,
			Value:  values[0],
			Value2: values[1],
		}
	}

}
func (c *WinControlT) GetMessageChan() chan MessageT {
	if c.messageChan == nil {
		c.messageChan = make(chan MessageT, 100)
	}
	return c.messageChan
}

func (c *WinControlT) GetKeyList() (result []string) {
	return []string{
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
}
func (c *WinControlT) SetHotKey(hk HotKey, k string) {

	if c.status != CONTROL_TYPE_FREE {
		c.sendMessage(CONTROL_EVENT_ERROR, language.ErrorSetHotkeyInFreeStatusStr)
		return
	}
	if k == "" {
		c.sendMessage(CONTROL_EVENT_ERROR, language.ErrorSetHotkeyWithoutHotkeyStr)
		return
	}

	if err := c.record.SetHotKey(recordAndPlayBack.HotKey(hk), c.keyM[k]); err != nil {
		c.sendMessage(CONTROL_EVENT_ERROR, err.Error())
	}
}
func (c *WinControlT) SetFileName(fileName string) {

	if fileName != "" {
		fileName += FILE_EXT
	}

	c.fileName = fileName
}
func (c *WinControlT) SetSpeed(speed float64) {
	c.playBack.SetSpeed(speed)
}
func (c *WinControlT) SetPlaybackTimes(times int) {
	c.playBack.SetPlaybackTimes(times)
}
func (c *WinControlT) SetIfTrackMouseMove(sign bool) {
	c.record.SetIfTrackMouseMove(sign)
}

func (c *WinControlT) ScanFile() (result []string) {
	if fs, err := ioutil.ReadDir("./"); err != nil {
		c.sendMessage(CONTROL_EVENT_ERROR, err.Error())
		return
	} else {
		for _, v := range fs {
			if filepath.Ext(v.Name()) == FILE_EXT {
				baseName := filepath.Base(v.Name())
				baseName = baseName[:len(baseName)-len(FILE_EXT)]
				result = append(result, baseName)
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})
	return
}

func getKeyM() map[string]keyMouTool.VKCode {
	return map[string]keyMouTool.VKCode{
		"0":   keyMouTool.VK_0,
		"1":   keyMouTool.VK_1,
		"2":   keyMouTool.VK_2,
		"3":   keyMouTool.VK_3,
		"4":   keyMouTool.VK_4,
		"5":   keyMouTool.VK_5,
		"6":   keyMouTool.VK_6,
		"7":   keyMouTool.VK_7,
		"8":   keyMouTool.VK_8,
		"9":   keyMouTool.VK_9,
		"A":   keyMouTool.VK_A,
		"B":   keyMouTool.VK_B,
		"C":   keyMouTool.VK_C,
		"D":   keyMouTool.VK_D,
		"E":   keyMouTool.VK_E,
		"F":   keyMouTool.VK_F,
		"G":   keyMouTool.VK_G,
		"H":   keyMouTool.VK_H,
		"I":   keyMouTool.VK_I,
		"J":   keyMouTool.VK_J,
		"K":   keyMouTool.VK_K,
		"L":   keyMouTool.VK_L,
		"M":   keyMouTool.VK_M,
		"N":   keyMouTool.VK_N,
		"O":   keyMouTool.VK_O,
		"P":   keyMouTool.VK_P,
		"Q":   keyMouTool.VK_Q,
		"R":   keyMouTool.VK_R,
		"S":   keyMouTool.VK_S,
		"T":   keyMouTool.VK_T,
		"U":   keyMouTool.VK_U,
		"V":   keyMouTool.VK_V,
		"W":   keyMouTool.VK_W,
		"X":   keyMouTool.VK_X,
		"Y":   keyMouTool.VK_Y,
		"Z":   keyMouTool.VK_Z,
		"F1":  keyMouTool.VK_F1,
		"F2":  keyMouTool.VK_F2,
		"F3":  keyMouTool.VK_F3,
		"F4":  keyMouTool.VK_F4,
		"F5":  keyMouTool.VK_F5,
		"F6":  keyMouTool.VK_F6,
		"F7":  keyMouTool.VK_F7,
		"F8":  keyMouTool.VK_F8,
		"F9":  keyMouTool.VK_F9,
		"F10": keyMouTool.VK_F10,
		"F11": keyMouTool.VK_F11,
		"F12": keyMouTool.VK_F12,
	}
}
