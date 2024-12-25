package enum

type Button int32

const (
	RecordButton   Button = iota + 1 // 回放按钮
	PlaybackButton                   // 播放按钮
	PauseButton                      // 暂停按钮
	StopButton                       // 停止按钮
)
