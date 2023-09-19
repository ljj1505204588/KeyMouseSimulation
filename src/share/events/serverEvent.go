package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
)

const ServerChange eventCenter.Topic = "server_change"     // 服务变动
const PlayBackFinish eventCenter.Topic = "playback_finish" // 回放结束
const ServerError eventCenter.Topic = "server_error"       // 错误

type ServerChangeData struct {
	Status        enum.Status   // 服务状态
	CurrentTimes  int           // 当前回放次数
	FileNamesData FileNamesData // 文件名称结构体
}
type FileNamesData struct {
	Change    bool     // 是否变动
	FileNames []string // 文件名称
}

type ServerErrorData struct {
	ErrInfo string
}

type PlayBackFinishData struct{}
