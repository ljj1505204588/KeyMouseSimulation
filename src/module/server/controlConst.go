package server

import (
	"KeyMouseSimulation/module/language"
	"KeyMouseSimulation/module/server/recordAndPlayBack"
)

const (
	FILE_EXT = ".recordPlayback"
)

//ControlStatus -------------------------------------- 控制器状态 --------------------------------------
type ControlStatus int

func (t ControlStatus) String() string {
	switch t {
	case CONTROL_TYPE_FREE:
		return language.ControlTypeFreeStr
	case CONTROL_TYPE_RECORDING:
		return language.ControlTypeRecordingStr
	case CONTROL_TYPE_RECORD_PAUSE:
		return language.ControlTypeRecordPauseStr
	case CONTROL_TYPE_PLAYBACK:
		return language.ControlTypePlaybackStr
	case CONTROL_TYPE_PLAYBACK_PAUSE:
		return language.ControlTypePlaybackPauseStr
	}
	return ""
}

const (
	CONTROL_TYPE_FREE ControlStatus = 1 << iota
	CONTROL_TYPE_RECORDING
	CONTROL_TYPE_RECORD_PAUSE
	CONTROL_TYPE_PLAYBACK
	CONTROL_TYPE_PLAYBACK_PAUSE
)

//Event -------------------------------------- 控制器事件 --------------------------------------
type Event int8

const (
	CONTROL_EVENT_STATUS_CHANGE Event = 1 << iota
	CONTROL_EVENT_HOTKEY_DOWN
	CONTROL_EVENT_PLAYBACK_TIMES_CHANGE
	CONTROL_EVENT_SAVE_FILE_ERROR
	CONTROL_EVENT_ERROR
)

//HotKey -------------------------------------- 控制器热键 --------------------------------------
type HotKey recordAndPlayBack.HotKey

func (t HotKey) String() string {
	return t.String()
}

const (
	HOT_KEY_STOP           = HotKey(recordAndPlayBack.HOT_KEY_STOP)
	HOT_KEY_PUASE          = HotKey(recordAndPlayBack.HOT_KEY_PAUSE)
	HOT_KEY_RECORD_START   = HotKey(recordAndPlayBack.HOT_KEY_RECORD_START)
	HOT_KEY_PLAYBACK_START = HotKey(recordAndPlayBack.HOT_KEY_PLAYBACK_START)
)
