package recordAndPlayBack

import "KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"

type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   keyMouTool.KeyInputChanT
	MouseNote keyMouTool.MouseInputChanT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}

// ------------------ 消息 ------------------

type PlaybackMessageT struct {
	Event PlaybackEvent
	Value interface{}
}

type RecordMessageT struct {
	Event RecordEvent
	Value interface{}
}
