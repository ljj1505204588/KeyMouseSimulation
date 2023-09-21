package enum

type HotKey string

const (
	HotKeyStop          HotKey = "hot_key_stop"
	HotKeyPause         HotKey = "hot_key_pause"
	HotKeyRecordStart   HotKey = "hot_key_record_start"
	HotKeyPlayBackStart HotKey = "hot_key_playback_start"
)
