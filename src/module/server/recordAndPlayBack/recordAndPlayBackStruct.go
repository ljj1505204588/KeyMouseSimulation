package recordAndPlayBack

import "KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"

type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   keyBoardNoteT
	MouseNote mouseNoteT
	TimeGap   int64 //Nanosecond
}

type keyBoardNoteT struct {
	Vk      keyMouTool.VKCode
	DWFlags keyMouTool.KeyBoardInputDW
}

type mouseNoteT struct {
	X       int32
	Y       int32
	DWFlags keyMouTool.MouseInputDW
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
