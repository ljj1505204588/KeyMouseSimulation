package events

import (
	eventCenter "KeyMouseSimulation/common/Event"
	"KeyMouseSimulation/share/enum"
)

const ServerStatus eventCenter.Topic = "server_status"     // 服务状态
const PlayBackFinish eventCenter.Topic = "playback_finish" // 回放结束
const RecordFinish eventCenter.Topic = "record_finish"     // 记录结束

const ServerError eventCenter.Topic = "server_error" // 错误
const ButtonClick eventCenter.Topic = "button_click" // 按钮按下

type ServerStatusChangeData struct {
	Status enum.Status // 服务状态
}

type ServerErrorData struct {
	ErrInfo string
}

type PlayBackFinishData struct{}
type RecordFinishData struct{}

type ButtonClickData struct {
	Button enum.Button
	Name   string
}
