package topic

import (
	"KeyMouseSimulation/common/windowsApi/windowsInput/keyMouTool"
	"KeyMouseSimulation/share/enum"
)

type Topic string

const ServerStatus Topic = "server_status"     // 服务状态
const PlaybackFinish Topic = "playback_finish" // 回放结束
const RecordFinish Topic = "record_finish"     // 记录结束

const ServerError Topic = "server_error" // 错误
const ButtonClick Topic = "button_click" // 按钮按下

const FileListChange Topic = "file_list_change" // 文件列表变动
const ConfigChange Topic = "config_change"      // 配置改变
const LanguageChange Topic = "language_change"  // 语言改变
const HotKeySet Topic = "hot_key_set"           // 热键设置
const HotKeyEffect Topic = "hot_key_effect"     // 热键触发

// -------------------- 数据 --------------------

type ServerStatusChangeData struct {
	Status enum.Status // 服务状态
}

type ServerErrorData struct {
	ErrInfo string
}

type PlayBackFinishData struct{}

type RecordFinishData struct {
	Notes keyMouTool.MulNote
}

type ButtonClickData struct {
	Button enum.Button
	Name   string
}

type FileListChangeData struct {
	ChooseFile string
	Files      []string
}

type ConfigChangeData struct {
	Key   enum.ConfEnum
	Value any
}

type LanguageChangeData struct {
	Typ enum.LanguageType // 语言类型
}

type HotKeySetData struct {
	Set map[enum.HotKey]keyMouTool.VKCode
}

type HotKeyEffectData struct {
	HotKey enum.HotKey
}
