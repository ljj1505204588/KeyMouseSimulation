package recordAndPlayBack

//ServerStatus -------------------------------------- 服务状态 --------------------------------------
type ServerStatus int

const (
	SERVER_TYPE_FREE ServerStatus = 1 << iota
	SERVER_TYPE_RECORD
	SERVER_TYPE_RECORD_PAUSE
	SERVER_TYPE_PLAYBACK
	SERVER_TYPE_PLAYBACK_PAUSE
)

func (t ServerStatus) String() (str string) {
	switch t {
	case SERVER_TYPE_FREE:
		str = "SERVER_TYPE_FREE"
	case SERVER_TYPE_RECORD:
		str = "SERVER_TYPE_RECORD"
	case SERVER_TYPE_RECORD_PAUSE:
		str = "SERVER_TYPE_RECORD_PAUSE"
	case SERVER_TYPE_PLAYBACK:
		str = "SERVER_TYPE_PLAYBACK"
	case SERVER_TYPE_PLAYBACK_PAUSE:
		str = "SERVER_TYPE_PLAYBACK_PAUSE"
	}
	return
}

//HotKey -------------------------------------- 热键 --------------------------------------
type HotKey int

func (t HotKey) String() (str string) {
	switch t {
	case HOT_KEY_STOP:
		str = "HOT_KEY_STOP"
	case HOT_KEY_PUASE:
		str = "HOT_KEY_PUASE"
	case HOT_KEY_RECORD_START:
		str = "HOT_KEY_RECORD_START"
	case HOT_KEY_PLAYBACK_START:
		str = "HOT_KEY_PLAYBACK_START"
	}
	return
}

const (
	HOT_KEY_STOP HotKey = iota
	HOT_KEY_PUASE
	HOT_KEY_RECORD_START
	HOT_KEY_PLAYBACK_START
)

//PlaybackEvent -------------------------------------- 回放事件 --------------------------------------
type PlaybackEvent int8

func (t PlaybackEvent) String() (str string) {
	switch t {
	case PLAYBACK_EVENT_STATUS_CHANGE:
		str = "PLAYBACK_EVENT_STATUS_CHANGE"
	case PLAYBACK_EVENT_CURRENT_TIMES_CHANGE:
		str = "PLAYBACK_EVENT_CURRENT_TIMES_CHANGE"
	case PLAYBACK_EVENT_ERROR:
		str = "PLAYBACK_EVENT_ERROR"
	}
	return
}

const (
	PLAYBACK_EVENT_STATUS_CHANGE PlaybackEvent = 1 << iota
	PLAYBACK_EVENT_CURRENT_TIMES_CHANGE
	PLAYBACK_EVENT_ERROR
)

//RecordEvent -------------------------------------- 记录事件 --------------------------------------
type RecordEvent int8

func (t RecordEvent) String() (str string) {
	switch t {
	case RECORD_EVENT_STATUS_CHANGE:
		str = "RECORD_EVENT_STATUS_CHANGE"
	case RECORD_EVENT_HOTKEY_DOWN:
		str = "RECORD_EVENT_HOTKEY_DOWN"
	case RECORD_SAVE_FILE_ERROR:
		str = "RECORD_SAVE_FILE_ERROR"
	}
	return
}

const (
	RECORD_EVENT_STATUS_CHANGE RecordEvent = 1 << iota
	RECORD_EVENT_HOTKEY_DOWN
	RECORD_SAVE_FILE_ERROR
	RECORD_EVENT_ERROR
)
