package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
)

//CurrentTimesChange 回放次数变动
const ServerCurrentTimesChange eventCenter.Topic = "server_current_times_change"

type ServerCurrentTimesChangeData struct {
	CurrentTimes int
}

//StatusChange 回放状态变动
const ServerStatusChange eventCenter.Topic = "server_status_change"

type ServerStatusChangeData struct {
	Status enum.Status
}

//ServerError 回放错误
const ServerError eventCenter.Topic = "server_error"

type ServerErrorData struct {
	ErrInfo string
}

//ServerFileError 文件记录错误
const ServerFileError eventCenter.Topic = "server_file_error"

type ServerFileErrorData struct {
	ErrInfo string
}

//ServerHotKeyDown 热键按下
const ServerHotKeyDown eventCenter.Topic = "server_hot_key_down"

type ServerHotKeyDownData struct {
	HotKey enum.HotKey
}
