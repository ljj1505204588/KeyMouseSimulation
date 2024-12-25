package enum

type Status string

const (
	Free          Status = "Free"          // 空闲状态
	Recording     Status = "Recording"     // 记录中
	RecordPause   Status = "RecordPause"   // 记录暂停
	Playback      Status = "Playback"      // 播放中
	PlaybackPause Status = "PlaybackPause" // 播放暂停
)
