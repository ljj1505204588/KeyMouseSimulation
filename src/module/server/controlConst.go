package server

const (
	FILE_EXT = ".recordPlayback"
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
