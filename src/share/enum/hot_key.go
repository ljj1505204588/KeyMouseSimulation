package enum

type HotKey int32

const ( //
	HotKeyRecord   HotKey = iota + 1 // 记录热键
	HotKeyPlayBack                   // 播放热键
	HotKeyPause                      // 暂停热键
	HotKeyStop                       // 停止热键
)

func TotalHotkey() []HotKey {
	return []HotKey{
		HotKeyRecord,
		HotKeyPlayBack,
		HotKeyPause,
		HotKeyStop,
	}
}
