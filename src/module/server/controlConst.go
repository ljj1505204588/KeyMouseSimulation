package server

import "KeyMouseSimulation/module/server/recordAndPlayBack"

const (
	FILE_EXT = ".recordPlayback"
)

type ControlStatus int

const (
	CONTROL_TYPE_FREE ControlStatus = 1 << iota
	CONTROL_TYPE_RECORDING
	CONTROL_TYPE_RECORD_PAUSE
	CONTROL_TYPE_PLAYBACK
	CONTROL_TYPE_PLAYBACK_PAUSE
)

type Event int8

const (
	CONTROL_EVENT_STATUS_CHANGE Event = 1 << iota
	CONTROL_EVENT_HOTKEY_DOWN
	CONTROL_EVENT_PLAYBACK_TIMES_CHANGE
	CONTROL_EVENT_SAVE_FILE_ERROR
	CONTROL_EVENT_ERROR
)

type HotKey int

const (
	HOT_KEY_STOP           = HotKey(recordAndPlayBack.HOT_KEY_STOP)
	HOT_KEY_PUASE          = HotKey(recordAndPlayBack.HOT_KEY_PUASE)
	HOT_KEY_RECORD_START   = HotKey(recordAndPlayBack.HOT_KEY_RECORD_START)
	HOT_KEY_PLAYBACK_START = HotKey(recordAndPlayBack.HOT_KEY_PLAYBACK_START)
)
