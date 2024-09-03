package enum

import "KeyMouseSimulation/module/language"

type Status string

const (
	Free          Status = "Free"
	Recording     Status = "Recording"
	RecordPause   Status = "RecordPause"
	Playback      Status = "Playback"
	PlaybackPause Status = "PlaybackPause"
)

var statusLanguageM = map[Status]language.DisPlay{
	Free:          language.ControlTypeFreeStr,
	Recording:     language.ControlTypeRecordingStr,
	RecordPause:   language.ControlTypeRecordPauseStr,
	Playback:      language.ControlTypePlaybackStr,
	PlaybackPause: language.ControlTypePlaybackPauseStr,
}

func (s Status) Language() string {
	var no = statusLanguageM[s]
	return language.Center.Get(no)
}
