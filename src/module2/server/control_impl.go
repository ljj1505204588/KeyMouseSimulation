package server

type ControlI interface {
	Record()              // 记录
	Playback(name string) // 回放
	Pause()               // 暂停
	Stop()                // 停止
	Save(name string)     // 存储
}
