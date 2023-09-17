package server

import "KeyMouseSimulation/share/enum"

type ControlI interface {
	Record()   // 记录
	Playback() // 回放
	Pause()    // 暂停
	Stop()     // 停止
	Save()     // 存储

	GetKeyList() (hotKeyList [4]string, keyList []string) // 获取热键列表
	SetHotKey(k enum.HotKey, key string)                  // 设置热键

	SetFileName(fileName string)   // 设置文件名称
	SetSpeed(speed float64)        // 设置播放速度
	SetPlaybackTimes(times int)    // 设置回放次数
	SetIfTrackMouseMove(sign bool) // 设置是否记录鼠标移动
}
