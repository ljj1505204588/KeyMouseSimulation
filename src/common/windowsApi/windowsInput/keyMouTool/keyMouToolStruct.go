package keyMouTool

type KeyInputT struct {
	VK      VKCode
	DwFlags KeyBoardInputDW
}

type MouseInputT struct {
	X         int32
	Y         int32
	DWFlags   MouseInputDW
	MouseData uint32
	Time      uint32
}
