package enum

type Status string

const (
	FREE           Status = "FREE"
	RECORDING      Status = "RECORDING"
	RECORD_PAUSE   Status = "RECORD_PAUSE"
	PLAYBACK       Status = "PLAYBACK"
	PLAYBACK_PAUSE Status = "PLAYBACK_PAUSE"
)
