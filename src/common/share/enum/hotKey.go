package enum

import (
	"KeyMouseSimulation/module/baseComponent"
)

type HotKey int

const ( //
	HotKeyRecord   HotKey = iota + 1 //"hot_key_record"
	HotKeyPlayBack                   //"hot_key_playback"
	HotKeyPause                      //"hot_key_pause"
	HotKeyStop                       //"hot_key_stop"
)

func TotalHotkey() []HotKey {
	return []HotKey{
		HotKeyStop,
		HotKeyPause,
		HotKeyRecord,
		HotKeyPlayBack,
	}
}

var hkDisplayM = map[HotKey]component.DisPlay{
	HotKeyStop:     component.StopStr,
	HotKeyPause:    component.PauseStr,
	HotKeyRecord:   component.RecordStr,
	HotKeyPlayBack: component.PlaybackStr,
}

func (h HotKey) Language() component.DisPlay {
	return hkDisplayM[h]
}

var hkDefaultKey = map[HotKey]string{
	HotKeyStop:     "F10",
	HotKeyPause:    "F9",
	HotKeyRecord:   "F7",
	HotKeyPlayBack: "F8",
}

func (h HotKey) DefKey() string {
	return hkDefaultKey[h]
}
