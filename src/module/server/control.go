package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/server/recordAndPlayBack"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"KeyMouseSimulation/share/language"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

/*
* 对接UI层
 */

type ControlI interface {
	StartRecord() error
	StartPlayback() error
	Pause() error
	Stop() error

	GetKeyList() (hotKeyList [4]string, keyList []string)
	SetHotKey(k enum.HotKey, key string)
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
		playBack: recordAndPlayBack.GetPlaybackServer(),
		record:   recordAndPlayBack.GetRecordServer(),
		keyM:     getKeyM(),
		status:   enum.FREE,
		fileName: "",
	}

	c.SetHotKey(enum.HOT_KEY_RECORD_START, "F7")
	c.SetHotKey(enum.HOT_KEY_PLAYBACK_START, "F8")
	c.SetHotKey(enum.HOT_KEY_PAUSE, "F9")
	c.SetHotKey(enum.HOT_KEY_STOP, "F10")

	c.changeStatus(c.status)

	return &c
}

type WinControlT struct {
	//绑定 playback record 实例
	playBack recordAndPlayBack.PlayBackServerI
	record   recordAndPlayBack.RecordServerI

	keyM map[string]keyMouTool.VKCode

	status   enum.Status
	fileName string
}

func (c *WinControlT) checkStatusChange(s enum.Status) error {

	ns := c.status

	switch ns {
	case enum.FREE:
		if s == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorFreeToRecordPause)
		} else if s == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorFreeToPlaybackPause)
		} else if s == enum.FREE {
			return fmt.Errorf(language.ErrorFreeToFree)
		}
	case enum.PLAYBACK:
		if s == enum.RECORDING || s == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackToRecordOrRecordPause)
		}
	case enum.PLAYBACK_PAUSE:
		if s == enum.RECORDING || s == enum.RECORD_PAUSE {
			return fmt.Errorf(language.ErrorPlaybackPauseToRecordOrRecordPause)
		}
	case enum.RECORDING:
		if s == enum.PLAYBACK || s == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordToPlaybackOrPlaybackPause)
		}
	case enum.RECORD_PAUSE:
		if s == enum.PLAYBACK || s == enum.PLAYBACK_PAUSE {
			return fmt.Errorf(language.ErrorRecordPauseToPlaybackOrPlaybackPause)
		}
	}
	fmt.Println("control ", s)
	return nil
}
func (c *WinControlT) changeStatus(s enum.Status) {
	c.status = s

	_ = eventCenter.Event.Publish(events.ServerStatusChange, events.ServerStatusChangeData{Status: s})
}

func (c *WinControlT) StartRecord() error {
	if err := c.checkStatusChange(enum.RECORDING); err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
		return err
	}

	err := c.record.Start()
	if err != nil {
		return err
	}

	c.changeStatus(enum.RECORDING)
	return nil
}
func (c *WinControlT) StartPlayback() error {
	if err := c.checkStatusChange(enum.PLAYBACK); err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})

		return err
	}

	err := c.playBack.Start(c.fileName)
	if err != nil {
		return err
	}

	c.changeStatus(enum.PLAYBACK)
	return nil
}
func (c *WinControlT) Pause() error {

	var status = enum.FREE
	if c.status == enum.RECORDING {
		status = enum.RECORD_PAUSE
	} else if c.status == enum.PLAYBACK {
		status = enum.PLAYBACK_PAUSE
	} else {
		return nil
	}

	var err error
	if err = c.checkStatusChange(status); err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})

		return err
	}

	switch status {
	case enum.PLAYBACK_PAUSE:
		err = c.playBack.Pause()
	case enum.RECORD_PAUSE:
		err = c.record.Pause()
	}

	if err != nil {
		return err
	}

	c.changeStatus(status)
	return nil
}
func (c *WinControlT) Stop() error {

	if err := c.checkStatusChange(enum.FREE); err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})

		return err
	}

	switch {
	case c.status == enum.RECORDING || c.status == enum.RECORD_PAUSE:
		if err := c.record.Stop(c.fileName); err != nil {
			return err
		}
	case c.status == enum.PLAYBACK || c.status == enum.PLAYBACK_PAUSE:
		if err := c.playBack.Stop(); err != nil {
			return err
		}
	}

	c.changeStatus(enum.FREE)
	return nil
}

func (c *WinControlT) GetKeyList() (hotKeyList [4]string, keyList []string) {
	return [4]string{"F7", "F8", "F9", "F10"}, []string{
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
}
func (c *WinControlT) SetHotKey(hk enum.HotKey, k string) {

	if c.status != enum.FREE {
		return
	}
	if k == "" {
		return
	}

	if err := c.record.SetHotKey(hk, c.keyM[k]); err != nil {
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
