package enum

type ConfEnum string

const (
	RecordHotKeyConf     ConfEnum = "record.hotkey"     // 记录
	PlaybackHotKeyConf   ConfEnum = "playback.hotkey"   // 回放
	PauseHotKeyConf      ConfEnum = "pause.hotkey"      // 暂停
	StopHotKeyConf       ConfEnum = "stop.hotkey"       // 停止
	RecordMouseTrackConf ConfEnum = "record.mouseTrack" // 记录鼠标路径
	PlaybackSpeedConf    ConfEnum = "playback.speed"    // 回放
	PlaybackTimesConf    ConfEnum = "playback.times"    // 回放次数
	LanguageConf         ConfEnum = "system.language"   // 语言设置
)
