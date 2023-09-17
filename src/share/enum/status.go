package enum

type Status string

const (
	Free          Status = "Free"
	Recording     Status = "Recording"
	RecordPause   Status = "RecordPause"
	Playback      Status = "Playback"
	PlaybackPause Status = "PlaybackPause"
)
