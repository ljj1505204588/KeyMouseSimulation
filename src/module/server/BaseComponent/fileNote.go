package recordAndPlayBack

import (
	"KeyMouseSimulation/common/commonTool"
	"KeyMouseSimulation/common/windowsApiTool/windowsHook"
	"KeyMouseSimulation/common/windowsApiTool/windowsInput/keyMouTool"
)

// todo 看看能不能补充享元
type noteT struct {
	NoteType  keyMouTool.InputType
	KeyNote   *keyMouTool.KeyInputT
	MouseNote *keyMouTool.MouseInputT
	TimeGap   int64 //Nanosecond
	timeGap   float64
}

type MulNote []noteT

// 添加记录
func (m *MulNote) appendMouseNote(startTime int64, event *windowsHook.MouseEvent) {
	var dw, exist = mouseDwMap[event.Message]
	if exist {
		var note = noteT{
			NoteType: keyMouTool.TYPE_INPUT_MOUSE,
			MouseNote: &keyMouTool.MouseInputT{X: event.X, Y: event.Y,
				DWFlags:   dw,
				MouseData: event.MouseData,
			},
			TimeGap: commonTool.Max(event.RecordTime-startTime, 0),
		}

		*m = append(*m, note)
	}
}

// 添加记录
func (m *MulNote) appendKeyBoardNote(startTime int64, event *windowsHook.KeyboardEvent) {
	var dw, exist = keyDwMap[event.Message]
	if exist {
		if _, ok := t.m[keyMouTool.VKCode(event.VkCode)]; ok {
			return
		}

		var note = noteT{
			NoteType: keyMouTool.TYPE_INPUT_KEYBOARD,
			KeyNote: &keyMouTool.KeyInputT{VK: keyMouTool.VKCode(event.VkCode),
				DwFlags: dw,
			},
			TimeGap: commonTool.Max(event.RecordTime-startTime, 0),
		}

		*m = append(*m, note)
	}
}

// 适应窗口大小
func (m *MulNote) adaptWindow(x, y int) {

	for nodePos := range *m {
		(*m)[nodePos].timeGap = float64((*m)[nodePos].TimeGap)
		if (*m)[nodePos].NoteType == keyMouTool.TYPE_INPUT_MOUSE {
			(*m)[nodePos].MouseNote.X = (*m)[nodePos].MouseNote.X * 65535 / int32(x)
			(*m)[nodePos].MouseNote.Y = (*m)[nodePos].MouseNote.Y * 65535 / int32(y)
		}
	}
}

var mouseDwMap = map[windowsHook.Message]keyMouTool.MouseInputDW{
	windowsHook.WM_MOUSEMOVE:         keyMouTool.DW_MOUSEEVENTF_MOVE | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTDOWN:     keyMouTool.DW_MOUSEEVENTF_LEFTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTUP:       keyMouTool.DW_MOUSEEVENTF_LEFTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSERIGHTDOWN:    keyMouTool.DW_MOUSEEVENTF_RIGHTDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSERIGHTUP:      keyMouTool.DW_MOUSEEVENTF_RIGHTUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEMIDDLEDOWN:   keyMouTool.DW_MOUSEEVENTF_MIDDLEDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEMIDDLEUP:     keyMouTool.DW_MOUSEEVENTF_MIDDLEUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTSICEDOWN: keyMouTool.DW_MOUSEEVENTF_XDOWN | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSELEFTSICEUP:   keyMouTool.DW_MOUSEEVENTF_XUP | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEWHEEL:        keyMouTool.DW_MOUSEEVENTF_WHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE,
	windowsHook.WM_MOUSEHWHEEL:       keyMouTool.DW_MOUSEEVENTF_HWHEEL | keyMouTool.DW_MOUSEEVENTF_ABSOLUTE, //这个暂时不知道是啥，
}
var keyDwMap = map[windowsHook.Message]keyMouTool.KeyBoardInputDW{
	windowsHook.WM_KEYDOWN: keyMouTool.DW_KEYEVENTF_KEYDown,
	windowsHook.WM_KEYUP:   keyMouTool.DW_KEYEVENTF_KEYUP,
}
