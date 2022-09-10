package enum

type HotKey int

const (
	HOT_KEY_STOP HotKey = iota
	HOT_KEY_PAUSE
	HOT_KEY_RECORD_START
	HOT_KEY_PLAYBACK_START
)
