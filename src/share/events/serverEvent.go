package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
)

//ServerChange 服务变动
const ServerChange eventCenter.Topic = "server_change"

type ServerChangeData struct {
	Status        enum.Status   // 服务状态
	CurrentTimes  int           // 当前回放次数
	FileNamesData FileNamesData // 文件名称结构体
}
type FileNamesData struct {
	Change    bool     // 是否变动
	FileNames []string // 文件名称
}

//ServerError 错误
const ServerError eventCenter.Topic = "server_error"

type ServerErrorData struct {
	ErrInfo string
}

//ServerHotKeyDown 热键按下
const ServerHotKeyDown eventCenter.Topic = "server_hot_key_down"

type ServerHotKeyDownData struct {
	HotKey enum.HotKey
}

//PlayBackFinish 回放结束
const PlayBackFinish eventCenter.Topic = "playback_finish"

type PlayBackFinishData struct{}

// FileScanNewFile 扫描到新文件
const FileScanNewFile eventCenter.Topic = "file_scan_new_file"

type FileScanNewFileData struct {
	NewFile []string
	FileList []string
}
