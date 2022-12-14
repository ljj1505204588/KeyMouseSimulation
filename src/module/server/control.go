package server

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/common/logTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/module/server/recordAndPlayBack"
	"KeyMouseSimulation/share/enum"
	"KeyMouseSimulation/share/events"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"time"
)

/*
* 对接UI层
 */

type ControlI interface {
	Record()
	Playback()
	Pause() (save bool)
	Stop(save bool)

	GetKeyList() (hotKeyList [4]string, keyList []string)
	SetHotKey(k enum.HotKey, key string)
	SetFileName(fileName string)
	SetSpeed(speed float64)
	SetPlaybackTimes(times int)
	SetIfTrackMouseMove(sign bool)
}

/*
	windows 实现接口
*/

type WinControlT struct {
	playBack recordAndPlayBack.PlayBackServerI // playback 实例
	record   recordAndPlayBack.RecordServerI   // record 实例

	keyM          map[string]keyMouTool.VKCode // 热键
	status        statusT                      // 状态
	fileName      string                       // 存储文件名称
	playBackTimes int                          // 回放次数
	lastTimes     int                          // 当前剩余次数
}

func NewWinControl() *WinControlT {
	c := WinControlT{
		playBack: recordAndPlayBack.GetPlaybackServer(),
		record:   recordAndPlayBack.GetRecordServer(),
		keyM:     getKeyM(),
		status:   statusT{statusEnum: enum.FREE},
		fileName: "",
	}

	c.SetHotKey(enum.HOT_KEY_RECORD_START, "F7")
	c.SetHotKey(enum.HOT_KEY_PLAYBACK_START, "F8")
	c.SetHotKey(enum.HOT_KEY_PAUSE, "F9")
	c.SetHotKey(enum.HOT_KEY_STOP, "F10")

	c.scanFile()
	eventCenter.Event.Register(events.PlayBackFinish, c.SubPlayBackFinish)

	return &c
}

// --------------------------------------- 基础功能 ---------------------------------------

// Record 记录
func (c *WinControlT) Record() {
	if err := c.changeStatus(enum.RECORDING); err == nil {
		c.record.Start()
	}

	return
}

// Playback 回放
func (c *WinControlT) Playback() {
	logTool.DebugAJ("Control 回放点击")

	if err := c.changeStatus(enum.PLAYBACK); err == nil {
		c.playBack.Start(c.fileName)
	}

	return
}

// Pause 暂停
func (c *WinControlT) Pause() (save bool) {
	logTool.DebugAJ("Control 暂停点击")

	//获取暂停状态
	status, err := c.status.getAfterPauseStatus()
	if err != nil {
		c.tryPublishServerErr(err)
		return
	}

	if err = c.changeStatus(status); err != nil {
		return
	}

	if status == enum.PLAYBACK_PAUSE {
		c.playBack.Pause()
	} else if status == enum.RECORD_PAUSE {
		save = true
		c.record.Pause()
	}

	return
}

// Stop 停止
func (c *WinControlT) Stop(save bool) {
	logTool.DebugAJ("Control 停止点击")

	status := c.status.statusEnum

	//校验 & 改动
	if err := c.changeStatus(enum.FREE); err != nil {
		return
	}

	//修改回放 & 记录状态
	if status == enum.PLAYBACK || status == enum.PLAYBACK_PAUSE {
		c.playBack.Stop()
	}
	if status == enum.RECORDING || status == enum.RECORD_PAUSE {
		c.record.Stop(c.fileName, save)
	}

}

// --------------------------------------- 额外功能 ---------------------------------------

// GetKeyList 获取热键数组
func (c *WinControlT) GetKeyList() (hotKeyList [4]string, keyList []string) {
	return [4]string{"F7", "F8", "F9", "F10"}, []string{
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
}

// SetHotKey 设置热键
func (c *WinControlT) SetHotKey(hk enum.HotKey, k string) {
	if key, exist := c.keyM[k]; exist {
		c.tryPublishServerErr(c.record.SetHotKey(hk, key))
	}
}

// SetFileName 设置存储文件名称
func (c *WinControlT) SetFileName(fileName string) {

	if fileName != "" {
		fileName += FILE_EXT
	}

	c.fileName = fileName
}

// SetSpeed 设置回放速度
func (c *WinControlT) SetSpeed(speed float64) {
	c.playBack.SetSpeed(speed)
}

// SetPlaybackTimes 设置回放次数
func (c *WinControlT) SetPlaybackTimes(times int) {
	c.playBackTimes = times
	c.lastTimes = times
	c.publishServerChange()
}

// SetIfTrackMouseMove 设置是否记录鼠标移动记录
func (c *WinControlT) SetIfTrackMouseMove(sign bool) {
	c.record.SetIfTrackMouseMove(sign)
}

func (c *WinControlT) scanFile() {
	go func() {
		var lastTimeNames []string
		for {
			var names []string
			//遍历存储当前文件名字
			if fs, err := ioutil.ReadDir("./"); err == nil {
				for _, f := range fs {
					if filepath.Ext(f.Name()) == FILE_EXT {
						name := filepath.Base(f.Name())
						name = name[:len(name)-len(FILE_EXT)]
						names = append(names, name)
					}
				}
			}
			if !reflect.DeepEqual(lastTimeNames, names) {
				lastTimeNames = names
				c.publishServerChange(names...)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// --------------------------------------- publishEvent ---------------------------------------

func (c *WinControlT) publishServerChange(fileNames ...string) {
	var fileNamesData events.FileNamesData
	if len(fileNames) != 0 {
		fileNamesData.Change = true
		fileNamesData.FileNames = fileNames
	}

	err := eventCenter.Event.Publish(events.ServerChange, events.ServerChangeData{
		Status:        c.status.statusEnum,
		CurrentTimes:  c.lastTimes,
		FileNamesData: fileNamesData,
	})
	c.tryPublishServerErr(err)
}

func (c *WinControlT) tryPublishServerErr(err error) {
	if err != nil {
		_ = eventCenter.Event.Publish(events.ServerError, events.ServerErrorData{ErrInfo: err.Error()})
	}
}

// --------------------------------------- util ---------------------------------------

func (c *WinControlT) changeStatus(s enum.Status) (err error) {
	if err = c.status.changeStatus(s); err != nil {
		c.tryPublishServerErr(err)
	} else {
		c.publishServerChange()
	}
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

// --------------------------------------- Sub ---------------------------------------

// SubPlayBackFinish 订阅回放结束
func (c *WinControlT) SubPlayBackFinish(_ interface{}) (err error) {
	//无限循环
	if c.playBackTimes == -1 {
		return
	}

	c.lastTimes -= 1
	if c.lastTimes > 0 {
		c.playBack.Start(c.fileName)
	} else {
		c.tryPublishServerErr(c.changeStatus(enum.FREE))
		c.lastTimes = c.playBackTimes
	}

	c.publishServerChange()

	return
}
