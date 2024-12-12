package enum

import (
	"KeyMouseSimulation/module/baseComponent"
)

type Status string

const (
	Free          Status = "Free"
	Recording     Status = "Recording"
	RecordPause   Status = "RecordPause"
	Playback      Status = "Playback"
	PlaybackPause Status = "PlaybackPause"
)

var statusLanguageM = map[Status]component.DisPlay{
	Free:          component.ControlTypeFreeStr,
	Recording:     component.ControlTypeRecordingStr,
	RecordPause:   component.ControlTypeRecordPauseStr,
	Playback:      component.ControlTypePlaybackStr,
	PlaybackPause: component.ControlTypePlaybackPauseStr,
}

func (s Status) Language() string {
	var no = statusLanguageM[s]
	return component.Center.Get(no)
}
