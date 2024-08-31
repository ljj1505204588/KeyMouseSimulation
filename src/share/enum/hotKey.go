package enum

type HotKey string

const (
	HotKeyStop     HotKey = "hot_key_stop"
	HotKeyPause    HotKey = "hot_key_pause"
	HotKeyRecord   HotKey = "hot_key_record"
	HotKeyPlayBack HotKey = "hot_key_playback"
)

func TotalHotkey() []HotKey {
	return []HotKey{
		HotKeyStop,
		HotKeyPause,
		HotKeyRecord,
		HotKeyPlayBack,
	}
}
