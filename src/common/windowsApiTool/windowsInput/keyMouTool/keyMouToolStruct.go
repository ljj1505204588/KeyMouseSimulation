package keyMouTool

type KeyInputChanT struct {
	VK      VKCode
	DwFlags KeyBoardInputDW
}

type MouseInputChanT struct {
	X int32
	Y int32
	DWFlags MouseInputDW
}