package recordAndPlayBack

import "KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"

type BaseT struct {
	PlayBack PlayBackServerI
	Record   RecordServerI
}

type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   *keyMouTool.KeyInputT
	MouseNote *keyMouTool.MouseInputT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}
