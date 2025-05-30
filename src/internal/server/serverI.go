package server

import (
	"KeyMouseSimulation/internal/server/status"
	"KeyMouseSimulation/share/enum"
)

var Svc SvcI = &serverT{
	control: status.NewKmStatusI(),
}

type SvcI interface {
	StatusShow(status enum.Status) string
	Record()           // 记录
	PlayBack()         // 回放
	Pause()            // 暂停
	Stop() (save bool) // 停止
	Save(name string)  // 存储文件
}
