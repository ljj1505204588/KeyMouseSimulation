package enum

type Button int32

const (
	RecordButton Button = iota + 1
	PlaybackButton
	PauseButton
	StopButton
)
