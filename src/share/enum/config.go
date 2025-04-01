package enum

type ConfEnum string

const (
	RecordHotKeyConf     ConfEnum = "record.hotkey"
	PlaybackHotKeyConf   ConfEnum = "playback.hotkey"
	PauseHotKeyConf      ConfEnum = "pause.hotkey"
	StopHotKeyConf       ConfEnum = "stop.hotkey"
	RecordMouseTrackConf ConfEnum = "record.mouseTrack"
	PlaybackSpeedConf    ConfEnum = "playback.speed"
	PlaybackTimesConf    ConfEnum = "playback.times"
	LanguageConf         ConfEnum = "system.language"
)
