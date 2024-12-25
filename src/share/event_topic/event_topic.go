package event_topic

import (
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
	"KeyMouseSimulation/pkg/event"
	enum2 "KeyMouseSimulation/share/enum"
)

const ServerStatus eventCenter.Topic = "server_status"     // 服务状态
const PlaybackFinish eventCenter.Topic = "playback_finish" // 回放结束
const RecordFinish eventCenter.Topic = "record_finish"     // 记录结束

const ServerError eventCenter.Topic = "server_error" // 错误
const ButtonClick eventCenter.Topic = "button_click" // 按钮按下

const ConfigChange eventCenter.Topic = "config_change"     // 配置改变
const LanguageChange eventCenter.Topic = "language_change" // 语言改变

// -------------------- 数据 --------------------

type ServerStatusChangeData struct {
	Status enum2.Status // 服务状态
}

type ServerErrorData struct {
	ErrInfo string
}

type PlayBackFinishData struct{}
type RecordFinishData struct {
	Notes keyMouTool.MulNote
}

type ButtonClickData struct {
	Button enum2.Button
	Name   string
}

type ConfigChangeData struct {
	Key string
}

type LanguageChangeData struct {
}
