package enum

type ConfEnum string

const (
	RecordHotKeyConf    ConfEnum = "record.hotkey"    // 记录
	PlaybackHotKeyConf  ConfEnum = "playback.hotkey"  // 回放
	PauseHotKeyConf     ConfEnum = "pause.hotkey"     // 暂停
	StopHotKeyConf      ConfEnum = "stop.hotkey"      // 停止
	SpeedUpHotKeyConf   ConfEnum = "speedUp.hotkey"   // 加速
	SpeedDownHotKeyConf ConfEnum = "speedDown.hotkey" // 减速

	RecordMouseTrackConf    ConfEnum = "record.mouseTrack"    // 记录鼠标路径
	RecordLenConf           ConfEnum = "record.len"           // 记录长度
	PlaybackSpeedConf       ConfEnum = "playback.speed"       // 回放
	PlaybackTimesConf       ConfEnum = "playback.times"       // 回放次数
	PlaybackRemainTimesConf ConfEnum = "playback.remainTimes" // 剩余回放次数
	LanguageConf            ConfEnum = "system.language"      // 语言设置
)
