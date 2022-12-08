package recordAndPlayBack

import (
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
)

type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   *keyMouTool.KeyInputT
	MouseNote *keyMouTool.MouseInputT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}
