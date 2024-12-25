package keyMouTool

import (
	"KeyMouseSimulation/pkg/common/commonTool"
	"KeyMouseSimulation/pkg/common/windowsApiTool/windowsHook"
)

type NoteT struct {
	NoteType  InputType
	KeyNote   *KeyInputT
	MouseNote *MouseInputT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}

type MulNote []NoteT

func (m *MulNote) AppendMouseNote(startTime int64, event *windowsHook.MouseEvent) {
	var dw, exist = mouseDwMap[event.Message]
	if exist {
		var note = NoteT{
			NoteType: TYPE_INPUT_MOUSE,
			MouseNote: &MouseInputT{X: event.X, Y: event.Y,
				DWFlags:   dw,
				MouseData: event.MouseData,
			},
			TimeGap: commonTool.Max(event.RecordTime-startTime, 0),
		}

		*m = append(*m, note)
	}
}

func (m *MulNote) AppendKeyBoardNote(startTime int64, event *windowsHook.KeyboardEvent) {
	var dw, exist = keyDwMap[event.Message]
	if exist {
		var code = VKCode(event.VkCode)

		var note = NoteT{
			NoteType: TYPE_INPUT_KEYBOARD,
			KeyNote:  &KeyInputT{VK: code, DwFlags: dw},
			TimeGap:  commonTool.Max(event.RecordTime-startTime, 0),
		}

		*m = append(*m, note)
	}
}

// AdaptWindow 适应窗口大小
func (m *MulNote) AdaptWindow(x, y int) {

	for nodePos := range *m {
		(*m)[nodePos].timeGap = float64((*m)[nodePos].TimeGap)
		if (*m)[nodePos].NoteType == TYPE_INPUT_MOUSE {
			(*m)[nodePos].MouseNote.X = (*m)[nodePos].MouseNote.X * 65535 / int32(x)
			(*m)[nodePos].MouseNote.Y = (*m)[nodePos].MouseNote.Y * 65535 / int32(y)
		}
	}
}

var mouseDwMap = map[windowsHook.Message]MouseInputDW{
	windowsHook.WM_MOUSEMOVE:         DW_MOUSEEVENTF_MOVE | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTDOWN:     DW_MOUSEEVENTF_LEFTDOWN | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTUP:       DW_MOUSEEVENTF_LEFTUP | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSERIGHTDOWN:    DW_MOUSEEVENTF_RIGHTDOWN | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSERIGHTUP:      DW_MOUSEEVENTF_RIGHTUP | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEMIDDLEDOWN:   DW_MOUSEEVENTF_MIDDLEDOWN | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEMIDDLEUP:     DW_MOUSEEVENTF_MIDDLEUP | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTSICEDOWN: DW_MOUSEEVENTF_XDOWN | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTSICEUP:   DW_MOUSEEVENTF_XUP | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEWHEEL:        DW_MOUSEEVENTF_WHEEL | DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEHWHEEL:       DW_MOUSEEVENTF_HWHEEL | DW_MOUSEEVENTF_ABSOLUTE, //这个暂时不知道是啥，
}
var keyDwMap = map[windowsHook.Message]KeyBoardInputDW{
	windowsHook.WM_KEYDOWN: DW_KEYEVENTF_KEYDown,
	windowsHook.WM_KEYUP:   DW_KEYEVENTF_KEYUP,
}
