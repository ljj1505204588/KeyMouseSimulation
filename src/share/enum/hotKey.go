package enum

import "KeyMouseSimulation/module/language"

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

var hkDisplayM = map[HotKey]language.DisPlay{
	HotKeyStop:     language.StopStr,
	HotKeyPause:    language.PauseStr,
	HotKeyRecord:   language.RecordStr,
	HotKeyPlayBack: language.PlaybackStr,
}

func (h HotKey) Language() language.DisPlay {
	return hkDisplayM[h]
}
