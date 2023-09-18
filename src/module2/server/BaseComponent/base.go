package recordAndPlayBack

import "KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"

const FileExt = ".recordPlayback"

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

type mulNote []noteT

func (m mulNote)adaptWindow(x,y int ){
	for nodePos := range m {
		m[nodePos].timeGap = float64(m[nodePos].TimeGap)
		if m[nodePos].NoteType == keyMouTool.TYPE_INPUT_MOUSE {
			m[nodePos].MouseNote.X = m[nodePos].MouseNote.X * 65535 / int32(x)
			m[nodePos].MouseNote.Y = m[nodePos].MouseNote.Y * 65535 / int32(y)
		}
	}
}