package server

import (
	"KeyMouseSimulation/module2/server/status"
	"KeyMouseSimulation/share/enum"
)

type control struct {
	status.KMStatusI
}

// GetKeyList 获取热键列表
func (c *control) GetKeyList() (hotKeyList [4]string, keyList []string) {
	return
}

// SetHotKey 设置热键
func (c *control) SetHotKey(k enum.HotKey, key string) {
	return
}

// SetFileName 设置文件名称
func (c *control) SetFileName(fileName string) {
	return
}

// SetSpeed 设置播放速度
func (c *control) SetSpeed(speed float64) {
	return
}

// SetPlaybackTimes 设置回放次数
func (c *control) SetPlaybackTimes(times int) {
	return
}

// SetIfTrackMouseMove 设置是否记录鼠标移
// return 动
func (c *control) SetIfTrackMouseMove(sign bool) {}
